package app

import (
	"net/http"
	"strings"

	"github.com/ezcnrmn/vaito/services/gateway/internal/handler"
	"github.com/ezcnrmn/vaito/services/gateway/internal/lib/contextutil"
	"github.com/ezcnrmn/vaito/services/gateway/internal/lib/jsonutil"
	"github.com/julienschmidt/httprouter"
)

const apiV1 = "/api/v1"

func (a *App) routes() http.Handler {
	handler := handler.New(a.log, a.services.user, a.services.listing, a.services.health.user, a.services.health.listing)
	routes := httprouter.New()

	routes.NotFound = http.HandlerFunc(handler.NotFound)
	routes.MethodNotAllowed = http.HandlerFunc(handler.MethodNotAllowed)

	routes.HandlerFunc(http.MethodGet, apiV1+"/healthcheck", handler.Healthcheck)

	routes.HandlerFunc(http.MethodPost, apiV1+"/users", handler.CreateUser)
	routes.HandlerFunc(http.MethodPatch, apiV1+"/users/:userID", handler.UpdateUser)
	routes.HandlerFunc(http.MethodPut, apiV1+"/users/:userID/update-password", handler.UpdateUserPassword)
	routes.HandlerFunc(http.MethodGet, apiV1+"/users/:userID", handler.GetUser)
	routes.HandlerFunc(http.MethodPost, apiV1+"/login", handler.AuthenticateUser)

	routes.HandlerFunc(http.MethodGet, apiV1+"/users/:userID/listings", handler.GetUserListings)
	routes.HandlerFunc(http.MethodGet, apiV1+"/users/:userID/listings/:id", handler.GetUserListing)

	routes.HandlerFunc(http.MethodPost, apiV1+"/listings", handler.CreateListing)
	routes.HandlerFunc(http.MethodPatch, apiV1+"/listings/:id", handler.UpdateListing)
	routes.HandlerFunc(http.MethodGet, apiV1+"/listings/:id", handler.GetListing)
	routes.HandlerFunc(http.MethodDelete, apiV1+"/listings/:id", handler.DeleteListing)
	routes.HandlerFunc(http.MethodPost, apiV1+"/listings/:id/moderation", handler.SendListingToModeration)
	routes.HandlerFunc(http.MethodPost, apiV1+"/listings/:id/activate", handler.ActivateListing)
	routes.HandlerFunc(http.MethodPost, apiV1+"/listings/:id/deactivate", handler.DeactivateListing)

	routes.HandlerFunc(http.MethodGet, apiV1+"/listings", handler.GetListings)
	routes.HandlerFunc(http.MethodGet, apiV1+"/categories", handler.GetListingCategories)

	routes.HandlerFunc(http.MethodGet, apiV1+"/moderation/listings", handler.ModerationListings)
	routes.HandlerFunc(http.MethodPost, apiV1+"/moderation/listings/:id/activate", handler.ModerationActivateListing)
	routes.HandlerFunc(http.MethodPost, apiV1+"/moderation/listings/:id/deactivate", handler.ModerationDeactivateListing)

	return a.recoverPanic(validateToken(routes))
}

func (a *App) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				a.log.Error("panic", "method", r.Method, "path", r.URL.Path, "err", r)
				w.Header().Set("Connection", "close")
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func validateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			r = contextutil.SetToken(r, "")
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		// 32byte token -> 44 length
		if (len(headerParts) != 2) || (len(headerParts) == 2 && (headerParts[0] != "Bearer" || len(headerParts[1]) != 44)) {
			w.Header().Set("WWW-Authenticate", "Bearer")
			msg := jsonutil.Envelope{"error": "invalid authentication token"}
			jsonutil.WriteJSON(w, http.StatusUnauthorized, msg)
			return
		}

		token := headerParts[1]

		r = contextutil.SetToken(r, token)
		next.ServeHTTP(w, r)
	})
}
