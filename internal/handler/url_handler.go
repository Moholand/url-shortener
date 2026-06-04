package handler

import (
	"encoding/json"
	"net/http"
	"url-shortener/internal/service"
	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortCode string `json:"short_code"`
	ShortURL  string `json:"short_url"`
}

func ShortenURL(service *service.URLService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var req ShortenRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		if !govalidator.IsURL(req.URL) {
            http.Error(w, "Invalid URL format", http.StatusBadRequest)
            return
        }

		urlData, err := service.Create(r.Context(), req.URL)
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}

		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}

		shortURL := scheme + "://" + r.Host + "/" + urlData.ShortCode

		res := ShortenResponse{
			ShortCode: urlData.ShortCode,
			ShortURL:  shortURL,
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
