package handler

import (
	"fmt"
	"net/http"
	"strings"

	"gostatus/internal/badge"
	"gostatus/internal/store"
)

type Handler struct {
	store *store.Store
}

func New(s *store.Store) *Handler {
	return &Handler{store: s}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/badge/status/", h.Status)
	mux.HandleFunc("/badge/playing/", h.Playing)
	mux.HandleFunc("/badge/crunchyroll/", h.CrunchyRoll)
	mux.HandleFunc("/badge/code/", h.Code)
	mux.HandleFunc("/badge/music/", h.Music)
	mux.HandleFunc("/presence/", h.Presence)
}

func (h *Handler) id(r *http.Request, prefix string) string {
	return strings.TrimPrefix(r.URL.Path, prefix)
}

// builds a badge.Options from the request, merging caller supplied
// defaults with any query-param overrides the user may have provided
func (h *Handler) badgeOpts(r *http.Request, label, message, labelColor, color, logo string) badge.Options {
	hideLogo := r.URL.Query().Get("hideLogo") == "true"
	return badge.Options{
		Label:      qp(r, "label", label),
		Message:    message,
		LabelColor: qp(r, "labelColor", labelColor),
		Color:      qp(r, "color", color),
		Logo:       logoOrEmpty(logo, hideLogo),
		Style:      badge.ParseStyle(r.URL.Query().Get("style")),
	}
}

func (h *Handler) renderBadge(w http.ResponseWriter, r *http.Request, label, message, labelColor, color, logo string) {
	_, _ = fmt.Fprint(w, badge.Render(h.badgeOpts(r, label, message, labelColor, color, logo)))
}

func svgHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	w.Header().Set("Cache-Control", "max-age=0, no-cache, no-store, must-revalidate")
}

// the query param value for key, falling back to fallback
func qp(r *http.Request, key, fallback string) string {
	if v := r.URL.Query().Get(key); v != "" {
		return v
	}
	return fallback
}

func logoOrEmpty(logo string, hide bool) string {
	if hide {
		return ""
	}
	return logo
}
