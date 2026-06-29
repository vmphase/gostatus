package handler

import (
	"net/http"
	"strings"

	"gostatus/internal/gateway"
	"gostatus/internal/badge"
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
		color:  badge.ColorVSCode,
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
		color:  badge.ColorZed,
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
	{
		name:   "Visual Studio",
		prefer: "visualstudio",
		label:  "visual studio",
		color:  badge.ColorVisualStudio,
		logo:   "visualstudio",
		message: func(a *gateway.Activity) string {
			if a.Details == "" {
				return ""
			}
			file := strings.TrimPrefix(a.Details, "File ")
			solution := strings.TrimPrefix(a.State, "Solution ")
			if solution == "" {
				return file
			}
			return file + " in " + solution
		},
	},
}

func findEditor(prefer string) editor {
	if len(editors) == 0 {
		return editor{label: "code", color: badge.ColorDiscord}
	}
	for _, ed := range editors {
		if ed.prefer == prefer {
			return ed
		}
	}
	return editors[0]
}

func (h *Handler) Code(w http.ResponseWriter, r *http.Request) {
	svgHeaders(w)

	fallback := qp(r, "fallback", "nothing")
	prefer := strings.ToLower(r.URL.Query().Get("prefer"))

	renderEditor := func(ed editor, message string) {
		h.renderBadge(w, r,
			qp(r, "label", ed.label),
			message,
			qp(r, "labelColor", badge.ColorDark),
			qp(r, "color", ed.color),
			ed.logo,
		)
	}

	p, ok := h.store.Get(h.id(r, "/badge/code/"))
	if !ok {
		renderEditor(findEditor(prefer), fallback)
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
		if msg := ed.message(a); msg != "" {
			matches = append(matches, match{ed, msg})
		}
	}

	if len(matches) == 0 {
		renderEditor(findEditor(prefer), fallback)
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
	renderEditor(chosen.ed, chosen.msg)
}
