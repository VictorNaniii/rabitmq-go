package messaging

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMq struct {
	conn *amqp091.Connection
}

func NewRabbitMq(uri string) (*RabbitMq, error) {
	if uri == "" {
		return nil, errors.New("rabbitmq uri is empty")
	}
	// trim spaces
	uri = strings.TrimSpace(uri)
	// If scheme is missing but contains @ or :, try to auto-prefix amqp://
	if !strings.HasPrefix(uri, "amqp://") && !strings.HasPrefix(uri, "amqps://") {
		// helpful fallback: if it looks like host[:port] or user:pass@host format, prepend amqp://
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
	return &RabbitMq{
		conn: conn,
	}, nil
}

func (r *RabbitMq) Close() {
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			log.Printf("warning: error closing rabbitmq connection: %v", err)
		}
	}
}
