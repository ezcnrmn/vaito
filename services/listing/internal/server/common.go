package server

import (
	"context"
	"errors"
	"slices"

	pbUser "github.com/ezcnrmn/vaito/gen/go/user"
	"github.com/ezcnrmn/vaito/services/listing/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) validateToken(ctx context.Context, token string) (userID int64, err error) {
	userResp, err := s.service.user.GetUserIDByToken(
		ctx,
		&pbUser.GetUserIDByTokenRequest{Token: &pbUser.Token{Token: token}},
	)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			code := s.Code()
			msg := s.Message()
			switch code {
			case codes.Unauthenticated:
				return 0, status.Error(codes.Unauthenticated, msg)
			default:
				return 0, status.Error(codes.Internal, err.Error())
			}
		}
	}

	return userResp.GetUserId(), nil
}

const permissionDeniedMsg = "you do not have permission to perform this action"

func (s *Server) validatePermission(ctx context.Context, token, permission string) (userID int64, err error) {
	userPermissionResp, err := s.service.user.GetUserPermissionsByToken(
		ctx,
		&pbUser.GetUserPermissionsByTokenRequest{Token: &pbUser.Token{Token: token}},
	)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			code := s.Code()
			msg := s.Message()
			switch code {
			case codes.Unauthenticated:
				return 0, status.Error(codes.Unauthenticated, msg)
			default:
				return 0, status.Error(codes.Internal, err.Error())
			}
		}
	}

	if i := slices.Index(userPermissionResp.Permissions, permission); i == -1 {
		return 0, status.Error(codes.PermissionDenied, permissionDeniedMsg)
	}

	return userPermissionResp.GetUserId(), nil
}

func (s *Server) getListing(ctx context.Context, id, userID int64) (*model.Listing, error) {
	listing, err := s.model.listing.GetListing(ctx, id, nil, nil)
	if err != nil {
		if errors.Is(err, model.ErrListingNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	if listing.UserID != userID {
		return nil, status.Error(codes.PermissionDenied, permissionDeniedMsg)
	}

	return listing, nil
}

func (s *Server) updateListingStatus(ctx context.Context, listing *model.Listing, newStatus string) error {
	err := s.model.listing.UpdateListingStatus(ctx, listing, newStatus)
	if err != nil {
		if errors.Is(err, model.ErrListingNotFound) {
			return status.Error(codes.NotFound, err.Error())
		} else {
			return status.Error(codes.Internal, err.Error())
		}
	}

	return nil
}
