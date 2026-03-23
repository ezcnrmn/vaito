package server

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	pb "github.com/ezcnrmn/vaito/gen/go/user"
	"github.com/ezcnrmn/vaito/services/user/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedUserServer

	model model.UserModel
	log   *slog.Logger
}

func NewServer(db *sql.DB, logger *slog.Logger) *Server {
	return &Server{model: model.UserModel{DB: db}, log: logger}
}

func (s *Server) CreateUser(_ context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	user := model.User{
		Name:         req.GetName(),
		Email:        req.GetEmail(),
		PasswordHash: req.GetPasswordHash(),
	}

	err := s.model.Insert(&user)
	if err != nil {
		if errors.Is(err, model.ErrDuplicateEmail) {
			return nil, status.Error(codes.AlreadyExists, model.ErrDuplicateEmail.Error())
		} else {
			return nil, err
		}
	}

	return &pb.UserResponse{Id: user.ID}, nil
}
