package user_handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Romasmi/s-shop-microservices/internal/domain/profile"
	"github.com/Romasmi/s-shop-microservices/internal/repository"
	"github.com/Romasmi/s-shop-microservices/internal/utils/http_utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type userSearchParameters struct {
	firstName  string
	secondName string
}

func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := uuid.Parse(vars["userId"])
	if err != nil {
		http_utils.ErrorInvalidRequestBody(w, fmt.Errorf("invalid user id"))
		return
	}

	profile, err := h.userService.GetUserByProfileId(r.Context(), userId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http_utils.JsonErrorNotFound(w)
			return
		}

		fmt.Printf("error while retreiving a user: %v\n", err)
		http_utils.JsonInternalServerError(w)
		return
	}

	http_utils.SuccessJsonResponse(w, profile)
}

func (h *UserHandler) SearchUserHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	userSearchData := &profile.UserSearchParams{
		FirstName:  queryParams.Get("first_name"),
		SecondName: queryParams.Get("last_name"),
	}

	profiles, err := h.userService.SearchUsers(r.Context(), userSearchData)
	if err != nil {
		fmt.Printf("error while searching users: %v\n", err)
		http_utils.JsonInternalServerError(w)
		return
	}

	http_utils.SuccessJsonResponse(w, profiles)
}
