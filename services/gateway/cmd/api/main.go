package main

import (
	"fmt"
	"log/slog"
	"os"

	pb "github.com/ezcnrmn/vaito/gen/go/storage"
	"github.com/ezcnrmn/vaito/services/gateway/internal/app"
	"github.com/ezcnrmn/vaito/services/gateway/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	storageUrl := fmt.Sprintf("%s:%s", os.Getenv("STORAGE_HOST"), os.Getenv("STORAGE_PORT"))
	port := os.Getenv("GATEWAY_PORT")
	host := os.Getenv("GATEWAY_HOST")
	grpcClientPort := os.Getenv("GATEWAY_CLIENT_PORT")

	config := &config.Config{
		Port:           port,
		GrpcClientPort: grpcClientPort,
		Host:           host,
		StorageUrl:     storageUrl,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// TODO: fix, если указать http://localhost:4001, он автоматические добавляет порт еще раз
	grpcClient, err := grpc.NewClient(storageUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer grpcClient.Close()

	userConn := pb.NewUserClient(grpcClient)
	listingConn := pb.NewListingClient(grpcClient)

	app := app.New(config, logger, userConn, listingConn)

	err = app.Run()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
