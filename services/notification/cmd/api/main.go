package main

import (
	"log/slog"
	"os"

	"github.com/ezcnrmn/vaito/services/notification/internal/app"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	rabbitmqUrl := os.Getenv("RABBITMQ_URL")
	if rabbitmqUrl == "" {
		logger.Error("empty rabbitmq url")
		os.Exit(1)
	}

	conn, err := amqp.Dial(rabbitmqUrl)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer channel.Close()

	app := app.New(logger, channel)

	err = app.DeclareQueues()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	err = app.Run()
	if err != nil {
		logger.Error(err.Error())
	}
}
