package notification

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EmailType string

const (
	ListingVisibilityChanged EmailType = "listing-visibility-changed"
)

type VisibilityEmailMessage struct {
	Type       EmailType `json:"type"`
	UserEmail  string    `json:"userEmail"`
	ListingID  int64     `json:"listingId"`
	Visibility bool      `json:"visibility"`
}

func (n *Notification) PublishVisibilityChanged(ctx context.Context, userEmail string, listingID int64, visibility bool) error {
	body, err := json.Marshal(VisibilityEmailMessage{Type: ListingVisibilityChanged, UserEmail: userEmail, ListingID: listingID, Visibility: visibility})
	if err != nil {
		return err
	}

	err = n.channel.PublishWithContext(ctx,
		"",
		n.emailQueue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	return err
}
