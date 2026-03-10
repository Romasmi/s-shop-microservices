package user_handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Romasmi/s-shop-microservices/internal/repository"
	"github.com/Romasmi/s-shop-microservices/internal/utils/http_utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIdStr := vars["userId"]
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		http_utils.JsonError(w, http.StatusBadRequest, fmt.Errorf("invalid user id"))
		return
	}

	err = h.userService.DeleteUser(r.Context(), userId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http_utils.JsonErrorNotFound(w)
			return
		}

		fmt.Printf("error while user deletion: %v\n", err)
		http_utils.JsonInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
