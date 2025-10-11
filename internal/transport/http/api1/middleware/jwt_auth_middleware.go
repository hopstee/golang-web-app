package middleware

import (
	"context"
	"mobile-backend-boilerplate/internal/service"
	"net/http"
	"strings"
)

type contextKey string

const userIDKey = contextKey("userID")

func JWTAuthMiddleware(authService *service.MobileAuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "missing or invalid authorization header", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")

			user, err := authService.Me(token)
			if err != nil {
				http.Error(w, "invalid access token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, user.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) int64 {
	if val, ok := ctx.Value(userIDKey).(int64); ok {
		return val
	}
	return 0
}
