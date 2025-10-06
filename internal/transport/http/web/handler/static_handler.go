package handler

import (
	"net/http"

	templ "github.com/a-h/templ"
)

func HandleStaticPage(w http.ResponseWriter, r *http.Request, component templ.Component, content templ.Component) {
	if r.Header.Get("HX-Request") != "" && content != nil {
		err := content.Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if component != nil {
		err := component.Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
