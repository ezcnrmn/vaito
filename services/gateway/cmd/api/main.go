package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/ezcnrmn/vaito/services/gateway/internal/app"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// @title						Vaito Gateway API
// @version					1.0
//
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Введите токен в формате: Bearer ...
func main() {
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

	port := os.Getenv("GATEWAY_PORT")

	userURL := fmt.Sprintf("%s:%s", os.Getenv("USER_HOST"), os.Getenv("USER_PORT"))
	userGrpcClient, err := grpc.NewClient(userURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer userGrpcClient.Close()

	listingURL := fmt.Sprintf("%s:%s", os.Getenv("LISTING_HOST"), os.Getenv("LISTING_PORT"))
	listingGrpcClient, err := grpc.NewClient(listingURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer listingGrpcClient.Close()

	app := app.New(port, logger, userGrpcClient, listingGrpcClient)

	err = app.Run()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
