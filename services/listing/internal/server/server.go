package server

import (
	"context"
	"database/sql"
	"log/slog"

	pb "github.com/ezcnrmn/vaito/gen/go/listing"
	pbUser "github.com/ezcnrmn/vaito/gen/go/user"
	"github.com/ezcnrmn/vaito/services/listing/internal/model"
)

type notification interface {
	PublishVisibilityChanged(ctx context.Context, userEmail string, listingID int64, visibility bool) error
}

type Server struct {
	pb.UnimplementedListingServiceServer

	model struct {
		listing model.ListingModel
	}
	service struct {
		user pbUser.UserServiceClient
	}
	log *slog.Logger

	notification notification
}

func NewServer(db *sql.DB, logger *slog.Logger, user pbUser.UserServiceClient, notification notification) *Server {
	return &Server{
		model: struct{ listing model.ListingModel }{
			listing: model.ListingModel{DB: db},
		},

		service: struct{ user pbUser.UserServiceClient }{
			user: user,
		},

		log: logger,

		notification: notification,
	}
}
