package consumer

import "log/slog"

type Consumer struct {
	log *slog.Logger
}

func New(logger *slog.Logger) *Consumer {
	return &Consumer{log: logger}
}
