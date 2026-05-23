package handler

import (
	"fmt"
	"net/http"
	"strings"

	"gostatus/internal/badge"
	"gostatus/internal/gateway"
)

// single code editor activity
type editor struct {
	// discord activity name to match against
	name string
	// ?prefer= to choose this editor, when none the first match is used
	prefer string
	// left-hand side text
	label string
	// right-hand side color
	color string
	// logo is the shields.io logo slug (empty: no logo)
	logo string
	// builds the badge text from a matched activity
	// "" when the activity lacks the required fields
	message func(a *gateway.Activity) string
}

// ordered registry of supported code editors
var editors = []editor{
	{
		name:   "Visual Studio Code",
		prefer: "vscode",
		label:  "vscode",
		color:  "#23a7f2",
		logo:   "vscode",
		message: func(a *gateway.Activity) string {
			if a.Details == "" {
				return ""
			}
			file := strings.TrimPrefix(a.Details, "Editing ")
			ws := strings.ReplaceAll(strings.ReplaceAll(a.State, "Workspace: ", ""), " (Workspace)", "")
			if ws == "" {
				return file
			}
			return file + " in " + ws
		},
	},
	{
		name:   "Zed",
		prefer: "zed",
		label:  "zed",
		color:  "#7c6df2",
		logo:   "zed",
		message: func(a *gateway.Activity) string {
			if a.Details == "" {
				return ""
			}
			file := strings.TrimPrefix(a.State, "Working on ")
			workspace := strings.TrimPrefix(a.Details, "In ")
			if file == "" {
				return workspace
			}
			return file + " in " + workspace
		},
	},
}

func (h *Handler) Code(w http.ResponseWriter, r *http.Request) {
	svgHeaders(w)

	fallback := qp(r, "fallback", "nothing")
	prefer := strings.ToLower(r.URL.Query().Get("prefer"))

	defaultEditor := func() editor {
		if prefer != "" {
			for _, ed := range editors {
				if ed.prefer == prefer {
					return ed
				}
			}
		}
		return editors[0]
	}

	p, ok := h.store.Get(h.id(r, "/badge/code/"))
	if !ok {
		def := defaultEditor()
		h.renderBadge(w, r, qp(r, "label", def.label), fallback, qp(r, "labelColor", "#1e1e2e"), qp(r, "color", def.color), def.logo)
		return
	}

	type match struct {
		ed  editor
		msg string
	}

	var matches []match
	for _, ed := range editors {
		a := FindActivity(p, gateway.ActivityTypePlaying, ed.name)
		if a == nil {
			continue
		}
		msg := ed.message(a)
		if msg == "" {
			continue
		}
		matches = append(matches, match{ed, msg})
	}

	if len(matches) == 0 {
		def := defaultEditor()
		h.renderBadge(w, r, qp(r, "label", def.label), fallback, qp(r, "labelColor", "#1e1e2e"), qp(r, "color", def.color), def.logo)
		return
	}

	chosen := matches[0]
	if prefer != "" {
		for _, m := range matches {
			if m.ed.prefer == prefer {
				chosen = m
				break
			}
		}
	}

	label := qp(r, "label", chosen.ed.label)
	color := qp(r, "color", chosen.ed.color)
	labelColor := qp(r, "labelColor", "#1e1e2e")
	fmt.Fprint(w, badge.MakeWithLogo(label, chosen.msg, labelColor, color, chosen.ed.logo))
}
