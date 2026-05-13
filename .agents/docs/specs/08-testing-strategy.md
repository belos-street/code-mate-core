# 08 — 测试策略

> **所属 Phase**：Phase 7
> **参考**：Go `testing` 标准库 + Table-Driven Tests

---

## 1. 单元测试

### core_test.go

覆盖范围：
- `NewParserContext` → 初始上下文创建
- `PushState` / `PopState` → 状态栈操作 + 错误边界
- `CurrentState` → 栈顶获取
- `SplitTokenByLineBreak` → `\n` / `\r\n` / 无换行 / 连续换行
- `MatchToken` → 规则匹配 / 优先级 / fallback
- `Parse` → 空输入 / 单行 / 多行 / 跨行 token

### language_test.go

覆盖范围：
- `RegisterLanguage` → 注册 + 别名 + 重复注册报错
- `GetLanguage` → 查找 + 大小写不敏感 + 未注册
- `ListLanguages` → 去重
- `Tokenize` → 正常解析 + 语言未注册报错

### api_test.go

覆盖范围：
- `CodeToTokens` → 正常返回 + 未知语言报错 + 空 lang 报错
- `CodeToHtml` → 完整 HTML 结构 + 主题别名
- `UpdateTheme` → 主题切换后样式变化 + token text/scope 不变
- 缓存 → 命中后不重复解析

## 2. 一致性验证

验证 Go 版本与 TS 版本对相同输入产生相同的 token 流。

### 验证方法

1. 为每种语言选取 3-5 个代表性代码片段（覆盖该语言的主要语法结构）
2. Go 版本解析后输出 JSON：`[{Text, Scope, Line, Col}, ...]`（不包含 Style）
3. TS 版本用 `codeToTokens` 解析相同代码，输出同样格式 JSON
4. 用 diff 工具逐行比较两个 JSON 输出，验证 `text`、`scope`、`line`、`col` 完全一致
5. 对 HTML 输出，验证 `<pre><code>` 结构一致，每个 token 的 text 内容正确转义

### 验证片段选取原则

- 单行代码
- 多行代码（含空行）
- 包含字符串（单/双引号、转义）
- 包含注释（单行/多行）
- 包含该语言特有语法结构
- 边界情况：空字符串、单字符、超长行

## 3. 性能基准测试

```go
func BenchmarkParseJavaScript(b *testing.B) {
    code := "const x = 1;\nfunction foo() { return x; }\nfoo();\n"
    spec := &javascript.Spec
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        core.Parse(code, spec)
    }
}

func BenchmarkParsePython(b *testing.B) {
    code := "def foo():\n    return 1\n\nx = foo()\n"
    spec := &python.Spec
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        core.Parse(code, spec)
    }
}

func BenchmarkCodeToHtml(b *testing.B) {
    hl := NewHighlighter()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        hl.CodeToHtml("const x = 1;", "javascript")
    }
}
```

运行：`go test -bench=. -benchmem ./...`
