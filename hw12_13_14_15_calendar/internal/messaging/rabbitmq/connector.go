package rabbitmq

import (
	"io"

	"github.com/streadway/amqp"
)

type Connector interface {
	Consumer
	Producer
	io.Closer
}

type connector struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	exchange   string
	queue      amqp.Queue
	consumeCh  <-chan amqp.Delivery
}

func NewConnector(connectionString string, exchange string, queueName string) (Connector, error) {
	connection, err := amqp.Dial(connectionString)
	if err != nil {
		return nil, err
	}
	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}
	err = channel.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	queue, err := channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &connector{exchange: exchange, connection: connection, channel: channel, queue: queue}, nil
}

func (c *connector) Consume() (<-chan amqp.Delivery, error) {
	err := c.channel.QueueBind(c.queue.Name, "", c.exchange, false, nil)
	if err != nil {
		return nil, err
	}

	c.consumeCh, err = c.channel.Consume(c.queue.Name, "", false, false, false, false, nil)

	return c.consumeCh, err
}

func (c *connector) Send(data []byte) error {
	err := c.channel.Publish(c.exchange, "", false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         data,
	})

	return err
}

func (c *connector) Close() error {
	if c.connection != nil {
		return c.connection.Close()
	}
	return nil
}
