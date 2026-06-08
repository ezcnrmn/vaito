package notification

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Notification struct {
	channel    *amqp.Channel
	emailQueue amqp.Queue
}

func NewNotification(channel *amqp.Channel) (*Notification, error) {
	n := &Notification{channel: channel}

	err := n.declareQueues()
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (n *Notification) declareQueues() error {
	emailQueue, err := n.channel.QueueDeclare(
		"email",
		true,
		false,
		false,
		false,
		amqp.Table{
			amqp.QueueTypeArg: amqp.QueueTypeQuorum,
		},
	)
	if err != nil {
		return err
	}

	n.emailQueue = emailQueue

	return nil
}
