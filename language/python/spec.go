package python

import "code-mate-core/core"

var Spec = core.TokenizerSpec{
	InitialState:  StateGlobal,
	Rules:         GrammarRules,
	FallbackScope: ScopeDefault,
}
