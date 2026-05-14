# Go + WASM 复刻实现计划

> **目标**：将 `ts/lib/` 语法高亮引擎复刻为 Go，编译为 WASM
> **策略**：按 Phase 顺序，每轮对话完成一个子任务，逐步推进
> **参考**：specs/01-10, `ts/lib/` 源码

---

## 对话约定

每轮对话开始时提供：
1. **本轮目标**：要完成什么
2. **输入文件**：需要参考的 TS 源码和 spec
3. **输出文件**：本轮要创建的 Go 文件
4. **验收标准**：完成的判断条件

每轮对话结束时确认：
- 文件是否创建完成
- `go build ./...` 是否通过
- 下轮需要的上下文



## Phase 1：核心分词引擎

> **参考 spec**：[03-core-engine.md](../specs/03-core-engine.md)
> **参考 TS 源码**：`ts/lib/core/types.ts`, `ts/lib/core/tokenizer.ts`, `ts/lib/core/registry.ts`

### 任务 1.1 — 核心类型定义

**目标**：创建 `core/types.go`

**输出文件**：`core/types.go`

**内容**：
- `TokenStyle` = `map[string]string`
- `Token` 结构体（Text, Scope, Line, Col [2]int, Style）
- `GrammarRule` 结构体（Regex *regexp.Regexp, Scope, PushState, PopState, Skip）
- `ParserContext` 结构体（StateStack []string, Line, Col）
- `TokenizerSpec` 结构体（InitialState, Rules map[string][]GrammarRule, FallbackScope）
- `TokenStream` = `[][]Token`

**关键约束**：
- `Col` 是 `[2]int`，闭区间 `[startCol, endCol]`
- `PushState` 空字符串 = 不压栈
- `PopState` true = 匹配后弹栈
- `Skip` 保留字段，当前不处理

**验收**：`go build ./core` 通过

### 任务 1.2 — FSM 分词引擎（上下文操作）

**目标**：创建 `core/tokenizer.go` 的上下文管理函数

**输出文件**：`core/tokenizer.go`

**先实现 4 个函数**：
1. `NewParserContext(initialState string) *ParserContext`
   - StateStack 初始化为 `[initialState]`，Line=1, Col=1
2. `PushState(ctx *ParserContext, state string) *ParserContext`
   - 返回新上下文（值拷贝），StateStack 追加 state
3. `PopState(ctx *ParserContext) (*ParserContext, error)`
   - 返回新上下文，StateStack 弹出末尾。栈长度 <= 1 时返回 error
4. `CurrentState(ctx *ParserContext) string`
   - 返回 StateStack 末尾元素

**实现要点**：
- 值拷贝语义：`newCtx := *ctx`
- StateStack 拷贝：`newStack := make([]string, len(ctx.StateStack)); copy(newStack, ctx.StateStack)`

**验收**：`go build ./core` 通过

### 任务 1.3 — FSM 分词引擎（分词核心）

**目标**：继续 `core/tokenizer.go`，实现分词逻辑

**参考 TS 源码**：`ts/lib/core/tokenizer.ts` 的 `splitTokenByLineBreak`, `matchToken`, `parse`

**实现 3 个函数**：

1. `SplitTokenByLineBreak(text, tokenScope, fallbackScope string, startLine, startCol int) []Token`
   - 处理 `\n` 和 `\r\n`
   - 换行符 token 使用 fallbackScope
   - 正确维护 line/col 计数

2. `MatchToken(code string, ctx *ParserContext, spec *TokenizerSpec) (*Token, *ParserContext)`
   - 获取当前状态的规则列表
   - 按顺序遍历规则，用正则从代码头部匹配（`^` 锚定）
   - 第一个匹配胜出
   - 根据 PushState/PopState 修改上下文
   - 所有规则不匹配时，消耗一个字符作为 fallback（使用 spec.FallbackScope）
   - 仅当 code 长度为 0 时返回 nil

3. `Parse(code string, spec *TokenizerSpec) TokenStream`
   - 维护 remainingCode 和 context
   - 每次迭代：先检查换行符 → 调用 MatchToken → 处理跨行 token
   - 最后一行无换行符也要推入结果

**关键对照**：与 TS 版本的 `parse` 函数逻辑必须一致（对照 tokenizer.ts:169-263）

**验收**：`go build ./core` 通过

### 任务 1.4 — 语言注册表

**目标**：创建 `core/registry.go`

**输出文件**：`core/registry.go`

**内容**：
```go
type LanguageAdapter interface {
    ID() string
    Aliases() []string
    Parse(code string) TokenStream
}
```

**实现 4 个函数**：
1. `RegisterLanguage(lang LanguageAdapter) error` — 注册 id + aliases，已注册报错
2. `GetLanguage(id string) (LanguageAdapter, bool)` — 大小写不敏感查找
3. `ListLanguages() []LanguageAdapter` — 去重列出
4. `Tokenize(code, langID string) (TokenStream, error)` — 便捷函数

**语言 id 规范化**：`strings.TrimSpace(strings.ToLower(id))`

**验收**：`go build ./core` 通过，Phase 1 全部完成

---

## Phase 2：JavaScript 语言适配器

> **参考 spec**：[04-language-adapters.md](../specs/04-language-adapters.md)
> **参考 TS 源码**：`ts/lib/language/javascript/` (type.ts, rule.ts, spec.ts, engine.ts, index.ts)

### 任务 2.1 — JavaScript 类型定义

**目标**：创建 `language/javascript/types.go`

**输出文件**：`language/javascript/types.go`

**内容**：
- `TokenScope` 类型（string）+ 所有 scope 常量（从 TS 的 type.ts 翻译）
- `GrammarState` 类型（string）+ 所有状态常量（global, multiline-comment, string-double, string-single, string-backtick, template-interpolation, import-dynamic）

**参考**：`ts/lib/language/javascript/type.ts`

**验收**：`go build ./language/javascript` 通过

### 任务 2.2 — JavaScript 语法规则

**目标**：创建 `language/javascript/rules.go`

**输出文件**：`language/javascript/rules.go`

**内容**：
- `var GRAMMAR_RULES = map[string][]core.GrammarRule{...}`
- 所有正则必须预编译：`regexp.MustCompile(...)`
- 规则顺序必须与 TS 版本一致（影响匹配优先级）
- 状态：global, multiline-comment, string-double, string-single, string-backtick, template-interpolation, import-dynamic

**参考**：`ts/lib/language/javascript/rule.ts`

**关键**：这是最复杂的文件，规则数量多，需要逐条对照翻译

**验收**：`go build ./language/javascript` 通过

### 任务 2.3 — JavaScript Spec + Adapter

**目标**：创建 `language/javascript/spec.go` 和 `language/javascript/adapter.go`

**输出文件**：
- `language/javascript/spec.go` — TokenizerSpec 实例
- `language/javascript/adapter.go` — LanguageAdapter 实现

**spec.go 内容**：
```go
var Spec = core.TokenizerSpec{
    InitialState:  "global",
    Rules:         GRAMMAR_RULES,
    FallbackScope: "default",
}
```

**adapter.go 内容**：
```go
type adapter struct{}

func (a adapter) ID() string       { return "javascript" }
func (a adapter) Aliases() []string { return []string{"js"} }
func (a adapter) Parse(code string) core.TokenStream {
    return core.Parse(code, &Spec)
}

var Language core.LanguageAdapter = adapter{}
```

**验收**：`go build ./language/javascript` 通过，Phase 2 全部完成

---

## Phase 3：主题系统

> **参考 spec**：[05-theme-system.md](../specs/05-theme-system.md)
> **参考 TS 源码**：`ts/lib/themes/` (types.ts, index.ts, theme-variant.ts, presets.ts, dark-plus/)

### 任务 3.1 — 主题类型 + 注册表

**目标**：创建 `theme/types.go` 和 `theme/registry.go`

**输出文件**：
- `theme/types.go` — HighlightTheme, ThemeStyleMap
- `theme/registry.go` — RegisterTheme, GetTheme, ListThemes, ResolveTheme, THEME_ALIAS_MAP

**types.go 内容**：
```go
type ThemeStyleMap map[string]string

type HighlightTheme struct {
    ID           string
    DisplayName  string
    DefaultStyle string
    PreStyle     string
    Styles       ThemeStyleMap
}
```

**registry.go 内容**：
- `RegisterTheme(theme *HighlightTheme) error`
- `GetTheme(id string) (*HighlightTheme, error)`
- `ListThemes() []*HighlightTheme`
- `ResolveTheme(theme interface{}) *HighlightTheme` — 接收 string 或 *HighlightTheme
- `THEME_ALIAS_MAP`：github → github-light, dark → dark-plus
- 默认主题 ID：`"dark-plus"`

**验收**：`go build ./theme` 通过

### 任务 3.2 — 主题变体生成器

**目标**：创建 `theme/variant.go`

**输出文件**：`theme/variant.go`

**内容**：
```go
func CreateThemeVariant(base *HighlightTheme, opts VariantOptions) *HighlightTheme
```

**算法**（对照 theme-variant.ts）：
1. 遍历 base.Styles 的每个 scope → style
2. 用正则 `#[0-9A-Fa-f]{6}` 找到所有颜色值
3. 将颜色值大写后在 colorMap 中查找替换
4. 生成新的 HighlightTheme

**验收**：`go build ./theme` 通过

### 任务 3.3 — Dark+ 基础主题（base + shared）

**目标**：创建 `theme/darkplus/base.go` 和 `theme/darkplus/shared.go`

**输出文件**：
- `theme/darkplus/base.go` — 主题常量（ID, DisplayName, DefaultStyle, PreStyle）
- `theme/darkplus/shared.go` — 通用样式（`"default": "color: #D4D4D4;"`）

**验收**：`go build ./theme/darkplus` 通过

### 任务 3.4 — Dark+ JavaScript 样式

**目标**：创建 `theme/darkplus/js.go`

**输出文件**：`theme/darkplus/js.go`

**内容**：
- `var darkPlusJsStyles = map[string]string{...}`
- 从 `ts/lib/themes/dark-plus/js.ts` 翻译所有 scope → CSS style 映射

**验收**：`go build ./theme/darkplus` 通过

### 任务 3.5 — Dark+ 其余语言样式

**目标**：创建剩余 16 个语言样式文件

**输出文件**（每个文件导出 `map[string]string`）：
- `theme/darkplus/typescript.go`
- `theme/darkplus/python.go`
- `theme/darkplus/css.go`
- `theme/darkplus/bash.go`
- `theme/darkplus/sql.go`
- `theme/darkplus/yaml.go`
- `theme/darkplus/markdown.go`
- `theme/darkplus/java.go`
- `theme/darkplus/c.go`
- `theme/darkplus/cpp.go`
- `theme/darkplus/csharp.go`
- `theme/darkplus/go.go`
- `theme/darkplus/rust.go`
- `theme/darkplus/php.go`
- `theme/darkplus/html.go`
- `theme/darkplus/json.go`

**参考**：`ts/lib/themes/dark-plus/` 下对应文件

**验收**：`go build ./theme/darkplus` 通过

### 任务 3.6 — Dark+ 组装 + 预设主题

**目标**：创建 `theme/darkplus/index.go` 和 `theme/presets.go`

**输出文件**：
- `theme/darkplus/index.go` — 组装所有语言样式为完整 HighlightTheme
- `theme/presets.go` — 9 个预设主题（通过 CreateThemeVariant 生成）

**presets.go 包含**：
- `VariantPalette` 结构体（12 种颜色）
- 9 个预设定义（github-light, dracula, one-dark-pro, nord, monokai, material-ocean, tokyo-night, solarized-dark, solarized-light）
- 每个预设的调色板数据（从 presets.ts 翻译）

**验收**：`go build ./theme` 通过，Phase 3 全部完成

---

## Phase 4：公共 API

> **参考 spec**：[06-public-api.md](../specs/06-public-api.md)
> **参考 TS 源码**：`ts/lib/api.ts`, `ts/lib/language/manager.ts`, `ts/lib/language/builtins.ts`

### 任务 4.1 — 语言管理器 + 内置语言列表

**目标**：创建 `language/manager.go` 和 `language/builtins.go`

**输出文件**：
- `language/manager.go` — 懒加载管理器
- `language/builtins.go` — 内置语言注册

**manager.go 内容**（对照 manager.ts）：
- 封装 `core.RegisterLanguage`, `core.GetLanguage`, `core.ListLanguages`, `core.Tokenize`
- 懒加载逻辑：首次调用时自动注册所有内置语言
- 包 init 时默认注册（与 TS 一致，避免行为不一致）

**builtins.go 内容**：
- `var BUILTIN_LANGUAGES = []core.LanguageAdapter{...}` — 17 种语言
- 目前只有 JavaScript，其余 Phase 5 补充

**验收**：`go build ./language` 通过

### 任务 4.2 — Highlighter 接口 + 实现

**目标**：创建 `api.go` 和 `highlighter.go`

**输出文件**：
- `api.go` — 接口定义
- `highlighter.go` — 实现

**api.go 内容**：
```go
type Highlighter interface {
    CodeToTokens(code string, lang string) (core.TokenStream, error)
    CodeToHtml(code string, lang string, opts ...HtmlOption) (string, error)
    UpdateTheme(theme interface{}) error
}

type HighlighterOptions struct {
    Theme           interface{}
    PreStyle        string
    LineClassPrefix string
}

func NewHighlighter(opts ...HighlighterOptions) Highlighter

type HtmlOption func(*htmlOptions)
func WithPreStyle(style string) HtmlOption
func WithLineClassPrefix(prefix string) HtmlOption
```

**highlighter.go 内容**（对照 api.ts）：
1. `CodeToTokens` — 规范化语言 ID → 查缓存 → tokenize → applyThemeStyles → 存缓存 → 深拷贝返回
2. `CodeToHtml` — CodeToTokens → escapeHtml → span/div/pre 拼接
3. `UpdateTheme` — 切换主题，遍历缓存重算 styledRows
4. 内部工具函数：`parseInlineStyle`, `stringifyInlineStyle`, `applyThemeStyles`, `cloneTokenStream`, `escapeHtml`

**escapeHtml 映射**（必须精确一致）：
| 字符 | 转义 |
|------|------|
| `&` | `&amp;` |
| `<` | `&lt;` |
| `>` | `&gt;` |
| `"` | `&quot;` |
| `'` | `&#39;` |
| `` ` `` | `&#96;` |
| `$` | `&#36;` |
| `\t` | `&#9;` |

**缓存 key**：`lang + "\x00" + code`

**验收**：`go build .` 通过，Phase 4 全部完成

---

## Phase 5：其余 16 种语言适配器

> **参考 spec**：[04-language-adapters.md](../specs/04-language-adapters.md)
> **参考 TS 源码**：`ts/lib/language/{lang}/` (type.ts, rule.ts, spec.ts, engine.ts, index.ts)

每种语言 4 个文件，结构与 JavaScript 相同：
- `types.go` — TokenScope + GrammarState 常量
- `rules.go` — GRAMMAR_RULES map（预编译正则）
- `spec.go` — TokenizerSpec 实例
- `adapter.go` — LanguageAdapter 实现

### 任务 5.1 — TypeScript

**目录**：`language/typescript/`
**参考**：`ts/lib/language/typescript/`
**注意**：TypeScript 在 JavaScript 基础上增加类型相关 scope

### 任务 5.2 — Python

**目录**：`language/python/`
**参考**：`ts/lib/language/python/`
**注意**：Python 有缩进状态、三引号字符串

### 任务 5.3 — CSS

**目录**：`language/css/`
**参考**：`ts/lib/language/css/`

### 任务 5.4 — Bash

**目录**：`language/bash/`
**参考**：`ts/lib/language/bash/`

### 任务 5.5 — SQL

**目录**：`language/sql/`
**参考**：`ts/lib/language/sql/`

### 任务 5.6 — YAML

**目录**：`language/yaml/`
**参考**：`ts/lib/language/yaml/`

### 任务 5.7 — Markdown

**目录**：`language/markdown/`
**参考**：`ts/lib/language/markdown/`

### 任务 5.8 — Java

**目录**：`language/java/`
**参考**：`ts/lib/language/java/`

### 任务 5.9 — C

**目录**：`language/c/`
**参考**：`ts/lib/language/c/`

### 任务 5.10 — C++

**目录**：`language/cpp/`
**参考**：`ts/lib/language/cpp/`

### 任务 5.11 — C#

**目录**：`language/csharp/`
**参考**：`ts/lib/language/csharp/`

### 任务 5.12 — Go

**目录**：`language/gorule/`（避免与标准库 `go` 冲突）
**参考**：`ts/lib/language/go/`
**注意**：adapter 的 ID 仍为 `"go"`，目录名是 `gorule`

### 任务 5.13 — Rust

**目录**：`language/rust/`
**参考**：`ts/lib/language/rust/`

### 任务 5.14 — PHP

**目录**：`language/php/`
**参考**：`ts/lib/language/php/`

### 任务 5.15 — HTML

**目录**：`language/html/`
**参考**：`ts/lib/language/html/`

### 任务 5.16 — JSON

**目录**：`language/json/`
**参考**：`ts/lib/language/json/`

### Phase 5 收尾

完成所有语言后：
1. 更新 `language/builtins.go`，将 17 种语言全部加入 `BUILTIN_LANGUAGES`
2. 运行 `go build ./...` 确认全部通过

**验收**：`go build ./...` 通过，17 种语言全部注册，Phase 5 全部完成

---

## Phase 6：WASM 入口

> **参考 spec**：[07-wasm-entry.md](../specs/07-wasm-entry.md)
> **参考**：Go `syscall/js` 标准库

### 任务 6.1 — WASM 入口实现

**目标**：创建 `cmd/wasm/main.go`

**输出文件**：`cmd/wasm/main.go`

**构建约束**：
```go
//go:build js && wasm
```

**内容**：
1. 全局变量：`highlighters map[int]*highlighterImpl`, `nextID int`, `defaultHighlighter Highlighter`
2. `main()` 函数：
   - 注册所有内置语言和主题
   - 创建默认 Highlighter（id=0, dark-plus 主题）
   - 挂载全局函数到 `js.Global()`
3. 导出的 JS API：
   - `coderMateHighlight(code, lang)` — 使用默认实例
   - `coderMate.createHighlighter(options)` — 返回 highlighterId
   - `coderMate.highlight(id, code, lang)` — 指定实例
   - `coderMate.updateTheme(id, theme)` — 切换主题
   - `coderMate.dispose(id)` — 释放实例
   - `coderMate.getThemes()` — 列出主题
   - `coderMate.getLanguages()` — 列出语言

**验收**：`GOOS=js GOARCH=wasm go build -o coder-mate.wasm ./cmd/wasm/` 通过

---

## Phase 7：测试与优化

> **参考 spec**：[08-testing-strategy.md](../specs/08-testing-strategy.md)
> **参考**：Go `testing` 标准库 + Table-Driven Tests

### 任务 7.1 — 核心引擎单元测试

**目标**：创建 `core/core_test.go`

**覆盖范围**：
- `NewParserContext` → 初始上下文
- `PushState` / `PopState` → 状态栈操作 + 错误边界
- `CurrentState` → 栈顶获取
- `SplitTokenByLineBreak` → `\n` / `\r\n` / 无换行 / 连续换行
- `MatchToken` → 规则匹配 / 优先级 / fallback
- `Parse` → 空输入 / 单行 / 多行 / 跨行 token

**验收**：`go test ./core/` 全部通过

### 任务 7.2 — 语言注册表测试

**目标**：创建 `language/language_test.go`

**覆盖范围**：
- `RegisterLanguage` → 注册 + 别名 + 重复注册报错
- `GetLanguage` → 查找 + 大小写不敏感 + 未注册
- `ListLanguages` → 去重
- `Tokenize` → 正常解析 + 语言未注册报错

**验收**：`go test ./language/` 全部通过

### 任务 7.3 — 公共 API 测试

**目标**：创建 `api_test.go`（根目录）

**覆盖范围**：
- `CodeToTokens` → 正常返回 + 未知语言报错 + 空 lang 报错
- `CodeToHtml` → 完整 HTML 结构 + 主题别名
- `UpdateTheme` → 主题切换后样式变化 + token text/scope 不变
- 缓存 → 命中后不重复解析

**验收**：`go test .` 全部通过

### 任务 7.4 — Go ↔ TS 一致性验证

**目标**：验证 Go 版本与 TS 版本对相同输入产生相同 token 流

**方法**：
1. 为每种语言选取 3-5 个代表性代码片段
2. Go 版本解析后输出 JSON：`[{Text, Scope, Line, Col}, ...]`
3. TS 版本用 `codeToTokens` 解析相同代码，输出同样格式
4. diff 比较 `text`, `scope`, `line`, `col` 完全一致

**验收**：所有语言的 token 流输出一致

### 任务 7.5 — 性能基准测试

**目标**：创建 `bench_test.go`（根目录）

**内容**：
- `BenchmarkParseJavaScript` — JS 解析性能
- `BenchmarkParsePython` — Python 解析性能
- `BenchmarkCodeToHtml` — 完整 HTML 生成性能

**运行**：`go test -bench=. -benchmem ./...`

**验收**：基准测试可运行，记录性能数据

---

## 快速参考

### 每轮对话启动模板

```
本轮目标：Phase X — 任务 X.N — [任务名]
参考 spec：specs/0X-xxx.md
参考 TS 源码：ts/lib/xxx/
输出文件：xxx.go
```

### 依赖关系

```
Phase 0 → Phase 1 → Phase 2 ─┐
                    Phase 3 ─┤→ Phase 4 → Phase 5 ─┐
                                              Phase 6 ─┤→ Phase 7
```

### 关键规则

1. **零外部依赖**：只用标准库
2. **正则预编译**：包级 `regexp.MustCompile`
3. **值拷贝语义**：ParserContext 操作返回新对象
4. **规则顺序**：与 TS 版本一致（影响匹配优先级）
5. **escapeHtml**：8 个字符映射必须精确一致
6. **缓存深拷贝**：防止调用方污染缓存

---

*本计划按 Phase 顺序执行，每轮对话完成一个子任务。*
---

## Phase 0：工程基建

### 任务 0.1 — 初始化 Go 模块

**目标**：创建 `go.mod`，确认项目可编译

**步骤**：
1. 在项目根目录执行 `go mod init code-mate-core`
2. 确认 `go.mod` 内容：
   ```
   module code-mate-core
   go 1.21
   ```
3. 创建空的 `core/doc.go`（包级文档，使 `go build ./...` 不报错）：
   ```go
   // Package core provides the language-agnostic tokenizer engine.
   package core
   ```
4. 运行 `go build ./...` 确认通过

**验收**：`go build ./...` 零报错，`go.mod` 存在

---