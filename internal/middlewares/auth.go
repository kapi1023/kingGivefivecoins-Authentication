package middlewares

import (
	"net/http"
	"strings"

	"github.com/gorilla/context"
	"github.com/kapi1023/kingGivefivecoins-Authentication/internal/services"
)

func AuthMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := services.ValidateToken(token, secretKey)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			context.Set(r, "user", claims.Email)
			next.ServeHTTP(w, r)
		})
	}
}
