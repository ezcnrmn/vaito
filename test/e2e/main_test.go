package e2e

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/lib/pq"
)

var envConfig = map[string]string{
	"POSTGRES_PASSWORD": "pa55word",
	"POSTGRES_USER":     "vaito",

	"USER_DB_NAME": "users",
	"USER_DB_USER": "vaito_user",
	"USER_DB_PASS": "pa55word",
	"USER_DB_DSN":  "postgres://vaito_user:pa55word@postgres_db:5432/users?sslmode=disable",

	"LISTING_DB_NAME": "listings",
	"LISTING_DB_USER": "vaito_listing",
	"LISTING_DB_PASS": "pa55word",
	"LISTING_DB_DSN":  "postgres://vaito_listing:pa55word@postgres_db:5432/listings?sslmode=disable",

	"GATEWAY_HOST": "gateway",
	"GATEWAY_PORT": "4000",

	"USER_HOST": "user",
	"USER_PORT": "4001",

	"LISTING_HOST": "listing",
	"LISTING_PORT": "4002",
}

var (
	gatewayAddr string
	db          *sql.DB
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	composeFile := "../../docker-compose.test.yaml"

	stack, err := compose.NewDockerCompose(composeFile)
	if err != nil {
		log.Fatal(err)
	}

	err = stack.WithEnv(envConfig).WaitForService("postgres_db", wait.ForHealthCheck().WithPollInterval(2*time.Second)).Up(ctx, compose.Wait(true))
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = stack.Down(
			context.Background(),
			compose.RemoveOrphans(true),
			compose.RemoveVolumes(true),
			compose.RemoveImagesLocal,
		)
		if err != nil {
			log.Fatal(err)
		}
	}()

	postgresContainer, err := stack.ServiceContainer(ctx, "postgres_db")
	if err != nil {
		log.Fatal(err)
	}

	dockerEndpoint, err := postgresContainer.Endpoint(ctx, "")
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("postgres://vaito_user:pa55word@%s/users?sslmode=disable", dockerEndpoint)
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = seedDatabase(db)
	if err != nil {
		log.Fatal(err)
	}

	gatewayContainer, err := stack.ServiceContainer(ctx, "gateway")
	if err != nil {
		log.Fatal(err)
	}

	gatewayEndpoint, err := gatewayContainer.Endpoint(ctx, "")
	if err != nil {
		log.Fatal(err)
	}

	gatewayAddr = "http://" + gatewayEndpoint

	m.Run()
}
