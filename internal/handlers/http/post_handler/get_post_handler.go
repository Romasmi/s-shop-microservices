package post_handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Romasmi/s-shop-microservices/internal/repository"
	"github.com/Romasmi/s-shop-microservices/internal/utils/http_utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (h *PostHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postIdStr := vars["postId"]
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		http_utils.ErrorInvalidRequestBody(w, fmt.Errorf("invalid post id"))
		return
	}

	post, err := h.postService.GetPost(r.Context(), postId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http_utils.JsonErrorNotFound(w)
			return
		}
		fmt.Printf("error while getting post: %v\n", err)
		http_utils.JsonInternalServerError(w)
		return
	}

	http_utils.SuccessJsonResponse(w, post)
}

func (h *PostHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	profileId := uuid.MustParse(r.Context().Value("profileId").(string))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	posts, err := h.postService.GetFeed(r.Context(), profileId, limit, offset)
	if err != nil {
		fmt.Printf("error while getting feed: %v\n", err)
		http_utils.JsonInternalServerError(w)
		return
	}
	http_utils.SuccessJsonResponse(w, posts)
}
