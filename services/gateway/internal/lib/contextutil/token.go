package contextutil

import (
	"context"
	"net/http"
)

type CustomContextKey string

const tokenKey CustomContextKey = "token"

func SetToken(r *http.Request, token string) *http.Request {
	ctx := context.WithValue(r.Context(), tokenKey, token)
	return r.WithContext(ctx)
}

func GetToken(r *http.Request) string {
	token, ok := r.Context().Value(tokenKey).(string)
	if !ok {
		panic("missing token in request context")
	}
	return token
}
