package user_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Romasmi/s-shop-microservices/internal/domain/user"
	"github.com/Romasmi/s-shop-microservices/internal/repository"
	"github.com/Romasmi/s-shop-microservices/internal/utils/http_utils"
)

type CreateUserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	CreatedAt string `json:"createdAt"`
}

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload user.User
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http_utils.ErrorInvalidRequestBody(w, err)
		return
	}

	newUser, err := h.userService.CreateUser(r.Context(), &payload)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			http_utils.JsonError(w, http.StatusBadRequest, fmt.Errorf("user already exists"))
			return
		}
		slog.Error("error while user creation", "error", err)
		http_utils.JsonInternalServerError(w)
		return
	}

	http_utils.SuccessJsonResponse(w, CreateUserResponse{
		ID:        newUser.ID.String(),
		Username:  newUser.Username,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Email:     newUser.Email,
		Phone:     newUser.Phone,
		Password:  newUser.Password,
		CreatedAt: newUser.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}
