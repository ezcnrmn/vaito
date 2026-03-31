package handler

import "net/http"

func (h *Handler) CreateListing(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		CategoryID  int    `json:"categoryId"`
		Price       int    `json:"price"`
	}
	_ = payload
}

func (h *Handler) UpdateListing(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) GetListing(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) DeleteListing(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) SendListingToModeration(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) ActivateListing(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) DeactivateListing(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) GetListings(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) GetListingCategories(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) GetUserListings(w http.ResponseWriter, r *http.Request) {
}
