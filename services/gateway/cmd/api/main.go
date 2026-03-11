package main

import (
	"fmt"
	"gateway/internal/app"
	"gateway/internal/config"
	"log/slog"
	"os"
)

func main() {
	storageUrl := fmt.Sprintf("%s:%s", os.Getenv("STORAGE_HOST"), os.Getenv("STORAGE_PORT"))
	port := os.Getenv("GATEWAY_PORT")

	config := &config.Config{
		Port:       port,
		StorageUrl: storageUrl,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := app.New(config, logger)
	err := app.Run()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
