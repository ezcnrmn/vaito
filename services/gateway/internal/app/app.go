package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/ezcnrmn/vaito/gen/go/storage"
	"github.com/ezcnrmn/vaito/services/gateway/internal/config"
)

type App struct {
	cfg        *config.Config
	log        *slog.Logger
	httpServer *http.Server
	grpc       struct {
		user    pb.UserClient
		listing pb.ListingClient
	}
}

func New(
	config *config.Config,
	logger *slog.Logger,
	userConn pb.UserClient,
	listingConn pb.ListingClient,
) *App {
	app := &App{
		cfg: config,
		log: logger,
		grpc: struct {
			user    pb.UserClient
			listing pb.ListingClient
		}{
			user:    userConn,
			listing: listingConn,
		},
	}

	routes := app.routes()

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: routes,
	}
	app.httpServer = httpServer

	return app
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

		err := a.httpServer.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		shutdownError <- nil
	}()

	a.log.Info("starting gateway app", "port", a.cfg.Port)

	err := a.httpServer.ListenAndServe()
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
