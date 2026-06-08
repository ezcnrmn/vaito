package app

import (
	"database/sql"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/ezcnrmn/vaito/gen/go/listing"
	pbUser "github.com/ezcnrmn/vaito/gen/go/user"
	"github.com/ezcnrmn/vaito/services/listing/internal/notification"
	"github.com/ezcnrmn/vaito/services/listing/internal/server"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	services   struct {
		user pbUser.UserServiceClient
	}
	notification *notification.Notification
}

func New(logger *slog.Logger, db *sql.DB, userClientConn *grpc.ClientConn, amqpChannel *amqp.Channel) (*App, error) {
	app := &App{log: logger}

	notification, err := notification.NewNotification(amqpChannel)
	if err != nil {
		return nil, err
	}

	user := pbUser.NewUserServiceClient(userClientConn)

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(app.recoverPanic),
	)
	pb.RegisterListingServiceServer(s, server.NewServer(db, logger, user, notification))

	health := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, health)

	app.gRPCServer = s
	app.services.user = user

	return app, nil
}

func (a *App) Run(listener net.Listener) error {
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		s := <-sigChan

		a.log.Info("shutting down listing service", "signal", s.String())

		a.gRPCServer.GracefulStop()
	}()

	err := a.gRPCServer.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}
