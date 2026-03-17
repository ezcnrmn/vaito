package handler

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var lettersAndDigitsRX = regexp.MustCompile("^[a-zA-Z0-9]*$")

var humanReadableMessages = map[string]string{
	"lettersAndDigits": "letters and digits only",
}

func lettersAndDigits(fl validator.FieldLevel) bool {
	return lettersAndDigitsRX.MatchString(fl.Field().String())
}
