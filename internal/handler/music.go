package handler

import (
	"net/http"
	"strings"

	"gostatus/internal/badge"
	"gostatus/internal/gateway"
)

// single music-service activity
type musicProvider struct {
	// discord activity name to match against
	name string
	// left-hand side text
	label string
	// right-hand side color
	color string
	// shields.io logo slug
	logo string
	// builds the badge text from a matched activity
	// "" when the activity lacks the required fields
	message func(a *gateway.Activity) string
}

// ordered registry of supported music services
var musicProviders = []musicProvider{
	{
		name:  "Spotify",
		label: "listening to",
		color: badge.ColorSpotify,
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
	def := musicProviders[0]

	renderProvider := func(p musicProvider, message string) {
		h.renderBadge(w, r,
			qp(r, "label", p.label),
			message,
			qp(r, "labelColor", badge.ColorLabel),
			qp(r, "color", p.color),
			p.logo,
		)
	}

	pres, ok := h.store.Get(h.id(r, "/badge/music/"))
	if !ok {
		renderProvider(def, fallback)
		return
	}

	for _, provider := range musicProviders {
		a := FindActivity(pres, gateway.ActivityTypeListening, provider.name)
		if a == nil {
			continue
		}
		if msg := provider.message(a); msg != "" {
			renderProvider(provider, msg)
			return
		}
	}

	renderProvider(def, fallback)
}
