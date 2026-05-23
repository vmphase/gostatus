package handler

import (
	"net/http"
	"strings"

	"gostatus/internal/gateway"
)

// single music-service activity.
type musicProvider struct {
	// discord activity name to match against
	name string
	// right-hand side color.
	color string
	// shields.io logo slug
	logo string
	// builds the badge text from a matched activity
	// "" when the activity lacks the required fields
	message func(a *gateway.Activity) string
}

// ordered registry of supported music services.
var musicProviders = []musicProvider{
	{
		name:  "Spotify",
		color: "#1db954",
		logo:  "spotify",
		message: func(a *gateway.Activity) string {
			if a.Details == "" {
				return ""
			}
			return a.Details + " by " + strings.ReplaceAll(a.State, "; ", ", ")
		},
	},
}

func (h *Handler) Music(w http.ResponseWriter, r *http.Request) {
	svgHeaders(w)

	fallback := qp(r, "fallback", "nothing")

	p, ok := h.store.Get(h.id(r, "/badge/music/"))
	if !ok {
		h.renderBadge(w, r, "listening to", fallback, "#555", "#1db954", "")
		return
	}

	for _, provider := range musicProviders {
		a := FindActivity(p, gateway.ActivityTypeListening, provider.name)
		if a == nil {
			continue
		}
		msg := provider.message(a)
		if msg == "" {
			continue
		}

		label := qp(r, "label", "listening to")
		color := qp(r, "color", provider.color)
		labelColor := qp(r, "labelColor", "#555")
		h.renderBadge(w, r, label, msg, labelColor, color, provider.logo)
		return
	}

	h.renderBadge(w, r,
		qp(r, "label", "listening to"), fallback,
		qp(r, "labelColor", "#555"), qp(r, "color", "#1db954"),
		"",
	)
}
