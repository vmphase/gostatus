package badge

import (
	"bytes"
	"html/template"
	"os"
)

var StatusColors = map[string]string{
	"online":  "#44b700",
	"idle":    "#faa61a",
	"dnd":     "#f04747",
	"offline": "#747f8d",
}

var tpl = func() *template.Template {
	b, err := os.ReadFile("assets/badge.svg")
	if err != nil {
		panic(err)
	}
	return template.Must(template.New("badge").Parse(string(b)))
}()

func Make(label, message, labelColor, color string) string {
	lw := len(label)*7 + 10
	mw := len(message)*7 + 10

	var buf bytes.Buffer
	
	tpl.Execute(&buf, map[string]any{
		"Total":      lw + mw,
		"LW":         lw,
		"MW":         mw,
		"LabelColor": labelColor,
		"Color":      color,
		"Label":      label,
		"Message":    message,
		"LX":         lw / 2,
		"LX1":        lw/2 + 1,
		"MX":         lw + mw/2,
		"MX1":        lw + mw/2 + 1,
	})
	return buf.String()
}
