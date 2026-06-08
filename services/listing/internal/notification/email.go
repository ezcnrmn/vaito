package notification

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (n *Notification) PublishVisibilityChanged(ctx context.Context, visibility bool) error {
	body := fmt.Sprintf(`{"visibility":%t}`, visibility)

	err := n.channel.PublishWithContext(ctx,
		"",
		n.emailQueue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})

	return err
}
