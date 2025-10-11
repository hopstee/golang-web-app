package handler

import (
	"encoding/json"
	"mobile-backend-boilerplate/internal/service"
	"mobile-backend-boilerplate/internal/transport/http/api1/middleware"
	"net/http"
	"time"
)

type AdminAuthHandler struct {
	webAuthService *service.AdminAuthService
	adminService   *service.AdminService
}

func NewAdminAuthHandler(webAuthService *service.AdminAuthService, adminService *service.AdminService) *AdminAuthHandler {
	return &AdminAuthHandler{
		webAuthService: webAuthService,
		adminService:   adminService,
	}
}

type adminAuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email,omitempty"`
	DeviceID string `json:"device_id"`
}

func (h *AdminAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req adminAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.webAuthService.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
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

	w.WriteHeader(http.StatusOK)
}

func (h *AdminAuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "admin_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	w.WriteHeader(http.StatusOK)
}

func (h *AdminAuthHandler) MeWeb(w http.ResponseWriter, r *http.Request) {
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
