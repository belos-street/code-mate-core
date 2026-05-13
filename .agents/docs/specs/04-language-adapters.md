# 04 — 语言适配器系统

> **所属 Phase**：Phase 2（JavaScript）+ Phase 5（其余 16 种）
> **参考源码**：`ts/lib/language/manager.ts`, `ts/lib/language/builtins.ts`

---

## 1. 管理器

实现文件：`language/manager.go`

- 封装懒加载逻辑：首次调用时自动注册所有内置语言
- 提供与 `core/registry` 相同的接口（`GetLanguage` / `ListLanguages` / `Tokenize`），但自动确保内置语言已注册
- 模块加载时默认注册内置语言（与 TS 版本一致，避免延迟注册导致行为不一致）

## 2. 内置语言列表

实现文件：`language/builtins.go`

必须实现以下 **17 种语言**的适配器：

| 序号 | 语言 | 目录名 | 别名 |
|------|------|--------|------|
| 1 | JavaScript | javascript | js |
| 2 | TypeScript | typescript | ts |
| 3 | Python | python | py |
| 4 | CSS | css | - |
| 5 | Bash | bash | sh, shell |
| 6 | SQL | sql | - |
| 7 | YAML | yaml | yml |
| 8 | Markdown | markdown | md |
| 9 | Java | java | - |
| 10 | C | c | - |
| 11 | C++ | cpp | c++ |
| 12 | C# | csharp | cs |
| 13 | Go | gorule | go |
| 14 | Rust | rust | rs |
| 15 | PHP | php | - |
| 16 | HTML | html | - |
| 17 | JSON | json | - |

## 3. 每种语言的文件结构

每种语言包含 4 个文件：

| 文件 | 职责 |
|------|------|
| `types.go` | TokenScope 常量 + GrammarState 常量 |
| `rules.go` | GRAMMAR_RULES map（状态 → 正则规则列表） |
| `spec.go` | TokenizerSpec 实例（InitialState + Rules + FallbackScope） |
| `adapter.go` | 实现 LanguageAdapter 接口 |

### 关键约束

- 所有正则必须在包级变量中预编译（`regexp.MustCompile`）
- 规则顺序必须与 TS 版本完全一致（顺序影响匹配优先级）
- Go 目录名 `gorule` 避免与标准库 `go` 冲突

### adapter.go 模板

```go
package javascript

import "code-mate-core/core"

type adapter struct{}

func (a adapter) ID() string              { return "javascript" }
func (a adapter) Aliases() []string        { return []string{"js"} }
func (a adapter) Parse(code string) core.TokenStream {
    return core.Parse(code, &Spec)
}

var Language core.LanguageAdapter = adapter{}
```

## 4. 正则表达式兼容性

Go 的 `regexp` 包基于 RE2 语法，本项目的 TS 正则**不使用**以下不兼容特性：
- `(?<=...)` 后行断言
- `(?=...)` 前行断言
- `\1` 反向引用

因此所有正则可以直接从 TS 移植到 Go。

## 5. JavaScript 适配器参考

作为 Phase 2 的首个实现，JavaScript 的规则需要完全从 [ts/lib/language/javascript/rule.ts](../../../ts/lib/language/javascript/rule.ts) 翻译。规则包含以下状态：

| 状态 | 用途 |
|------|------|
| `global` | 顶层代码解析 |
| `multiline-comment` | `/* ... */` |
| `string-double` | `"..."` |
| `string-single` | `'...'` |
| `string-backtick` | `` `...` ``（含 `${}` 插值） |
| `template-interpolation` | `${...}` 内部代码 |
| `import-dynamic` | `import(...)` |
