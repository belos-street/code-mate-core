package theme

import (
	"fmt"
	"strings"
)

var DefaultThemeID = "dark-plus"

var THEME_ALIAS_MAP = map[string]string{
	"dark":   "dark-plus",
	"github": "github-light",
}

var themeRegistry = make(map[string]*HighlightTheme)

func normalizeThemeID(id string) string {
	return strings.TrimSpace(strings.ToLower(id))
}

func RegisterTheme(theme *HighlightTheme) error {
	normalizedID := normalizeThemeID(theme.ID)
	if normalizedID == "" {
		return fmt.Errorf("theme id cannot be empty")
	}

	if existing, ok := themeRegistry[normalizedID]; ok && existing != theme {
		return fmt.Errorf("theme %q is already registered", theme.ID)
	}

	themeRegistry[normalizedID] = theme
	return nil
}

func GetTheme(themeID string) (*HighlightTheme, error) {
	if themeID == "" {
		themeID = DefaultThemeID
	}

	// 先检查别名
	if resolvedID, ok := THEME_ALIAS_MAP[normalizeThemeID(themeID)]; ok {
		themeID = resolvedID
	}

	normalizedID := normalizeThemeID(themeID)
	if normalizedID == "" {
		return nil, fmt.Errorf("theme id cannot be empty")
	}

	theme, ok := themeRegistry[normalizedID]
	if !ok {
		return nil, fmt.Errorf("theme %q is not registered", normalizedID)
	}

	return theme, nil
}

func ListThemes() []*HighlightTheme {
	themes := make([]*HighlightTheme, 0, len(themeRegistry))
	for _, theme := range themeRegistry {
		themes = append(themes, theme)
	}
	return themes
}

func ResolveAlias(theme interface{}) interface{} {
	if s, ok := theme.(string); ok {
		normalized := normalizeThemeID(s)
		if resolved, ok := THEME_ALIAS_MAP[normalized]; ok {
			return resolved
		}
		return normalized
	}
	return theme
}

func ResolveTheme(theme interface{}) (*HighlightTheme, error) {
	switch v := theme.(type) {
	case nil:
		return GetTheme("")
	case string:
		return GetTheme(v)
	case *HighlightTheme:
		return v, nil
	default:
		return nil, fmt.Errorf("invalid theme type: %T", theme)
	}
}

// ResolveScopeStyle 查找 scope 对应的样式：
// 1) 精确匹配
// 2) 前缀回退（a.b.c -> a.b -> a）
// 3) defaultStyle
func ResolveScopeStyle(scope string, theme *HighlightTheme) string {
	if style, ok := theme.Styles[scope]; ok {
		return style
	}

	candidate := scope
	for strings.Contains(candidate, ".") {
		idx := strings.LastIndex(candidate, ".")
		candidate = candidate[:idx]
		if style, ok := theme.Styles[candidate]; ok {
			return style
		}
	}

	if style, ok := theme.Styles["default"]; ok {
		return style
	}
	return theme.DefaultStyle
}

func ClearRegistry() {
	themeRegistry = make(map[string]*HighlightTheme)
}
