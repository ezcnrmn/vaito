package handler

import (
	"context"
	"net/http"
	"time"

	pbListing "github.com/ezcnrmn/vaito/gen/go/listing"
	"github.com/ezcnrmn/vaito/services/gateway/internal/lib/contextutil"
	"github.com/ezcnrmn/vaito/services/gateway/internal/lib/jsonutil"
	"google.golang.org/grpc/codes"
)

type CreateListingPayload struct {
	Title       string `json:"title" validate:"required,min=10,max=300" example:"Продам iPhone 3G"`
	Description string `json:"description" validate:"required" example:"Полный комплект: Зарядное устройство. Коробка. Наушники. Чехол в подарок."`
	CategoryID  int64  `json:"category_id" validate:"required" example:"2"`
	Price       int64  `json:"price" validate:"required,min=1" example:"1500"`
}

// CreateListing - Создание объявления
//
//	@summary	Создание объявления
//	@tags		listings
//	@param		input	body	CreateListingPayload	true	"Данные нового объявления"
//	@accept		json
//	@produce	json
//	@success	200	{object}	ListingResponse
//	@failure	400	{object}	ErrorResponse
//	@failure	401	{object}	ErrorResponse
//	@failure	403	{object}	ErrorResponse
//	@security	BearerAuth
//	@router		/listings [post]
func (h *Handler) CreateListing(w http.ResponseWriter, r *http.Request) {
	token := contextutil.GetToken(r)
	if token == "" {
		sendMissingTokenError(w)
		return
	}

	var payload CreateListingPayload

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

	resp, err := h.listingConn.CreateListing(ctx, &pbListing.CreateListingRequest{
		Title:       payload.Title,
		Description: payload.Description,
		CategoryId:  payload.CategoryID,
		Price:       payload.Price,

		Authentication: &pbListing.Authentication{
			Token: token,
		},
	})
	if err != nil {
		h.handleGRPCError(w, err, func(code codes.Code, msg string) {
			switch code {
			case codes.Unauthenticated:
				sendError(w, http.StatusUnauthorized, msg)
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

type UpdateListingPayload struct {
	Title       *string `json:"title" validate:"omitempty,min=10,max=300" example:"Продам iPhone 3G"`
	Description *string `json:"description" validate:"omitempty" example:"Полный комплект: Зарядное устройство. Коробка. Наушники. Чехол в подарок."`
	CategoryID  *int64  `json:"category_id" validate:"omitempty" example:"2"`
	Price       *int64  `json:"price" validate:"omitempty,min=1" example:"1500"`
}

// UpdateListing - Обновление данных объявления
//
//	@summary	Обновление данных объявления
//	@tags		listings
//	@param		id		path	int						true	"Идентификатор объявления"
//	@param		input	body	UpdateListingPayload	true	"Данные объявления"
//	@accept		json
//	@produce	json
//	@success	200	{object}	ListingResponse
//	@failure	400	{object}	ErrorResponse
//	@failure	401	{object}	ErrorResponse
//	@failure	403	{object}	ErrorResponse
//	@failure	404	{object}	ErrorResponse
//	@failure	423	{object}	ErrorResponse
//	@security	BearerAuth
//	@router		/listings/{id} [patch]
func (h *Handler) UpdateListing(w http.ResponseWriter, r *http.Request) {
	token := contextutil.GetToken(r)
	if token == "" {
		sendMissingTokenError(w)
		return
	}

	id, err := readIDParam(r)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	var payload UpdateListingPayload

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
	if payload.Title == nil && payload.Description == nil && payload.CategoryID == nil && payload.Price == nil {
		sendError(w, http.StatusBadRequest, "you must specify at least one field")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.listingConn.UpdateListing(ctx, &pbListing.UpdateListingRequest{
		Id:          id,
		Title:       payload.Title,
		Description: payload.Description,
		CategoryId:  payload.CategoryID,
		Price:       payload.Price,

		Authentication: &pbListing.Authentication{
			Token: token,
		},
	})
	if err != nil {
		h.handleGRPCError(w, err, func(code codes.Code, msg string) {
			switch code {
			case codes.Unauthenticated:
				sendUnauthorizedError(w)
			case codes.PermissionDenied:
				sendForbiddenError(w)
			case codes.NotFound:
				sendError(w, http.StatusNotFound, msg)
			case codes.Aborted:
				sendError(w, http.StatusUnprocessableEntity, msg)
			default:
				h.log.Error(msg, "code", code)
				sendInternalError(w)
			}
		})
		return
	}

	writeListingResponse(w, resp.GetListing())
}

// GetListing - Получение данных объявления (только активные)
//
//	@summary	Получение данных объявления (только активные)
//	@tags		listings
//	@param		id	path	int	true	"Идентификатор объявления"
//	@produce	json
//	@success	200	{object}	ListingResponse
//	@failure	400	{object}	ErrorResponse
//	@failure	404	{object}	ErrorResponse
//	@router		/listings/{id} [get]
func (h *Handler) GetListing(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.listingConn.GetListing(ctx, &pbListing.GetListingRequest{Id: id})
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

	writeListingResponse(w, resp.GetListing())
}

// DeleteListing - Удаление объявления
//
//	@summary	Удаление объявления
//	@tags		listings
//	@param		id	path	int	true	"Идентификатор объявления"
//	@produce	json
//	@success	200	{object}	MessageResponse
//	@failure	400	{object}	ErrorResponse
//	@failure	401	{object}	ErrorResponse
//	@failure	403	{object}	ErrorResponse
//	@failure	404	{object}	ErrorResponse
//	@security	BearerAuth
//	@router		/listings/{id} [delete]
func (h *Handler) DeleteListing(w http.ResponseWriter, r *http.Request) {
	token := contextutil.GetToken(r)
	if token == "" {
		sendMissingTokenError(w)
		return
	}

	id, err := readIDParam(r)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = h.listingConn.DeleteListing(ctx, &pbListing.DeleteListingRequest{
		Id: id,
		Authentication: &pbListing.Authentication{
			Token: token,
		},
	})
	if err != nil {
		h.handleGRPCError(w, err, func(code codes.Code, msg string) {
			switch code {
			case codes.PermissionDenied:
				sendForbiddenError(w)
			case codes.NotFound:
				sendError(w, http.StatusNotFound, msg)
			default:
				h.log.Error(msg, "code", code)
				sendInternalError(w)
			}
		})
		return
	}

	sendSuccessMessage(w, "listing successfully deleted")
}

// SendListingToModeration - Отправка объявления на модерацию
//
//	@summary	Отправка объявления на модерацию
//	@tags		listings
//	@param		id	path	int	true	"Идентификатор объявления"
//	@produce	json
//	@success	200	{object}	MessageResponse
//	@failure	400	{object}	ErrorResponse
//	@failure	401	{object}	ErrorResponse
//	@failure	403	{object}	ErrorResponse
//	@failure	404	{object}	ErrorResponse
//	@security	BearerAuth
//	@router		/listings/{id}/moderation [post]
func (h *Handler) SendListingToModeration(w http.ResponseWriter, r *http.Request) {
	token := contextutil.GetToken(r)
	if token == "" {
		sendMissingTokenError(w)
		return
	}

	id, err := readIDParam(r)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = h.listingConn.SendListingToModeration(ctx, &pbListing.SendListingToModerationRequest{
		Id: id,
		Authentication: &pbListing.Authentication{
			Token: token,
		},
	})
	if err != nil {
		h.handleGRPCError(w, err, func(code codes.Code, msg string) {
			switch code {
			case codes.Unauthenticated:
				sendUnauthorizedError(w)
			case codes.PermissionDenied:
				sendForbiddenError(w)
			case codes.NotFound:
				sendError(w, http.StatusNotFound, msg)
			case codes.InvalidArgument:
				sendError(w, http.StatusBadRequest, msg)
			default:
				h.log.Error(msg, "code", code)
				sendInternalError(w)
			}
		})
		return
	}

	sendSuccessMessage(w, "listing successfully sent to moderation")
}

// ActivateListing - Активация объявления (только неактивные)
//
//	@summary	Активация объявления (только неактивные)
//	@tags		listings
//	@param		id	path	int	true	"Идентификатор объявления"
//	@produce	json
//	@success	200	{object}	MessageResponse
//	@failure	400	{object}	ErrorResponse
//	@failure	401	{object}	ErrorResponse
//	@failure	403	{object}	ErrorResponse
//	@failure	404	{object}	ErrorResponse
//	@security	BearerAuth
//	@router		/listings/{id}/activate [post]
func (h *Handler) ActivateListing(w http.ResponseWriter, r *http.Request) {
	token := contextutil.GetToken(r)
	if token == "" {
		sendMissingTokenError(w)
		return
	}

	id, err := readIDParam(r)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = h.listingConn.ActivateListing(ctx, &pbListing.ActivateListingRequest{
		Id: id,
		Authentication: &pbListing.Authentication{
			Token: token,
		},
	})
	if err != nil {
		h.handleGRPCError(w, err, func(code codes.Code, msg string) {
			switch code {
			case codes.Unauthenticated:
				sendUnauthorizedError(w)
			case codes.PermissionDenied:
				sendForbiddenError(w)
			case codes.NotFound:
				sendError(w, http.StatusNotFound, msg)
			case codes.InvalidArgument:
				sendError(w, http.StatusBadRequest, msg)
			default:
				h.log.Error(msg, "code", code)
				sendInternalError(w)
			}
		})
		return
	}

	sendSuccessMessage(w, "listing successfully activated")
}

// DeactivateListing - Деактивация объявления
//
//	@summary	Деактивация объявления
//	@tags		listings
//	@param		id	path	int	true	"Идентификатор объявления"
//	@produce	json
//	@success	200	{object}	MessageResponse
//	@failure	400	{object}	ErrorResponse
//	@failure	401	{object}	ErrorResponse
//	@failure	403	{object}	ErrorResponse
//	@failure	404	{object}	ErrorResponse
//	@security	BearerAuth
//	@router		/listings/{id}/deactivate [post]
func (h *Handler) DeactivateListing(w http.ResponseWriter, r *http.Request) {
	token := contextutil.GetToken(r)
	if token == "" {
		sendMissingTokenError(w)
		return
	}

	id, err := readIDParam(r)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = h.listingConn.DeactivateListing(ctx, &pbListing.DeactivateListingRequest{
		Id: id,
		Authentication: &pbListing.Authentication{
			Token: token,
		},
	})
	if err != nil {
		h.handleGRPCError(w, err, func(code codes.Code, msg string) {
			switch code {
			case codes.Unauthenticated:
				sendUnauthorizedError(w)
			case codes.PermissionDenied:
				sendForbiddenError(w)
			case codes.NotFound:
				sendError(w, http.StatusNotFound, msg)
			case codes.InvalidArgument:
				sendError(w, http.StatusBadRequest, msg)
			default:
				h.log.Error(msg, "code", code)
				sendInternalError(w)
			}
		})
		return
	}

	sendSuccessMessage(w, "listing successfully deactivated")
}

// GetListings - Получение объявлений с пагинацией (только активные)
//
//	@summary	Получение объявлений с пагинацией (только активные)
//	@tags		listings
//	@param		page	query	int	true	"Страница"							default(1)
//	@param		size	query	int	true	"Количество объявлений на страницу"	default(25)
//	@produce	json
//	@success	200	{object}	PaginatedListingResponse
//	@failure	400	{object}	ErrorResponse
//	@router		/listings [get]
func (h *Handler) GetListings(w http.ResponseWriter, r *http.Request) {
	page, size, err := readPaginationParams(r)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.listingConn.GetActiveListings(ctx, &pbListing.GetActiveListingsRequest{
		Pagination: &pbListing.PaginationRequest{
			Page: page,
			Size: size,
		},
	})
	if err != nil {
		h.handleGRPCError(w, err, func(code codes.Code, msg string) {
			switch code {
			default:
				h.log.Error(msg, "code", code)
				sendInternalError(w)
			}
		})
		return
	}

	writePaginatedListingResponse(w, resp.Items, resp.Pagination)
}
