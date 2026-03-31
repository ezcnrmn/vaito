package server

import (
	"database/sql"
	"log/slog"

	pb "github.com/ezcnrmn/vaito/gen/go/user"
	"github.com/ezcnrmn/vaito/services/user/internal/model"
)

type Server struct {
	pb.UnimplementedUserServiceServer

	model struct {
		user  model.UserModel
		token model.TokenModel
	}

	log *slog.Logger
}

func NewServer(db *sql.DB, logger *slog.Logger) *Server {
	return &Server{model: struct {
		user  model.UserModel
		token model.TokenModel
	}{
		user:  model.UserModel{DB: db},
		token: model.TokenModel{DB: db},
	}, log: logger}
}
