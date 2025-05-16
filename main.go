package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"github.com/mattzech/cloudcaddie-mapper/handlers"
)

func main() {
	funcs := template.FuncMap{"seq": handlers.HoleCount}
	var err error
	templates, err := template.New("home").Funcs(funcs).ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}

	http.HandleFunc("/", handlers.IndexHandler(templates))
	http.HandleFunc("/load-form", handlers.LoadFormHandler(templates))
	http.HandleFunc("/generate", handlers.GenerateHandler(templates))

	fmt.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
