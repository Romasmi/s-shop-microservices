package user_handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Romasmi/s-shop-microservices/internal/domain/profile"
	"github.com/Romasmi/s-shop-microservices/internal/utils/http_utils"
)

type CreateUserRequest struct {
	FirstName  string         `json:"first_name"`
	SecondName string         `json:"second_name"`
	Birthdate  string         `json:"birthdate"` // format 2017-02-01
	Biography  string         `json:"biography"`
	Gender     profile.Gender `json:"gender"`
	City       string         `json:"city"`
	Password   string         `json:"password"`
}

func (payload *CreateUserRequest) toModel() (*profile.CreateProfileModel, error) {
	birthdate, err := time.Parse("2006-01-02", payload.Birthdate)
	if err != nil {
		return nil, err
	}

	return &profile.CreateProfileModel{
		FirstName:  payload.FirstName,
		SecondName: payload.SecondName,
		Birthdate:  birthdate,
		Biography:  payload.Biography,
		Gender:     payload.Gender,
		City:       payload.City,
		Password:   payload.Password,
	}, nil
}

// RegisterUserHandler Note: according to the specification, there is not unique identifier for the user like email
// so on each request it creates a new user
func (h *UserHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http_utils.ErrorInvalidRequestBody(w, err)
		return
	}

	creatProfileModel, err := payload.toModel()
	if err != nil {
		http_utils.ErrorInvalidRequestBody(w, err)
		return
	}

	profile, err := h.userService.RegisterUser(r.Context(), creatProfileModel)
	if err != nil {
		fmt.Printf("error while user registration: %v\n", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	http_utils.SuccessJsonResponse(w, profile)
}
