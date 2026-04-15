package handler

import (
	"context"
	"net/http"
	"time"

	pbListing "github.com/ezcnrmn/vaito/gen/go/listing"
	pbUser "github.com/ezcnrmn/vaito/gen/go/user"
	"github.com/ezcnrmn/vaito/services/gateway/internal/lib/contextutil"
	"github.com/ezcnrmn/vaito/services/gateway/internal/lib/jsonutil"
	"google.golang.org/grpc/codes"
)

type CreateUserPayload struct {
	Name     string `json:"name" validate:"required,min=3,max=50,lettersAndDigits" example:"User"`
	Email    string `json:"email" validate:"required,email" example:"user@test.com"`
	Password string `json:"password" validate:"required,min=8,max=100" example:"pa55word"`
}

// CreateUser - Создание пользователя
//
//	@summary	Создание пользователя
//	@tags		users
//	@param		input	body	CreateUserPayload	true	"Данные нового пользователя"
//	@accept		json
//	@produce	json
//	@success	200	{object}	UserResponse
//	@failure	400	{object}	ErrorResponse
//	@failure	409	{object}	ErrorResponse
//	@router		/users [post]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserPayload

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

type UpdateUserPayload struct {
	Name  *string `json:"name" validate:"omitempty,min=3,max=50,lettersAndDigits" example:"User"`
	Email *string `json:"email" validate:"omitempty,email" example:"user@test.com"`
}

// UpdateUser - Обновление данных пользователя (name, email)
//
//	@summary	Обновление данных пользователя
//	@tags		users
//	@param		userID	path	int					true	"Идентификатор пользователя"
//	@param		input	body	UpdateUserPayload	true	"Данные пользователя"
//	@accept		json
//	@produce	json
//	@success	200	{object}	UserResponse
//	@failure	400	{object}	ErrorResponse
//	@failure	401	{object}	ErrorResponse
//	@failure	403	{object}	ErrorResponse
//	@failure	404	{object}	ErrorResponse
//	@failure	422	{object}	ErrorResponse
//	@security	BearerAuth
//	@router		/users/{userID} [patch]
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

	var payload UpdateUserPayload

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

type UpdateUserPasswordPayload struct {
	Password string `json:"password" validate:"required,min=8,max=100" example:"pa55word"`
}

// UpdateUserPassword - Обновление пароля пользователя (password)
//
//	@summary	Обновление пароля пользователя
//	@tags		users
//	@param		userID	path	int							true	"Идентификатор пользователя"
//	@param		input	body	UpdateUserPasswordPayload	true	"Пароль пользователя"
//	@accept		json
//	@produce	json
//	@success	200	{object}	MessageResponse
//	@failure	400	{object}	ErrorResponse
//	@failure	401	{object}	ErrorResponse
//	@failure	403	{object}	ErrorResponse
//	@failure	404	{object}	ErrorResponse
//	@failure	422	{object}	ErrorResponse
//	@security	BearerAuth
//	@router		/users/{userID}/update-password [put]
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

	var payload UpdateUserPasswordPayload
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

// GetUser - Получение данных пользователя
//
//	@summary	Получение данных пользователя
//	@tags		users
//	@param		userID	path	int	true	"Идентификатор пользователя"
//	@produce	json
//	@success	200	{object}	UserResponse
//	@failure	404	{object}	ErrorResponse
//	@router		/users/{userID} [get]
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

type AuthenticateUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
type TokenResponse struct {
	Token  string `json:"token"`
	UserID int64  `json:"userID"`
}

// AuthenticateUser - Аутентификация пользователя
//
//	@summary	Аутентификация пользователя
//	@tags		users
//	@param		input	body	AuthenticateUserPayload	true	"Email и пароль пользователя"
//	@accept		json
//	@produce	json
//	@success	200	{object}	TokenResponse
//	@failure	400	{object}	ErrorResponse
//	@failure	401	{object}	ErrorResponse
//	@router		/login [post]
func (h *Handler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var payload AuthenticateUserPayload
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

// GetUserListings - Получение объявлений пользователя с пагинацией
//
//	@summary	Получение объявлений пользователя с пагинацией
//	@tags		users
//	@param		userID	path	int	true	"Идентификатор пользователя"
//	@param		page	query	int	true	"Страница"							default(1)
//	@param		size	query	int	true	"Количество объявлений на страницу"	default(25)
//	@produce	json
//	@success	200	{object}	PaginatedListingResponse
//	@failure	400	{object}	ErrorResponse
//	@failure	401	{object}	ErrorResponse
//	@failure	403	{object}	ErrorResponse
//	@security	BearerAuth
//	@router		/users/{userID}/listings [get]
func (h *Handler) GetUserListings(w http.ResponseWriter, r *http.Request) {
	token := contextutil.GetToken(r)
	if token == "" {
		sendMissingTokenError(w)
		return
	}

	userID, err := readUserIDParam(r)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	page, size, err := readPaginationParams(r)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.listingConn.GetListingsByUser(ctx, &pbListing.GetListingsByUserRequest{
		Pagination: &pbListing.PaginationRequest{
			Page: page,
			Size: size,
		},
		Authentication: &pbListing.Authentication{
			Token: token,
		},
		UserId: userID,
	})
	if err != nil {
		h.handleGRPCError(w, err, func(code codes.Code, msg string) {
			switch code {
			case codes.Unauthenticated:
				sendUnauthorizedError(w)
			case codes.PermissionDenied:
				sendForbiddenError(w)
			default:
				h.log.Error(msg, "code", code)
				sendInternalError(w)
			}
		})
		return
	}

	writePaginatedListingResponse(w, resp.Items, resp.Pagination)
}

// GetUserListing - Получение данных объявления пользователя
//
//	@summary	Получение данных объявления пользователя
//	@tags		users
//	@param		userID	path	int	true	"Идентификатор пользователя"
//	@param		id		path	int	true	"Идентификатор объявления"
//	@produce	json
//	@success	200	{object}	ListingResponse
//	@failure	400	{object}	ErrorResponse
//	@failure	401	{object}	ErrorResponse
//	@failure	403	{object}	ErrorResponse
//	@failure	404	{object}	ErrorResponse
//	@security	BearerAuth
//	@router		/users/{userID}/listings/{id} [get]
func (h *Handler) GetUserListing(w http.ResponseWriter, r *http.Request) {
	token := contextutil.GetToken(r)
	if token == "" {
		sendMissingTokenError(w)
		return
	}

	userID, err := readUserIDParam(r)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := readIDParam(r)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.listingConn.GetUserListing(
		ctx,
		&pbListing.GetUserListingRequest{Id: id, UserId: userID, Authentication: &pbListing.Authentication{Token: token}},
	)
	if err != nil {
		h.handleGRPCError(w, err, func(code codes.Code, msg string) {
			switch code {
			case codes.Unauthenticated:
				sendUnauthorizedError(w)
			case codes.NotFound:
				sendError(w, http.StatusNotFound, msg)
			case codes.PermissionDenied:
				sendForbiddenError(w)
			default:
				h.log.Error(msg, "code", code)
				sendInternalError(w)
			}
		})
		return
	}

	writeListingResponse(w, resp.GetListing())
}
