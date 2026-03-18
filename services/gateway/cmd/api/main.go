package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ezcnrmn/vaito/services/gateway/internal/app"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

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
