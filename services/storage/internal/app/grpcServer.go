package app

import (
	"context"
	"fmt"

	pb "github.com/ezcnrmn/vaito/gen/go/storage"
)

type userServer struct {
	pb.UnimplementedUserServer
}

func (s *userServer) CreateUser(_ context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
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

type listingServer struct {
	pb.UnimplementedListingServer
}

func (s *listingServer) CreateListing(_ context.Context, req *pb.CreateListingRequest) (*pb.CreateListingResponse, error) {
	listing := struct {
		Title       string
		Description string
		CategoryId  int32
		Price       int32
	}{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		CategoryId:  req.GetCategoryId(),
		Price:       req.GetPrice(),
	}
	fmt.Printf("Got: %+v", listing)
	return &pb.CreateListingResponse{Id: 1}, nil
}
