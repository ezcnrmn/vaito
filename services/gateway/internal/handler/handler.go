package handler

import (
	"log/slog"
	"reflect"
	"strings"

	pb "github.com/ezcnrmn/vaito/gen/go/storage"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	log         *slog.Logger
	validator   *validator.Validate
	userConn    pb.UserClient
	listingConn pb.ListingClient
}

func New(logger *slog.Logger, userConn pb.UserClient, listingConn pb.ListingClient) *Handler {
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
		log:         logger,
		validator:   validator,
		userConn:    userConn,
		listingConn: listingConn,
	}
}
