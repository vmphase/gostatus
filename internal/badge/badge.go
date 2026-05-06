package badge

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"text/template"
)

const logoShift = 20

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

type params struct {
	Total, LW, MW     int
	LX, LX1, MX, MX1  int
	LabelColor, Color string
	Label, Message    string
	LogoB64           string
}

func Make(label, message, labelColor, color string) string {
	return render(label, message, labelColor, color, "")
}

func MakeWithLogo(label, message, labelColor, color, logoName string) string {
	logo := ""
	if data, err := os.ReadFile(fmt.Sprintf("assets/logos/%s.svg", logoName)); err == nil {
		logo = base64.StdEncoding.EncodeToString(data)
	}
	return render(label, message, labelColor, color, logo)
}

func render(label, message, labelColor, color, logoB64 string) string {
	lw := len(label)*7 + 10
	mw := len(message)*7 + 10

	shift := 0
	if logoB64 != "" {
		shift = logoShift
	}

	total := lw + mw + shift

	var buf bytes.Buffer
	tpl.Execute(&buf, params{
		Total:      total,
		LW:         lw + shift,
		MW:         mw,
		LX:         lw/2 + shift,
		LX1:        lw/2 + shift + 1,
		MX:         lw + shift + mw/2,
		MX1:        lw + shift + mw/2 + 1,
		LabelColor: labelColor,
		Color:      color,
		Label:      label,
		Message:    message,
		LogoB64:    logoB64,
	})
	return buf.String()
}
