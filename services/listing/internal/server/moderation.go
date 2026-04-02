package server

import (
	"context"

	pb "github.com/ezcnrmn/vaito/gen/go/listing"
)

func (s *Server) GetModerationListings(_ context.Context, req *pb.GetModerationListingsRequest) (*pb.GetModerationListingsResponse, error) {
	return &pb.GetModerationListingsResponse{}, nil
}

func (s *Server) ActivateListingByModeration(_ context.Context, req *pb.ActivateListingByModerationRequest) (*pb.ActivateListingByModerationResponse, error) {
	return &pb.ActivateListingByModerationResponse{}, nil
}

func (s *Server) DeactivateListingByModeration(_ context.Context, req *pb.DeactivateListingByModerationRequest) (*pb.DeactivateListingByModerationResponse, error) {
	return &pb.DeactivateListingByModerationResponse{}, nil
}
