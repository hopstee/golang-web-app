package web

import (
	"mobile-backend-boilerplate/internal/transport/http/options"
	"mobile-backend-boilerplate/internal/transport/http/web/handler"
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

	r.Get("/", opts.StaticPageHandler.RenderStaticPage)
	r.Get("/about", opts.StaticPageHandler.RenderStaticPage)
	r.Get("/projects", opts.StaticPageHandler.RenderStaticPage)

	r.Route("/contact", func(r chi.Router) {
		r.Get("/", opts.ContactHandler.Show)
		r.Post("/", opts.ContactHandler.Submit)
	})

	r.Get("/blog", func(w http.ResponseWriter, r *http.Request) {
		data := layouts.NewPublicLayoutProps(r)
		data.WideWrapper = true

		handler.HandleStaticPage(w, r, pages.BlogPage(data), pages.BlogPageContent(data))
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

		data := layouts.NewPublicLayoutProps(r)
		data.Centered = true

		handler.HandleStaticPage(w, r, pages.NotFoundPage(data), pages.NotFoundPageContent(data))
	})

	return r
}
