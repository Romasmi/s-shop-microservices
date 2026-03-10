package auth_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Romasmi/s-shop-microservices/internal/services"
	"github.com/Romasmi/s-shop-microservices/internal/utils/http_utils"
	"github.com/google/uuid"
)

type AuthHandler struct {
	AuthService *services.AuthService
}

type LoginRequest struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

func CreateAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var payload LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http_utils.ErrorInvalidRequestBody(w, err)
		return
	}
	profileId, err := uuid.Parse(payload.ID)
	if err != nil {
		http_utils.ErrorInvalidRequestBody(w, err)
		return
	}

	session, err := h.AuthService.LoginUser(r.Context(), profileId, payload.Password)
	if err != nil {
		if errors.Is(err, services.InvalidCredentialsError) {
			http_utils.JsonError(w, http.StatusForbidden, err)
			return
		}

		fmt.Printf("error while login user: %v\n", err)
		http_utils.JsonInternalServerError(w)
		return
	}
	http_utils.SuccessJsonResponse(w, session)
}
