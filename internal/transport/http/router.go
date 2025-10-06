package http

import (
	"mobile-backend-boilerplate/internal/transport/http/api1"
	"mobile-backend-boilerplate/internal/transport/http/options"
	"mobile-backend-boilerplate/internal/transport/http/web"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(opts options.Options) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	apiRoutes := api1.NewRouter(opts)
	r.Mount("/api", apiRoutes)

	webRoutes := web.NewRouter(opts)
	r.Mount("/", webRoutes)

	return r
}
