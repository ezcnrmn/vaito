package server

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	pb "github.com/ezcnrmn/vaito/gen/go/storage"
)

type ListingServer struct {
	pb.UnimplementedListingServer
	db  *sql.DB
	log *slog.Logger
}

func NewListingServer(db *sql.DB, logger *slog.Logger) *ListingServer {
	return &ListingServer{db: db, log: logger}
}

func (ls *ListingServer) CreateListing(_ context.Context, req *pb.CreateListingRequest) (*pb.CreateListingResponse, error) {
	listing := struct {
		Title       string
		Description string
		CategoryId  int32
		Price       int32
	}{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		CategoryId:  req.GetCategoryId(),
		Price:       req.GetPrice(),
	}
	ls.log.Debug("Got message", "message", fmt.Sprintf("%+v", listing))
	return &pb.CreateListingResponse{Id: 1}, nil
}
