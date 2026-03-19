package handler

import (
	"context"
	"net/http"
	"time"

	pb "github.com/ezcnrmn/vaito/gen/go/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name     string `json:"name" validate:"required,min=3,max=50,lettersAndDigits"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=100"`
	}

	err := readJSON(w, r, &payload)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.validator.Struct(payload)
	if err != nil {
		sendValidateError(w, err)
		return
	}

	hash, err := hashPassword(payload.Password)
	if err != nil {
		sendInternalError(w)
		h.log.Error(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := h.userConn.CreateUser(ctx, &pb.CreateUserRequest{
		Name:         payload.Name,
		Email:        payload.Email,
		PasswordHash: hash,
		CreatedAt:    timestamppb.New(time.Now()),
	})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			msg := s.Message()
			code := s.Code()
			switch code {
			case codes.AlreadyExists:
				sendError(w, http.StatusConflict, msg)
			default:
				h.log.Error(msg, "code", code)
				sendInternalError(w)
			}
		} else {
			h.log.Error(err.Error())
			sendInternalError(w)
		}
		return
	}

	writeJSON(w, http.StatusOK, envelope{"userId": resp.String()})
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) Authenticate(w http.ResponseWriter, r *http.Request) {}
