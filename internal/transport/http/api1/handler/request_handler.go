package handler

import (
	"encoding/json"
	"mobile-backend-boilerplate/internal/notifier"
	"mobile-backend-boilerplate/internal/repository"
	"mobile-backend-boilerplate/internal/service"
	"net/http"
)

type RequestHandler struct {
	requestService *service.RequestService
	notifier       notifier.Notifier
}

func NewRequestHandler(requestService *service.RequestService, notifier notifier.Notifier) *RequestHandler {
	return &RequestHandler{
		requestService: requestService,
		notifier:       notifier,
	}
}

func (h *RequestHandler) Submit(w http.ResponseWriter, r *http.Request) {
	var req repository.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.requestService.Create(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.notifier.SendMessage(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(id)
}
