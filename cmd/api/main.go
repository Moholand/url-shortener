package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"url-shortener/internal/handler"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"
	"url-shortener/pkg/db"

	"github.com/go-chi/chi/v5"
)

func main() {

	database, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	rdb := db.NewRedisClient()
	
	urlRepo := repository.NewURLRepository(database)
	clickRepo := repository.NewClickRepository(database)
	urlService := service.NewURLService(urlRepo, clickRepo, rdb)

	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("failed to parse templates:", err)
	}

	panelHandler := handler.NewPanelHandler(urlService, tmpl)

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Route("/panel", func(r chi.Router) {
		r.Get("/", panelHandler.PanelList)
		r.Get("/create", panelHandler.PanelCreate)
		r.Post("/create", panelHandler.PanelStoreShortCode)
		r.Get("/{shortCode}", panelHandler.PanelDetail)
		r.Post("/{shortCode}/delete", panelHandler.PanelDelete)
	})

	r.Post("/shorten", handler.ShortenURL(urlService))

	r.Get("/analytics/{shortCode}", handler.GetAnalytics(urlService))

	r.Delete("/{shortCode}", handler.DeleteURL(urlService))

	r.Get("/{shortCode}", handler.RedirectURL(urlService))

	fmt.Println("Server running on :8080")

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}

}
