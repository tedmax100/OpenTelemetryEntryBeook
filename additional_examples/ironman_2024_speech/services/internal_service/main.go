package main

import (
	"context"
	"errors"
	"time"

	"net/http"
	"os"
	"os/signal"

	otel "demo/internal/otel"

	"math/rand"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
)

var SERVICE_NAME = "internal-service"
var conn *grpc.ClientConn

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	otel.InitOtel(ctx, SERVICE_NAME)

	r := gin.Default()
	r.Use(otelgin.Middleware("api"))
	r.GET("/echo", func(c *gin.Context) {
		tracer := otel.GetTracer(SERVICE_NAME)
		_, span := tracer.Start(c.Request.Context(), c.HandlerName())
		defer span.End()
		bagMsg := baggage.FromContext(c.Request.Context())
		baggageItems := make(map[string]string)
		for _, member := range bagMsg.Members() {
			// baggageItems[member.Key()] = member.Value()
			span.SetAttributes(attribute.String(member.Key(), member.Value()))
		}

		rand.Seed(time.Now().UnixNano())
		randomValue := rand.Intn(10)

		if randomValue < 3 {
			span.SetStatus(codes.Error, "random error occurred")
			otel.SpanSetError(span, errors.New("random error"))
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "random error occurred",
			})
			return
		}

		otel.SpanSetOk(span, "")
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
			"baggage": baggageItems,
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
