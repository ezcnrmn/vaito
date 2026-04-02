package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	pbListing "github.com/ezcnrmn/vaito/gen/go/listing"
	pbUser "github.com/ezcnrmn/vaito/gen/go/user"
	"github.com/ezcnrmn/vaito/services/gateway/internal/lib/jsonutil"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (h *Handler) handleGRPCError(w http.ResponseWriter, err error, handler func(code codes.Code, msg string)) {
	if s, ok := status.FromError(err); ok {
		code := s.Code()
		msg := s.Message()
		handler(code, msg)
		return
	}
	h.log.Error(err.Error())
	sendInternalError(w)
}

type user struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	RoleName string `json:"role_name"`
}

func writeUserResponse(w http.ResponseWriter, userResponse *pbUser.User) {
	user := user{
		ID:       userResponse.GetId(),
		Name:     userResponse.GetName(),
		Email:    userResponse.GetEmail(),
		RoleName: userResponse.GetRoleName(),
	}

	jsonutil.WriteJSON(w, http.StatusOK, jsonutil.Envelope{"user": user})
}

type category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func writeCategoriesResponse(w http.ResponseWriter, categoryResponse []*pbListing.Category) {
	categories := make([]category, 0, len(categoryResponse))

	for _, c := range categoryResponse {
		categories = append(categories, category{ID: c.GetId(), Name: c.GetName()})
	}

	jsonutil.WriteJSON(w, http.StatusOK, jsonutil.Envelope{"categories": categories})
}

type listing struct {
	ID           int64      `json:"id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	CategoryID   int64      `json:"category_id"`
	CategoryName string     `json:"category_name"`
	UserID       int64      `json:"user_id"`
	Status       string     `json:"status"`
	Price        int64      `json:"price"`
	CreatedAt    *time.Time `json:"created_at"`
	PublishedAt  *time.Time `json:"published_at"`
}

func writeListingResponse(w http.ResponseWriter, listingResponse *pbListing.Listing) {
	listing := listing{
		ID:           listingResponse.GetId(),
		Title:        listingResponse.GetTitle(),
		Description:  listingResponse.GetDescription(),
		CategoryID:   listingResponse.Category.GetId(),
		CategoryName: listingResponse.Category.GetName(),
		UserID:       listingResponse.GetUserId(),
		Status:       listingResponse.GetStatus(),
		Price:        listingResponse.GetPrice(),
	}
	if listingResponse.CreatedAt != nil {
		t := listingResponse.GetCreatedAt().AsTime()
		listing.CreatedAt = &t
	}
	if listingResponse.PublishedAt != nil {
		t := listingResponse.GetPublishedAt().AsTime()
		listing.PublishedAt = &t
	}

	jsonutil.WriteJSON(w, http.StatusOK, jsonutil.Envelope{"listing": listing})
}

func sendSuccessMessage(w http.ResponseWriter, message string) {
	data := jsonutil.Envelope{
		"message": message,
	}
	jsonutil.WriteJSON(w, http.StatusOK, data)
}
