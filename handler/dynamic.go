package handler

import (
	"html/template"
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {

	lessen := []les{
		les{Id: 0, Name: "test les 1"},
		les{Id: 1, Name: "test les 2"},
		les{Id: 2, Name: "test les 3"},
	}

	renderTemplate(w, "home", lessen)
}

func Lesson(w http.ResponseWriter, r *http.Request) {
}

// TODO: optimize this by loading templates once on startup instead of on every request.
func renderTemplate(w http.ResponseWriter, name string, data any) {
	templateBaseDir := "./www/templates"

	tmpl, err := template.ParseGlob(templateBaseDir + "/*.html")
	if err != nil {
		log.Print("failed to load template", name, ":", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl, err = tmpl.ParseFiles(templateBaseDir + "/pages/" + name + ".html")

	if err != nil {
		log.Print("failed to load template ", name, " :", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Print("failed to execute template ", name, " :", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
