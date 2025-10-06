package web

import (
	"mobile-backend-boilerplate/internal/transport/http/options"
	"mobile-backend-boilerplate/internal/transport/http/web/handler"
	middleware "mobile-backend-boilerplate/internal/transport/http/web/middleware"
	"mobile-backend-boilerplate/internal/view/layouts"
	"mobile-backend-boilerplate/internal/view/pages"
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

		handler.HandleStaticPage(w, r, pages.IndexPage(data), pages.IndexPageContent(data))
	})

	r.Get("/about", func(w http.ResponseWriter, r *http.Request) {
		data := layouts.NewBaseLayoutProps(r)
		handler.HandleStaticPage(w, r, pages.AboutPage(data), pages.AboutPageContent(data))
	})

	r.Get("/projects", func(w http.ResponseWriter, r *http.Request) {
		data := layouts.NewBaseLayoutProps(r)
		handler.HandleStaticPage(w, r, pages.ProjectsPage(data), pages.ProjectsPageContent(data))
	})

	r.Route("/contact", func(r chi.Router) {
		r.Get("/", opts.ContactHandler.Show)
		r.Post("/", opts.ContactHandler.Submit)
	})

	r.Get("/blog", func(w http.ResponseWriter, r *http.Request) {
		data := layouts.NewBaseLayoutProps(r)
		data.WideWrapper = true

		handler.HandleStaticPage(w, r, pages.BlogPage(data), pages.BlogPageContent(data))
	})

	r.With(middleware.CorsMiddleware).Group(func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", opts.WebAuthHandler.Login)
			r.Post("/login", opts.WebAuthHandler.Login)
			r.Post("/logout", opts.WebAuthHandler.Logout)

			r.With(middleware.AuthMiddleware(opts.WebAuthService)).Group(func(r chi.Router) {
				r.Get("/me", opts.WebAuthHandler.MeWeb)
			})
		})

		r.Route("/post", func(r chi.Router) {
			r.Get("/all/public", opts.PostHandler.GetAllPublic)
			r.Get("/public/:id", opts.PostHandler.GetPublic)

			r.With(middleware.AuthMiddleware(opts.WebAuthService)).Group(func(r chi.Router) {
				r.Get("/all", opts.PostHandler.GetAllPosts)
				r.Get("/{id}", opts.PostHandler.GetPost)
				r.Post("/", opts.PostHandler.CreatePost)
				r.Put("/{id}", opts.PostHandler.UpdatePost)
				r.Delete("/{id}", opts.PostHandler.DeletePost)
			})
		})

		r.Post("/submit", opts.RequestHandler.Submit)
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

		handler.HandleStaticPage(w, r, pages.NotFoundPage(data), pages.NotFoundPageContent(data))
	})

	return r
}
