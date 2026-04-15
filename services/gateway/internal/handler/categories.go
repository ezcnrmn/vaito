package handler

import (
	"context"
	"net/http"
	"time"

	pbListing "github.com/ezcnrmn/vaito/gen/go/listing"
	"google.golang.org/grpc/codes"
)

// GetListingCategories - Получение категорий
//
//	@summary	Получение категорий
//	@tags		categories
//	@produce	json
//	@success	200	{object}	CategoriesResponse
//	@router		/categories [get]
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
