package middleware

import (
	"context"
	"mobile-backend-boilerplate/internal/service"
	"net/http"
)

const adminIDKey = contextKey("adminID")

func AdminAuthMiddleware(authService *service.AdminAuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("admin_token")
			if err != nil {
				http.Error(w, "missing or invalid authorization header", http.StatusUnauthorized)
				return
			}

			user, err := authService.Me(c.Value)
			if err != nil {
				http.Error(w, "invalid access token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), adminIDKey, user.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetAdminID(ctx context.Context) int64 {
	if val, ok := ctx.Value(adminIDKey).(int64); ok {
		return val
	}
	return 0
}
