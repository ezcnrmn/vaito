package server

import (
	"context"
	"errors"

	pb "github.com/ezcnrmn/vaito/gen/go/listing"
	"github.com/ezcnrmn/vaito/services/listing/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateListing(ctx context.Context, req *pb.CreateListingRequest) (*pb.CreateListingResponse, error) {
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

func (s *Server) UpdateListing(ctx context.Context, req *pb.UpdateListingRequest) (*pb.UpdateListingResponse, error) {
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

func (s *Server) DeleteListing(ctx context.Context, req *pb.DeleteListingRequest) (*pb.DeleteListingResponse, error) {
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

func (s *Server) GetListing(ctx context.Context, req *pb.GetListingRequest) (*pb.GetListingResponse, error) {
	id, statusName := req.GetId(), "Active"
	listing, err := s.model.listing.GetListing(ctx, id, nil, &statusName)
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

func (s *Server) GetUserListing(ctx context.Context, req *pb.GetUserListingRequest) (*pb.GetUserListingResponse, error) {
	token := req.Authentication.GetToken()
	userID, err := s.validateToken(ctx, token)
	if err != nil {
		return nil, err
	}

	requestedUserID := req.GetUserId()
	if userID != requestedUserID {
		return nil, status.Error(codes.PermissionDenied, permissionDeniedMsg)
	}

	id := req.GetId()
	listing, err := s.model.listing.GetListing(ctx, id, &userID, nil)
	if err != nil {
		if errors.Is(err, model.ErrListingNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &pb.GetUserListingResponse{
		Listing: listingToProtobufListing(listing),
	}, nil
}

func (s *Server) GetActiveListings(ctx context.Context, req *pb.GetActiveListingsRequest) (*pb.GetActiveListingsResponse, error) {
	statusName := "Active"
	pagination := model.Pagination{
		Page:          req.GetPagination().GetPage(),
		Size:          req.GetPagination().GetSize(),
		Sort:          "published_at",
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

	return &pb.GetActiveListingsResponse{
		Items:      pbListings,
		Pagination: paginationToProtobufPagination(&pagination),
	}, nil
}

func (s *Server) GetListingsByUser(ctx context.Context, req *pb.GetListingsByUserRequest) (*pb.GetListingsByUserResponse, error) {
	token := req.Authentication.GetToken()
	userID, err := s.validateToken(ctx, token)
	if err != nil {
		return nil, err
	}

	reqUserID := req.GetUserId()
	if reqUserID != userID {
		return nil, status.Error(codes.PermissionDenied, permissionDeniedMsg)
	}

	pagination := model.Pagination{
		Page:          req.GetPagination().GetPage(),
		Size:          req.GetPagination().GetSize(),
		Sort:          "created_at",
		SortDirection: "DESC",
		Filter: struct {
			Status *string
			UserID *int64
		}{
			UserID: &userID,
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

	return &pb.GetListingsByUserResponse{
		Items:      pbListings,
		Pagination: paginationToProtobufPagination(&pagination),
	}, nil
}

func (s *Server) SendListingToModeration(ctx context.Context, req *pb.SendListingToModerationRequest) (*pb.SendListingToModerationResponse, error) {
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

func (s *Server) ActivateListing(ctx context.Context, req *pb.ActivateListingRequest) (*pb.ActivateListingResponse, error) {
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

func (s *Server) DeactivateListing(ctx context.Context, req *pb.DeactivateListingRequest) (*pb.DeactivateListingResponse, error) {
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

func (s *Server) GetCategories(ctx context.Context, req *pb.GetCategoriesRequest) (*pb.GetCategoriesResponse, error) {
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
