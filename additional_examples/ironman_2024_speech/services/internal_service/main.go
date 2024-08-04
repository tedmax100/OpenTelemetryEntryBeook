package main

import (
	"context"
	"log/slog"

	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	logger "demo/internal/logger"
	otel "demo/internal/otel"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/baggage"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

var SERVICE_NAME = "internal-service"
var serviceName = semconv.ServiceNameKey.String(SERVICE_NAME)
var conn *grpc.ClientConn

var _logger *slog.Logger

func init() {
	_logger = logger.GetLogger().With("gin_mode", gin.EnvGinMode).With("service", SERVICE_NAME)
}

func main() {
	_logger.InfoContext(context.Background(), "init grpc connection...")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	otel.InitOtel(ctx, SERVICE_NAME, "otelcol:4317", _logger)

	r := gin.Default()
	r.Use(otelgin.Middleware(SERVICE_NAME))
	r.GET("/echo", func(c *gin.Context) {
		tracer := otel.GetTracer(SERVICE_NAME)
		_, span := tracer.Start(c.Request.Context(), c.HandlerName())
		defer span.End()
		bagMsg := baggage.FromContext(c.Request.Context())
		baggageItems := make(map[string]string)
		for _, member := range bagMsg.Members() {
			baggageItems[member.Key()] = member.Value()
		}

		otel.SpanSetOk(span, "")
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
			"baggage": baggageItems,
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
