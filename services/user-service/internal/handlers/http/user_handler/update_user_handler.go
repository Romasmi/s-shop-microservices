package user_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Romasmi/s-shop-microservices/user-service/internal/domain/user"
	"github.com/Romasmi/s-shop-microservices/user-service/internal/repository"
	http2 "github.com/Romasmi/s-shop-microservices/user-service/internal/transport/http"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (h *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIdStr := vars["userId"]
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		http2.JsonError(w, http.StatusBadRequest, fmt.Errorf("invalid user id"))
		return
	}

	var payload user.User
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http2.ErrorInvalidRequestBody(w, err)
		return
	}
	payload.ID = userId

	updatedUser, err := h.userService.UpdateUser(r.Context(), &payload)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http2.JsonErrorNotFound(w)
			return
		}
		if errors.Is(err, repository.ErrDuplicate) {
			http2.JsonError(w, http.StatusBadRequest, fmt.Errorf("user already exists"))
			return
		}

		slog.Error("error while user update", "error", err)
		http2.JsonInternalServerError(w)
		return
	}

	http2.SuccessJsonResponse(w, updatedUser)
}
