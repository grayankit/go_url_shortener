package handler

import (
	"encoding/base64"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {
	records, err := database.GetAll()
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	tmpl, _ := template.ParseFiles(filepath.Join("templates", "dashboard.html"))
	tmpl.Execute(w, records)
}
func DeleteURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := r.FormValue("code")
	if code == "" {
		http.Error(w, "Missing code", http.StatusBadRequest)
		return
	}

	if err := database.Delete(code); err != nil {
		http.Error(w, "Failed to delete", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
func BasicAuth(username, password string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const basicPrefix = "Basic "

		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, basicPrefix) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		payload, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(auth, basicPrefix))
		if err != nil {
			http.Error(w, "Invalid auth", http.StatusUnauthorized)
			return
		}

		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 || pair[0] != username || pair[1] != password {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		handler(w, r)
	}
}
