package infrastructure

import (
	apiHandler "mobile-backend-boilerplate/internal/transport/http/api1/handler"
	webHandler "mobile-backend-boilerplate/internal/transport/http/web/handler"
)

func (d *Dependencies) InitHandlers() {
	d.MobileAuthHandler = apiHandler.NewMobileAuthHandler(d.MobileAuthService, d.UserService)
	d.AdminAuthHandler = apiHandler.NewAdminAuthHandler(d.AdminAuthService, d.AdminService)
	d.RequestHandler = apiHandler.NewRequestHandler(d.RequestService, d.TelegramNotifier)
	d.SchemaEntityHandler = apiHandler.NewSchemaEntityHandler(d.SchemaEntityService)
	d.FilesHandler = apiHandler.NewFilesHandler(d.FileStorage)

	d.PostHandler = webHandler.NewPostHandler(d.PostService)
	d.ContactHandler = webHandler.NewContactHandler(d.RequestService, d.SchemaEntityService, d.TelegramNotifier)
	d.StaticPageHandler = webHandler.NewStaticPageHandler(d.SchemaEntityService)
}
