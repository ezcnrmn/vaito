package app

import (
	"context"
	"errors"
	"fmt"
	"gateway/internal/config"
	"gateway/internal/handler"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
)

type App struct {
	cfg    *config.Config
	log    *slog.Logger
	server *http.Server
}

func New(config *config.Config, logger *slog.Logger) *App {
	routes := getRoutes(config, logger)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: routes,
	}

	return &App{
		cfg:    config,
		log:    logger,
		server: server,
	}
}

func (a *App) Run() error {
	shutdownError := make(chan error)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		s := <-sigChan

		a.log.Info("shutting down gateway app", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := a.server.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		shutdownError <- nil
	}()

	a.log.Info("starting gateway app", "port", a.cfg.Port)

	err := a.server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	a.log.Error("gateway app stopped")
	return nil
}

func getRoutes(config *config.Config, logger *slog.Logger) *httprouter.Router {

	handler := handler.New(config, logger)
	routes := httprouter.New()

	routes.HandlerFunc(http.MethodGet, "/v1/healthcheck", handler.Healthcheck)

	routes.HandlerFunc(http.MethodPost, "/v1/users", handler.CreateUser)
	routes.HandlerFunc(http.MethodPut, "/v1/users/:id", handler.UpdateUser)
	routes.HandlerFunc(http.MethodGet, "/v1/users/authenticate", handler.Authenticate)

	routes.HandlerFunc(http.MethodPost, "/v1/listings", handler.CreateListings)
	routes.HandlerFunc(http.MethodGet, "/v1/listings/:id", handler.ShowListings)
	routes.HandlerFunc(http.MethodPut, "/v1/listings/:id", handler.UpdateListings)
	routes.HandlerFunc(http.MethodDelete, "/v1/listings/:id", handler.DeleteListings)
	routes.HandlerFunc(http.MethodGet, "/v1/listings", handler.ListListings)

	return routes
}
