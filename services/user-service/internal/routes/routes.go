package routes

import (
	"fmt"
	"net/http"

	"github.com/Romasmi/s-shop-microservices/internal/config"
	"github.com/Romasmi/s-shop-microservices/internal/infra/database"
	"github.com/Romasmi/s-shop-microservices/internal/middleware"
	"github.com/Romasmi/s-shop-microservices/internal/utils/http_utils"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type App interface {
	GetDB() *database.Connection
	GetConfig() *config.Config
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
	router.StrictSlash(true)
	router.Use(middleware.MetricsMiddleware)
	router.Use(middleware.ResponseHeadersMiddleware)
	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(MethodNotAllowedHandler)

	router.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		http_utils.SuccessJsonResponse(w, map[string]string{"status": "OK"})
	}).Methods(http.MethodGet)

	RegisterUserRoutes(router, app.GetDB().DB, app.GetConfig())
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http_utils.JsonError(w, http.StatusNotFound, fmt.Errorf("undefined route"))
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	http_utils.JsonError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
}
