package handler

import (
	"html/template"
	"net/http"
	"path/filepath"
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
