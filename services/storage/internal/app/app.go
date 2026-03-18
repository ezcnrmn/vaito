package app

import (
	"database/sql"
	"log/slog"
	"net"

	pb "github.com/ezcnrmn/vaito/gen/go/storage"
	"github.com/ezcnrmn/vaito/services/storage/internal/server"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
}

func New(logger *slog.Logger, db *sql.DB) *App {
	s := grpc.NewServer()
	pb.RegisterListingServer(s, server.NewListingServer(db))
	pb.RegisterUserServer(s, server.NewUserServer(db))

	app := &App{
		log:        logger,
		gRPCServer: s,
	}

	return app
}

func (a *App) Run(listener net.Listener) error {
	err := a.gRPCServer.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}
