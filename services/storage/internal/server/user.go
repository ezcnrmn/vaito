package server

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	pb "github.com/ezcnrmn/vaito/gen/go/storage"
	"github.com/ezcnrmn/vaito/services/storage/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	pb.UnimplementedUserServer

	model model.UserModel
	log   *slog.Logger
}

func NewUserServer(db *sql.DB, logger *slog.Logger) *UserServer {
	return &UserServer{model: model.UserModel{DB: db}, log: logger}
}

func (us *UserServer) CreateUser(_ context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := model.User{
		Name:         req.GetName(),
		Email:        req.GetEmail(),
		PasswordHash: req.GetPasswordHash(),
		CreatedAt:    req.GetCreatedAt().AsTime(),
	}

	err := us.model.Insert(&user)
	if err != nil {
		if errors.Is(err, model.ErrDuplicateEmail) {
			return nil, status.Error(codes.AlreadyExists, model.ErrDuplicateEmail.Error())
		} else {
			return nil, err
		}
	}

	return &pb.CreateUserResponse{Id: user.ID}, nil
}
