package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	sloggin "github.com/samber/slog-gin"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	//amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	logger "demo/internal/logger"
	otel "demo/internal/otel"

	"go.opentelemetry.io/otel/baggage"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

var serviceName = semconv.ServiceNameKey.String("api-server")
var conn *grpc.ClientConn

var (
	uri          = flag.String("uri", "amqp://guest:guest@localhost:5672/", "AMQP URI")
	exchange     = flag.String("exchange", "test-exchange", "Durable AMQP exchange name")
	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
	queue        = flag.String("queue", "test-queue", "Ephemeral AMQP queue name")
	routingKey   = flag.String("key", "test-key", "AMQP routing key")
	body         = flag.String("body", "foobar", "Body of message")
	continuous   = flag.Bool("continuous", false, "Keep publishing messages at a 1msg/sec rate")
)

var _logger *slog.Logger

func init() {
	_logger = logger.GetLogger().With("gin_mode", gin.EnvGinMode).With("service", "api")
}

func main() {
	_logger.InfoContext(context.Background(), "init grpc connection...")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	otel.InitOtel(ctx, "api", "otelcol:4317", _logger)

	r := gin.Default()
	r.Use(otelgin.Middleware("api"))
	config := sloggin.Config{
		WithRequestID: true,
		WithSpanID:    true,
		WithTraceID:   true,
	}
	r.Use(sloggin.NewWithConfig(_logger, config))

	r.GET("/hello", func(c *gin.Context) {
		tracer := otel.GetTracer("api-server")
		ctx, span := tracer.Start(c.Request.Context(), "hello")
		defer span.End()

		member, err := baggage.NewMember("user", "nathan")
		if err != nil {
			_logger.Error("failed to NewMember", slog.Any("error", err))
			otel.SpanSetError(span, err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		bag, err := baggage.New(member)
		if err != nil {
			_logger.Error("failed to create baggage", slog.Any("error", err))
			otel.SpanSetError(span, err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

		ctx = baggage.ContextWithBaggage(ctx, bag)
		req, err := http.NewRequestWithContext(ctx, "GET", "http://internal_service:3000/echo?message=Hello_Nathan", nil)
		if err != nil {
			_logger.Error("failed to create a request", slog.Any("error", err))
			otel.SpanSetError(span, err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		_, err = client.Do(req)
		if err != nil {
			var opErr *net.OpError
			if errors.As(err, &opErr) {
				_logger.Error("failed to link to server", slog.Any("error", err))
			} else {
				_logger.Error("failed to request to echo", slog.Any("error", err))
			}

			otel.SpanSetError(span, err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		otel.SpanSetOk(span, "")
		c.JSON(http.StatusOK, gin.H{
			"message": "success	",
		})
	})
	/*
		r.GET("gentasks", func(c *gin.Context) {
			message := []byte("Hello, 雷N!")
			msg := sarama.ProducerMessage{
				Topic: "tasks", // 替換成你的 Kafka 主題名稱
				Value: sarama.ByteEncoder(message),
			}

			// Inject tracing info into message
			span := createProducerSpan(c.Request.Context(), &msg)
			defer span.End()

			config := amqp.Config{
				Vhost:      "/",
				Properties: amqp.NewConnectionProperties(),
			}
			config.Properties.SetClientConnectionName("producer-with-confirms")
			conn, err := amqp.DialConfig(*uri, config)
			if err != nil {
				ErrLog.Fatalf("producer: error in dial: %s", err)
			}
			defer conn.Close()

			// Send message and handle response
			startTime := time.Now()
			select {
			case KafkaProducerClient.Input() <- &msg:
				log.Infof("Message sent to Kafka: %v", msg)
				select {
				case successMsg := <-KafkaProducerClient.Successes():
					span.SetAttributes(
						attribute.Bool("messaging.kafka.producer.success", true),
						attribute.Int("messaging.kafka.producer.duration_ms", int(time.Since(startTime).Milliseconds())),
						attribute.KeyValue(semconv.MessagingKafkaMessageOffset(int(successMsg.Offset))),
					)
					log.Infof("Successful to write message. offset: %v, duration: %v", successMsg.Offset, time.Since(startTime))
					span.SetStatus(codes.Ok, "")
					c.JSON(http.StatusOK, gin.H{
						"message": "success",
					})
					return
				case errMsg := <-KafkaProducerClient.Errors():
					span.SetAttributes(
						attribute.Bool("messaging.kafka.producer.success", false),
						attribute.Int("messaging.kafka.producer.duration_ms", int(time.Since(startTime).Milliseconds())),
					)
					span.SetStatus(codes.Error, errMsg.Err.Error())
					log.Errorf("Failed to write message: %v", errMsg.Err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": errMsg.Err.Error(),
					})
					return
				case <-ctx.Done():
					span.SetAttributes(
						attribute.Bool("messaging.kafka.producer.success", false),
						attribute.Int("messaging.kafka.producer.duration_ms", int(time.Since(startTime).Milliseconds())),
					)
					span.SetStatus(codes.Error, "Context cancelled: "+ctx.Err().Error())
					log.Warnf("Context canceled before success message received: %v", ctx.Err())
					c.JSON(http.StatusRequestTimeout, gin.H{
						"message": "Request timeout",
					})
					return
				}
			case <-ctx.Done():
				span.SetAttributes(
					attribute.Bool("messaging.kafka.producer.success", false),
					attribute.Int("messaging.kafka.producer.duration_ms", int(time.Since(startTime).Milliseconds())),
				)
				span.SetStatus(codes.Error, "Failed to send: "+ctx.Err().Error())
				log.Errorf("Failed to send message to Kafka within context deadline: %v", ctx.Err())
				c.JSON(http.StatusRequestTimeout, gin.H{
					"message": "Request timeout",
				})
				return
			}
		})
	*/
	go func() {
		r.Run()
	}()
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	<-signalCh

	//r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

/* // Initializes an OTLP exporter, and configures the corresponding trace provider.
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
*/
