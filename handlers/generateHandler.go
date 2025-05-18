package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type Coordinate struct {
    Lat float64 `json:"lat"`
    Lng float64 `json:"lng"`
}

type Hole struct {
    TeeBox Coordinate `json:"tee_box"`
    Green  Coordinate `json:"green"`
}

type Course struct {
    CourseName string           `json:"course_name"`
    Holes      map[string]Hole  `json:"holes"`
}

func GenerateHandler(templates *template.Template) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
        if err := r.ParseForm(); err != nil {
            http.Error(w, "Failed to parse form", http.StatusBadRequest)
            return
        }

        holes := make(map[string]Hole)

        for i := 1; i <= 18; i++ {
            prefix := fmt.Sprintf("hole%d_", i)
        
            latTee, err1 := strconv.ParseFloat(r.FormValue(prefix+"lat_tee"), 64)
            lngTee, err2 := strconv.ParseFloat(r.FormValue(prefix+"lng_tee"), 64)
            latGreen, err3 := strconv.ParseFloat(r.FormValue(prefix+"lat_green"), 64)
            lngGreen, err4 := strconv.ParseFloat(r.FormValue(prefix+"lng_green"), 64)
        
            if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
                continue // Skip incomplete or malformed holes
            }
        
            holes[strconv.Itoa(i)] = Hole{
                TeeBox: Coordinate{Lat: latTee, Lng: lngTee},
                Green:  Coordinate{Lat: latGreen, Lng: lngGreen},
            }
        }

        course := Course{
            CourseName: "My Course",
            Holes:      holes,
        }

        w.Header().Set("Content-Disposition", "attachment; filename=course.json")
	    w.Header().Set("Content-Type", "application/json")

        if err := json.NewEncoder(w).Encode(course); err != nil {
            http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
        }
    }       
}
