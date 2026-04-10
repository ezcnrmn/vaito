package server

import (
	"context"
	"errors"
	"time"

	pb "github.com/ezcnrmn/vaito/gen/go/listing"
	"github.com/ezcnrmn/vaito/services/listing/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetModerationListings(_ context.Context, req *pb.GetModerationListingsRequest) (*pb.GetModerationListingsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	token := req.Authentication.GetToken()
	_, err := s.validatePermission(ctx, token, "listing:moderate")
	if err != nil {
		return nil, err
	}

	statusName := "Moderation"
	pagination := model.Pagination{
		Page:          req.GetPagination().GetPage(),
		Size:          req.GetPagination().GetSize(),
		Sort:          "created_at",
		SortDirection: "DESC",
		Filter: struct {
			Status *string
			UserID *int64
		}{
			Status: &statusName,
		},
	}

	listings, err := s.model.listing.GetListings(ctx, &pagination)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbListings := make([]*pb.Listing, 0, len(listings))
	for _, l := range listings {
		pbListings = append(pbListings, listingToProtobufListing(&l))
	}

	return &pb.GetModerationListingsResponse{
		Items:      pbListings,
		Pagination: paginationToProtobufPagination(&pagination),
	}, nil
}

func (s *Server) ActivateListingByModeration(_ context.Context, req *pb.ActivateListingByModerationRequest) (*pb.ActivateListingByModerationResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	token := req.Authentication.GetToken()
	_, err := s.validatePermission(ctx, token, "listing:moderate")
	if err != nil {
		return nil, err
	}

	id := req.GetId()
	listing, err := s.model.listing.GetListing(ctx, id, nil, nil)
	if err != nil {
		if errors.Is(err, model.ErrListingNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	if listing.Status != "Moderation" {
		return nil, status.Error(codes.InvalidArgument, "you can only activate listings with 'Moderation' status")
	}

	err = s.updateListingStatus(ctx, listing, "Active")
	if err != nil {
		return nil, err
	}

	return &pb.ActivateListingByModerationResponse{}, nil
}

func (s *Server) DeactivateListingByModeration(_ context.Context, req *pb.DeactivateListingByModerationRequest) (*pb.DeactivateListingByModerationResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	token := req.Authentication.GetToken()
	_, err := s.validatePermission(ctx, token, "listing:moderate")
	if err != nil {
		return nil, err
	}

	id := req.GetId()
	listing, err := s.model.listing.GetListing(ctx, id, nil, nil)
	if err != nil {
		if errors.Is(err, model.ErrListingNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	if listing.Status != "Active" {
		return nil, status.Error(codes.InvalidArgument, "you can only deactivate listings with 'Active' status")
	}

	err = s.updateListingStatus(ctx, listing, "Inactive")
	if err != nil {
		return nil, err
	}

	return &pb.DeactivateListingByModerationResponse{}, nil
}
