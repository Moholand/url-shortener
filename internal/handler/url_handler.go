package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"url-shortener/internal/service"
	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
)

type AnalyticsResponse struct {
	ShortCode   string                      `json:"short_code"`
	ExpiresAt   *string                     `json:"expires_at,omitempty"`
	TotalClicks int                         `json:"total_clicks"`
	Clicks      []AnalyticsClick            `json:"clicks"`
}

type AnalyticsClick struct {
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	Referer   string `json:"referer"`
	ClickedAt string `json:"clicked_at"`
}

type ShortenRequest struct {
	URL       string `json:"url"`
	ExpiresAt string `json:"expires_at,omitempty"`
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

		var expiresAt *time.Time
		if req.ExpiresAt != "" {
			t, err := time.Parse(time.RFC3339, req.ExpiresAt)
			if err != nil {
				http.Error(w, "invalid expires_at format, use RFC 3339 (e.g. 2026-12-31T23:59:59Z)", http.StatusBadRequest)
				return
			}
			expiresAt = &t
		}

		urlData, err := service.Create(r.Context(), req.URL, expiresAt)
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

		go service.RecordClick(context.Background(), shortCode, r.RemoteAddr, r.UserAgent(), r.Referer())

		http.Redirect(w, r, originalURL, http.StatusFound)
	}
}

func GetAnalytics(service *service.URLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortCode := chi.URLParam(r, "shortCode")

		urlData, err := service.GetURLInfo(r.Context(), shortCode)
		if err != nil || urlData == nil {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		total, clicks, err := service.GetAnalytics(r.Context(), shortCode)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		analyticsClicks := make([]AnalyticsClick, len(clicks))
		for i, c := range clicks {
			analyticsClicks[i] = AnalyticsClick{
				IPAddress: c.IPAddress,
				UserAgent: c.UserAgent,
				Referer:   c.Referer,
				ClickedAt: c.ClickedAt.Format("2006-01-02T15:04:05Z07:00"),
			}
		}

		var expiresAtStr *string
		if urlData.ExpiresAt != nil {
			s := urlData.ExpiresAt.Format("2006-01-02T15:04:05Z07:00")
			expiresAtStr = &s
		}

		res := AnalyticsResponse{
			ShortCode:   shortCode,
			ExpiresAt:   expiresAtStr,
			TotalClicks: total,
			Clicks:      analyticsClicks,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}
