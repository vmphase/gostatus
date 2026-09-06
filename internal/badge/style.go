package badge

// Style is a badge visual style.
type Style string

// Supported badge styles.
const (
	StyleFlat        Style = "flat"
	StyleFlatSquare  Style = "flat-square"
	StyleForTheBadge Style = "for-the-badge"
)

// ParseStyle converts a style string into a Style, defaulting to flat for unknown values.
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
