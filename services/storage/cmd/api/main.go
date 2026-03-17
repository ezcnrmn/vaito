package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/ezcnrmn/vaito/services/storage/internal/app"
	"github.com/ezcnrmn/vaito/services/storage/internal/config"
	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("no DB_DSN was provided")
	}

	gatewayUrl := fmt.Sprintf("%s:%s", os.Getenv("GATEWAY_HOST"), os.Getenv("GATEWAY_PORT"))
	port := os.Getenv("STORAGE_PORT")

	config := &config.Config{
		Port:       port,
		GatewayUrl: gatewayUrl,
	}
	config.DB.DSN = dsn

	// TODO: вынести в env
	config.DB.MaxOpenConns = 25
	config.DB.MaxIdleConns = 25
	config.DB.MaxIdleTime = 15 * time.Minute

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := app.New(config, logger)
	err := app.Run()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
