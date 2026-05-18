package handler

import (
	"fmt"
	"net/http"
	"strings"

	"gostatus/internal/badge"
	"gostatus/internal/gateway"
)

func (h *Handler) CrunchyRoll(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/badge/crunchyroll/")
	svgHeaders(w)

	message := qp(r, "fallback", "nothing")

	if p, ok := h.store.Get(id); ok {
		for _, a := range p.Activities {
			if a.Type == gateway.ActivityTypeWatching && a.Name == "Crunchyroll" && a.Details != "" {
				message = fmt.Sprintf("%s – %s", a.Details, a.State)
				break
			}
		}
	}

	make := badge.Make
	if qp(r, "hideLogo", "false") != "true" {
		make = func(label, msg, lc, c string) string {
			return badge.MakeWithLogo(label, msg, lc, c, "crunchyroll")
		}
	}

	fmt.Fprint(w, make(
		qp(r, "label", "watching"),
		message,
		qp(r, "labelColor", "#555"),
		qp(r, "color", "#5865f2"),
	))
}
