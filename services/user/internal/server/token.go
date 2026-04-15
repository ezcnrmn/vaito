package server

import (
	"context"
	"errors"
	"time"

	pb "github.com/ezcnrmn/vaito/gen/go/user"
	"github.com/ezcnrmn/vaito/services/user/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const unauthenticatedMsg = "invalid authentication credentials"

func (s *Server) AuthenticateUser(ctx context.Context, req *pb.AuthenticateUserRequest) (*pb.AuthenticateUserResponse, error) {
	email, password := req.GetEmail(), req.GetPassword()

	user, err := s.model.user.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, status.Error(codes.Unauthenticated, unauthenticatedMsg)
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	user.Password.Plaintext = password
	if !user.Password.CompareHash() {
		return nil, status.Error(codes.Unauthenticated, unauthenticatedMsg)
	}

	token, err := model.NewToken(user.ID, model.ScopeAuthentication, 24*time.Hour)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = s.model.token.Insert(ctx, token)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.AuthenticateUserResponse{Token: &pb.Token{Token: token.Text}, UserId: user.ID}, nil
}
