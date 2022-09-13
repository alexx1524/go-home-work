package rabbitmq

import "github.com/streadway/amqp"

type Consumer interface {
	Consume() (<-chan amqp.Delivery, error)
	Close() error
}
