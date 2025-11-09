package messaging

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMq struct {
	conn   *amqp091.Connection
	Chanel *amqp091.Channel
}
type MessageHandlers func(context.Context, amqp091.Delivery) error

func NewRabbitMq(uri string) (*RabbitMq, error) {
	if uri == "" {
		return nil, errors.New("rabbitmq uri is empty")
	}
	// trim spaces
	uri = strings.TrimSpace(uri)
	if !strings.HasPrefix(uri, "amqp://") && !strings.HasPrefix(uri, "amqps://") {
		if strings.Contains(uri, "@") || strings.Contains(uri, ":") {
			uri = "amqp://" + uri
		} else {
			return nil, fmt.Errorf("AMQP scheme must be either 'amqp://' or 'amqps://' - got: %q", uri)
		}
	}

	conn, err := amqp091.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("failed connect to RabbitMq: %v", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("fail to create a chanel: %v", err)
	}
	rmq := &RabbitMq{
		conn:   conn,
		Chanel: ch,
	}
	if err := rmq.setupExchangesAndQueues(); err != nil {
		//Clean up if setup fails
		rmq.Close()
		return nil, fmt.Errorf("failed to setup exchanges and queues: %v", err)
	}
	return rmq, nil
}
func (r *RabbitMq) ConsumeMessages(queueName string, handler MessageHandlers) error {
	msgs, err := r.Chanel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	ctx := context.Background()

	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %s", msg.Body)

			if err := handler(ctx, msg); err != nil {
				log.Fatalf("failed to handle the message: %v", err)
			}
		}
	}()

	return nil
}

func (r *RabbitMq) setupExchangesAndQueues() error {
	_, err := r.Chanel.QueueDeclare(
		"hello",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *RabbitMq) PublishMessage(ctx context.Context, routingKey string, message string) error {
	return r.Chanel.PublishWithContext(ctx,
		"",         // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp091.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(message),
			DeliveryMode: amqp091.Persistent,
		})
}

func (r *RabbitMq) Close() {
	if r.Chanel != nil {
		if err := r.Chanel.Close(); err != nil {
			log.Printf("warning: error closing rabbitmq channel: %v", err)
		}
	}
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			log.Printf("warning: error closing rabbitmq connection: %v", err)
		}
	}
}
