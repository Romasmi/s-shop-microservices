package routes

import (
	"net/http"

	"github.com/Romasmi/s-shop-microservices/internal/handlers/http/user_handler"
	"github.com/Romasmi/s-shop-microservices/internal/middleware"
	"github.com/Romasmi/s-shop-microservices/internal/repository"
	"github.com/Romasmi/s-shop-microservices/internal/services"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserRoutes(r *mux.Router, db *pgxpool.Pool) {

	userService := services.CreateUserService(cityRepo, profileRepo, profileFriendsRepo, uow, publisher)
	sessionService := services.CreateSessionService(sessionRepo)

	userHandler := user_handler.CreateUserHandler(userService)

	r.HandleFunc("/user/register", userHandler.RegisterUserHandler).Methods(http.MethodPost)

	privateRoute := r.PathPrefix("/").Subrouter()

	authMiddleware := middleware.CreateAuthMiddleware(sessionService)

	privateRoute.Use(authMiddleware.Process)
	privateRoute.HandleFunc("/user/get/{userId}", userHandler.GetUserHandler).Methods(http.MethodGet)
	privateRoute.HandleFunc("/user/search", userHandler.SearchUserHandler).Methods(http.MethodGet)
	privateRoute.HandleFunc("/friend/set/{userId}", userHandler.SetFriend).Methods(http.MethodPut)
	privateRoute.HandleFunc("/friend/delete/{userId}", userHandler.DeleteFriend).Methods(http.MethodPut)
}
