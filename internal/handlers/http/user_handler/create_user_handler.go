package user_handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Romasmi/s-shop-microservices/internal/domain/user"
	"github.com/Romasmi/s-shop-microservices/internal/utils/http_utils"
)

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload user.User
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http_utils.ErrorInvalidRequestBody(w, err)
		return
	}

	newUser, err := h.userService.CreateUser(r.Context(), &payload)
	if err != nil {
		fmt.Printf("error while user creation: %v\n", err)
		http_utils.JsonInternalServerError(w)
		return
	}

	http_utils.SuccessJsonResponse(w, newUser)
}
