package markdown

const (
	ScopeHeading        = "markup.heading.markdown"
	ScopeSetextHeading  = "markup.heading.setext.markdown"
	ScopeHr             = "markup.hr.markdown"
	ScopeQuote          = "markup.quote.markdown"
	ScopeList           = "markup.list.markdown"
	ScopeTask           = "markup.task.markdown"
	ScopeTableSep       = "markup.table.separator.markdown"
	ScopeTableRow       = "markup.table.row.markdown"
	ScopeBold           = "markup.bold.markdown"
	ScopeItalic         = "markup.italic.markdown"
	ScopeStrikethrough  = "markup.strikethrough.markdown"
	ScopeInlineRaw      = "markup.inline.raw.markdown"
	ScopeIndentedCode   = "markup.code.indented.markdown"
	ScopeFencedBegin    = "markup.fenced-code.block.begin.markdown"
	ScopeFencedEnd      = "markup.fenced-code.block.end.markdown"
	ScopeLinkLabel      = "markup.link.label.markdown"
	ScopeLinkUrl        = "markup.link.url.markdown"
	ScopeLinkPunct      = "markup.link.punctuation.markdown"
	ScopeImageLabel     = "markup.image.label.markdown"
	ScopeRefLinkLabel   = "markup.reference.link.label.markdown"
	ScopeRefLinkUrl     = "markup.reference.link.url.markdown"
	ScopeComment        = "comment.block.markdown"
	ScopeStringUnquoted = "string.unquoted.markdown"
	ScopeDefault        = "default"
)

const (
	GrammarStateGlobal        = "global"
	GrammarStateInlineCode    = "inline-code"
	GrammarStateFencedCode    = "fenced-code"
	GrammarStateComment       = "comment-block"
	GrammarStateStrong        = "strong"
	GrammarStateEmphasis      = "emphasis"
	GrammarStateStrikethrough = "strikethrough"
)
