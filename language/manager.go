package language

import "code-mate-core/core"

var builtinsRegistered bool

func EnsureBuiltinsRegistered() {
	if builtinsRegistered {
		return
	}
	RegisterBuiltins()
	builtinsRegistered = true
}

func GetLanguage(id string) (core.LanguageAdapter, bool) {
	EnsureBuiltinsRegistered()
	return core.GetLanguage(id)
}

func ListLanguages() []core.LanguageAdapter {
	EnsureBuiltinsRegistered()
	return core.ListLanguages()
}

func Tokenize(code, langID string) (core.TokenStream, error) {
	EnsureBuiltinsRegistered()
	return core.Tokenize(code, langID)
}
