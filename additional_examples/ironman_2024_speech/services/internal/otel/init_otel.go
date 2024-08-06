package otel

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitOtel(ctx context.Context, serviceName string) {
	conn, err := initConn()
	if err != nil {
		log.Fatal("failed to create gRPC connection to collector")
		os.Exit(1)
	}

	res := newResource(ctx, serviceName)
	_, err = initTracerProvider(ctx, res, conn)
	if err != nil {
		//logger.Error("failed to initTracerProvider", slog.Any("error", err))
		os.Exit(1)
	}

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

func newResource(ctx context.Context, serviceName string) *resource.Resource {
	hostName, _ := os.Hostname()
	//serviceSemconv := semconv.ServiceNameKey.String(serviceName)

	res, err := resource.New(ctx,
		resource.WithContainer(),
		resource.WithOS(),
		resource.WithProcessRuntimeVersion(),
		resource.WithTelemetrySDK(),
		resource.WithProcessCommandArgs(),
		resource.WithAttributes(
			// The service name used to display traces in backends
			//serviceSemconv,
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersion("1.0.0"),
			semconv.HostName(hostName),
		),
	)
	if err != nil {
		//logger.Error("failed to new resource", slog.Any("error", err))
		os.Exit(1)
	}

	return res
}

// Initialize a gRPC connection to be used by both the tracer and meter
// providers.
func initConn() (*grpc.ClientConn, error) {
	// It connects the OpenTelemetry Collector through local gRPC connection.
	// You may replace `localhost:4317` with your endpoint.
	conn, err := grpc.NewClient("otelcol:4317",
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	return conn, err
}
