package server

import (
	"context"
	"database/sql"
	"fmt"

	pb "github.com/ezcnrmn/vaito/gen/go/storage"
)

type UserServer struct {
	pb.UnimplementedUserServer
	db *sql.DB
}

func NewUserServer(db *sql.DB) *UserServer {
	return &UserServer{db: db}
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
	fmt.Printf("Got: %+v", user)
	return &pb.CreateUserResponse{Id: 1}, nil
}
