package handler

import (
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/grayankit/go_url_shortener/internal/db"
	"github.com/grayankit/go_url_shortener/internal/logger"
	"github.com/grayankit/go_url_shortener/internal/shortener"
)

var (
	database *db.DB
	log      *slog.Logger
)

func InitHandlers(d *db.DB) {
	database = d
	log = logger.NewLogger()
}

func Home(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Error("failed to parse template", "err", err, "path", tmplPath)
		http.Error(w, "Unable to load page", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Error("failed to execute template", "err", err, "path", tmplPath)
		http.Error(w, "Unable to load page", http.StatusInternalServerError)
	}
}
func Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Warn("method not allowed", "method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	longURL := r.FormValue("url")
	if longURL == "" {
		log.Warn("missing url in request")
		http.Error(w, "Missing URL", http.StatusBadRequest)
		return
	}
	existingCode, errDouble := database.GetCodeByURL(longURL)
	if errDouble == nil {
		shortURL := "http://localhost:8080/u/" + existingCode
		tmpl, err := template.ParseFiles(filepath.Join("templates", "result.html"))
		if err != nil {
			log.Error("failed to parse template", "err", err)
			http.Error(w, "Failed to generate short URL", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, shortURL)
		if err != nil {
			log.Error("failed to execute template", "err", err)
			http.Error(w, "Failed to generate short URL", http.StatusInternalServerError)
		}
		return
	}

	log.Info("could not find existing code for url, will generate a new one", "err", errDouble, "url", longURL)

	code := shortener.GenerateCode()
	err := database.Save(code, longURL)
	if err != nil {
		log.Error("failed to save url", "err", err, "code", code, "url", longURL)
		http.Error(w, "Failed to save", http.StatusInternalServerError)
		return
	}

	shortURL := "http://localhost:8080/u/" + code
	tmpl, err := template.ParseFiles(filepath.Join("templates", "result.html"))
	if err != nil {
		log.Error("failed to parse template", "err", err)
		http.Error(w, "Failed to generate short URL", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, shortURL)
	if err != nil {
		log.Error("failed to execute template", "err", err)
		http.Error(w, "Failed to generate short URL", http.StatusInternalServerError)
	}
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/u/")
	longURL, err := database.GetLongURL(code)
	if err != nil {
		log.Error("failed to get long url", "err", err, "code", code)
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, longURL, http.StatusFound)
}