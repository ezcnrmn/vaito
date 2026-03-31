package server

import (
	"context"
	"database/sql"
	"log/slog"

	pb "github.com/ezcnrmn/vaito/gen/go/listing"
	"github.com/ezcnrmn/vaito/services/listing/internal/model"
)

type Server struct {
	pb.UnimplementedListingServiceServer

	model model.ListingModel
	log   *slog.Logger
}

func NewServer(db *sql.DB, logger *slog.Logger) *Server {
	return &Server{model: model.ListingModel{DB: db}, log: logger}
}

func (s *Server) CreateListing(_ context.Context, req *pb.CreateListingRequest) (*pb.CreateListingResponse, error) {
	return &pb.CreateListingResponse{}, nil
}
