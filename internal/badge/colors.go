package badge

// Badge colors for the supported services and statuses.
const (
	ColorVSCode       = "#23a7f2"
	ColorZed          = "#7c6df2"
	ColorVisualStudio = "#5c2d91"

	ColorSpotify = "#1db954"
	ColorTIDAL   = "#000000"

	ColorCrunchyroll = "#f47521"

	ColorDiscord = "#5865f2"

	ColorOnline  = "#44b700"
	ColorIdle    = "#faa61a"
	ColorDnd     = "#f04747"
	ColorOffline = "#747f8d"

	ColorLabel = "#555"
	ColorDark  = "#1e1e2e"
)

// StatusColors maps Discord presence status names to their badge colors.
var StatusColors = map[string]string{
	"online":  ColorOnline,
	"idle":    ColorIdle,
	"dnd":     ColorDnd,
	"offline": ColorOffline,
}
