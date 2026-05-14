package javascript

import "code-mate-core/core"

var Spec = core.TokenizerSpec{
	InitialState:  StateGlobal,
	Rules:         GRAMMAR_RULES,
	FallbackScope: ScopeDefault,
}
