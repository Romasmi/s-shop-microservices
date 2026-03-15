package routes

import (
	"net/http"

	"github.com/Romasmi/s-shop-microservices/internal/handlers/http/user_handler"
	"github.com/Romasmi/s-shop-microservices/internal/repository"
	"github.com/Romasmi/s-shop-microservices/internal/services"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserRoutes(r *mux.Router, db *pgxpool.Pool) {
	userRepo := repository.CreateUserRepository(db)
	userService := services.CreateUserService(userRepo)
	userHandler := user_handler.NewUserHandler(userService)

	r.HandleFunc("/user", userHandler.CreateUserHandler).Methods(http.MethodPost)
	r.HandleFunc("/user/{userId}", userHandler.GetUserHandler).Methods(http.MethodGet)
	r.HandleFunc("/user/{userId}", userHandler.UpdateUserHandler).Methods(http.MethodPut)
	r.HandleFunc("/user/{userId}", userHandler.DeleteUserHandler).Methods(http.MethodDelete)
}
