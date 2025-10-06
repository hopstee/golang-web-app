package web

import (
	"log/slog"
	"mobile-backend-boilerplate/internal/service"

	"github.com/go-chi/chi/v5"
)

func InitPostRoutes(r chi.Router, staticDir string, logger *slog.Logger, postService *service.PostService) {
	// postHandler := handler.NewPostHandler(
	// 	staticDir,
	// 	logger,
	// 	postService,
	// )

	r.Route("/blog/post", func(r chi.Router) {
		// r.Get("/{id}", postHandler.RenderPost)
	})
}
