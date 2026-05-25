package handler

import (
	"encoding/json"
	"net/http"
	"url-shortener/internal/service"
	"github.com/go-chi/chi/v5"
	"url-shortener/pkg/code"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortCode string `json:"short_code"`
}

func ShortenURL(service *service.URLService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var req ShortenRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		shortCode := code.Generate(6)

		urlData, err := service.Create(shortCode, req.URL)
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}

		res := ShortenResponse{
			ShortCode: urlData.ShortCode,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}

func RedirectURL(service *service.URLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		shortCode := chi.URLParam(r, "shortCode")
		
		originalURL, err := service.GetOriginalURL(r.Context(), shortCode)
		if err != nil || originalURL == "" {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, originalURL, http.StatusFound)
	}
}
