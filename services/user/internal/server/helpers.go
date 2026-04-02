package server

import (
	pb "github.com/ezcnrmn/vaito/gen/go/user"
	"github.com/ezcnrmn/vaito/services/user/internal/model"
)

func userToProtobufUser(user *model.User) *pb.User {
	return &pb.User{
		Id:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		RoleName: user.Role.Name,
	}
}
