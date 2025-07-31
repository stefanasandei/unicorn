package model

import (
	"sync"

	"github.com/common/broker"
)

type ReplyQueue struct {
	Broker        broker.MessageBroker
	msgs          <-chan broker.DeliveryMessage
	messageBuffer []string
	msgMutex      sync.Mutex
	cond          *sync.Cond
}

func NewReplyQueue(b broker.MessageBroker) (*ReplyQueue, error) {
	r := &ReplyQueue{
		Broker:        b,
		messageBuffer: make([]string, 0), // Pass the mutex to NewCond
	}

	r.cond = sync.NewCond(&r.msgMutex)

	if err := r.Setup(); err != nil {
		return nil, err
	}

	go r.consumeMessages()

	return r, nil
}

func (r *ReplyQueue) Setup() error {
	queueName := "reply"

	// the broker should be already connected

	err := r.Broker.CreateQueue(queueName)
	if err != nil {
		return err
	}

	r.msgs, err = r.Broker.Consume(queueName)
	if err != nil {
		return err
	}

	return nil
}

func (r *ReplyQueue) consumeMessages() {
	for msg := range r.msgs {
		r.msgMutex.Lock()
		r.messageBuffer = append(r.messageBuffer, msg.Body)
		r.msgMutex.Unlock()
		r.cond.Signal() // Signal that a new message is available
	}
}

func (r *ReplyQueue) GetNewMessage() (string, error) {
	r.msgMutex.Lock()

	for len(r.messageBuffer) == 0 {
		// TODO: don't block the main thread
		r.cond.Wait() // Wait for a new message
	}

	// Return the latest unhandled message (first in the buffer)
	msg := r.messageBuffer[0]

	// Remove the message from the buffer
	r.messageBuffer = r.messageBuffer[1:]

	r.msgMutex.Unlock()

	return msg, nil
}

func (r *ReplyQueue) Close() error {
	return r.Broker.Close()
}
