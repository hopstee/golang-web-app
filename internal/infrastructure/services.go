package infrastructure

import "mobile-backend-boilerplate/internal/service"

func (d *Dependencies) InitServices() {
	d.MobileAuthService = service.NewMobileAuthService(
		d.AuthRepo,
		d.UserRepo,
		[]byte(d.Config.Authentication.JWT.Secret),
		d.Logger,
	)
	d.AdminAuthService = service.NewAdminAuthService(
		d.AdminRepo,
		[]byte(d.Config.Authentication.JWT.Secret),
		d.Logger,
	)
	d.UserService = service.NewUserService(d.UserRepo, d.Logger)
	d.AdminService = service.NewAdminService(d.AdminRepo, d.Logger)
	d.RequestService = service.NewRequestService(d.RequestRepo, d.Logger)
	d.PostService = service.NewPostService(d.PostRepo, d.Logger)
}
