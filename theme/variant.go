package theme

import (
	"regexp"
	"strings"
)

var hexColorPattern = regexp.MustCompile(`#[0-9A-Fa-f]{6}`)

type VariantOptions struct {
	ID           string
	DisplayName  string
	DefaultStyle string
	PreStyle     string
	ColorMap     map[string]string
}

func normalizeColor(color string) string {
	return strings.ToUpper(color)
}

func normalizeColorMap(colorMap map[string]string) map[string]string {
	normalized := make(map[string]string, len(colorMap))
	for sourceColor, targetColor := range colorMap {
		normalized[normalizeColor(sourceColor)] = targetColor
	}
	return normalized
}

func replaceStyleColors(style string, colorMap map[string]string) string {
	return hexColorPattern.ReplaceAllStringFunc(style, func(color string) string {
		if mappedColor, ok := colorMap[normalizeColor(color)]; ok {
			return mappedColor
		}
		return color
	})
}

func CreateThemeVariant(base *HighlightTheme, opts VariantOptions) *HighlightTheme {
	normalizedColorMap := normalizeColorMap(opts.ColorMap)
	styles := make(ThemeStyleMap, len(base.Styles))

	for scope, style := range base.Styles {
		styles[scope] = replaceStyleColors(style, normalizedColorMap)
	}

	return &HighlightTheme{
		ID:           opts.ID,
		DisplayName:  opts.DisplayName,
		DefaultStyle: opts.DefaultStyle,
		PreStyle:     opts.PreStyle,
		Styles:       styles,
	}
}
