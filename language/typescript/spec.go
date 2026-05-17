package typescript

import (
	"sync"

	"code-mate-core/core"
)

var (
	spec     *core.TokenizerSpec
	specOnce sync.Once
)

func GetSpec() *core.TokenizerSpec {
	specOnce.Do(func() {
		s := core.TokenizerSpec{
			InitialState:  "global",
			Rules:         getGrammarRules(),
			FallbackScope: "default",
		}
		spec = &s
	})
	return spec
}
