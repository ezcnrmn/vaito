package contextutil

import (
	"context"
	"net/http"
)

const tokenKey = "token"

func SetToken(r *http.Request, token []byte) *http.Request {
	ctx := context.WithValue(r.Context(), tokenKey, token)
	return r.WithContext(ctx)
}

func GetToken(r *http.Request) []byte {
	token, ok := r.Context().Value(tokenKey).([]byte)
	if !ok {
		panic("missing token in request context")
	}
	return token
}
