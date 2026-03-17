package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"time"

	pb "github.com/ezcnrmn/vaito/gen/go/storage"
	"github.com/ezcnrmn/vaito/services/storage/internal/config"
	"google.golang.org/grpc"
)

type App struct {
	cfg    *config.Config
	log    *slog.Logger
	db     *sql.DB
	server *grpc.Server
}

func New(config *config.Config, logger *slog.Logger) *App {
	s := grpc.NewServer()
	pb.RegisterListingServer(s, &listingServer{})
	pb.RegisterUserServer(s, &userServer{})

	app := &App{
		cfg:    config,
		log:    logger,
		server: s,
	}

	return app
}

func (a *App) Run() error {
	// TODO: переосмыслить такой запуск/создание клиента (возможно перенести в main)
	db, err := a.openDB()
	if err != nil {
		return err
	}

	defer db.Close()
	a.log.Info("database connection pool established")
	a.db = db

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", a.cfg.Port))
	if err != nil {
		return err
	}

	a.log.Info("starting storage app", "port", a.cfg.Port)
	err = a.server.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) openDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", a.cfg.DB.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(a.cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(a.cfg.DB.MaxIdleConns)
	db.SetConnMaxIdleTime(a.cfg.DB.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
