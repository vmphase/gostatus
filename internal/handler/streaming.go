package handler

import (
	"net/http"
	"strings"

	"gostatus/internal/badge"
	"gostatus/internal/store"
)

// Streaming renders the user's current STREAMING activity (activity type 1),
// per Discord docs the only fields guaranteed are name and url.
func (h *Handler) Streaming(w http.ResponseWriter, r *http.Request) {
	svgHeaders(w)

	message := qp(r, "fallback", "nothing")

	if p, ok := h.store.Get(h.id(r, "/badge/streaming/")); ok {
		if a := FindActivity(p, store.ActivityTypeStreaming, ""); a != nil {
			title := strings.TrimSpace(a.Details)
			if title == "" {
				title = strings.TrimSpace(a.State)
			}
			if title != "" {
				message = title + " on " + a.Name
			} else if a.Name != "" {
				message = "on " + a.Name
			}
		}
	}

	h.renderBadge(w, r,
		qp(r, "label", "streaming"), message,
		qp(r, "labelColor", badge.ColorLabel), qp(r, "color", badge.ColorDiscord),
		"",
	)
}
