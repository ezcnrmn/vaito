package handler

import "net/http"

func (h *Handler) Healthcheck(w http.ResponseWriter, r *http.Request) {
	data := envelope{
		"status": "available",
	}

	err := writeJSON(w, http.StatusOK, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
