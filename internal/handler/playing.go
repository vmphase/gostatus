package handler

import (
	"fmt"
	"net/http"
	"strings"

	"gostatus/internal/badge"
	"gostatus/internal/gateway"
)

func (h *Handler) Playing(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/badge/playing/")
	svgHeaders(w)

	var games []string
	if p, ok := h.store.Get(id); ok {
		for _, a := range p.Activities {
			if a.Type == gateway.ActivityTypePlaying &&
				a.Name != "Visual Studio Code" &&
				a.Name != "Zed" {
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
