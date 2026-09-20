package handler

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"

	"url-shortener/internal/service"

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
