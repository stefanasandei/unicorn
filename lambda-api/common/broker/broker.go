package broker

// a common interface for a message broker
// to be implemented for rabbitmq, sqs of whatever we may need

type DeliveryMessage struct {
	Body string
}

type MessageBroker interface {
	Connect(string) error

	CreateQueue(string) error

	Consume(string) (<-chan DeliveryMessage, error)

	SendMessageToQueue(string, string) error

	Close() error
}
