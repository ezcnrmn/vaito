package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/ezcnrmn/vaito/services/notification/internal/consumer"
)

type App struct {
	log    *slog.Logger
	amqpCh *amqp.Channel
	queues struct {
		email *amqp.Queue
	}
	consumer *consumer.Consumer
}

func New(logger *slog.Logger, amqpChannel *amqp.Channel) *App {
	app := &App{
		log:      logger,
		amqpCh:   amqpChannel,
		consumer: consumer.New(logger),
	}

	return app
}

func (a *App) DeclareQueues() error {
	emailQueue, err := a.amqpCh.QueueDeclare(
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

	a.queues.email = &emailQueue

	return nil
}

func (a *App) Run() error {
	consumerTag := "notification-consumers"

	emails, err := a.amqpCh.Consume(
		a.queues.email.Name,
		consumerTag,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	wg := sync.WaitGroup{}

	a.log.Info("starting notification service")

	go a.consumer.ConsumeEmails(&wg, emails)

	<-ctx.Done()
	a.log.Info("shutting down notification service")

	err = a.amqpCh.Cancel(consumerTag, false)
	if err != nil {
		return nil
	}

	wg.Wait()

	return nil
}
