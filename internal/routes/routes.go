package routes

import (
	"fmt"
	"net/http"

	"github.com/Romasmi/s-shop-microservices/internal/events/publisher"
	"github.com/Romasmi/s-shop-microservices/internal/infra/database"
	"github.com/Romasmi/s-shop-microservices/internal/infra/redis_client"
	"github.com/Romasmi/s-shop-microservices/internal/middleware"
	"github.com/Romasmi/s-shop-microservices/internal/utils/http_utils"
	"github.com/gorilla/mux"
)

type App interface {
	GetDB() *database.Connection
	GetRedis() *redis_client.Connection
}

type NotFoundResponse struct {
	Error string `json:"error"`
}

func RegisterRoutes(
	router *mux.Router,
	app App,
) {
	if router == nil {
		panic("router must be initialized before routes registration")
	}
	router.Use(middleware.ResponseHeadersMiddleware)
	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	RegisterUserRoutes(router, app.GetDB().DB)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	http_utils.JsonError(w, http.StatusNotFound, fmt.Errorf("undefined route"))
}
