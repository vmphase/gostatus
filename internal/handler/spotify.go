package handler

import (
	"fmt"
	"net/http"
	"strings"

	"gostatus/internal/badge"
	"gostatus/internal/gateway"
)

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
