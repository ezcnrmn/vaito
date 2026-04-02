package server

import (
	"context"
	"errors"
	"slices"
	"time"

	pb "github.com/ezcnrmn/vaito/gen/go/listing"
	pbUser "github.com/ezcnrmn/vaito/gen/go/user"
	"github.com/ezcnrmn/vaito/services/listing/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const permissionDeniedMsg = "you do not have permission to perform this action"

func (s *Server) CreateListing(_ context.Context, req *pb.CreateListingRequest) (*pb.CreateListingResponse, error) {
	token := req.Authentication.GetToken()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	userPermissionResp, err := s.service.user.GetUserPermissionsByToken(ctx, &pbUser.GetUserPermissionsByTokenRequest{Token: &pbUser.Token{Token: token}})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			code := s.Code()
			msg := s.Message()
			switch code {
			case codes.Unauthenticated:
				return nil, status.Error(codes.Unauthenticated, msg)
			default:
				return nil, status.Error(codes.Internal, err.Error())
			}
		}
	}

	if i := slices.Index(userPermissionResp.Permissions, "listing:create"); i == -1 {
		return nil, status.Error(codes.PermissionDenied, permissionDeniedMsg)
	}

	listing := &model.Listing{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Cetegory:    model.Category{ID: req.GetCategoryId()},
		UserID:      userPermissionResp.GetUserId(),
		Price:       req.GetPrice(),
	}

	err = s.model.listing.Insert(listing)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateListingResponse{
		Listing: listingToProtobufListing(listing),
	}, nil
}

func (s *Server) UpdateListing(_ context.Context, req *pb.UpdateListingRequest) (*pb.UpdateListingResponse, error) {
	token := req.Authentication.GetToken()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	userPermissionResp, err := s.service.user.GetUserPermissionsByToken(ctx, &pbUser.GetUserPermissionsByTokenRequest{Token: &pbUser.Token{Token: token}})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			code := s.Code()
			msg := s.Message()
			switch code {
			case codes.Unauthenticated:
				return nil, status.Error(codes.Unauthenticated, msg)
			default:
				return nil, status.Error(codes.Internal, err.Error())
			}
		}
	}

	if i := slices.Index(userPermissionResp.Permissions, "listing:edit"); i == -1 {
		return nil, status.Error(codes.PermissionDenied, permissionDeniedMsg)
	}

	id := req.GetId()
	listing, err := s.model.listing.GetListing(id)
	if err != nil {
		if errors.Is(err, model.ErrListingNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	if listing.UserID != userPermissionResp.GetUserId() {
		return nil, status.Error(codes.PermissionDenied, permissionDeniedMsg)
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

	err = s.model.listing.UpdateListing(listing)
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
	token := req.Authentication.GetToken()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	userPermissionResp, err := s.service.user.GetUserPermissionsByToken(ctx, &pbUser.GetUserPermissionsByTokenRequest{Token: &pbUser.Token{Token: token}})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			code := s.Code()
			msg := s.Message()
			switch code {
			case codes.Unauthenticated:
				return nil, status.Error(codes.Unauthenticated, msg)
			default:
				return nil, status.Error(codes.Internal, err.Error())
			}
		}
	}

	if i := slices.Index(userPermissionResp.Permissions, "listing:delete"); i == -1 {
		return nil, status.Error(codes.PermissionDenied, permissionDeniedMsg)
	}

	id, userID := req.GetId(), userPermissionResp.GetUserId()
	err = s.model.listing.DeleteListing(id, userID)
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
	id := req.GetId()

	listing, err := s.model.listing.GetListing(id)
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
	return &pb.SendListingToModerationResponse{}, nil
}

func (s *Server) ActivateListing(_ context.Context, req *pb.ActivateListingRequest) (*pb.ActivateListingResponse, error) {
	return &pb.ActivateListingResponse{}, nil
}

func (s *Server) DeactivateListing(_ context.Context, req *pb.DeactivateListingRequest) (*pb.DeactivateListingResponse, error) {
	return &pb.DeactivateListingResponse{}, nil
}

func (s *Server) GetCategories(_ context.Context, req *pb.GetCategoriesRequest) (*pb.GetCategoriesResponse, error) {
	categories, err := s.model.listing.GetCategories()
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to get categories")
	}

	pbCategories := make([]*pb.Category, 0, len(*categories))
	for _, c := range *categories {
		pbCategories = append(pbCategories, &pb.Category{Id: c.ID, Name: c.Name})
	}

	return &pb.GetCategoriesResponse{Categories: pbCategories}, nil
}
