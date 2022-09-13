package rabbitmq

type Producer interface {
	Send(data []byte) error
	Close() error
}
