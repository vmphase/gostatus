package badge

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type Style string

const (
	StyleFlat        Style = "flat"
	StyleFlatSquare  Style = "flat-square"
	StyleForTheBadge Style = "for-the-badge"
)

func ParseStyle(s string) Style {
	switch s {
	case "flat-square":
		return StyleFlatSquare
	case "for-the-badge":
		return StyleForTheBadge
	default:
		return StyleFlat
	}
}

var tpl = func() *template.Template {
	b, err := os.ReadFile("assets/badge.svg")
	if err != nil {
		panic(err)
	}
	return template.Must(template.New("badge").Parse(string(b)))
}()

type Options struct {
	Label      string
	Message    string
	LabelColor string
	Color      string
	Logo       string
	Style      Style
}

type templateParams struct {
	Total, LW, MW     int
	LX, LX1, MX, MX1  int
	LabelColor, Color string
	Label, Message    string
	LogoB64           string
	Style             Style
}

const (
	logoShift = 15
	flatCharW = 7
	flatPad   = 10
	ftbCharW  = 8
	ftbPad    = 20
)

func Render(o Options) string {
	logoB64 := loadLogo(o.Logo)

	label, message := o.Label, o.Message
	charW, pad := flatCharW, flatPad
	if o.Style == StyleForTheBadge {
		label = strings.ToUpper(label)
		message = strings.ToUpper(message)
		charW, pad = ftbCharW, ftbPad
	}

	lw := len(label)*charW + pad
	mw := len(message)*charW + pad
	shift := 0
	if logoB64 != "" {
		shift = logoShift
	}
	total := lw + mw + shift

	var buf bytes.Buffer
	tpl.Execute(&buf, templateParams{
		Total:      total,
		LW:         lw + shift,
		MW:         mw,
		LX:         lw/2 + shift,
		LX1:        lw/2 + shift + 1,
		MX:         lw + shift + mw/2,
		MX1:        lw + shift + mw/2 + 1,
		LabelColor: o.LabelColor,
		Color:      o.Color,
		Label:      label,
		Message:    message,
		LogoB64:    logoB64,
		Style:      o.Style,
	})
	return buf.String()
}

func loadLogo(name string) string {
	if name == "" {
		return ""
	}
	data, err := os.ReadFile(fmt.Sprintf("assets/logos/%s.svg", name))
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(data)
}
