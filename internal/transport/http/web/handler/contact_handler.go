package handler

import (
	"fmt"
	"log"
	"mobile-backend-boilerplate/internal/notifier"
	"mobile-backend-boilerplate/internal/repository"
	"mobile-backend-boilerplate/internal/service"
	"mobile-backend-boilerplate/internal/view/layouts"
	"mobile-backend-boilerplate/internal/view/pages"
	"mobile-backend-boilerplate/pkg/helper/conversion"
	"mobile-backend-boilerplate/pkg/helper/json"
	"mobile-backend-boilerplate/pkg/helper/markdown"
	"net/http"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
)

type ContactHandler struct {
	requestService        *service.RequestService
	schemaEntitiesService *service.SchemaEntityService
	notifier              notifier.Notifier
}

func NewContactHandler(
	requestService *service.RequestService,
	schemaEntitiesService *service.SchemaEntityService,
	notifier notifier.Notifier,
) *ContactHandler {
	return &ContactHandler{
		requestService:        requestService,
		schemaEntitiesService: schemaEntitiesService,
		notifier:              notifier,
	}
}

func (h *ContactHandler) Submit(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Incorrect data", http.StatusBadRequest)
		return
	}

	values := make(map[string]string)
	for key := range r.Form {
		values[key] = r.FormValue(key)
	}

	pageData, layoutName, err := h.schemaEntitiesService.CollectFullEntityData(r.Context(), "contact")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	layoutData := layouts.NewPublicLayoutProps(r)
	metaData := layoutData.HeadData
	if pageMetaData, ok := pageData[layoutName]; ok {
		if err := json.MapToStruct(pageMetaData, &metaData); err != nil {
			http.Error(w, fmt.Sprintf("failed to decode head data: %v", err), http.StatusInternalServerError)
			return
		}
	}
	layoutData.Centered = true
	layoutData.HeadData = metaData

	var contactContent pages.ContactPagePartialProps
	if err := json.MapToStruct(pageData[repository.MainContentKey], &contactContent); err != nil {
		http.Error(w, fmt.Sprintf("failed to decode page data: %v", err), http.StatusInternalServerError)
		return
	}

	contactPageData := pages.ContactPageProps{
		LayoutContent: layoutData,
		PageContent:   contactContent,
		State:         *pages.NewContactFormState(),
	}

	errors := h.validateDynamicForm(pageData, values)
	if len(errors) > 0 {
		contactPageData.State.Errors = errors
		contactPageData.State.Values = values
		HandleStaticPage(
			w,
			r,
			pages.ContactPage(contactPageData),
			pages.ContactPagePartialForm(contactPageData.PageContent, contactPageData.State),
		)
		return
	}

	msg := markdown.EscapeMarkdownV2("🚨*Новый запрос!*🚨\n")
	for key, value := range values {
		switch {
		case strings.Contains(key, "name"):
			msg += fmt.Sprintf("🙋‍♂️ %s\n", markdown.EscapeMarkdownV2(value))
		case strings.Contains(key, "phone"):
			msg += fmt.Sprintf("📞 %s\n", markdown.EscapeMarkdownV2(value))
		case strings.Contains(key, "email"):
			msg += fmt.Sprintf("📫 %s\n", markdown.EscapeMarkdownV2(value))
		case strings.Contains(key, "message"):
			msg += fmt.Sprintf("✉️ %s\n", markdown.EscapeMarkdownV2(value))
		case strings.Contains(key, "amount"):
			msg += fmt.Sprintf("💰 %s\n", markdown.EscapeMarkdownV2(value))
		default:
			msg += fmt.Sprintf("▶️ %s\n", markdown.EscapeMarkdownV2(value))
		}
	}
	log.Printf("Notification message: %+v", msg)
	err = h.notifier.SendMessage(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	contactPageData.State.Success = true

	HandleStaticPage(
		w,
		r,
		pages.ContactPage(contactPageData),
		pages.ContactPagePartialForm(contactPageData.PageContent, contactPageData.State),
	)
}

func (h *ContactHandler) validateDynamicForm(pageData map[string]interface{}, values map[string]string) map[string]string {
	errorsMap := make(map[string]string)

	content, ok := pageData[repository.MainContentKey].(map[string]interface{})
	if !ok {
		errorsMap["form"] = "Неверная структура данных (content)"
		return errorsMap
	}

	formFields, ok := content[repository.FormFieldsKey].([]interface{})
	if !ok {
		errorsMap["form"] = "Нет полей формы"
		return errorsMap
	}

	for _, f := range formFields {
		field, ok := f.(map[string]interface{})
		if !ok {
			continue
		}

		id := fmt.Sprint(field["id"])
		val := values[id]

		if req := conversion.ParseBool(field["required"]); req && val == "" {
			errorsMap[id] = "Поле обязательно"
			continue
		}

		if minL, ok := conversion.ParseInt(field["min_length"]); ok && minL != 0 {
			if len(val) < minL {
				errorsMap[id] = fmt.Sprintf("Поле должно содержать минимум %d символов", minL)
				continue
			}
		}

		if maxL, ok := conversion.ParseInt(field["max_length"]); ok && maxL != 0 {
			if len(val) > maxL {
				errorsMap[id] = fmt.Sprintf("Поле должно быть максимум %d символов", maxL)
				continue
			}
		}

		if min, ok := conversion.ParseInt(field["min"]); ok && min != 0 {
			if v, err := strconv.Atoi(val); err == nil && v < min {
				errorsMap[id] = fmt.Sprintf("Значение должно быть не меньше %v", min)
				continue
			}
		}

		if max, ok := conversion.ParseInt(field["max"]); ok && max != 0 {
			if v, err := strconv.Atoi(val); err == nil && v > max {
				errorsMap[id] = fmt.Sprintf("Значение должно быть не больше %v", max)
				continue
			}
		}

		if fType := conversion.ParseString(field["type"]); fType == "email" {
			_, err := mail.ParseAddress(val)
			if err != nil {
				errorsMap[id] = "Неверный формат почты"
				continue
			}
		}

		if pattern := conversion.ParseString(field["pattern"]); pattern != "" {
			re, err := regexp.Compile(pattern)
			if err == nil && !re.MatchString(val) {
				errorsMap[id] = "Неверный формат"
				continue
			}
		}

		delete(errorsMap, id)
	}

	return errorsMap
}
