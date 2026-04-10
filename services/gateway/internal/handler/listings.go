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

func (h *Handler) CreateListing(w http.ResponseWriter, r *http.Request) {
	token := contextutil.GetToken(r)
	if token == "" {
		sendMissingTokenError(w)
		return
	}

	var payload struct {
		Title       string `json:"title" validate:"required,min=10,max=300"`
		Description string `json:"description" validate:"required"`
		CategoryID  int64  `json:"category_id" validate:"required"`
		Price       int64  `json:"price" validate:"required,min=1"`
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

	var payload struct {
		Title       *string `json:"title" validate:"omitempty,min=10,max=300"`
		Description *string `json:"description" validate:"omitempty"`
		CategoryID  *int64  `json:"category_id" validate:"omitempty"`
		Price       *int64  `json:"price" validate:"omitempty,min=1"`
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

func (h *Handler) GetListingCategories(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.listingConn.GetCategories(ctx, &pbListing.GetCategoriesRequest{})
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

	writeCategoriesResponse(w, resp.GetCategories())
}

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
