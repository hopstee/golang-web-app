package infrastructure

import (
	"log/slog"
	"mobile-backend-boilerplate/internal/config"
	"mobile-backend-boilerplate/internal/filestorage"
	"mobile-backend-boilerplate/internal/kvstore"
	"mobile-backend-boilerplate/internal/notifier"
	"mobile-backend-boilerplate/internal/repository"
	"mobile-backend-boilerplate/internal/service"
	apiHandler "mobile-backend-boilerplate/internal/transport/http/api1/handler"
	webHandler "mobile-backend-boilerplate/internal/transport/http/web/handler"
)

type Dependencies struct {
	Logger *slog.Logger
	Config *config.Config

	// repo
	Repository repository.Repository

	// kvstore
	KVStore kvstore.KVStore

	// filestorage
	FileStorage filestorage.FileStorage

	// notifiers
	TelegramNotifier notifier.Notifier

	// services
	MobileAuthService   *service.MobileAuthService
	AdminAuthService    *service.AdminAuthService
	UserService         *service.UserService
	AdminService        *service.AdminService
	RequestService      *service.RequestService
	PostService         *service.PostService
	SchemaEntityService *service.SchemaEntityService

	// api handlers
	MobileAuthHandler   *apiHandler.MobileAuthHandler
	AdminAuthHandler    *apiHandler.AdminAuthHandler
	RequestHandler      *apiHandler.RequestHandler
	SchemaEntityHandler *apiHandler.SchemaEntityHandler
	FilesHandler        *apiHandler.FilesHandler

	// web handlers
	PostHandler       *webHandler.PostHandler
	ContactHandler    *webHandler.ContactHandler
	StaticPageHandler *webHandler.StaticPageHandler
	AdminHandler      *webHandler.AdminHandler
}
