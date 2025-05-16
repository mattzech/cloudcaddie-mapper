package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var templates *template.Template

func seq(end int) []int {
	s := make([]int, end)
	for i := range s {
		s[i] = i+1
	}
	return s
}

func main() {
	funcs := template.FuncMap{"seq": seq}
	var err error
	templates, err = template.New("home").Funcs(funcs).ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/load-form", loadFormHandler)
	http.HandleFunc("/generate", generateHandler)

	fmt.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Rendering index.html")
	err := templates.ExecuteTemplate(w, "home", map[string][]int{"HoleCount": seq(18)})
	if err != nil {
		log.Printf("Template render error: %v\n", err)
		http.Error(w, "Template rendering failed", http.StatusInternalServerError)
	}
}

func loadFormHandler(w http.ResponseWriter, r *http.Request) {
	holesStr := r.URL.Query().Get("holes")
	holes, err := strconv.Atoi(holesStr)
	fmt.Println("holes " + holesStr)
	if err != nil || (holes != 9 && holes != 18) {
		http.Error(w, "Invalid hole count", http.StatusBadRequest)
		return
	}

	err = templates.ExecuteTemplate(w, "form", map[string][]int{"HoleCount": seq(holes)})
	if err != nil {
		http.Error(w, "Form rendering failed", http.StatusInternalServerError)
	}
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	type Hole struct {
		Hole int     `json:"hole"`
		Lat  float64 `json:"lat"`
		Lon  float64 `json:"lon"`
	}

	var holes []Hole
	for i := 1; ; i++ {
		latStr := r.FormValue(fmt.Sprintf("lat%d", i))
		lonStr := r.FormValue(fmt.Sprintf("lon%d", i))
		if latStr == "" || lonStr == "" {
			break
		}

		lat, err1 := strconv.ParseFloat(latStr, 64)
		lon, err2 := strconv.ParseFloat(lonStr, 64)
		if err1 != nil || err2 != nil {
			http.Error(w, fmt.Sprintf("Invalid coordinates for hole %d", i), http.StatusBadRequest)
			return
		}

		holes = append(holes, Hole{Hole: i, Lat: lat, Lon: lon})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(holes)
}
