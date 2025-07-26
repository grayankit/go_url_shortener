package main

import (
	"log"
	"net/http"
	"os"

	"github.com/grayankit/go_url_shortener/internal/db"
	"github.com/grayankit/go_url_shortener/internal/handler"
	"github.com/grayankit/go_url_shortener/internal/logger"
	"github.com/grayankit/go_url_shortener/internal/middleware"
	"github.com/joho/godotenv"
)

func main() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("Error Loading .env file")
	}

	adminUser := os.Getenv("ADMIN_USER")
	adminPass := os.Getenv("ADMIN_PASS")
	log.Println(adminPass, adminUser)
	database := db.New("shortener.db")
	handler.InitHandlers(database)
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/favicon.ico", fs)
	mux.HandleFunc("/", handler.Home)
	mux.HandleFunc("/shorten", handler.Shorten)
	mux.HandleFunc("/u/", handler.Redirect)
	mux.HandleFunc("/dashboard", handler.BasicAuth(adminUser, adminPass, handler.Dashboard))
	mux.HandleFunc("/delete", handler.BasicAuth(adminUser, adminPass, handler.DeleteURL))

	//Logging every request
	logger := logger.NewLogger()
	loggingMiddleware := middleware.LoggingMiddleware(logger)
	wrappedMux := loggingMiddleware(mux)

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", wrappedMux)

	if err != nil {
		log.Fatal(err)
	}
}
