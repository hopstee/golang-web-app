package handler

import (
	"encoding/json"
	"fmt"
	"mobile-backend-boilerplate/internal/notifier"
	"mobile-backend-boilerplate/internal/repository"
	"mobile-backend-boilerplate/internal/service"
	"mobile-backend-boilerplate/pkg/helper/markdown"
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

	msg := fmt.Sprintf("*Имя:* %s\n", markdown.EscapeMarkdownV2(req.Name))
	msg += fmt.Sprintf("*Сообщение:* %s\n", markdown.EscapeMarkdownV2(req.Message))
	msg += fmt.Sprintf("*Способ связи:* %s\n", markdown.EscapeMarkdownV2(req.ContactType))
	msg += fmt.Sprintf("*Телефон:* %s\n", markdown.EscapeMarkdownV2(req.Phone))
	msg += fmt.Sprintf("*Почта:* %s\n", markdown.EscapeMarkdownV2(req.Email))
	err = h.notifier.SendMessage(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(id)
}
