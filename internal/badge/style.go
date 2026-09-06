package badge

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
