# 02 — 架构设计

> **所属 Phase**：全部（参考用）
> **依赖**：无

---

## 1. TS 源码架构参考

```
ts/lib/
├── core/                    # 语言无关的核心引擎
│   ├── types.ts             # Token, GrammarRule, ParserContext, TokenizerSpec
│   ├── tokenizer.ts         # 通用 FSM 分词引擎（状态机 + 规则匹配 + 跨行拆分）
│   └── registry.ts          # 语言注册表（Map<id, LanguageAdapter>）
│
├── language/                # 17 种语言适配器
│   ├── builtins.ts          # 内置语言数组
│   ├── manager.ts           # 懒加载注册管理器
│   ├── javascript/          # 每种语言：type.ts, rule.ts, spec.ts, engine.ts, index.ts
│   ├── typescript/
│   ├── python/ ... (共17种)
│
├── themes/                  # 主题系统
│   ├── types.ts             # HighlightTheme, ThemeStyleMap
│   ├── theme-variant.ts     # 主题变体生成器（颜色替换）
│   ├── presets.ts           # 9 个预设主题
│   ├── index.ts             # 主题注册表 + resolveTheme
│   └── dark-plus/           # 基础主题（每种语言单独样式文件）
│
├── api.ts                   # 公共 API（createHighlighter, codeToTokens, codeToHtml）
└── index.ts                 # 顶层导出
```

## 2. 核心数据流

```
输入代码 string
    │
    ▼
┌─────────────────────────────┐
│  tokenize(code, langId)     │  ← language/manager.ts
│  1. 从 registry 查找语言    │
│  2. 调用 language.parse()   │
│  3. parse() 内部调用         │
│     core/tokenizer.parse()  │
│     遍历正则规则，FSM 状态机 │
└─────────────┬───────────────┘
              │
              ▼
        TokenStream ([][]Token)
        每个 Token: { text, scope, line, col, style:{} }
              │
              ▼
┌─────────────────────────────┐
│  applyThemeStyles(rows, theme) │  ← api.ts
│  对每个 token 的 scope       │
│  查找对应的 CSS style 字符串 │
│  写入 token.style            │
└─────────────┬───────────────┘
              │
              ▼
┌─────────────────────────────┐
│  codeToHtml()               │
│  1. escapeHtml(token.text)  │
│  2. 拼接 <span style="..."> │
│  3. 每行包裹 <div>          │
│  4. 外层 <pre><code>        │
└─────────────┬───────────────┘
              │
              ▼
        HTML string
```

## 3. Go 项目目录结构

```
code-mate-core/
├── go.mod
├── go.sum
├── cmd/
│   └── wasm/
│       └── main.go              # WASM 入口
│
├── core/
│   ├── types.go                 # 核心类型定义
│   ├── tokenizer.go             # FSM 分词引擎
│   └── registry.go              # 语言注册表
│
├── language/
│   ├── manager.go               # 懒加载管理器
│   ├── builtins.go              # 内置语言注册
│   ├── javascript/              # 每种语言：types.go, rules.go, spec.go, adapter.go
│   ├── typescript/
│   ├── python/
│   ├── css/
│   ├── bash/
│   ├── sql/
│   ├── yaml/
│   ├── markdown/
│   ├── java/
│   ├── c/
│   ├── cpp/
│   ├── csharp/
│   ├── gorule/                  # Go 语言（避免与标准库冲突）
│   ├── rust/
│   ├── php/
│   ├── html/
│   └── json/
│
├── theme/
│   ├── types.go                 # 主题类型
│   ├── variant.go               # 主题变体生成器
│   ├── registry.go              # 主题注册表
│   ├── presets.go               # 9 个预设主题
│   └── darkplus/                # Dark+ 基础主题
│       ├── base.go
│       ├── shared.go
│       ├── js.go
│       ├── typescript.go
│       └── ...（每种语言一个文件）
│
├── api.go                       # Highlighter 接口定义
├── highlighter.go               # Highlighter 实现
│
├── core_test.go                 # 核心引擎测试
├── language_test.go             # 语言注册表测试
├── api_test.go                  # 公共 API 测试
└── bench_test.go                # 性能基准测试
```

## 4. TS → Go 映射速查表

| TS 文件 | Go 文件 | 说明 |
|---------|---------|------|
| `lib/core/types.ts` | `core/types.go` | 类型定义，去掉泛型参数 |
| `lib/core/tokenizer.ts` | `core/tokenizer.go` | 核心算法，不可变 → 值拷贝 |
| `lib/core/registry.ts` | `core/registry.go` | 注册表，Map → map |
| `lib/language/manager.ts` | `language/manager.go` | 懒加载管理 |
| `lib/language/builtins.ts` | `language/builtins.go` | 内置语言列表 |
| `lib/language/{lang}/type.ts` | `language/{lang}/types.go` | 常量定义 |
| `lib/language/{lang}/rule.ts` | `language/{lang}/rules.go` | 正则规则 |
| `lib/language/{lang}/spec.ts` | `language/{lang}/spec.go` | TokenizerSpec |
| `lib/language/{lang}/engine.ts` | `language/{lang}/adapter.go` | 适配器实现 |
| `lib/themes/types.ts` | `theme/types.go` | 主题类型 |
| `lib/themes/theme-variant.ts` | `theme/variant.go` | 变体生成器 |
| `lib/themes/index.ts` | `theme/registry.go` | 主题注册表 |
| `lib/themes/presets.ts` | `theme/presets.go` | 预设主题 |
| `lib/themes/dark-plus/*.ts` | `theme/darkplus/*.go` | Dark+ 主题样式 |
| `lib/api.ts` | `api.go` + `highlighter.go` | 公共 API |
| N/A | `cmd/wasm/main.go` | WASM 入口（新增） |
