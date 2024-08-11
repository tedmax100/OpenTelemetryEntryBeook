package main

import (
	"context"
	"fmt"
	"log"

	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	//amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"demo/internal/amqp"
	otel "demo/internal/otel"

	"go.opentelemetry.io/otel/baggage"
)

var SERVICE_NAME = "api-server"
var conn *grpc.ClientConn

var (
	uri          = flag.String("uri", "amqp://guest:guest@rabbitmq:5672/", "AMQP URI")
	exchange     = flag.String("exchange", "test-exchange", "Durable AMQP exchange name")
	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
	queue        = flag.String("queue", "test-queue", "Ephemeral AMQP queue name")
	routingKey   = flag.String("key", "test-key", "AMQP routing key")
	body         = flag.String("body", "foobar", "Body of message")
	continuous   = flag.Bool("continuous", false, "Keep publishing messages at a 1msg/sec rate")
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	channel, deferFuncs, err := amqp.NewAqmpConn("amqp://demo:demo@localhost:5672/", SERVICE_NAME)
	if err != nil {
		for _, deferFunc := range deferFuncs {
			deferFunc()
		}
		log.Fatal(err)
	}
	otel.InitOtel(ctx, SERVICE_NAME)

	r := gin.Default()
	r.Use(otelgin.Middleware("api"))

	r.GET("/hello", func(c *gin.Context) {
		tracer := otel.GetTracer("api-server")
		ctx, span := tracer.Start(c.Request.Context(), "hello")
		defer func() {
			span.End()
		}()

		member, err := baggage.NewMember("user", "nathan")
		if err != nil {
			otel.SpanSetError(span, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		bag, err := baggage.New(member)
		if err != nil {
			otel.SpanSetError(span, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

		ctx = baggage.ContextWithBaggage(ctx, bag)
		req, err := http.NewRequestWithContext(ctx, "GET", "http://internal_service:8080/echo?message=Hello_Nathan", nil)
		if err != nil {
			otel.SpanSetError(span, err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		res, err := client.Do(req)
		if err != nil {
			otel.SpanSetError(span, err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			otel.SpanSetError(span, fmt.Errorf("unexpected status code: %d", res.StatusCode))

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("unexpected status code: %d", res.StatusCode),
			})
			return
		}

		err = amqp.PublishWithCtx(ctx, channel, "demo")
		if err != nil {
			otel.SpanSetError(span, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		otel.SpanSetOk(span, "")
		c.JSON(http.StatusOK, gin.H{
			"message": "success	",
		})
	})

	go func() {
		r.Run()
	}()
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	<-signalCh

}
