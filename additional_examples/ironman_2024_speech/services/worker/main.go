package main

import (
	"context"
	"errors"
	"log"
	"time"

	"os"
	"os/signal"
	"syscall"

	"demo/internal/amqp"
	otel "demo/internal/otel"

	"math/rand"

	"github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var SERVICE_NAME = "worker"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	otel.InitOtel(ctx, SERVICE_NAME)

	channel, deferFuncs, err := amqp.NewAqmpConn("amqp://demo:demo@rabbitmq:5672/", SERVICE_NAME)
	if err != nil {
		for _, deferFunc := range deferFuncs {
			deferFunc()
		}
		log.Fatal(err)
	}

	meter := otel.GetMeter(SERVICE_NAME)
	taskCount, err := meter.Int64Counter("task_total", metric.WithDescription("The number of task"))
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := channel.Consume(
		"demo", // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal("Failed to register a consumer")
	}

	go func() {
		for d := range msgs {
			go func(msg amqp091.Delivery) {
				// Extract headers
				workCtx := otel.ExtractAMQPHeaders(context.Background(), msg.Headers)
				taskCount.Add(workCtx, 1)
				bagMsg := baggage.FromContext(workCtx)

				tracer := otel.GetTracer(SERVICE_NAME)
				_, workSpan := tracer.Start(workCtx, "work", trace.WithSpanKind(trace.SpanKindConsumer))
				for _, member := range bagMsg.Members() {
					workSpan.SetAttributes(attribute.String(member.Key(), member.Value()))
				}

				randomValue := rand.Intn(200)

				time.Sleep(time.Duration(randomValue) * time.Millisecond)
				msg.Ack(false)

				randomErrorValue := rand.Intn(10)

				if randomErrorValue < 3 {
					otel.SpanSetError(workSpan, errors.New("worker error"))
				} else {
					otel.SpanSetOk(workSpan, "")
				}

				workSpan.End()
			}(d)

		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	<-signalCh

}
