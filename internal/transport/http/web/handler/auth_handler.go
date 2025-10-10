package handler

import (
	"encoding/json"
	"errors"
	"mobile-backend-boilerplate/internal/service"
	"mobile-backend-boilerplate/internal/transport/http/web/middleware"
	"mobile-backend-boilerplate/internal/view/layouts"
	pages "mobile-backend-boilerplate/internal/view/pages/public"
	"net/http"
	"time"
)

type WebAuthHandler struct {
	webAuthService *service.WebAuthService
	adminService   *service.AdminService
}

func NewWebAuthHandler(webAuthService *service.WebAuthService, adminService *service.AdminService) *WebAuthHandler {
	return &WebAuthHandler{
		webAuthService: webAuthService,
		adminService:   adminService,
	}
}

func (h *WebAuthHandler) Show(w http.ResponseWriter, r *http.Request) {
	data := layouts.NewBaseLayoutProps(r)
	data.Centered = true
	data.WithNavigation = false
	data.WithTopPadding = false

	cookie, err := r.Cookie("admin_token")
	if err == nil && cookie.Value != "" {
		_, err := h.webAuthService.Me(cookie.Value)
		if err == nil {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}
	}

	formState := pages.NewLoginFormState()
	HandleStaticPage(w, r, pages.LoginPage(data, formState), pages.LoginPageContent(data, formState))
}

func (h *WebAuthHandler) Submit(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Incorrect data", http.StatusBadRequest)
		return
	}

	adminData := pages.LoginFormData{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	data := layouts.NewBaseLayoutProps(r)
	data.Centered = true

	formErrors := h.validateForm(&adminData)
	if len(formErrors) > 0 {
		formState := pages.NewLoginFormStateWithErrors(adminData, formErrors)

		HandleStaticPage(w, r, pages.LoginPage(data, formState), pages.LoginPagePartialForm(formState))
		return
	}

	token, err := h.webAuthService.Login(adminData.Username, adminData.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			formErrors["username"] = "Пользователь не найден"
		case errors.Is(err, service.ErrInvalidPassword):
			formErrors["password"] = "Неверный пароль"
		default:
			formErrors["form"] = "Ошибка авторизации, попробуйте позже"
		}

		formState := pages.NewLoginFormStateWithErrors(adminData, formErrors)
		HandleStaticPage(w, r, pages.LoginPage(data, formState), pages.LoginPagePartialForm(formState))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "admin_token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(30 * 24 * time.Hour),
	})

	w.Header().Set("HX-Redirect", "/dashboard")
	w.WriteHeader(http.StatusOK)
}

func (h *WebAuthHandler) validateForm(data *pages.LoginFormData) map[string]string {
	errors := make(map[string]string)

	if data.Username == "" {
		errors["username"] = "Обязательное поле"
	}

	if data.Password == "" {
		errors["password"] = "Обязательное поле"
	}

	return errors
}

func (h *WebAuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "admin_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, "/auth/login", http.StatusPermanentRedirect)
}

func (h *WebAuthHandler) MeWeb(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetAdminID(r.Context())
	if userID == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.adminService.GetAdminById(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
