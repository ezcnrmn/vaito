package handler

import (
	"context"
	"net/http"
	"time"

	pbUser "github.com/ezcnrmn/vaito/gen/go/user"
	"github.com/ezcnrmn/vaito/services/gateway/internal/lib/contextutil"
	"github.com/ezcnrmn/vaito/services/gateway/internal/lib/jsonutil"
	"google.golang.org/grpc/codes"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.userConn.CreateUser(ctx, &pbUser.CreateUserRequest{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		h.handleGRPCError(w, err, func(code codes.Code, msg string) {
			switch code {
			case codes.AlreadyExists:
				sendError(w, http.StatusConflict, msg)
			default:
				h.log.Error(msg, "code", code)
				sendInternalError(w)
			}
		})
		return
	}

	writeUserResponse(w, resp.GetUser())
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	token := contextutil.GetToken(r)
	if token == "" {
		sendMissingTokenError(w)
		return
	}

	id, err := readUserIDParam(r)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	var payload struct {
		Name  *string `json:"name" validate:"omitempty,min=3,max=50,lettersAndDigits"`
		Email *string `json:"email" validate:"omitempty,email"`
	}

	err = jsonutil.ReadJSON(w, r, &payload)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.validator.Struct(payload)
	if err != nil {
		sendValidateError(w, err)
		return
	}

	if payload.Name == nil && payload.Email == nil {
		sendError(w, http.StatusBadRequest, "you must specify at least one field")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.userConn.UpdateUser(ctx, &pbUser.UpdateUserRequest{
		Id:    id,
		Name:  payload.Name,
		Email: payload.Email,
		Token: &pbUser.Token{
			Token: token,
		},
	})
	if err != nil {
		h.handleGRPCError(w, err, func(code codes.Code, msg string) {
			switch code {
			case codes.Unauthenticated:
				sendError(w, http.StatusUnauthorized, msg)
			case codes.NotFound:
				sendUnauthorizedError(w)
			case codes.PermissionDenied:
				sendForbiddenError(w)
			case codes.Aborted:
				sendError(w, http.StatusUnprocessableEntity, msg)
			default:
				h.log.Error(msg, "code", code)
				sendInternalError(w)
			}
		})
		return
	}

	writeUserResponse(w, resp.GetUser())
}

func (h *Handler) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	token := contextutil.GetToken(r)
	if token == "" {
		sendMissingTokenError(w)
		return
	}

	id, err := readUserIDParam(r)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	var payload struct {
		Password string `json:"password" validate:"required,min=8,max=100"`
	}

	err = jsonutil.ReadJSON(w, r, &payload)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.validator.Struct(payload)
	if err != nil {
		sendValidateError(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = h.userConn.UpdateUserPassword(ctx, &pbUser.UpdateUserPasswordRequest{
		Id:       id,
		Password: payload.Password,
		Token: &pbUser.Token{
			Token: token,
		},
	})
	if err != nil {
		h.handleGRPCError(w, err, func(code codes.Code, msg string) {
			switch code {
			case codes.Unauthenticated:
				sendError(w, http.StatusUnauthorized, msg)
			case codes.NotFound:
				sendUnauthorizedError(w)
			case codes.PermissionDenied:
				sendForbiddenError(w)
			case codes.Aborted:
				sendError(w, http.StatusUnprocessableEntity, msg)
			default:
				h.log.Error(msg, "code", code)
				sendInternalError(w)
			}
		})
		return
	}

	sendSuccessMessage(w, "password successfully changed")
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := readUserIDParam(r)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.userConn.GetUser(ctx, &pbUser.GetUserRequest{Id: id})
	if err != nil {
		h.handleGRPCError(w, err, func(code codes.Code, msg string) {
			switch code {
			case codes.NotFound:
				sendError(w, http.StatusNotFound, msg)
			default:
				h.log.Error(msg, "code", code)
				sendInternalError(w)
			}
		})
		return
	}

	writeUserResponse(w, resp.GetUser())
}

func (h *Handler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.userConn.AuthenticateUser(ctx, &pbUser.AuthenticateUserRequest{
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		h.handleGRPCError(w, err, func(code codes.Code, msg string) {
			switch code {
			case codes.Unauthenticated:
				sendError(w, http.StatusUnauthorized, msg)
			default:
				h.log.Error(msg, "code", code)
				sendInternalError(w)
			}
		})
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, jsonutil.Envelope{"token": resp.GetToken().GetToken(), "userID": resp.GetUserId()})
}
