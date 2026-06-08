package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/ezcnrmn/vaito/services/listing/internal/app"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	dbDsn := os.Getenv("LISTING_DB_DSN")
	dbMaxOpenConns := flag.Int("db-max-open-conns", 25, "PostgreSQL max open connections")
	dbMaxIdleConns := flag.Int("db-max-idle-conns", 25, "PostgreSQL max idle connections")
	dbMaxIdleTime := flag.Duration("db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	showDebug := flag.Bool("debug-log", false, "Sets log level to Debug and shows source of message")
	flag.Parse()

	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	if *showDebug {
		opts.Level = slog.LevelDebug
		opts.AddSource = true
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	db, err := openDB(dbDsn, *dbMaxOpenConns, *dbMaxIdleConns, *dbMaxIdleTime)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	port := os.Getenv("LISTING_PORT")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	userURL := fmt.Sprintf("%s:%s", os.Getenv("USER_HOST"), os.Getenv("USER_PORT"))
	userGrpcClient, err := grpc.NewClient(userURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer userGrpcClient.Close()

	rabbitmqUrl := os.Getenv("RABBITMQ_URL")
	conn, err := amqp.Dial(rabbitmqUrl)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer channel.Close()

	app, err := app.New(logger, db, userGrpcClient, channel)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("starting listing service", "port", port)
	err = app.Run(listener)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(dsn string, maxOpenConns, maxIdleConns int, maxIdleTime time.Duration) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
