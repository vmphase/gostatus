package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gostatus/internal/badge"
	"gostatus/internal/gateway"
)

// Full cached Discord presence payload for a user.
// Fallbacks to an offline/empty presence response if not found.
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

// The game user is currently playing, excluding editor activities
func (h *Handler) Playing(w http.ResponseWriter, r *http.Request) {
	svgHeaders(w)

	message := qp(r, "fallback", "nothing")

	if p, ok := h.store.Get(h.id(r, "/badge/playing/")); ok {
		// exclude editors from /badge/code/
		editorNames := make([]string, len(editors))
		for i, ed := range editors {
			editorNames[i] = ed.name
		}

		acts := FindAllActivities(p, gateway.ActivityTypePlaying, editorNames...)
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
