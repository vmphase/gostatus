package badge

const (
	// editors
	ColorVSCode       = "#23a7f2"
	ColorZed          = "#7c6df2"
	ColorVisualStudio = "#5c2d91"

	// music
	ColorSpotify = "#1db954"
	ColorTIDAL   = "#000000"

	// media
	ColorCrunchyroll = "#f47521"

	// discord
	ColorDiscord = "#5865f2"

	// status
	ColorOnline  = "#44b700"
	ColorIdle    = "#faa61a"
	ColorDnd     = "#f04747"
	ColorOffline = "#747f8d"

	// generic
	ColorLabel = "#555"
	ColorDark  = "#1e1e2e"
)

var StatusColors = map[string]string{
	"online":  ColorOnline,
	"idle":    ColorIdle,
	"dnd":     ColorDnd,
	"offline": ColorOffline,
}
