package handler

import (
	"fmt"
	"net/http"
	"strings"

	"gostatus/internal/badge"
	"gostatus/internal/gateway"
)

// Current Discord presence status (online, idle, dnd, offline).
// Accepts ?simple=true to collapse idle and dnd into "online".
func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	svgHeaders(w)
	status := "offline"
	if p, ok := h.store.Get(h.id(r, "/badge/status/")); ok {
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

// Spotify track and artist the user is currently listening to
func (h *Handler) Spotify(w http.ResponseWriter, r *http.Request) {
	svgHeaders(w)
	message := qp(r, "fallback", "nothing")
	if p, ok := h.store.Get(h.id(r, "/badge/spotify/")); ok {
		if a := FindActivity(p, gateway.ActivityTypeListening, "Spotify"); a != nil {
			message = a.Details + " by " + strings.ReplaceAll(a.State, "; ", ", ")
		}
	}
	h.renderBadge(w, r,
		qp(r, "label", "listening to"), message,
		qp(r, "labelColor", "#555"), qp(r, "color", "#1db954"),
		"spotify",
	)
}

// CrunchyRoll episode and series the user is currently watching
func (h *Handler) CrunchyRoll(w http.ResponseWriter, r *http.Request) {
	svgHeaders(w)
	message := qp(r, "fallback", "nothing")
	if p, ok := h.store.Get(h.id(r, "/badge/crunchyroll/")); ok {
		if a := FindActivity(p, gateway.ActivityTypeWatching, "Crunchyroll"); a != nil && a.Details != "" {
			message = fmt.Sprintf("%s – %s", a.Details, a.State)
		}
	}
	h.renderBadge(w, r,
		qp(r, "label", "watching"), message,
		qp(r, "labelColor", "#555"), qp(r, "color", "#5865f2"),
		"crunchyroll",
	)
}

// List of games the user is currently playing, excluding editor activities
func (h *Handler) Playing(w http.ResponseWriter, r *http.Request) {
	svgHeaders(w)
	message := qp(r, "fallback", "nothing")
	if p, ok := h.store.Get(h.id(r, "/badge/playing/")); ok {
		acts := FindAllActivities(p, gateway.ActivityTypePlaying, "Visual Studio Code", "Zed")
		if len(acts) > 0 {
			names := make([]string, len(acts))
			for i, a := range acts {
				names[i] = a.Name
			}
			message = strings.Join(names, ", ")
		}
	}
	h.renderBadge(w, r,
		qp(r, "label", "playing"), message,
		qp(r, "labelColor", "#555"), qp(r, "color", "#5865f2"),
		"",
	)
}

// VSCode file and workspace the user is currently editing
func (h *Handler) VSCode(w http.ResponseWriter, r *http.Request) {
	svgHeaders(w)
	message := qp(r, "fallback", "nothing")
	if p, ok := h.store.Get(h.id(r, "/badge/vscode/")); ok {
		if a := FindActivity(p, gateway.ActivityTypePlaying, "Visual Studio Code"); a != nil && a.Details != "" {
			file := strings.TrimPrefix(a.Details, "Editing ")
			ws := strings.ReplaceAll(strings.ReplaceAll(a.State, "Workspace: ", ""), " (Workspace)", "")
			message = file + " in " + ws
		}
	}
	h.renderBadge(w, r,
		qp(r, "label", "vscode"), message,
		qp(r, "labelColor", "#555"), qp(r, "color", "#23a7f2"),
		"vscode",
	)
}

// Zed serves a badge showing the file and workspace the user is currently working on in Zed
func (h *Handler) Zed(w http.ResponseWriter, r *http.Request) {
	svgHeaders(w)
	message := qp(r, "fallback", "nothing")
	if p, ok := h.store.Get(h.id(r, "/badge/zed/")); ok {
		if a := FindActivity(p, gateway.ActivityTypePlaying, "Zed"); a != nil && a.Details != "" {
			file := strings.TrimPrefix(a.State, "Working on ")
			workspace := strings.TrimPrefix(a.Details, "In ")
			message = file + " in " + workspace
		}
	}
	h.renderBadge(w, r,
		qp(r, "label", "zed"), message,
		qp(r, "labelColor", "#1e1e2e"), qp(r, "color", "#7c6df2"),
		"zed",
	)
}
