package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func sendError(w http.ResponseWriter, err error) {
	data := envelope{
		"error": err.Error(),
	}
	writeJSON(w, http.StatusBadRequest, data)
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
	}

	writeJSON(w, http.StatusBadRequest, data)
}
