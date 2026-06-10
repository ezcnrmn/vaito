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

	pbListing "github.com/ezcnrmn/vaito/gen/go/listing"
	pbUser "github.com/ezcnrmn/vaito/gen/go/user"
	"github.com/ezcnrmn/vaito/services/gateway/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type App struct {
	cfg config.Config

	log        *slog.Logger
	httpServer *http.Server

	services struct {
		user    pbUser.UserServiceClient
		listing pbListing.ListingServiceClient

		health struct {
			user    grpc_health_v1.HealthClient
			listing grpc_health_v1.HealthClient
		}
	}
}

func New(config config.Config, logger *slog.Logger, userClient, listingClient *grpc.ClientConn) *App {
	userConn := pbUser.NewUserServiceClient(userClient)
	listingConn := pbListing.NewListingServiceClient(listingClient)

	app := &App{
		cfg: config,
		log: logger,
	}
	app.services.user = userConn
	app.services.listing = listingConn
	app.services.health.user = grpc_health_v1.NewHealthClient(userClient)
	app.services.health.listing = grpc_health_v1.NewHealthClient(listingClient)

	routes := app.routes()

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.cfg.Port),
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

		a.log.Info("shutting down gateway service", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := a.httpServer.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		shutdownError <- nil
	}()

	a.log.Info("starting gateway service", "port", a.cfg.Port)

	err := a.httpServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	a.log.Error("gateway service stopped")
	return nil
}
