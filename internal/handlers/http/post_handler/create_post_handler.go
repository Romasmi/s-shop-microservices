package post_handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Romasmi/s-shop-microservices/internal/utils/http_utils"
	"github.com/google/uuid"
)

type CreatePostRequest struct {
	Text string `json:"text"`
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	profileId, _ := uuid.Parse(r.Context().Value("profileId").(string))

	var payload CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http_utils.ErrorInvalidRequestBody(w, err)
		return
	}
	if payload.Text == "" {
		http_utils.ErrorInvalidRequestBody(w, fmt.Errorf("text is required"))
		return
	}

	post, err := h.postService.CreatePost(r.Context(), profileId, payload.Text)
	if err != nil {
		fmt.Printf("error while creating post: %v\n", err)
		http_utils.JsonInternalServerError(w)
		return
	}
	http_utils.SuccessJsonResponse(w, post)
}
