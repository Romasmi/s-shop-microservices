package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Romasmi/s-shop-microservices/internal/utils/http_utils"
	"github.com/gorilla/mux"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authUserID := r.Header.Get("X-Auth-User-ID")
		if authUserID == "" {
			http_utils.JsonError(w, http.StatusUnauthorized, fmt.Errorf("missing X-Auth-User-ID header"))
			return
		}

		vars := mux.Vars(r)
		pathUserID, ok := vars["userId"]
		if ok && !strings.EqualFold(pathUserID, authUserID) {
			http_utils.JsonError(w, http.StatusForbidden, fmt.Errorf("access denied: you can only access your own data"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
