package handler

import (
	"context"
	"net/http"
	"time"

	pbUser "github.com/ezcnrmn/vaito/gen/go/user"
	"github.com/ezcnrmn/vaito/services/gateway/internal/lib/jsonutil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name     string `json:"name" validate:"required,min=3,max=50,lettersAndDigits"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=100"`
	}

	err := jsonutil.ReadJSON(w, r, &payload)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.validator.Struct(payload)
	if err != nil {
		sendValidateError(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := h.userConn.CreateUser(ctx, &pbUser.CreateUserRequest{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
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

	user := struct {
		ID       int64
		Name     string
		Email    string
		RoleName string
	}{
		ID:       resp.Id,
		Name:     resp.Name,
		Email:    resp.Email,
		RoleName: resp.RoleName,
	}
	jsonutil.WriteJSON(w, http.StatusOK, jsonutil.Envelope{"user": user})
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name  *string `json:"name" validate:"min=3,max=50,lettersAndDigits"`
		Email *string `json:"email" validate:"email"`
	}

	err := jsonutil.ReadJSON(w, r, &payload)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.validator.Struct(payload)
	if err != nil {
		sendValidateError(w, err)
		return
	}

	var token []byte

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := h.userConn.UpdateUser(ctx, &pbUser.UpdateUserRequest{
		Name:  payload.Name,
		Email: payload.Email,
		Token: &pbUser.TokenRequest{
			Token: token,
		},
	})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			msg := s.Message()
			code := s.Code()
			switch code {
			case codes.Unauthenticated:
				sendError(w, http.StatusUnauthorized, msg)
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

	user := struct {
		ID       int64
		Name     string
		Email    string
		RoleName string
	}{
		ID:       resp.Id,
		Name:     resp.Name,
		Email:    resp.Email,
		RoleName: resp.RoleName,
	}
	jsonutil.WriteJSON(w, http.StatusOK, jsonutil.Envelope{"user": user})
}

func (h *Handler) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := h.userConn.GetUser(ctx, &pbUser.GetUserRequest{Id: id})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			msg := s.Message()
			code := s.Code()
			switch code {
			case codes.NotFound:
				sendError(w, http.StatusNotFound, msg)
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

	user := struct {
		ID       int64
		Name     string
		Email    string
		RoleName string
	}{
		ID:       resp.Id,
		Name:     resp.Name,
		Email:    resp.Email,
		RoleName: resp.RoleName,
	}
	jsonutil.WriteJSON(w, http.StatusOK, jsonutil.Envelope{"user": user})
}

func (h *Handler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=100"`
	}

	err := jsonutil.ReadJSON(w, r, &payload)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.validator.Struct(payload)
	if err != nil {
		sendValidateError(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := h.userConn.AuthenticateUser(ctx, &pbUser.AuthenticateUserRequest{
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			msg := s.Message()
			code := s.Code()
			switch code {
			case codes.Unauthenticated:
				sendError(w, http.StatusUnauthorized, msg)
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

	jsonutil.WriteJSON(w, http.StatusOK, jsonutil.Envelope{"token": string(resp.Token)})
}

func (h *Handler) GetUserListings(w http.ResponseWriter, r *http.Request) {
}
