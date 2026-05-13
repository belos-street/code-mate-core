# 05 — 主题系统

> **所属 Phase**：Phase 3
> **参考源码**：`ts/lib/themes/`

---

## 1. 主题类型

实现文件：`theme/types.go`

```go
type ThemeStyleMap map[string]string

type HighlightTheme struct {
    ID           string
    DisplayName  string
    DefaultStyle string       // scope 未匹配时的兜底样式
    PreStyle     string        // <pre> 标签的默认样式
    Styles       ThemeStyleMap // scope -> CSS style 字符串
}
```

## 2. 主题注册表

实现文件：`theme/registry.go`

- `RegisterTheme(theme *HighlightTheme) error`
- `GetTheme(id string) (*HighlightTheme, error)`
- `ListThemes() []*HighlightTheme`
- `ResolveTheme(theme interface{}) *HighlightTheme` — 接收 string 或 *HighlightTheme
- 默认主题 ID：`"dark-plus"`

### 主题别名映射（`THEME_ALIAS_MAP`）

```go
var THEME_ALIAS_MAP = map[string]string{
    "github": "github-light",
    "dark":   "dark-plus",
}
```

`ResolveTheme`、`NewHighlighter`、`UpdateTheme` 接收字符串主题时都必须先经过 alias 映射转换。

## 3. 主题变体生成器

实现文件：`theme/variant.go`

```go
func CreateThemeVariant(base *HighlightTheme, opts VariantOptions) *HighlightTheme
```

算法：
1. 遍历 base.Styles 的每个 scope → style
2. 用正则 `#[0-9A-Fa-f]{6}` 找到所有颜色值
3. 将颜色值大写后在 colorMap 中查找替换
4. 生成新的 HighlightTheme

## 4. Dark+ 基础主题

目录：`theme/darkplus/`

| 文件 | 职责 |
|------|------|
| `base.go` | 主题 ID (`dark-plus`)、显示名 (`Dark+`)、默认样式 (`color: #D4D4D4;`)、pre 样式 |
| `shared.go` | 跨语言通用样式（如 `"default": "color: #D4D4D4;"`） |
| `js.go` | JavaScript 专用样式映射 |
| `typescript.go` | TypeScript 专用样式映射 |
| `python.go` | Python 专用样式映射 |
| `css.go` | CSS 专用样式映射 |
| `bash.go` | Bash 专用样式映射 |
| `sql.go` | SQL 专用样式映射 |
| `yaml.go` | YAML 专用样式映射 |
| `markdown.go` | Markdown 专用样式映射 |
| `java.go` | Java 专用样式映射 |
| `c.go` | C 专用样式映射 |
| `cpp.go` | C++ 专用样式映射 |
| `csharp.go` | C# 专用样式映射 |
| `go.go` | Go 专用样式映射 |
| `rust.go` | Rust 专用样式映射 |
| `php.go` | PHP 专用样式映射 |
| `html.go` | HTML 专用样式映射 |
| `json.go` | JSON 专用样式映射 |

**格式**：每个文件导出 `map[string]string`，key 为完整 scope（如 `"comment.block.js"`），value 为 CSS 样式字符串。

## 5. 预设主题（9 个）

实现文件：`theme/presets.go`

通过 `CreateThemeVariant(darkPlusTheme, ...)` 生成：

| 主题 ID | 显示名 |
|---------|--------|
| github-light | GitHub Light |
| dracula | Dracula |
| one-dark-pro | One Dark Pro |
| nord | Nord |
| monokai | Monokai |
| material-ocean | Material Ocean |
| tokyo-night | Tokyo Night |
| solarized-dark | Solarized Dark |
| solarized-light | Solarized Light |

每个主题定义一个调色板（12 种颜色映射），自动从 dark-plus 生成所有 scope 的样式。

**调色板结构**：
```go
type VariantPalette struct {
    foreground string
    comment    string
    keyword    string
    accent     string
    typeName   string
    callable   string
    warning    string
    number     string
    symbol     string
    variable   string
    string     string
    muted      string
}
```
