package infrastructure

import "mobile-backend-boilerplate/internal/service"

func (d *Dependencies) InitServices() {
	d.MobileAuthService = service.NewMobileAuthService(
		d.Repository.Auth(),
		d.Repository.User(),
		[]byte(d.Config.Authentication.JWT.Secret),
		d.Logger,
	)
	d.AdminAuthService = service.NewAdminAuthService(
		d.Repository.Admin(),
		[]byte(d.Config.Authentication.JWT.Secret),
		d.Logger,
	)
	d.UserService = service.NewUserService(d.Repository.User(), d.Logger)
	d.AdminService = service.NewAdminService(d.Repository.Admin(), d.Logger)
	d.RequestService = service.NewRequestService(d.Repository.Request(), d.Logger)
	d.PostService = service.NewPostService(d.Repository.Post(), d.Logger)
}
