package handler

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/grayankit/go_url_shortener/internal/shortener"
	"github.com/grayankit/go_url_shortener/internal/storage"
)

var store = storage.NewStore()

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
	code := shortener.GenerateCode()
	store.Save(code, longURL)

	shortURL := "http://localhost:8080/u/" + code
	tmpl, _ := template.ParseFiles(filepath.Join("templates", "result.html"))
	tmpl.Execute(w, shortURL)
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/u/")
	longURL, ok := store.Get(code)
	if !ok {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, longURL, http.StatusFound)
}
