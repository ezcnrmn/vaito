package handler

import "net/http"

func (h *Handler) ListListings(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) ShowListings(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) CreateListings(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		CategoryID  int    `json:"categoryId"`
		Price       int    `json:"price"`
	}
	_ = payload
}

func (h *Handler) UpdateListings(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) DeleteListings(w http.ResponseWriter, r *http.Request) {}
