package handler

import (
	"log/slog"
	"reflect"
	"strings"

	pbListing "github.com/ezcnrmn/vaito/gen/go/listing"
	pbUser "github.com/ezcnrmn/vaito/gen/go/user"
	"github.com/ezcnrmn/vaito/services/gateway/internal/config"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Handler struct {
	cfg       config.Config
	log       *slog.Logger
	validator *validator.Validate

	userConn    pbUser.UserServiceClient
	listingConn pbListing.ListingServiceClient
	health      struct {
		user    grpc_health_v1.HealthClient
		listing grpc_health_v1.HealthClient
	}
}

func New(
	config config.Config,
	logger *slog.Logger,
	userConn pbUser.UserServiceClient,
	listingConn pbListing.ListingServiceClient,
	userHealthConn,
	listingHealthConn grpc_health_v1.HealthClient,
) *Handler {
	validator := validator.New()
	validator.RegisterValidation("lettersAndDigits", lettersAndDigits)
	validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Handler{
		cfg:         config,
		log:         logger,
		validator:   validator,
		userConn:    userConn,
		listingConn: listingConn,
		health: struct {
			user    grpc_health_v1.HealthClient
			listing grpc_health_v1.HealthClient
		}{
			user:    userHealthConn,
			listing: listingHealthConn,
		},
	}
}
