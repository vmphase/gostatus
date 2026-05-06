package handler

import (
	"fmt"
	"net/http"
	"strings"

	"gostatus/internal/badge"
)

func (h *Handler) VSCode(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/badge/vscode/")
	svgHeaders(w)

	message := qp(r, "fallback", "nothing rn")
	if p, ok := h.store.Get(id); ok {
		for _, a := range p.Activities {
			if a.Name == "Visual Studio Code" && a.Details != "" && a.State != "" {
				file := strings.TrimPrefix(a.Details, "Editing ")
				workspace := strings.ReplaceAll(a.State, "Workspace: ", "")
				workspace = strings.ReplaceAll(workspace, " (Workspace)", "")
				if strings.HasPrefix(workspace, "Glitch:") {
					workspace = strings.Replace(workspace, "Glitch:", "🎏", 1)
				}
				message = file + " in " + workspace
				break
			}
		}
	}

	fmt.Fprint(w, badge.Make(
		qp(r, "label", "vscode"),
		message,
		qp(r, "labelColor", "#555"),
		qp(r, "color", "#23a7f2"),
	))
}
