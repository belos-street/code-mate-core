package sql

import "code-mate-core/core"

var Spec = core.TokenizerSpec{
	InitialState:  GrammarStateGlobal,
	Rules:         GrammarRules,
	FallbackScope: ScopeDefault,
}
