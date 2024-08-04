package otel

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitOtel(ctx context.Context, serviceName string, collectorTarget string, logger *slog.Logger) {
	conn, err := grpc.NewClient(collectorTarget,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error("failed to create gRPC connection to collector", slog.Any("error", err))
		os.Exit(1)
	}

	serviceSemconv := semconv.ServiceNameKey.String(serviceName)

	res, err := resource.New(ctx,
		resource.WithContainer(),
		resource.WithOS(),
		resource.WithContainerID(),
		resource.WithTelemetrySDK(),
		resource.WithProcessCommandArgs(),
		resource.WithAttributes(
			// The service name used to display traces in backends
			serviceSemconv,
		),
	)
	if err != nil {
		logger.Error("failed to new resource", slog.Any("error", err))
		os.Exit(1)
	}

	shutdownTracerProvider, err := initTracerProvider(ctx, res, conn)
	if err != nil {
		logger.Error("failed to initTracerProvider", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		if err := shutdownTracerProvider(ctx); err != nil {
			logger.Error("failed to shutdown TracerProvider", slog.Any("error", err))
		}
	}()
}

func GetTracer(name string) trace.Tracer {
	return otel.Tracer(name)
}

func SpanSetError(span trace.Span, err error) {
	span.SetStatus(codes.Error, err.Error())
}

func SpanSetOk(span trace.Span, description string) {
	span.SetStatus(codes.Ok, description)
}
