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

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := model.User{
		Name:      req.GetName(),
		Email:     req.GetEmail(),
		Password:  model.Password{Plaintext: req.GetPassword()},
		CreatedAt: time.Now(),
	}

	err := user.Password.GenerateHash()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = s.model.user.Insert(ctx, &user)
	if err != nil {
		if errors.Is(err, model.ErrDuplicateEmail) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &pb.CreateUserResponse{
		User: userToProtobufUser(&user),
	}, nil
}

const invalidTokenMsg = "invalid token"

func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	token := &model.Token{
		Scope: model.ScopeAuthentication,
		Text:  req.Token.GetToken(),
	}
	err := token.DecodeBytes()
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, invalidTokenMsg)
	}

	user, err := s.model.user.GetUserByToken(ctx, token)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	userID := req.GetId()
	if user.ID != userID {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	if req.Email != nil {
		user.Email = req.GetEmail()
	}
	if req.Name != nil {
		user.Name = req.GetName()
	}

	err = s.model.user.UpdateUser(ctx, user)
	if err != nil {
		if errors.Is(err, model.ErrEditConflict) {
			return nil, status.Error(codes.Aborted, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &pb.UpdateUserResponse{
		User: userToProtobufUser(user),
	}, nil
}

func (s *Server) UpdateUserPassword(ctx context.Context, req *pb.UpdateUserPasswordRequest) (*pb.UpdateUserPasswordResponse, error) {
	token := &model.Token{
		Scope: model.ScopeAuthentication,
		Text:  req.Token.GetToken(),
	}
	err := token.DecodeBytes()
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, invalidTokenMsg)
	}

	user, err := s.model.user.GetUserByToken(ctx, token)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	userID := req.GetId()
	if user.ID != userID {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	user.Password.Plaintext = req.GetPassword()
	err = user.Password.GenerateHash()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = s.model.user.UpdateUser(ctx, user)
	if err != nil {
		if errors.Is(err, model.ErrEditConflict) {
			return nil, status.Error(codes.Aborted, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &pb.UpdateUserPasswordResponse{}, nil
}

func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	id := req.GetId()

	user, err := s.model.user.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &pb.GetUserResponse{
		User: userToProtobufUser(user),
	}, nil
}

func (s *Server) GetUserIDByToken(ctx context.Context, req *pb.GetUserIDByTokenRequest) (*pb.GetUserIDByTokenResponse, error) {
	token := &model.Token{
		Scope: model.ScopeAuthentication,
		Text:  req.Token.GetToken(),
	}
	err := token.DecodeBytes()
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, invalidTokenMsg)
	}

	userID, err := s.model.user.GetUserIDByToken(ctx, token)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, status.Error(codes.Unauthenticated, invalidTokenMsg)
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &pb.GetUserIDByTokenResponse{
		UserId: userID,
	}, nil
}

func (s *Server) GetUserPermissionsByToken(ctx context.Context, req *pb.GetUserPermissionsByTokenRequest) (*pb.GetUserPermissionsByTokenResponse, error) {
	token := &model.Token{
		Scope: model.ScopeAuthentication,
		Text:  req.Token.GetToken(),
	}
	err := token.DecodeBytes()
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, invalidTokenMsg)
	}

	userID, permissions, err := s.model.user.GetUserPermissionsByToken(ctx, token)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, status.Error(codes.Unauthenticated, invalidTokenMsg)
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &pb.GetUserPermissionsByTokenResponse{
		UserId:      userID,
		Permissions: permissions,
	}, nil
}
