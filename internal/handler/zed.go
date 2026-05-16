package handler

import (
	"fmt"
	"net/http"
	"strings"

	"gostatus/internal/badge"
)

func (h *Handler) Zed(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/badge/zed/")
	svgHeaders(w)

	message := qp(r, "fallback", "nothing")
	if p, ok := h.store.Get(id); ok {
		for _, a := range p.Activities {
			if a.Name == "Zed" && a.Details != "" && a.State != "" {
				file := strings.TrimPrefix(a.State, "Working on ")
				workspace := strings.TrimPrefix(a.Details, "In ")
				message = file + " in " + workspace
				break
			}
		}
	}

	make := badge.Make
	if qp(r, "hideLogo", "false") != "true" {
		make = func(label, msg, lc, c string) string {
			return badge.MakeWithLogo(label, msg, lc, c, "zed")
		}
	}

	fmt.Fprint(w, make(
		qp(r, "label", "zed"),
		message,
		qp(r, "labelColor", "#1e1e2e"),
		qp(r, "color", "#7c6df2"),
	))
}
