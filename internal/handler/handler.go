package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gostatus/internal/badge"
	"gostatus/internal/gateway"
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
	mux.HandleFunc("/badge/playing/", h.Playing)
	mux.HandleFunc("/badge/vscode/", h.VSCode)
	mux.HandleFunc("/presence/", h.Presence)
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/badge/status/")
	svgHeaders(w)

	status := "offline"
	if p, ok := h.store.Get(id); ok {
		status = p.Status
	}
	if r.URL.Query().Get("simple") == "true" && (status == "idle" || status == "dnd") {
		status = "online"
	}

	fmt.Fprint(w, badge.Make(
		qp(r, "label", "currently"),
		status,
		qp(r, "labelColor", "#555"),
		qp(r, "color", badge.StatusColors[status]),
	))
}

func (h *Handler) Spotify(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/badge/spotify/")
	svgHeaders(w)

	message := qp(r, "fallback", "nothing")
	if p, ok := h.store.Get(id); ok {
		for _, a := range p.Activities {
			if a.Type == gateway.ActivityTypeListening && a.Name == "Spotify" && a.Details != "" && a.State != "" {
				artists := strings.Join(strings.Split(a.State, "; "), ", ")
				message = a.Details + " by " + artists
				break
			}
		}
	}

	fmt.Fprint(w, badge.Make(
		qp(r, "label", "listening to"),
		message,
		qp(r, "labelColor", "#555"),
		qp(r, "color", "#1db954"),
	))
}

func (h *Handler) Playing(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/badge/playing/")
	svgHeaders(w)

	var games []string
	if p, ok := h.store.Get(id); ok {
		for _, a := range p.Activities {
			if a.Type == gateway.ActivityTypePlaying &&
				a.Name != "Visual Studio Code" &&
				a.Name != "IntelliJ IDEA Ultimate" {
				games = append(games, a.Name)
			}
		}
	}

	message := qp(r, "fallback", "nothing")
	if len(games) > 0 {
		message = strings.Join(games, ", ")
	}

	fmt.Fprint(w, badge.Make(
		qp(r, "label", "playing"),
		message,
		qp(r, "labelColor", "#555"),
		qp(r, "color", "#5865f2"),
	))
}

func (h *Handler) Presence(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/presence/")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if p, ok := h.store.Get(id); ok {
		json.NewEncoder(w).Encode(p)
	} else {
		json.NewEncoder(w).Encode(map[string]any{
			"status":        "offline",
			"client_status": map[string]any{},
			"activities":    []any{},
		})
	}
}


func (h *Handler) VSCode(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/badge/vscode/")
	svgHeaders(w)

	message := qp(r, "fallback", "nothing rn")
	if p, ok := h.store.Get(id); ok {
		for _, a := range p.Activities {
			if a.Name == "Visual Studio Code" && a.Details != "" && a.State != "" {
				file := strings.TrimPrefix(a.Details, "Editing ")
				workspace := a.State
				workspace = strings.ReplaceAll(workspace, "Workspace: ", "")
				workspace = strings.ReplaceAll(workspace, " (Workspace)", "")
				if strings.HasPrefix(workspace, "Glitch:") {
					workspace = strings.Replace(workspace, "Glitch:", "🎏", 1)
				}
				message = file + " in " + workspace
				break
			}
		}
	}

	svg := badge.Make(
		qp(r, "label", "vscode"),
		message,
		qp(r, "labelColor", "#555"),
		qp(r, "color", qp(r, "color", "#23a7f2")),
	)

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
