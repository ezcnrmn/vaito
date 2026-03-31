package handler

import (
	"errors"
	"net/http"
	"strconv"

	pbUser "github.com/ezcnrmn/vaito/gen/go/user"
	"github.com/ezcnrmn/vaito/services/gateway/internal/lib/jsonutil"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (h *Handler) handleGRPCError(w http.ResponseWriter, err error, handler func(code codes.Code, msg string)) {
	if s, ok := status.FromError(err); ok {
		code := s.Code()
		msg := s.Message()
		handler(code, msg)
		return
	}
	h.log.Error(err.Error())
	sendInternalError(w)
}

type user struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	RoleName string `json:"role_name"`
}

func writeUserResponse(w http.ResponseWriter, userResponse *pbUser.User) {
	user := user{
		ID:       userResponse.GetId(),
		Name:     userResponse.GetName(),
		Email:    userResponse.GetEmail(),
		RoleName: userResponse.GetRoleName(),
	}

	jsonutil.WriteJSON(w, http.StatusOK, jsonutil.Envelope{"user": user})
}

func sendSuccessMessage(w http.ResponseWriter, message string) {
	data := jsonutil.Envelope{
		"message": message,
	}
	jsonutil.WriteJSON(w, http.StatusOK, data)
}
