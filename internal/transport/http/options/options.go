package options

import (
	"log/slog"
	"mobile-backend-boilerplate/internal/service"
	apiHandlers "mobile-backend-boilerplate/internal/transport/http/api1/handler"
	webHandlers "mobile-backend-boilerplate/internal/transport/http/web/handler"
)

type Options struct {
	StaticDir         string
	MobileAuthHandler *apiHandlers.MobileAuthHandler
	WebAuthHandler    *webHandlers.WebAuthHandler
	RequestHandler    *apiHandlers.RequestHandler
	PostHandler       *webHandlers.PostHandler
	ContactHandler    *webHandlers.ContactHandler
	MobileAuthService *service.MobileAuthService
	WebAuthService    *service.WebAuthService
	PostService       *service.PostService
	Logger            *slog.Logger
}
