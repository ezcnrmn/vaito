package handler

import (
	"context"
	"net/http"
	"time"

	pbListing "github.com/ezcnrmn/vaito/gen/go/listing"
	"github.com/ezcnrmn/vaito/services/gateway/internal/lib/contextutil"
	"google.golang.org/grpc/codes"
)

// ModerationListings - Получение объявлений отправленных на модерацию
//
//	@summary	Получение объявлений отправленных на модерацию
//	@tags		moderation
//	@param		page	query	int	true	"Страница"							default(1)
//	@param		size	query	int	true	"Количество объявлений на страницу"	default(25)
//	@produce	json
//	@success	200	{object}	PaginatedListingResponse
//	@failure	400	{object}	ErrorResponse
//	@failure	401	{object}	ErrorResponse
//	@failure	403	{object}	ErrorResponse
//	@security	BearerAuth
//	@router		/moderation/listings [get]
func (h *Handler) ModerationListings(w http.ResponseWriter, r *http.Request) {
	token := contextutil.GetToken(r)
	if token == "" {
		sendMissingTokenError(w)
		return
	}

	page, size, err := readPaginationParams(r)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.listingConn.GetModerationListings(ctx, &pbListing.GetModerationListingsRequest{
		Pagination: &pbListing.PaginationRequest{
			Page: page,
			Size: size,
		},
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
			default:
				h.log.Error(msg, "code", code)
				sendInternalError(w)
			}
		})
		return
	}

	writePaginatedListingResponse(w, resp.Items, resp.Pagination)
}

// ModerationActivateListing - Активация объявления модерацией
//
//	@summary	Активация объявления модерацией
//	@tags		moderation
//	@param		id	path	int	true	"Идентификатор объявления"
//	@produce	json
//	@success	200	{object}	MessageResponse
//	@failure	400	{object}	ErrorResponse
//	@failure	401	{object}	ErrorResponse
//	@failure	403	{object}	ErrorResponse
//	@failure	404	{object}	ErrorResponse
//	@security	BearerAuth
//	@router		/moderation/listings/{id}/activate [post]
func (h *Handler) ModerationActivateListing(w http.ResponseWriter, r *http.Request) {
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

	_, err = h.listingConn.ActivateListingByModeration(ctx, &pbListing.ActivateListingByModerationRequest{
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

// ModerationDeactivateListing - Деактивация объявления модерацией
//
//	@summary	Деактивация объявления модерацией
//	@tags		moderation
//	@param		id	path	int	true	"Идентификатор объявления"
//	@produce	json
//	@success	200	{object}	MessageResponse
//	@failure	400	{object}	ErrorResponse
//	@failure	401	{object}	ErrorResponse
//	@failure	403	{object}	ErrorResponse
//	@failure	404	{object}	ErrorResponse
//	@security	BearerAuth
//	@router		/moderation/listings/{id}/deactivate [post]
func (h *Handler) ModerationDeactivateListing(w http.ResponseWriter, r *http.Request) {
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

	_, err = h.listingConn.DeactivateListingByModeration(ctx, &pbListing.DeactivateListingByModerationRequest{
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
