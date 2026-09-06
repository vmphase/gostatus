package handler

import (
	"fmt"
	"net/http"
	"strings"

	"gostatus/internal/badge"
	"gostatus/internal/gateway"
)

// Status renders the user's presence status.
func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	svgHeaders(w)
	status := "offline"
	if p, ok := h.store.Get(h.id(r, "/badge/status/")); ok {
		status = p.Status
	}
	if r.URL.Query().Get("simple") == "true" && (status == "idle" || status == "dnd") {
		status = "online"
	}
	color := badge.StatusColors[status]
	if color == "" {
		color = badge.ColorOffline
	}
	h.renderBadge(w, r, "currently", status, badge.ColorLabel, color, "")
}

// Playing renders the user's current game, excluding code editor activities.
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
		qp(r, "labelColor", badge.ColorLabel), qp(r, "color", badge.ColorDiscord),
		"",
	)
}

// CrunchyRoll renders the user's current Crunchyroll episode.
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
		qp(r, "labelColor", badge.ColorLabel), qp(r, "color", badge.ColorCrunchyroll),
		"crunchyroll",
	)
}
