package handler

import (
	"fmt"
	"net/http"
	"strings"

	"gostatus/internal/badge"
)

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/badge/status/")
	svgHeaders(w)

	status := "offline"
	if p, ok := h.store.Get(id); ok {
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
