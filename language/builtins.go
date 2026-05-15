package language

import (
	"code-mate-core/core"
	"code-mate-core/language/javascript"
)

var builtinLanguages = []core.LanguageAdapter{
	javascript.Language,
}

func RegisterBuiltins() {
	for _, lang := range builtinLanguages {
		core.RegisterLanguage(lang)
	}
}

func ResetBuiltins() {
	builtinsRegistered = false
}
