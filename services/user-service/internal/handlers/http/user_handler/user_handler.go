package user_handler

import "github.com/Romasmi/s-shop-microservices/internal/services"

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}
