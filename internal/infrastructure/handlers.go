package infrastructure

import (
	apiHandler "mobile-backend-boilerplate/internal/transport/http/api1/handler"
	webHandler "mobile-backend-boilerplate/internal/transport/http/web/handler"
)

func (d *Dependencies) InitHandlers() {
	d.MobileAuthHandler = apiHandler.NewMobileAuthHandler(d.MobileAuthService, d.UserService)
	d.RequestHandler = apiHandler.NewRequestHandler(d.RequestService, d.TelegramNotifier)

	d.WebAuthHandler = webHandler.NewWebAuthHandler(d.WebAuthService, d.AdminService)
	d.PostHandler = webHandler.NewPostHandler(d.PostService)
	d.ContactHandler = webHandler.NewContactHandler(d.RequestService, d.TelegramNotifier)
}
