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

func readUserIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("userID"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func readPaginationParams(r *http.Request) (page, size int32, err error) {
	pageParam, sizeParam := r.URL.Query().Get("page"), r.URL.Query().Get("size")
	if pageParam == "" {
		page = 1
	} else {
		parsed, err := strconv.ParseInt(pageParam, 10, 32)
		if err != nil {
			return 0, 0, errors.New("invalid page parameter")
		}
		page = int32(parsed)
		if page < 1 {
			return 0, 0, errors.New("invalid page parameter")
		}
	}

	if sizeParam == "" {
		size = 20
	} else {
		parsed, err := strconv.ParseInt(sizeParam, 10, 32)
		if err != nil {
			return 0, 0, errors.New("invalid size parameter")
		}
		size = int32(parsed)
		if size < 1 || size > 100 {
			return 0, 0, errors.New("invalid size parameter")
		}
	}

	return page, size, nil
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

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	RoleName string `json:"role_name"`
}
type UserResponse struct {
	User User `json:"user"`
}

func writeUserResponse(w http.ResponseWriter, userResponse *pbUser.User) {
	user := User{
		ID:       userResponse.GetId(),
		Name:     userResponse.GetName(),
		Email:    userResponse.GetEmail(),
		RoleName: userResponse.GetRoleName(),
	}

	jsonutil.WriteJSON(w, http.StatusOK, jsonutil.Envelope{"user": user})
}

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
type CategoriesResponse struct {
	Categories []Category `json:"categories"`
}

func writeCategoriesResponse(w http.ResponseWriter, categoryResponse []*pbListing.Category) {
	categories := make([]Category, 0, len(categoryResponse))

	for _, c := range categoryResponse {
		categories = append(categories, Category{ID: c.GetId(), Name: c.GetName()})
	}

	jsonutil.WriteJSON(w, http.StatusOK, jsonutil.Envelope{"categories": categories})
}

type Listing struct {
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
type ListingResponse struct {
	Listing Listing `json:"listing"`
}

func writeListingResponse(w http.ResponseWriter, listingResponse *pbListing.Listing) {
	listing := Listing{
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

type Pagination struct {
	Page  int32 `json:"page"`
	Size  int32 `json:"size"`
	Total int32 `json:"total"`
}

type PaginatedListingResponse struct {
	Items      []Listing  `json:"items"`
	Pagination Pagination `json:"pagination"`
}

func writePaginatedListingResponse(w http.ResponseWriter, listingsResponse []*pbListing.Listing, paginationResponse *pbListing.PaginationResponse) {
	listings := make([]Listing, 0, len(listingsResponse))
	for _, l := range listingsResponse {
		listing := Listing{
			ID:           l.GetId(),
			Title:        l.GetTitle(),
			Description:  l.GetDescription(),
			CategoryID:   l.Category.GetId(),
			CategoryName: l.Category.GetName(),
			UserID:       l.GetUserId(),
			Status:       l.GetStatus(),
			Price:        l.GetPrice(),
		}
		if l.CreatedAt != nil {
			t := l.GetCreatedAt().AsTime()
			listing.CreatedAt = &t
		}
		if l.PublishedAt != nil {
			t := l.GetPublishedAt().AsTime()
			listing.PublishedAt = &t
		}
		listings = append(listings, listing)
	}

	jsonutil.WriteJSON(w, http.StatusOK, jsonutil.Envelope{
		"items": listings,
		"pagination": Pagination{
			Page:  paginationResponse.GetPage(),
			Size:  paginationResponse.GetSize(),
			Total: paginationResponse.GetTotal(),
		},
	})
}

type MessageResponse struct {
	Message string `json:"message"`
}

func sendSuccessMessage(w http.ResponseWriter, message string) {
	data := jsonutil.Envelope{
		"message": message,
	}
	jsonutil.WriteJSON(w, http.StatusOK, data)
}
