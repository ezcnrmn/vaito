package server

import (
	"context"
	"database/sql"
	"fmt"

	pb "github.com/ezcnrmn/vaito/gen/go/storage"
)

type ListingServer struct {
	pb.UnimplementedListingServer
	db *sql.DB
}

func NewListingServer(db *sql.DB) *ListingServer {
	return &ListingServer{db: db}
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
	fmt.Printf("Got: %+v", listing)
	return &pb.CreateListingResponse{Id: 1}, nil
}
