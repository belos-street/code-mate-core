package core

import (
	"fmt"
	"reflect"
	"strings"
)

type LanguageAdapter interface {
	ID() string
	Aliases() []string
	Parse(code string) TokenStream
}

var languageRegistry = make(map[string]LanguageAdapter)

func normalizeID(id string) string {
	return strings.TrimSpace(strings.ToLower(id))
}

func ensureValidID(id string) string {
	normalized := normalizeID(id)
	if normalized == "" {
		panic("Language id cannot be empty")
	}
	return normalized
}

func RegisterLanguage(lang LanguageAdapter) error {
	canonicalID := ensureValidID(lang.ID())
	keys := []string{canonicalID}
	for _, alias := range lang.Aliases() {
		keys = append(keys, ensureValidID(alias))
	}

	for _, key := range keys {
		existing, ok := languageRegistry[key]
		if ok && !reflect.DeepEqual(existing, lang) {
			return fmt.Errorf("language id or alias %q is already registered", key)
		}
	}

	for _, key := range keys {
		languageRegistry[key] = lang
	}

	return nil
}

func GetLanguage(id string) (LanguageAdapter, bool) {
	normalized := ensureValidID(id)
	lang, ok := languageRegistry[normalized]
	return lang, ok
}

func ListLanguages() []LanguageAdapter {
	seen := make(map[string]bool)
	var languages []LanguageAdapter

	for _, lang := range languageRegistry {
		key := normalizeID(lang.ID())
		if seen[key] {
			continue
		}
		seen[key] = true
		languages = append(languages, lang)
	}

	return languages
}

func Tokenize(code, langID string) (TokenStream, error) {
	if langID == "" {
		return nil, fmt.Errorf("language id cannot be empty")
	}
	lang, ok := GetLanguage(langID)
	if !ok {
		return nil, fmt.Errorf("language %q is not registered", normalizeID(langID))
	}
	return lang.Parse(code), nil
}

func ClearRegistry() {
	languageRegistry = make(map[string]LanguageAdapter)
}
