package api1

import (
	customMiddleware "mobile-backend-boilerplate/internal/transport/http/api1/middleware"
	"mobile-backend-boilerplate/internal/transport/http/options"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(opts options.Options) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {
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

			r.Route("/admin", func(r chi.Router) {
				r.Route("/auth", func(r chi.Router) {
					r.Post("/login", opts.AdminAuthHandler.Login)
					r.Post("/logout", opts.AdminAuthHandler.Logout)

					r.With(customMiddleware.AdminAuthMiddleware(opts.AdminAuthService)).Group(func(r chi.Router) {
						r.Get("/me", opts.AdminAuthHandler.MeWeb)
					})
				})

				r.With(customMiddleware.AdminAuthMiddleware(opts.AdminAuthService)).Group(func(r chi.Router) {
					r.Route("/pages", func(r chi.Router) {
						r.Get("/", opts.PagesHandler.GetAllPagesSchemas)
						r.Get("/{slug}/schema", opts.PagesHandler.GetPageSchema)
						r.Get("/{slug}/data", opts.PagesHandler.GetPageData)
						r.Put("/{slug}/data", opts.PagesHandler.UpdatePageData)
						r.Delete("/{slug}", opts.PagesHandler.DeletePage)
					})

					r.Route("/files", func(r chi.Router) {
						r.Post("/", opts.FilesHandler.UploadFile)
						r.Delete("/", opts.FilesHandler.DeleteFile)
					})
				})
			})
		})

		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"status":"ok"}`))
		})
	})

	return r
}
