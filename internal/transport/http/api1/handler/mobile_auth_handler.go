package handler

import (
	"encoding/json"
	"mobile-backend-boilerplate/internal/service"
	"mobile-backend-boilerplate/internal/transport/http/api1/middleware"
	"net/http"
)

type MobileAuthHandler struct {
	authService *service.MobileAuthService
	userService *service.UserService
}

func NewMobileAuthHandler(authService *service.MobileAuthService, userService *service.UserService) *MobileAuthHandler {
	return &MobileAuthHandler{
		authService: authService,
		userService: userService,
	}
}

type userAuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email,omitempty"`
	DeviceID string `json:"device_id"`
}

type userRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
	DeviceID     string `json:"device_id"`
}

func (h *MobileAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req userAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokens, err := h.authService.Login(req.Username, req.Password, req.DeviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(tokens)
}

func (h *MobileAuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req userAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokens, err := h.authService.Register(req.Username, req.Password, req.Email, req.DeviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(tokens)
}

func (h *MobileAuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req userRefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokens, err := h.authService.Refresh(req.RefreshToken, req.DeviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(tokens)
}

func (h *MobileAuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req userRefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.authService.Logout(req.RefreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MobileAuthHandler) MeMobile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetUserById(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
