package handler

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/grayankit/go_url_shortener/internal/db"
	"github.com/grayankit/go_url_shortener/internal/shortener"
)

var database *db.DB

func InitHandlers(d *db.DB) {
	database = d
}

func Home(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		println(err.Error())
		http.Error(w, "Unable to load page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
func Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	longURL := r.FormValue("url")
	if longURL == "" {
		http.Error(w, "Missing URL", http.StatusBadRequest)
		return
	}
	existingCode, errDouble := database.GetCodeByURL(longURL)
	if errDouble == nil {
		shortURL := "http://localhost:8080/u/" + existingCode
		tmpl, _ := template.ParseFiles(filepath.Join("templates", "result.html"))
		tmpl.Execute(w, shortURL)
		return
	}
	code := shortener.GenerateCode()
	err := database.Save(code, longURL)
	if err != nil {
		http.Error(w, "Failed to save", http.StatusInternalServerError)
		return
	}

	shortURL := "http://localhost:8080/u/" + code
	tmpl, _ := template.ParseFiles(filepath.Join("templates", "result.html"))
	tmpl.Execute(w, shortURL)
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/u/")
	longURL, err := database.GetLongURL(code)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, longURL, http.StatusFound)
}
