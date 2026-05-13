# 09 — 实现阶段

> **说明**：按 Phase 1-7 顺序实现，每个 Phase 完成后验证再进入下一 Phase。

---

## Phase 1：核心引擎（必须）

| # | 文件 | 内容 |
|---|------|------|
| 1 | `core/types.go` | Token, GrammarRule, ParserContext, TokenizerSpec, TokenStream, TokenStyle |
| 2 | `core/tokenizer.go` | NewParserContext, PushState, PopState, CurrentState, SplitTokenByLineBreak, MatchToken, Parse |
| 3 | `core/registry.go` | LanguageAdapter 接口, RegisterLanguage, GetLanguage, ListLanguages, Tokenize |

## Phase 2：首个语言适配器

| # | 文件 | 内容 |
|---|------|------|
| 4 | `language/javascript/types.go` | TokenScope 常量 + GrammarState 常量 |
| 5 | `language/javascript/rules.go` | GRAMMAR_RULES map（global, multiline-comment, string-double, string-single, string-backtick, template-interpolation, import-dynamic） |
| 6 | `language/javascript/spec.go` | TokenizerSpec 实例 |
| 7 | `language/javascript/adapter.go` | LanguageAdapter 接口实现 |

## Phase 3：主题系统

| # | 文件 | 内容 |
|---|------|------|
| 8 | `theme/types.go` | HighlightTheme, ThemeStyleMap |
| 9 | `theme/darkplus/base.go` | 主题常量 |
| 10 | `theme/darkplus/shared.go` + 语言样式文件 | Dark+ 完整样式 |
| 11 | `theme/registry.go` | RegisterTheme, GetTheme, ResolveTheme, THEME_ALIAS_MAP |
| 12 | `theme/variant.go` | CreateThemeVariant |
| 13 | `theme/presets.go` | 9 个预设主题 |

## Phase 4：公共 API

| # | 文件 | 内容 |
|---|------|------|
| 14 | `api.go` | Highlighter 接口, HighlighterOptions, HtmlOption |
| 15 | `highlighter.go` | Highlighter 实现（CodeToTokens, CodeToHtml, UpdateTheme, 缓存） |
| 16 | `language/manager.go` | 懒加载管理器 |
| 17 | `language/builtins.go` | 内置语言列表 |

## Phase 5：其余 16 种语言适配器

每种语言 4 个文件（types.go, rules.go, spec.go, adapter.go），共 64 个文件。

| 语言 | 目录 |
|------|------|
| TypeScript | `language/typescript/` |
| Python | `language/python/` |
| CSS | `language/css/` |
| Bash | `language/bash/` |
| SQL | `language/sql/` |
| YAML | `language/yaml/` |
| Markdown | `language/markdown/` |
| Java | `language/java/` |
| C | `language/c/` |
| C++ | `language/cpp/` |
| C# | `language/csharp/` |
| Go | `language/gorule/` |
| Rust | `language/rust/` |
| PHP | `language/php/` |
| HTML | `language/html/` |
| JSON | `language/json/` |

## Phase 6：WASM 入口

| # | 文件 | 内容 |
|---|------|------|
| 18 | `cmd/wasm/main.go` | JS API 导出 + 初始化 |

## Phase 7：测试与优化

| # | 文件 | 内容 |
|---|------|------|
| 19 | `core_test.go` | 核心引擎测试 |
| 20 | `language_test.go` | 语言注册表测试 |
| 21 | `api_test.go` | 公共 API 测试 |
| 22 | `bench_test.go` | 性能基准测试 |
| 23 | 一致性验证 | Go ↔ TS 对比 |
