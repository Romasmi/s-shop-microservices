package post_handler

import "github.com/Romasmi/s-shop-microservices/internal/services"

type PostHandler struct {
	postService *services.PostService
}

func CreatePostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}
