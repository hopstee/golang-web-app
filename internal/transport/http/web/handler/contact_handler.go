package handler

import (
	"mobile-backend-boilerplate/internal/notifier"
	"mobile-backend-boilerplate/internal/repository"
	"mobile-backend-boilerplate/internal/service"
	"mobile-backend-boilerplate/internal/view/layouts"
	"mobile-backend-boilerplate/internal/view/pages"
	"net/http"
	"regexp"
	"strconv"
)

var phoneRegex = regexp.MustCompile(`^(\+7|8)\d{10}$`)

type ContactHandler struct {
	requestService *service.RequestService
	notifier       notifier.Notifier
}

func NewContactHandler(requestService *service.RequestService, notifier notifier.Notifier) *ContactHandler {
	return &ContactHandler{
		requestService: requestService,
		notifier:       notifier,
	}
}

func (h *ContactHandler) Show(w http.ResponseWriter, r *http.Request) {
	data := layouts.NewBaseLayoutProps(r)
	data.Centered = true

	formState := pages.NewFormState()

	HandleStaticPage(w, r, pages.ContactPage(data, formState), pages.ContactPageContent(data, formState))
}

func (h *ContactHandler) Submit(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Incorrect data", http.StatusBadRequest)
		return
	}

	state := repository.Request{
		Name:        r.FormValue("name"),
		Message:     r.FormValue("message"),
		Phone:       r.FormValue("phone"),
		Email:       r.FormValue("email"),
		ContactType: r.FormValue("contact_type"),
	}

	if amountStr := r.FormValue("amount"); amountStr != "" {
		if amount, err := strconv.Atoi(amountStr); err == nil {
			state.Amount = amount
		}
	}

	errors := h.validateForm(&state)
	if len(errors) > 0 {
		data := layouts.NewBaseLayoutProps(r)
		data.Centered = true
		formState := pages.NewFormStateWithErrors(state, errors)

		HandleStaticPage(w, r, pages.ContactPage(data, formState), pages.ContactPagePartialForm(formState))
		return
	}

	_, err = h.requestService.Create(state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.notifier.SendMessage(state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	data := layouts.NewBaseLayoutProps(r)
	data.Centered = true
	formState := pages.NewFormState()
	formState.Success = true

	HandleStaticPage(w, r, pages.ContactPage(data, formState), pages.ContactPagePartialForm(formState))
}

func (h *ContactHandler) validateForm(state *repository.Request) map[string]string {
	errors := make(map[string]string)

	if state.Name == "" {
		errors["name"] = "Обязательное поле"
	} else if len(state.Name) < 2 {
		errors["name"] = "Имя должно быть больше одного символа"
	}

	if state.Message == "" {
		errors["message"] = "Обязательное поле"
	}

	if state.Phone == "" {
		errors["phone"] = "Обязательное поле"
	} else if !phoneRegex.MatchString(state.Phone) {
		errors["phone"] = "Некорректный номер телефона"
	}

	if state.Email == "" {
		errors["email"] = "Обязательное поле"
	}

	if state.Amount <= 0 {
		errors["amount"] = "Обязательное поле"
	}

	if state.ContactType == "" {
		errors["contact_type"] = "Обязательное поле"
	}

	return errors
}
