package middleware

import (
	"context"
	"fmt"
	"mobile-backend-boilerplate/internal/service"
	"net/http"
)

type contextKey string

const adminIDKey = contextKey("adminID")

func AuthMiddleware(adminAuthService *service.WebAuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("admin_token")
			fmt.Print(c)
			if err != nil {
				http.Redirect(w, r, "/auth/login", http.StatusFound)
				return
			}

			user, err := adminAuthService.Me(c.Value)
			if err != nil {
				http.Redirect(w, r, "/auth/login", http.StatusFound)
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
