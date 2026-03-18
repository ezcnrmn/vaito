package handler

import (
	"context"
	"net/http"
	"time"

	pb "github.com/ezcnrmn/vaito/gen/go/storage"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name     string `json:"name" validate:"required,min=3,max=50,lettersAndDigits"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=100"`
	}

	err := readJSON(w, r, &payload)
	if err != nil {
		sendError(w, err)
		return
	}

	err = h.validator.Struct(payload)
	if err != nil {
		sendValidateError(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := h.userConn.CreateUser(ctx, &pb.CreateUserRequest{
		Name:         payload.Name,
		Email:        payload.Email,
		PasswordHash: payload.Password,
	})
	if err != nil {
		h.log.Debug(err.Error())
		writeJSON(w, http.StatusInternalServerError, envelope{"error": err.Error()})
		return
	}
	h.log.Debug(resp.String())
	writeJSON(w, http.StatusInternalServerError, envelope{"success": resp.String()})
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) Authenticate(w http.ResponseWriter, r *http.Request) {}
