package rust

const (
	ScopeCommentLineDS   = "comment.line.double-slash.rust"
	ScopeCommentBlock    = "comment.block.rust"
	ScopeMetaAttribute   = "meta.attribute.rust"
	ScopeKeywordControl  = "keyword.control.rust"
	ScopeKeywordDecl     = "keyword.declaration.rust"
	ScopeKeywordModifier = "keyword.modifier.rust"
	ScopeSupportBuiltin  = "support.type.builtin.rust"
	ScopeEntityType      = "entity.name.type.rust"
	ScopeEntityFunction  = "entity.name.function.rust"
	ScopeNamespace       = "entity.name.namespace.rust"
	ScopeEntityMacro     = "entity.name.macro.rust"
	ScopeVarLang         = "variable.language.rust"
	ScopeConstantBoolean = "constant.language.boolean.rust"
	ScopeConstantNull    = "constant.language.null.rust"
	ScopeConstantNumeric = "constant.numeric.rust"
	ScopeStringDouble    = "string.quoted.double.rust"
	ScopeStringSingle    = "string.quoted.single.rust"
	ScopeStringRaw       = "string.quoted.raw.rust"
	ScopeOperator        = "operator.rust"
	ScopeDefault         = "default"
)

const (
	GrammarStateGlobal       = "global"
	GrammarStateCommentBlock = "comment-block"
	GrammarStateStringDouble = "string-double"
	GrammarStateStringSingle = "string-single"
	GrammarStateStringRaw    = "string-raw"
)
