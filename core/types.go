package core

import "regexp"

type TokenStyle map[string]string

type Token struct {
	Text  string
	Scope string
	Line  int
	Col   [2]int
	Style TokenStyle
}

type GrammarRule struct {
	Regex     *regexp.Regexp
	Scope     string
	PushState string
	PopState  bool
	Skip      bool
}

type ParserContext struct {
	StateStack []string
	Line       int
	Col        int
}

type TokenizerSpec struct {
	InitialState  string
	Rules         map[string][]GrammarRule
	FallbackScope string
}

type TokenStream [][]Token
