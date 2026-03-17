package app

import (
	"net/http"

	"github.com/ezcnrmn/vaito/services/gateway/internal/handler"
	"github.com/julienschmidt/httprouter"
)

func (a *App) routes() http.Handler {
	handler := handler.New(a.cfg, a.log, a.grpc.user, a.grpc.listing)
	routes := httprouter.New()

	routes.HandlerFunc(http.MethodGet, "/v1/healthcheck", handler.Healthcheck)

	routes.HandlerFunc(http.MethodPost, "/v1/users", handler.CreateUser)
	routes.HandlerFunc(http.MethodPut, "/v1/users/:id", handler.UpdateUser)
	routes.HandlerFunc(http.MethodGet, "/v1/users/authenticate", handler.Authenticate)

	routes.HandlerFunc(http.MethodPost, "/v1/listings", handler.CreateListings)
	routes.HandlerFunc(http.MethodGet, "/v1/listings/:id", handler.ShowListings)
	routes.HandlerFunc(http.MethodPut, "/v1/listings/:id", handler.UpdateListings)
	routes.HandlerFunc(http.MethodDelete, "/v1/listings/:id", handler.DeleteListings)
	routes.HandlerFunc(http.MethodGet, "/v1/listings", handler.ListListings)

	return a.recoverPanic(routes)
}

func (a *App) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
