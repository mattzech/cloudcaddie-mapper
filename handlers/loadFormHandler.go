package handlers

import (
	"html/template"
	"net/http"
	"strconv"
)

func LoadFormHandler(templates *template.Template) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		holesStr := r.URL.Query().Get("holes")
		holes, err := strconv.Atoi(holesStr)
		if err != nil || (holes != 9 && holes != 18) {
			http.Error(w, "Invalid hole count", http.StatusBadRequest)
			return
		}

		err = templates.ExecuteTemplate(w, "form", map[string][]int{"HoleCount": HoleCount(holes)})
		if err != nil {
			http.Error(w, "Form rendering failed", http.StatusInternalServerError)
		}
	}
}