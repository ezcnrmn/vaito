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

func (s *Server) CreateListing(_ context.Context, req *pb.CreateListingRequest) (*pb.CreateListingResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	token := req.Authentication.GetToken()
	userID, err := s.validatePermission(ctx, token, "listing:create")
	if err != nil {
		return nil, err
	}

	listing := &model.Listing{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Cetegory:    model.Category{ID: req.GetCategoryId()},
		UserID:      userID,
		Price:       req.GetPrice(),
	}

	err = s.model.listing.Insert(ctx, listing)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateListingResponse{
		Listing: listingToProtobufListing(listing),
	}, nil
}

func (s *Server) UpdateListing(_ context.Context, req *pb.UpdateListingRequest) (*pb.UpdateListingResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	token := req.Authentication.GetToken()
	userID, err := s.validatePermission(ctx, token, "listing:edit")
	if err != nil {
		return nil, err
	}

	id := req.GetId()
	listing, err := s.getListing(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		listing.Title = req.GetTitle()
	}
	if req.Description != nil {
		listing.Description = req.GetDescription()
	}
	if req.Price != nil {
		listing.Price = req.GetPrice()
	}
	if req.CategoryId != nil {
		listing.Cetegory.ID = req.GetCategoryId()
	}

	err = s.model.listing.UpdateListing(ctx, listing)
	if err != nil {
		if errors.Is(err, model.ErrEditConflict) {
			return nil, status.Error(codes.Aborted, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &pb.UpdateListingResponse{
		Listing: listingToProtobufListing(listing),
	}, nil
}

func (s *Server) DeleteListing(_ context.Context, req *pb.DeleteListingRequest) (*pb.DeleteListingResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	token := req.Authentication.GetToken()
	userID, err := s.validatePermission(ctx, token, "listing:delete")
	if err != nil {
		return nil, err
	}

	id := req.GetId()
	err = s.model.listing.DeleteListing(ctx, id, userID)
	if err != nil {
		if errors.Is(err, model.ErrListingNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &pb.DeleteListingResponse{}, nil
}

func (s *Server) GetListing(_ context.Context, req *pb.GetListingRequest) (*pb.GetListingResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	id := req.GetId()
	listing, err := s.model.listing.GetListing(ctx, id)
	if err != nil {
		if errors.Is(err, model.ErrListingNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &pb.GetListingResponse{
		Listing: listingToProtobufListing(listing),
	}, nil
}

func (s *Server) GetUserListing(_ context.Context, req *pb.GetUserListingRequest) (*pb.GetUserListingResponse, error) {
	return &pb.GetUserListingResponse{}, nil
}

func (s *Server) GetActiveListings(_ context.Context, req *pb.GetActiveListingsRequest) (*pb.GetActiveListingsResponse, error) {
	return &pb.GetActiveListingsResponse{}, nil
}

func (s *Server) GetListingsByUser(_ context.Context, req *pb.GetListingsByUserRequest) (*pb.GetListingsByUserResponse, error) {
	return &pb.GetListingsByUserResponse{}, nil
}

func (s *Server) SendListingToModeration(_ context.Context, req *pb.SendListingToModerationRequest) (*pb.SendListingToModerationResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	token := req.Authentication.GetToken()
	userID, err := s.validatePermission(ctx, token, "listing:edit")
	if err != nil {
		return nil, err
	}

	id := req.GetId()
	listing, err := s.getListing(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	if listing.Status != "Draft" {
		return nil, status.Error(codes.InvalidArgument, "you can only send listings with 'Draft' status for moderation")
	}

	err = s.updateListingStatus(ctx, listing, "Moderation")
	if err != nil {
		return nil, err
	}

	return &pb.SendListingToModerationResponse{}, nil
}

func (s *Server) ActivateListing(_ context.Context, req *pb.ActivateListingRequest) (*pb.ActivateListingResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	token := req.Authentication.GetToken()
	userID, err := s.validatePermission(ctx, token, "listing:edit")
	if err != nil {
		return nil, err
	}

	id := req.GetId()
	listing, err := s.getListing(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	if listing.Status != "Inactive" {
		return nil, status.Error(codes.InvalidArgument, "you can only activate listings with 'Inactive' status")
	}

	err = s.updateListingStatus(ctx, listing, "Active")
	if err != nil {
		return nil, err
	}

	return &pb.ActivateListingResponse{}, nil
}

func (s *Server) DeactivateListing(_ context.Context, req *pb.DeactivateListingRequest) (*pb.DeactivateListingResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	token := req.Authentication.GetToken()
	userID, err := s.validatePermission(ctx, token, "listing:edit")
	if err != nil {
		return nil, err
	}

	id := req.GetId()
	listing, err := s.getListing(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	if listing.Status != "Active" {
		return nil, status.Error(codes.InvalidArgument, "you can only deactivate listings with 'Active' status")
	}

	err = s.updateListingStatus(ctx, listing, "Inactive")
	if err != nil {
		return nil, err
	}

	return &pb.DeactivateListingResponse{}, nil
}

func (s *Server) GetCategories(_ context.Context, req *pb.GetCategoriesRequest) (*pb.GetCategoriesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	categories, err := s.model.listing.GetCategories(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to get categories")
	}

	pbCategories := make([]*pb.Category, 0, len(*categories))
	for _, c := range *categories {
		pbCategories = append(pbCategories, &pb.Category{Id: c.ID, Name: c.Name})
	}

	return &pb.GetCategoriesResponse{Categories: pbCategories}, nil
}
