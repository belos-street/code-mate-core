package typescript

import "code-mate-core/core"

type adapter struct{}

func (a adapter) ID() string        { return "typescript" }
func (a adapter) Aliases() []string { return []string{"ts"} }
func (a adapter) Parse(code string) core.TokenStream {
	return core.Parse(code, GetSpec())
}

var Language core.LanguageAdapter = adapter{}
