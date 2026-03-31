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

func (s *Server) CreateUser(_ context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
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

	err = s.model.user.Insert(&user)
	if err != nil {
		if errors.Is(err, model.ErrDuplicateEmail) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &pb.CreateUserResponse{
		User: &pb.User{
			Id:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			RoleName: user.Role.Name,
		},
	}, nil
}

const invalidTokenMsg = "invalid token"

func (s *Server) UpdateUser(_ context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	token := &model.Token{
		Scope: model.ScopeAuthentication,
		Text:  req.Token.GetToken(),
	}
	err := token.DecodeBytes()
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, invalidTokenMsg)
	}

	user, err := s.model.user.GetUserByToken(token)
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

	err = s.model.user.UpdateUser(user)
	if err != nil {
		if errors.Is(err, model.ErrEditConflict) {
			return nil, status.Error(codes.Aborted, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &pb.UpdateUserResponse{
		User: &pb.User{
			Id:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			RoleName: user.Role.Name,
		},
	}, nil
}

func (s *Server) UpdateUserPassword(_ context.Context, req *pb.UpdateUserPasswordRequest) (*pb.UpdateUserPasswordResponse, error) {
	token := &model.Token{
		Scope: model.ScopeAuthentication,
		Text:  req.Token.GetToken(),
	}
	err := token.DecodeBytes()
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, invalidTokenMsg)
	}

	user, err := s.model.user.GetUserByToken(token)
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

	err = s.model.user.UpdateUser(user)
	if err != nil {
		if errors.Is(err, model.ErrEditConflict) {
			return nil, status.Error(codes.Aborted, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &pb.UpdateUserPasswordResponse{}, nil
}

func (s *Server) GetUser(_ context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	id := req.GetId()

	user, err := s.model.user.GetUser(id)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			Id:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			RoleName: user.Role.Name,
		},
	}, nil
}
