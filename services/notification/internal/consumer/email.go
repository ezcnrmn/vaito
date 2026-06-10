package consumer

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EmailType string

const (
	ListingVisibilityChanged EmailType = "listing-visibility-changed"
)

// SendEmail Симулирует отправку электронного сообщения, выводя в консоль
func (c Consumer) ConsumeEmails(wg *sync.WaitGroup, emails <-chan amqp.Delivery) {
	wg.Add(1)
	defer wg.Done()

	for msg := range emails {
		var email map[string]any

		err := json.Unmarshal(msg.Body, &email)
		if err != nil {
			c.log.Error("unable to decode message", "err", err.Error(), "msg", email)
			continue
		}

		rawType, ok := email["type"]
		if !ok {
			c.log.Error("message doesn't have `type` field", "msg", email)
			continue
		}

		emailType, ok := rawType.(string)
		if !ok {
			c.log.Error("unable to convert `type` field to `EmailType`", "type", rawType, "msg", email)
			continue
		}

		err = c.proceedEmailByType(EmailType(emailType), email)
		if err != nil {
			c.log.Error("unable to proceed message", "err", err.Error(), "msg", email)
		}
	}
}

func (c Consumer) proceedEmailByType(emailType EmailType, msg map[string]any) error {
	switch emailType {
	case ListingVisibilityChanged:
		return c.sendVisibilityEmail(msg)
	default:
		return errors.New("unknown message type")
	}
}

func (c Consumer) sendVisibilityEmail(msg map[string]any) error {
	rawUserEmail, emailOk := msg["userEmail"]
	rawListingID, listingIDOk := msg["listingId"]
	rawVisibility, visibilityOk := msg["visibility"]

	if !emailOk || !listingIDOk || !visibilityOk {
		return errors.New("message doesn't have one or many fields: `userEmail`, `listingId`, `visibility`")
	}

	userEmail, ok := rawUserEmail.(string)
	if !ok {
		return errors.New("unable to convert message's `userEmail` field to string")
	}

	floatListingID, ok := rawListingID.(float64)
	if !ok {
		return errors.New("unable to convert message's `listingId` field to float64")
	}

	listingID := int64(floatListingID)

	visibility, ok := rawVisibility.(bool)
	if !ok {
		return errors.New("unable to convert message's `visibility` field to bool")
	}

	c.log.Info(fmt.Sprintf("Sending email about changed listing visibility by moderation. To:%q ListingID:%d NewVisibility:%t", userEmail, listingID, visibility))

	return nil
}
