package infrastructure

import (
	"database/sql"
	"log/slog"
	"mobile-backend-boilerplate/internal/config"
	"mobile-backend-boilerplate/internal/notifier"
	"mobile-backend-boilerplate/internal/repository"
	"mobile-backend-boilerplate/internal/service"
	apiHandler "mobile-backend-boilerplate/internal/transport/http/api1/handler"
	webHandler "mobile-backend-boilerplate/internal/transport/http/web/handler"
)

type Dependencies struct {
	DB     *sql.DB
	Logger *slog.Logger
	Config *config.Config

	AuthRepo    repository.AuthRepository
	UserRepo    repository.UserRepository
	AdminRepo   repository.AdminRepository
	RequestRepo repository.RequestRepository
	PostRepo    repository.PostRepository

	TelegramNotifier notifier.Notifier

	MobileAuthService *service.MobileAuthService
	WebAuthService    *service.WebAuthService
	UserService       *service.UserService
	AdminService      *service.AdminService
	RequestService    *service.RequestService
	PostService       *service.PostService

	MobileAuthHandler *apiHandler.MobileAuthHandler
	WebAuthHandler    *webHandler.WebAuthHandler
	RequestHandler    *apiHandler.RequestHandler
	PostHandler       *webHandler.PostHandler
	ContactHandler    *webHandler.ContactHandler
}
