package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"os"
	"os/signal"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
)

var serviceName = semconv.ServiceNameKey.String("worker")
var conn *grpc.ClientConn

// var kafkaProducerClient  sarama.AsyncProducer
var log *logrus.Logger

func init() {
	log = logrus.New()
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	}
	log.Out = os.Stdout
}

func main() {
	log.Printf("Waiting for connection...")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Printf("init grpc connection...")
	conn, err := grpc.NewClient("otelcol:4317",
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create gRPC connection to collector: %w", err))
		// return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithContainer(),
		resource.WithOS(),
		resource.WithContainerID(),
		resource.WithTelemetrySDK(),
		resource.WithProcessCommandArgs(),
		resource.WithAttributes(
			// The service name used to display traces in backends
			serviceName,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	shutdownTracerProvider, err := initTracerProvider(ctx, res, conn)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdownTracerProvider(ctx); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %s", err)
		}
	}()

	brokerList := strings.Split("kafka:9092", ",")
	log.Printf("Kafka brokers: %s", strings.Join(brokerList, ", "))

	ctx, cancel = signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := startConsumerGroup(ctx, brokerList, log); err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()

}

// Initializes an OTLP exporter, and configures the corresponding trace provider.
func initTracerProvider(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) (func(context.Context) error, error) {
	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// Set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider.Shutdown, nil
}

func createProducerSpan(ctx context.Context, msg *sarama.ProducerMessage) trace.Span {
	tracer := otel.Tracer("workder")
	spanContext, span := tracer.Start(
		ctx,
		fmt.Sprintf("%s publish", msg.Topic),
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(
			semconv.PeerService("kafka"),
			semconv.NetworkTransportTCP,
			semconv.MessagingSystemKafka,
			semconv.MessagingDestinationName(msg.Topic),
			semconv.MessagingOperationPublish,
			semconv.MessagingKafkaDestinationPartition(int(msg.Partition)),
		),
	)

	carrier := propagation.MapCarrier{}
	propagator := otel.GetTextMapPropagator()
	propagator.Inject(spanContext, carrier)

	for key, value := range carrier {
		msg.Headers = append(msg.Headers, sarama.RecordHeader{Key: []byte(key), Value: []byte(value)})
	}

	return span
}

var (
	Topic           = "tasks"
	ProtocolVersion = sarama.V3_0_0_0
	GroupID         = "worker"
)

func startConsumerGroup(ctx context.Context, brokers []string, log *logrus.Logger) error {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = ProtocolVersion
	// So we can know the partition and offset of messages.
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Consumer.Interceptors = []sarama.ConsumerInterceptor{NewOTelInterceptor(GroupID)}

	consumerGroup, err := sarama.NewConsumerGroup(brokers, GroupID, saramaConfig)
	if err != nil {
		return err
	}

	handler := groupHandler{
		log: log,
	}

	err = consumerGroup.Consume(ctx, []string{Topic}, &handler)
	if err != nil {
		return err
	}
	return nil
}

type groupHandler struct {
	log *logrus.Logger
}

func (g *groupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (g *groupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (g *groupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			headerStrings := make([]string, len(message.Headers))
			for i, header := range message.Headers {
				headerStrings[i] = fmt.Sprintf("%s=%s", header.Key, header.Value)
			}
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s, headers = %v", message.Value, message.Timestamp, message.Topic, headerStrings)
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

type OTelInterceptor struct {
	tracer     trace.Tracer
	fixedAttrs []attribute.KeyValue
}

// NewOTelInterceptor processes span for intercepted messages and add some
// headers with the span data.
func NewOTelInterceptor(groupID string) *OTelInterceptor {
	oi := OTelInterceptor{}
	oi.tracer = otel.Tracer("github.com/open-telemetry/opentelemetry-demo/accountingservice/sarama")

	oi.fixedAttrs = []attribute.KeyValue{
		semconv.MessagingSystemKafka,
		semconv.MessagingKafkaConsumerGroup(groupID),
		semconv.NetTransportTCP,
	}
	return &oi
}

func (oi *OTelInterceptor) OnConsume(msg *sarama.ConsumerMessage) {
	headers := propagation.MapCarrier{}

	for _, recordHeader := range msg.Headers {
		headers[string(recordHeader.Key)] = string(recordHeader.Value)
	}

	propagator := otel.GetTextMapPropagator()
	ctx := propagator.Extract(context.Background(), headers)

	_, span := oi.tracer.Start(
		ctx,
		fmt.Sprintf("%s receive", msg.Topic),
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(oi.fixedAttrs...),
		trace.WithAttributes(
			semconv.MessagingDestinationName(msg.Topic),
			semconv.MessagingKafkaMessageOffset(int(msg.Offset)),
			semconv.MessagingMessageBodySize(len(msg.Value)),
			semconv.MessagingOperationReceive,
			semconv.MessagingKafkaDestinationPartition(int(msg.Partition)),
		),
	)
	defer span.End()
}
