package language

import (
	"testing"

	"code-mate-core/core"
)

func init() {
	core.ClearRegistry()
	ResetBuiltins()
	RegisterBuiltins()
}

func TestGetLanguage_Found(t *testing.T) {
	EnsureBuiltinsRegistered()
	lang, ok := GetLanguage("javascript")
	if !ok {
		t.Fatal("expected to find 'javascript'")
	}
	if lang.ID() != "javascript" {
		t.Errorf("expected ID 'javascript', got %q", lang.ID())
	}
}

func TestGetLanguage_ByAlias(t *testing.T) {
	EnsureBuiltinsRegistered()
	lang, ok := GetLanguage("js")
	if !ok {
		t.Fatal("expected to find 'js' alias")
	}
	if lang.ID() != "javascript" {
		t.Errorf("expected ID 'javascript' via alias, got %q", lang.ID())
	}
}

func TestGetLanguage_NotFound(t *testing.T) {
	EnsureBuiltinsRegistered()
	_, ok := GetLanguage("nonexistent")
	if ok {
		t.Fatal("expected not found for nonexistent language")
	}
}

func TestGetLanguage_CaseInsensitive(t *testing.T) {
	EnsureBuiltinsRegistered()
	_, ok := GetLanguage("PYTHON")
	if !ok {
		t.Fatal("expected case-insensitive lookup to succeed")
	}
	_, ok = GetLanguage("  python  ")
	if !ok {
		t.Fatal("expected trimmed lookup to succeed")
	}
}

func TestListLanguages_ContainsBuiltins(t *testing.T) {
	EnsureBuiltinsRegistered()
	list := ListLanguages()
	expectedIDs := []string{
		"javascript", "typescript", "python", "bash", "c",
		"cpp", "csharp", "css", "go", "html", "java", "json",
		"markdown", "php", "rust", "sql", "yaml",
	}
	foundIDs := make(map[string]bool)
	for _, lang := range list {
		foundIDs[lang.ID()] = true
	}
	for _, id := range expectedIDs {
		if !foundIDs[id] {
			t.Errorf("expected language %q not found in list", id)
		}
	}
}

func TestListLanguages_NoDuplicates(t *testing.T) {
	EnsureBuiltinsRegistered()
	list := ListLanguages()
	seen := make(map[string]bool)
	for _, lang := range list {
		id := lang.ID()
		if seen[id] {
			t.Errorf("duplicate language ID: %q", id)
		}
		seen[id] = true
	}
}

func TestTokenize_Success(t *testing.T) {
	EnsureBuiltinsRegistered()
	rows, err := Tokenize("const x = 1;", "javascript")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rows) == 0 {
		t.Fatal("expected at least 1 row")
	}
	if len(rows[0]) == 0 {
		t.Fatal("expected at least 1 token")
	}
}

func TestTokenize_UnknownLanguage(t *testing.T) {
	EnsureBuiltinsRegistered()
	_, err := Tokenize("hello", "unknown")
	if err == nil {
		t.Fatal("expected error for unknown language")
	}
}

func TestTokenize_EmptyCode(t *testing.T) {
	EnsureBuiltinsRegistered()
	rows, err := Tokenize("", "javascript")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rows) != 0 {
		t.Errorf("expected 0 rows for empty code, got %d", len(rows))
	}
}

func TestRegisteredAliases(t *testing.T) {
	EnsureBuiltinsRegistered()
	tests := []struct {
		alias string
		want  string
	}{
		{"js", "javascript"},
		{"ts", "typescript"},
		{"py", "python"},
		{"sh", "bash"},
		{"shell", "bash"},
		{"h", "c"},
		{"golang", "go"},
		{"rs", "rust"},
		{"md", "markdown"},
		{"yml", "yaml"},
		{"htm", "html"},
		{"cs", "csharp"},
		{"phtml", "php"},
		{"jav", "java"},
		{"c++", "cpp"},
	}
	for _, tt := range tests {
		lang, ok := GetLanguage(tt.alias)
		if !ok {
			t.Errorf("alias %q should resolve to a language", tt.alias)
			continue
		}
		if lang.ID() != tt.want {
			t.Errorf("alias %q resolved to %q, want %q", tt.alias, lang.ID(), tt.want)
		}
	}
}
