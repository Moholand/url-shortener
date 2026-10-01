package handler

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
	"time"

	"url-shortener/internal/service"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
)

type PanelHandler struct {
	Service  *service.URLService
	Template *template.Template
}

func NewPanelHandler(svc *service.URLService, tmpl *template.Template) *PanelHandler {
	return &PanelHandler{Service: svc, Template: tmpl}
}

func (h *PanelHandler) PanelList(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	search := r.URL.Query().Get("q")
	perPage := 20
	offset := (page - 1) * perPage

	items, total, err := h.Service.ListURLs(r.Context(), perPage, offset, search)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}
	nextPage := page + 1
	if nextPage > totalPages {
		nextPage = totalPages
	}

	var pageNumbers []int
	for i := 1; i <= totalPages; i++ {
		pageNumbers = append(pageNumbers, i)
	}

	data := struct {
		Items       interface{}
		Page        int
		TotalPages  int
		Total       int
		Search      string
		PrevPage    int
		NextPage    int
		PageNumbers []int
	}{
		Items:       items,
		Page:        page,
		TotalPages:  totalPages,
		Total:       total,
		Search:      search,
		PrevPage:    prevPage,
		NextPage:    nextPage,
		PageNumbers: pageNumbers,
	}

	h.Template.ExecuteTemplate(w, "panel.html", data)
}

func (h *PanelHandler) PanelDetail(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "shortCode")

	url, clicks, totalClicks, err := h.Service.GetURLDetail(r.Context(), shortCode)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if url == nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	shortURL := fmt.Sprintf("%s://%s/%s", scheme, r.Host, url.ShortCode)

	data := struct {
		URL         interface{}
		Clicks      interface{}
		TotalClicks int
		ShortURL    string
	}{
		URL:         url,
		Clicks:      clicks,
		TotalClicks: totalClicks,
		ShortURL:    shortURL,
	}

	h.Template.ExecuteTemplate(w, "detail.html", data)
}

func (h *PanelHandler) PanelDelete(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "shortCode")

	err := h.Service.Delete(r.Context(), shortCode)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/panel", http.StatusFound)
}

func (h *PanelHandler) PanelCreate(w http.ResponseWriter, r *http.Request) {
	h.renderCreateForm(w, "", "", "")
}

func (h *PanelHandler) PanelStoreShortCode(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	expiresAtRaw := r.FormValue("expires_at")

	if !govalidator.IsURL(url) {
		h.renderCreateForm(w, url, expiresAtRaw, "Invalid URL format")
		return
	}

	var expiresAt *time.Time
	if expiresAtRaw != "" {
		parsed, err := parseExpiresAt(expiresAtRaw)
		if err != nil {
			h.renderCreateForm(w, url, expiresAtRaw, "Invalid expires_at format")
			return
		}
		expiresAt = &parsed
	}

	urlData, err := h.Service.Create(r.Context(), url, expiresAt)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/panel/"+urlData.ShortCode, http.StatusFound)
}

func (h *PanelHandler) renderCreateForm(w http.ResponseWriter, url, expiresAt, errMsg string) {
	data := struct {
		URL       string
		ExpiresAt string
		Error     string
	}{
		URL:       url,
		ExpiresAt: expiresAt,
		Error:     errMsg,
	}

	h.Template.ExecuteTemplate(w, "create.html", data)
}

func parseExpiresAt(v string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, v); err == nil {
		return t, nil
	}
	return time.Parse("2006-01-02T15:04", v)
}
