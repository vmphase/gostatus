package handler

import (
	"net/http"

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
