package amqp

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewAqmpConn(uri, serviceName string) (*amqp.Channel, []func() error, error) {
	deferFuncs := make([]func() error, 0)
	config := amqp.Config{
		Vhost:      "/",
		Properties: amqp.NewConnectionProperties(),
	}
	config.Properties.SetClientConnectionName("producer-with-confirms")

	log.Printf("producer: dialing %s", uri)
	conn, err := amqp.DialConfig(uri, config)
	if err != nil {
		return nil, deferFuncs, fmt.Errorf("producer: error in dial: %w", err)
	}
	deferFuncs = append(deferFuncs, func() error {
		return conn.Close()
	})

	ch, err := conn.Channel()
	if err != nil {
		return nil, deferFuncs, fmt.Errorf("error getting a channel: %w", err)
	}
	deferFuncs = append(deferFuncs, func() error {
		return ch.Close()
	})

	_, err = ch.QueueDeclare(
		"demo", // name
		true,   // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		return nil, deferFuncs, fmt.Errorf("failed to declare a queue: %w", err)
	}

	return ch, deferFuncs, nil
}

func PublishWithCtx(ctx context.Context, channel *amqp.Channel, body string) error {
	return channel.PublishWithContext(ctx,
		"",     // exchange
		"demo", // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
}
