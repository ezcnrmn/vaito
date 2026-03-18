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
	"google.golang.org/grpc"
)

type App struct {
	port       string
	log        *slog.Logger
	httpServer *http.Server
	storage    struct {
		user    pb.UserClient
		listing pb.ListingClient
	}
}

func New(port string, logger *slog.Logger, grpcClient *grpc.ClientConn) *App {
	userConn := pb.NewUserClient(grpcClient)
	listingConn := pb.NewListingClient(grpcClient)

	app := &App{
		port: port,
		log:  logger,
	}
	app.storage.user = userConn
	app.storage.listing = listingConn

	routes := app.routes()

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
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

	a.log.Info("starting gateway app", "port", a.port)

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
