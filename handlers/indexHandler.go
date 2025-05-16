package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func IndexHandler(templates *template.Template) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		log.Println("Rendering index.html")
		err := templates.ExecuteTemplate(w, "home", map[string][]int{"HoleCount": HoleCount(18)})
		if err != nil {
			log.Printf("Template render error: %v\n", err)
			http.Error(w, "Template rendering failed", http.StatusInternalServerError)
		}
	}
}

func HoleCount(end int) []int {
	s := make([]int, end)
	for i := range s {
		s[i] = i+1
	}
	return s
}
