package broker

import (
	"context"
	"errors"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQMessageBroker struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func (broker *RabbitMQMessageBroker) Connect(connectionUrl string) error {
	var err error

	// start connection
	broker.conn, err = amqp.Dial(connectionUrl)
	if err != nil {
		log.Panicf("Failed to connect to RabbitMQ: %s", err)
		return err
	}

	// create a channel
	broker.channel, err = broker.conn.Channel()
	if err != nil {
		log.Panicf("Failed to open a channel: %s", err)
		return err
	}

	return nil
}

func (broker *RabbitMQMessageBroker) CreateQueue(queueName string) error {
	var err error

	if broker.channel == nil {
		return errors.New("channel is nil")
	}

	// create a queue
	broker.queue, err = broker.channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		log.Panicf("Failed to declare a queue: %s", err)
		return err
	}

	return nil
}

func (broker *RabbitMQMessageBroker) Consume(queueName string) (<-chan DeliveryMessage, error) {
	if broker.channel == nil {
		return make(<-chan DeliveryMessage), errors.New("channel is nil")
	}

	// create a queue consumer
	msgs, err := broker.channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)

	if err != nil {
		log.Panicf("Failed to register a consumer: %s", err)
		return make(<-chan DeliveryMessage), err
	}

	dMsgs := convertToDeliveryMessage(msgs)

	return dMsgs, err
}

func (broker *RabbitMQMessageBroker) SendMessageToQueue(queueName, message string) error {
	// TODO we shouldn't open and close the channel within this function
	ch, err := broker.conn.Channel()
	if err != nil {
		log.Panicf("Failed to open a channel: %s", err)
		return err
	}

	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			log.Panicf("Error closing the rabbitmq channel: %s", err)
		}
	}(ch)

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Panicf("Failed to declare a queue: %s", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

	if err != nil {
		log.Panicf("Failed to publish to the exchange: %s", err)
		return err
	}

	return nil
}

func (broker *RabbitMQMessageBroker) Close() error {
	defer func(channel *amqp.Channel) {
		err := channel.Close()
		if err != nil {
			log.Panicf("Error closing the rabbitmq channel: %s", err)
		}
	}(broker.channel)

	return broker.conn.Close()
}

func convertToDeliveryMessage(in <-chan amqp.Delivery) <-chan DeliveryMessage {
	out := make(chan DeliveryMessage)

	go func() {
		defer close(out)

		for msg := range in {
			deliveryMsg := DeliveryMessage{
				Body: string(msg.Body),
			}
			out <- deliveryMsg
		}
	}()

	return out
}
