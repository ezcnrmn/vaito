package handler

import (
	"net/http"

	"github.com/ezcnrmn/vaito/services/gateway/internal/lib/jsonutil"
)

func (h *Handler) Healthcheck(w http.ResponseWriter, r *http.Request) {
	data := jsonutil.Envelope{
		"status": "available",
	}

	err := jsonutil.WriteJSON(w, http.StatusOK, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
