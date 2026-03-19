package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func sendError(w http.ResponseWriter, code int, message string) {
	data := envelope{
		"error": message,
	}
	writeJSON(w, code, data)
}

func sendInternalError(w http.ResponseWriter) {
	data := envelope{
		"error": "an unexpected error occurred while processing your request",
	}
	writeJSON(w, http.StatusInternalServerError, data)
}

func sendValidateError(w http.ResponseWriter, err error) {
	data := envelope{}

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

	writeJSON(w, http.StatusBadRequest, data)
}
