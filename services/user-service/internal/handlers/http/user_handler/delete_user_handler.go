package user_handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Romasmi/s-shop-microservices/user-service/internal/repository"
	http2 "github.com/Romasmi/s-shop-microservices/user-service/internal/transport/http"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIdStr := vars["userId"]
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		http2.JsonError(w, http.StatusBadRequest, fmt.Errorf("invalid user id"))
		return
	}

	err = h.userService.DeleteUser(r.Context(), userId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http2.JsonErrorNotFound(w)
			return
		}

		slog.Error("error while user deletion", "error", err)
		http2.JsonInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
