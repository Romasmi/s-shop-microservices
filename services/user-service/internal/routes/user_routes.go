package routes

import (
	"net/http"

	"github.com/Romasmi/s-shop-microservices/user-service/internal/config"
	"github.com/Romasmi/s-shop-microservices/user-service/internal/handlers/http/user_handler"
	"github.com/Romasmi/s-shop-microservices/user-service/internal/middleware"
	"github.com/Romasmi/s-shop-microservices/user-service/internal/repository"
	"github.com/Romasmi/s-shop-microservices/user-service/internal/services"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserRoutes(r *mux.Router, db *pgxpool.Pool, cfg *config.Config) {
	userRepo := repository.CreateUserRepository(db)
	userService := services.CreateUserService(userRepo, cfg.AuthServiceURL)
	userHandler := user_handler.NewUserHandler(userService)

	r.HandleFunc("/user", userHandler.CreateUserHandler).Methods(http.MethodPost)

	authR := r.PathPrefix("/user").Subrouter()
	authR.Use(middleware.AuthMiddleware)
	authR.HandleFunc("/{userId}", userHandler.GetUserHandler).Methods(http.MethodGet)
	authR.HandleFunc("/{userId}", userHandler.UpdateUserHandler).Methods(http.MethodPut)
	authR.HandleFunc("/{userId}", userHandler.DeleteUserHandler).Methods(http.MethodDelete)
}
