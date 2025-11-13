package options

import (
	"log/slog"
	"mobile-backend-boilerplate/internal/service"
	apiHandlers "mobile-backend-boilerplate/internal/transport/http/api1/handler"
	webHandlers "mobile-backend-boilerplate/internal/transport/http/web/handler"
)

type Options struct {
	StaticDir string
	WebappDir string

	MobileAuthHandler   *apiHandlers.MobileAuthHandler
	AdminAuthHandler    *apiHandlers.AdminAuthHandler
	RequestHandler      *apiHandlers.RequestHandler
	SchemaEntityHandler *apiHandlers.SchemaEntityHandler
	FilesHandler        *apiHandlers.FilesHandler

	PostHandler       *webHandlers.PostHandler
	ContactHandler    *webHandlers.ContactHandler
	StaticPageHandler *webHandlers.StaticPageHandler
	AdminHandler      *webHandlers.AdminHandler

	MobileAuthService   *service.MobileAuthService
	AdminAuthService    *service.AdminAuthService
	PostService         *service.PostService
	SchemaEntityService *service.SchemaEntityService

	Logger *slog.Logger
}
