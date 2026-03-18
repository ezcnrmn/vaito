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

	storageUrl := fmt.Sprintf("%s:%s", os.Getenv("STORAGE_HOST"), os.Getenv("STORAGE_PORT"))
	grpcClient, err := grpc.NewClient(storageUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer grpcClient.Close()

	app := app.New(port, logger, grpcClient)

	err = app.Run()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
