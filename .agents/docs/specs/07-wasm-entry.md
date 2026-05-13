# 07 — WASM 入口

> **所属 Phase**：Phase 6
> **参考**：Go `syscall/js` 标准库

---

## 1. 构建约束

实现文件：`cmd/wasm/main.go`

```go
//go:build js && wasm
```

## 2. 导出的 JavaScript API

```javascript
// 全局便捷函数（使用默认 Highlighter 实例）
window.coderMateHighlight(code, lang) → { html: string } | { error: string }

// Highlighter 实例管理
window.coderMate.createHighlighter(options) → highlighterId  // int, 从 1 开始递增
window.coderMate.highlight(highlighterId, code, lang) → { html: string } | { error: string }
window.coderMate.updateTheme(highlighterId, theme) → void
window.coderMate.dispose(highlighterId) → void  // 释放 Go 端持有的实例和缓存

// 辅助查询 API（可选）
window.coderMate.getThemes() → [{id: string, displayName: string}, ...]
window.coderMate.getLanguages() → [{id: string, aliases: [string]}, ...]
```

### highlighterId 管理

- Go 端用 `map[int]*Highlighter` + `sync.Mutex` 管理实例
- `createHighlighter` 返回递增整数 ID
- `dispose(id)` 从 map 中删除实例，Go GC 回收缓存内存
- 默认实例 id=0，不需要 dispose

## 3. 初始化

- WASM 模块加载后 `main()` 自动执行
- 注册所有内置语言和主题
- 创建默认 Highlighter 实例（id=0，使用 dark-plus 主题）
- 将全局函数挂载到 `js.Global()`

## 4. WASM 编译命令

```bash
GOOS=js GOARCH=wasm go build -o coder-mate.wasm ./cmd/wasm/
```

## 5. 前端加载示例

```html
<script src="wasm_exec.js"></script>
<script>
  const go = new Go();
  WebAssembly.instantiateStreaming(fetch("coder-mate.wasm"), go.importObject)
    .then(result => {
      go.run(result.instance);
      // window.coderMateHighlight 现已可用
      const result = coderMateHighlight("const x = 1;", "javascript");
      console.log(result.html);
    });
</script>
```
