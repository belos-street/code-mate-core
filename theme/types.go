package theme

type HighlightTheme struct {
	ID           string
	DisplayName  string
	DefaultStyle string
	PreStyle     string
	Styles       ThemeStyleMap
}

type ThemeStyleMap map[string]string
