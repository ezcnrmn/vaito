package app

import (
	"database/sql"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/ezcnrmn/vaito/gen/go/user"
	"github.com/ezcnrmn/vaito/services/user/internal/server"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
}

func New(logger *slog.Logger, db *sql.DB) *App {
	s := grpc.NewServer()
	pb.RegisterUserServer(s, server.NewServer(db, logger))

	app := &App{
		log:        logger,
		gRPCServer: s,
	}

	return app
}

func (a *App) Run(listener net.Listener) error {
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		s := <-sigChan

		a.log.Info("shutting down user service", "signal", s.String())

		a.gRPCServer.GracefulStop()
	}()

	err := a.gRPCServer.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}
