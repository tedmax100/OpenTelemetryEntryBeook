package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var initResourcesOnce sync.Once
var r *sdkresource.Resource

type LogEntity struct {
	App      string        `json:"app"`
	Duration time.Duration `json:"duration"`
	Status   int32         `json:"status"`
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Set up OpenTelemetry.
	shutdown, err := initProvider()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	go func() {
		for {
			ctx, span := otel.Tracer("Ch11App").Start(context.Background(), "operation")

			duration := time.Duration(rand.Intn(300)+10) * time.Millisecond
			time.Sleep(duration)
			status := int32(200) // 預設為 2xx

			if duration.Milliseconds()%2 == 0 {
				entity := LogEntity{
					App:      "Ch11App",
					Duration: duration,
					Status:   status,
				}
				Logger.InfoContext(ctx, "response",
					zap.String("app", entity.App),
					zap.String("duration", formatDurationWithMS(entity.Duration)),
					zap.Int32("status", entity.Status),
				)
				span.SetAttributes(attribute.KeyValue{Key: "level", Value: attribute.StringValue("INFO")})
				span.SetAttributes(attribute.KeyValue{Key: "app", Value: attribute.StringValue(entity.App)})
				span.SetStatus(codes.Ok, "")
			} else {
				entity := LogEntity{
					App:      "Ch11App",
					Duration: duration,
					Status:   int32(400),
				}
				Logger.ErrorContext(ctx, "err response",
					zap.String("app", entity.App),
					zap.String("duration", formatDurationWithMS(entity.Duration)),
					zap.Int32("status", entity.Status),
				)
				span.SetAttributes(attribute.KeyValue{Key: "level", Value: attribute.StringValue("ERROR")})
				span.SetAttributes(attribute.KeyValue{Key: "app", Value: attribute.StringValue(entity.App)})
				span.SetStatus(codes.Error, "error")
			}
			span.End()
		}
	}()

	go func() {
		for {
			ctx, span := otel.Tracer("Ch11App2").Start(context.Background(), "operation")

			duration := time.Duration(rand.Intn(300)+10) * time.Millisecond
			time.Sleep(duration)
			status := int32(200) // 預設為 2xx

			if duration.Milliseconds()%2 == 0 {
				entity := LogEntity{
					App:      "Ch11App2",
					Duration: duration,
					Status:   status,
				}

				Logger.InfoContext(ctx, "response",
					zap.String("app", entity.App),
					zap.String("duration", formatDurationWithMS(entity.Duration)),
					zap.Int32("status", entity.Status),
				)
				span.SetAttributes(attribute.KeyValue{Key: "level", Value: attribute.StringValue("INFO")})
				span.SetAttributes(attribute.KeyValue{Key: "app", Value: attribute.StringValue(entity.App)})
				span.SetStatus(codes.Ok, "")
			} else {
				entity := LogEntity{
					App:      "Ch11App2",
					Duration: duration,
					Status:   int32(400),
				}

				Logger.ErrorContext(ctx, "err response",
					zap.String("app", entity.App),
					zap.String("duration", formatDurationWithMS(entity.Duration)),
					zap.Int32("status", entity.Status),
				)
				span.SetAttributes(attribute.KeyValue{Key: "level", Value: attribute.StringValue("ERROR")})
				span.SetAttributes(attribute.KeyValue{Key: "app", Value: attribute.StringValue(entity.App)})
				span.SetStatus(codes.Error, "error")
			}
			span.End()
		}
	}()

	<-sigs
}

func formatDurationWithMS(d time.Duration) string {
	return d.String() // `String` method already formats as "72ms", "2h45m", etc.
}

var Logger *otelzap.Logger

func init() {
	// Wrap zap logger to extend Zap with API that accepts a context.Context.
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = EncodeTime
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	logger, err := config.Build(zap.AddCaller(), zap.Fields(zap.String("service", "ch11")))
	if err != nil {
		slog.Error(err.Error())
	}
	logger.Sync()

	Logger = otelzap.New(logger, otelzap.WithTraceIDField(true), otelzap.WithMinLevel(zapcore.DebugLevel))

}

func EncodeTime(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendInt64(t.Unix())
}

func CapitalLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(level.String())
}

func initProvider() (func(context.Context) error, error) {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "otel-collector:4317",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(initResource()),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(newPropagator())

	return tracerProvider.Shutdown, nil
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func initResource() *sdkresource.Resource {
	initResourcesOnce.Do(func() {
		extraResources, _ := sdkresource.New(
			context.Background(),
			sdkresource.WithOS(),
			sdkresource.WithProcess(),
			sdkresource.WithContainer(),
			sdkresource.WithHost(),
			sdkresource.WithAttributes(
				semconv.ServiceName("ch11_app"),
				semconv.ServiceNamespace("ch11"),
			),
		)
		r, _ = sdkresource.Merge(
			sdkresource.Default(),
			extraResources,
		)
	})
	return r
}
