package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"demo/internal/amqp"
	otel "demo/internal/otel"

	"go.opentelemetry.io/otel/attribute"

	"math/rand"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/metric"
)

var SERVICE_NAME = "api-server"

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	channel, deferFuncs, err := amqp.NewAqmpConn("amqp://demo:demo@rabbitmq:5672/", SERVICE_NAME)
	if err != nil {
		for _, deferFunc := range deferFuncs {
			deferFunc()
		}
		log.Fatal(err)
	}
	otel.InitOtel(ctx, SERVICE_NAME)

	r := gin.Default()
	r.Use(otelgin.Middleware(SERVICE_NAME))

	meter := otel.GetMeter(SERVICE_NAME)
	reqCount, err := meter.Int64Counter("request_total", metric.WithDescription("The number of access API"))
	if err != nil {
		log.Fatal(err)
	}

	r.GET("/hello", func(c *gin.Context) {
		commonAttrs := []attribute.KeyValue{
			attribute.String("path", c.Request.RequestURI),
			attribute.String("method", c.Request.Method),
		}
		reqCount.Add(ctx, 1, metric.WithAttributes(commonAttrs...))

		tracer := otel.GetTracer(SERVICE_NAME)
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

		randomSleepValue := rand.Intn(50)

		time.Sleep(time.Duration(randomSleepValue) * time.Millisecond)

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
		headers := otel.InjectAMQPHeaders(ctx)
		err = amqp.PublishWithCtx(ctx, channel, "demo", headers)
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
