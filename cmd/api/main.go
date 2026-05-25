package main

import (
	"fmt"
	"net/http"

	"url-shortener/internal/handler"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"
	"url-shortener/pkg/db"

	"github.com/go-chi/chi/v5"
)

func main() {

	database := db.Connect()
	defer database.Close()

	rdb := db.NewRedisClient()
	
	urlRepo := repository.NewURLRepository(database)
	urlService := service.NewURLService(urlRepo, rdb)

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Get("/{shortCode}", handler.RedirectURL(urlService))

	r.Post("/shorten", handler.ShortenURL(urlService))

	fmt.Println("Server running on :8080")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}

}
