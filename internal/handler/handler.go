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
	mux.HandleFunc("/badge/spotify/", h.Spotify)
	mux.HandleFunc("/badge/crunchyroll/", h.CrunchyRoll)
	mux.HandleFunc("/badge/playing/", h.Playing)
	mux.HandleFunc("/badge/vscode/", h.VSCode)
	mux.HandleFunc("/badge/zed/", h.Zed)
	mux.HandleFunc("/presence/", h.Presence)
}

func (h *Handler) id(r *http.Request, prefix string) string {
	return strings.TrimPrefix(r.URL.Path, prefix)
}

func (h *Handler) renderBadge(w http.ResponseWriter, r *http.Request, label, message, labelColor, color, logo string) {
	var svg string
	if logo != "" && qp(r, "hideLogo", "false") != "true" {
		svg = badge.MakeWithLogo(label, message, labelColor, color, logo)
	} else {
		svg = badge.Make(label, message, labelColor, color)
	}
	fmt.Fprint(w, svg)
}

func svgHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	w.Header().Set("Cache-Control", "max-age=0, no-cache, no-store, must-revalidate")
}

func qp(r *http.Request, key, fallback string) string {
	if v := r.URL.Query().Get(key); v != "" {
		return v
	}
	return fallback
}
