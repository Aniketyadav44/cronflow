package config

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

func loadRabbitMQ(rUser, rPass, rHost, rPort, rVHost string) (*amqp091.Connection, *amqp091.Channel, error) {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s%s", rUser, rPass, rHost, rPort, rVHost)

	conn, err := amqp091.Dial(dsn)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, nil
}
