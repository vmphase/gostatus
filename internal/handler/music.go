package handler

import (
	"net/http"
	"strings"

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
	def := musicProviders[0]

	p, ok := h.store.Get(h.id(r, "/badge/music/"))
	if !ok {
		h.renderBadge(
			w, r,
			qp(r, "label", def.label),
			fallback,
			qp(r, "labelColor", "#555"),
			qp(r, "color", def.color),
			def.logo,
		)
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
		h.renderBadge(
			w, r,
			qp(r, "label", provider.label),
			msg,
			qp(r, "labelColor", "#555"),
			qp(r, "color", provider.color),
			provider.logo,
		)
		return
	}

	h.renderBadge(
		w, r,
		qp(r, "label", def.label),
		fallback,
		qp(r, "labelColor", "#555"),
		qp(r, "color", def.color),
		def.logo,
	)
}
