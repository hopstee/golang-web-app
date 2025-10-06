package api1

import (
	customMiddleware "mobile-backend-boilerplate/internal/transport/http/api1/middleware"
	"mobile-backend-boilerplate/internal/transport/http/options"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(opts options.Options) *chi.Mux {
	r := chi.NewRouter()

	r.With(customMiddleware.CorsMiddleware).Group(func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", opts.MobileAuthHandler.Login)
			r.Post("/register", opts.MobileAuthHandler.Register)
			r.Post("/refresh", opts.MobileAuthHandler.Refresh)
			r.Post("/logout", opts.MobileAuthHandler.Logout)

			r.With(customMiddleware.JWTAuthMiddleware(opts.MobileAuthService)).Group(func(r chi.Router) {
				r.Get("/me", opts.MobileAuthHandler.MeMobile)
			})
		})
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	return r
}
