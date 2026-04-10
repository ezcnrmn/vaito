package handler

import (
	"context"
	"net/http"
	"time"

	pbListing "github.com/ezcnrmn/vaito/gen/go/listing"
	"github.com/ezcnrmn/vaito/services/gateway/internal/lib/contextutil"
	"google.golang.org/grpc/codes"
)

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
