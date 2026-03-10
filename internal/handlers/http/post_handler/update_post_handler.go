package post_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Romasmi/s-shop-microservices/internal/repository"
	"github.com/Romasmi/s-shop-microservices/internal/utils/http_utils"
	"github.com/google/uuid"
)

type UpdatePostRequest struct {
	PostID string `json:"id"`
	Text   string `json:"text"`
}

func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	profileId, _ := uuid.Parse(r.Context().Value("profileId").(string))

	var payload UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http_utils.ErrorInvalidRequestBody(w, err)
		return
	}
	if payload.Text == "" || payload.PostID == "" {
		http_utils.ErrorInvalidRequestBody(w, fmt.Errorf("postId and text are required"))
		return
	}
	postId, err := uuid.Parse(payload.PostID)
	if err != nil {
		http_utils.ErrorInvalidRequestBody(w, fmt.Errorf("invalid postId"))
		return
	}

	post, err := h.postService.UpdatePost(r.Context(), profileId, postId, payload.Text)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http_utils.JsonErrorNotFound(w)
			return
		}
		fmt.Printf("error while updating post: %v\n", err)
		http_utils.JsonInternalServerError(w)
		return
	}
	http_utils.SuccessJsonResponse(w, post)
}
