package user_handler

import (
	"fmt"
	"net/http"

	"github.com/Romasmi/s-shop-microservices/internal/utils/http_utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (h *UserHandler) SetFriend(w http.ResponseWriter, r *http.Request) {
	profileId, _ := uuid.Parse(r.Context().Value("profileId").(string))

	vars := mux.Vars(r)
	friendId, err := uuid.Parse(vars["userId"])
	if err != nil {
		http_utils.ErrorInvalidRequestBody(w, fmt.Errorf("invalid user id"))
		return
	}

	if profileId == friendId {
		http_utils.JsonError(w, http.StatusBadRequest, fmt.Errorf("current userId equels to friendId"))
		return
	}

	err = h.userService.SetFriend(r.Context(), profileId, friendId)
	if err != nil {
		fmt.Printf("error settings a friend: %v\n", err)
		http_utils.JsonInternalServerError(w)
		return
	}

	http_utils.SuccessJsonResponse(w, struct{}{})
}

func (h *UserHandler) DeleteFriend(w http.ResponseWriter, r *http.Request) {
	profileId, _ := uuid.Parse(r.Context().Value("profileId").(string))

	vars := mux.Vars(r)
	friendId, err := uuid.Parse(vars["userId"])
	if err != nil {
		http_utils.ErrorInvalidRequestBody(w, fmt.Errorf("invalid user id"))
		return
	}

	if profileId == friendId {
		http_utils.JsonError(w, http.StatusBadRequest, fmt.Errorf("current userId equels to friendId"))
		return
	}

	err = h.userService.DeleteFriend(r.Context(), profileId, friendId)
	if err != nil {
		fmt.Printf("error deleting a friend: %v\n", err)
		http_utils.JsonInternalServerError(w)
		return
	}

	http_utils.SuccessJsonResponse(w, struct{}{})
}
