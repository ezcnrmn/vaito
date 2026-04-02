package server

import (
	"database/sql"
	"log/slog"

	pb "github.com/ezcnrmn/vaito/gen/go/listing"
	pbUser "github.com/ezcnrmn/vaito/gen/go/user"
	"github.com/ezcnrmn/vaito/services/listing/internal/model"
)

type Server struct {
	pb.UnimplementedListingServiceServer

	model struct {
		listing model.ListingModel
	}
	service struct {
		user pbUser.UserServiceClient
	}
	log *slog.Logger
}

func NewServer(db *sql.DB, logger *slog.Logger, user pbUser.UserServiceClient) *Server {
	return &Server{
		model: struct{ listing model.ListingModel }{
			listing: model.ListingModel{DB: db},
		},

		service: struct{ user pbUser.UserServiceClient }{
			user: user,
		},

		log: logger,
	}
}
