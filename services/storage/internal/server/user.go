package server

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	pb "github.com/ezcnrmn/vaito/gen/go/storage"
)

type UserServer struct {
	pb.UnimplementedUserServer
	db  *sql.DB
	log *slog.Logger
}

func NewUserServer(db *sql.DB, logger *slog.Logger) *UserServer {
	return &UserServer{db: db, log: logger}
}

func (us *UserServer) CreateUser(_ context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := struct {
		Name         string
		Email        string
		PasswordHash string
	}{
		Name:         req.GetName(),
		Email:        req.GetEmail(),
		PasswordHash: req.GetPasswordHash(),
	}
	us.log.Debug("Got message", "message", fmt.Sprintf("%+v", user))
	return &pb.CreateUserResponse{Id: 1}, nil
}
