package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ezcnrmn/vaito/services/gateway/internal/lib/jsonutil"
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) NotFound(w http.ResponseWriter, _ *http.Request) {
	msg := "the requested resource could not be found"
	sendError(w, http.StatusNotFound, msg)
}

func (h *Handler) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	sendError(w, http.StatusMethodNotAllowed, msg)
}

func sendError(w http.ResponseWriter, code int, message string) {
	data := jsonutil.Envelope{
		"error": message,
	}
	jsonutil.WriteJSON(w, code, data)
}

func sendInternalError(w http.ResponseWriter) {
	msg := "an unexpected error occurred while processing your request"
	sendError(w, http.StatusInternalServerError, msg)
}

func sendMissingTokenError(w http.ResponseWriter) {
	msg := "missing authentication token"
	sendError(w, http.StatusUnauthorized, msg)
}

func sendUnauthorizedError(w http.ResponseWriter) {
	msg := "you must be authenticated to access this resource"
	sendError(w, http.StatusUnauthorized, msg)
}

func sendForbiddenError(w http.ResponseWriter) {
	msg := "you don't have the necessary permissions to access this resource or to perform this action"
	sendError(w, http.StatusForbidden, msg)
}

func sendValidateError(w http.ResponseWriter, err error) {
	data := jsonutil.Envelope{}

	var validateErrs validator.ValidationErrors
	if errors.As(err, &validateErrs) {
		for _, e := range validateErrs {
			if mes, ok := humanReadableMessages[e.ActualTag()]; ok {
				data[e.Namespace()] = fmt.Sprintf("The %s field must satisfy the condition: %s", e.Field(), mes)
			} else if e.Param() != "" {
				data[e.Namespace()] = fmt.Sprintf("The %s field must satisfy the condition: %s=%s", e.Field(), e.ActualTag(), e.Param())
			} else {
				data[e.Namespace()] = fmt.Sprintf("The %s field must satisfy the condition: %s", e.Field(), e.ActualTag())
			}
		}
	} else {
		panic("wrong error type")
	}

	jsonutil.WriteJSON(w, http.StatusBadRequest, data)
}
