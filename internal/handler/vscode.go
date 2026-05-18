package handler

import (
	"fmt"
	"net/http"
	"strings"

	"gostatus/internal/badge"
	"gostatus/internal/gateway"
)

func (h *Handler) VSCode(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/badge/vscode/")
	svgHeaders(w)

	message := qp(r, "fallback", "nothing")
	if p, ok := h.store.Get(id); ok {
		for _, a := range p.Activities {
			if a.Type == gateway.ActivityTypePlaying && a.Name == "Visual Studio Code" && a.Details != "" {
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

	make := badge.Make
	if qp(r, "hideLogo", "false") != "true" {
		make = func(label, msg, lc, c string) string {
			return badge.MakeWithLogo(label, msg, lc, c, "vscode")
		}
	}

	fmt.Fprint(w, make(
		qp(r, "label", "vscode"),
		message,
		qp(r, "labelColor", "#555"),
		qp(r, "color", "#23a7f2"),
	))
}
