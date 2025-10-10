package web

import (
	"mobile-backend-boilerplate/internal/transport/http/options"
	"mobile-backend-boilerplate/internal/transport/http/web/handler"
	"mobile-backend-boilerplate/internal/transport/http/web/middleware"
	"mobile-backend-boilerplate/internal/view/layouts"
	public_pages "mobile-backend-boilerplate/internal/view/pages/public"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

func NewRouter(opts options.Options) *chi.Mux {
	r := chi.NewRouter()

	fs := http.FileServer(http.Dir(opts.StaticDir))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("User-agent: *\nAllow: /\n"))
	})

	r.Get("/unsupported_browser.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(opts.StaticDir, "unsupported_browser.js"))
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		data := layouts.NewBaseLayoutProps(r)
		data.Centered = true

		handler.HandleStaticPage(w, r, public_pages.IndexPage(data), public_pages.IndexPageContent(data))
	})

	r.Get("/about", func(w http.ResponseWriter, r *http.Request) {
		data := layouts.NewBaseLayoutProps(r)
		handler.HandleStaticPage(w, r, public_pages.AboutPage(data), public_pages.AboutPageContent(data))
	})

	r.Get("/projects", func(w http.ResponseWriter, r *http.Request) {
		data := layouts.NewBaseLayoutProps(r)
		handler.HandleStaticPage(w, r, public_pages.ProjectsPage(data), public_pages.ProjectsPageContent(data))
	})

	r.Route("/contact", func(r chi.Router) {
		r.Get("/", opts.ContactHandler.Show)
		r.Post("/", opts.ContactHandler.Submit)
	})

	r.Get("/blog", func(w http.ResponseWriter, r *http.Request) {
		data := layouts.NewBaseLayoutProps(r)
		data.WideWrapper = true

		handler.HandleStaticPage(w, r, public_pages.BlogPage(data), public_pages.BlogPageContent(data))
	})

	r.Route("/auth", func(r chi.Router) {
		r.Route("/login", func(r chi.Router) {
			r.Get("/", opts.WebAuthHandler.Show)
			r.Post("/", opts.WebAuthHandler.Submit)
		})

		r.Get("/logout", opts.WebAuthHandler.Logout)
	})

	r.With(middleware.AuthMiddleware(opts.WebAuthService)).Group(func(r chi.Router) {
		r.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte{'t', 'e', 's', 't'})
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		filePath := filepath.Join(opts.StaticDir, r.URL.Path)
		if _, err := os.Stat(filePath); err == nil {
			http.ServeFile(w, r, filePath)
			return
		}

		data := layouts.NewBaseLayoutProps(r)
		data.Centered = true

		handler.HandleStaticPage(w, r, public_pages.NotFoundPage(data), public_pages.NotFoundPageContent(data))
	})

	return r
}
