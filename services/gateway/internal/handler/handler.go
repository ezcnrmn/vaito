package handler

import (
	"gateway/internal/config"
	"log/slog"
)

type Handler struct {
	cfg *config.Config
	log *slog.Logger
}

func New(config *config.Config, logger *slog.Logger) *Handler {
	return &Handler{
		cfg: config,
		log: logger,
	}
}
