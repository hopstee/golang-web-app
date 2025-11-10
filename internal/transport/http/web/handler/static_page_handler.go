package handler

import (
	"fmt"
	"log"
	"mobile-backend-boilerplate/internal/repository"
	"mobile-backend-boilerplate/internal/service"
	"mobile-backend-boilerplate/internal/view/layouts"
	"mobile-backend-boilerplate/internal/view/pages"
	"mobile-backend-boilerplate/pkg/helper/json"
	"net/http"
	"strings"
)

type StaticPageHandler struct {
	schemaEntitiesService *service.SchemaEntityService
}

func NewStaticPageHandler(schemaEntitiesService *service.SchemaEntityService) *StaticPageHandler {
	return &StaticPageHandler{
		schemaEntitiesService: schemaEntitiesService,
	}
}

func (h *StaticPageHandler) RenderStaticPage(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	path = strings.TrimPrefix(path, "/")
	path = strings.ReplaceAll(path, "/", ":")

	if path == "" {
		path = "index"
	}

	pageData, layoutName, err := h.schemaEntitiesService.CollectFullEntityData(r.Context(), path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("PAGE DATA: %+v", pageData)

	layoutData := layouts.NewPublicLayoutProps(r)
	metaData := layoutData.HeadData
	if pageMetaData, ok := pageData[layoutName]; ok {
		if err := json.MapToStruct(pageMetaData, &metaData); err != nil {
			http.Error(w, fmt.Sprintf("failed to decode head data: %v", err), http.StatusInternalServerError)
			return
		}
	}

	layoutData.HeadData = metaData

	switch path {
	case "index":
		var indexContent pages.IndexPagePartialProps
		if err := json.MapToStruct(pageData[repository.MainContentKey], &indexContent); err != nil {
			http.Error(w, fmt.Sprintf("failed to decode page data: %v", err), http.StatusInternalServerError)
			return
		}

		layoutData.Centered = true

		indexPageData := pages.IndexPageProps{
			LayoutContent: layoutData,
			PageContent:   indexContent,
		}

		HandleStaticPage(w, r, pages.IndexPage(indexPageData), pages.IndexPageContent(indexPageData))
	case "about":
		var aboutContent pages.AboutPagePartialProps
		if err := json.MapToStruct(pageData[repository.MainContentKey], &aboutContent); err != nil {
			http.Error(w, fmt.Sprintf("failed to decode page data: %v", err), http.StatusInternalServerError)
			return
		}

		aboutPageData := pages.AboutPageProps{
			LayoutContent: layoutData,
			PageContent:   aboutContent,
		}

		HandleStaticPage(w, r, pages.AboutPage(aboutPageData), pages.AboutPageContent(aboutPageData))
	case "projects":
		var projectsContent pages.ProjectPagePartialProps
		if err := json.MapToStruct(pageData[repository.MainContentKey], &projectsContent); err != nil {
			http.Error(w, fmt.Sprintf("failed to decode page data: %v", err), http.StatusInternalServerError)
			return
		}

		projectPageData := pages.ProjectPageProps{
			LayoutContent: layoutData,
			PageContent:   projectsContent,
		}

		HandleStaticPage(w, r, pages.ProjectsPage(projectPageData), pages.ProjectsPageContent(projectPageData))
	case "contact":
		var contactContent pages.ContactPagePartialProps
		if err := json.MapToStruct(pageData[repository.MainContentKey], &contactContent); err != nil {
			http.Error(w, fmt.Sprintf("failed to decode page data: %v", err), http.StatusInternalServerError)
			return
		}

		layoutData.Centered = true

		contactPageData := pages.ContactPageProps{
			LayoutContent: layoutData,
			PageContent:   contactContent,
			State:         *pages.NewContactFormState(),
		}

		HandleStaticPage(w, r, pages.ContactPage(contactPageData), pages.ContactPageContent(contactPageData))
	}
}
