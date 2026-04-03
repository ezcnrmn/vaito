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
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
}

func New(logger *slog.Logger, db *sql.DB) *App {
	app := &App{log: logger}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(app.recoverPanic),
	)
	pb.RegisterUserServiceServer(s, server.NewServer(db, logger))

	health := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, health)

	app.gRPCServer = s

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
