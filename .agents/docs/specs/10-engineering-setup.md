# 10 — 工程基建

> **所属 Phase**：Phase 0（编码前准备）
> **参考**：`golang-best-practices` skill

---

## 1. Go Modules 初始化

```bash
cd code-mate-core
go mod init code-mate-core
go mod tidy
```

- Module path：`code-mate-core`
- Go 版本：`go 1.21`（go.mod 中声明）
- 零外部依赖，`go.sum` 为空或不存在

### go.mod 示例

```
module code-mate-core

go 1.21
```

## 2. 项目目录创建

按 Phase 顺序创建目录，不要一次性创建全部：

```bash
# Phase 1 — 核心引擎
mkdir -p core

# Phase 2 — 首个语言适配器
mkdir -p language/javascript

# Phase 3 — 主题系统
mkdir -p theme/darkplus

# Phase 4 — 公共 API（api.go 和 highlighter.go 在根目录）
mkdir -p language

# Phase 5 — 其余语言适配器（按需创建）
mkdir -p language/{typescript,python,css,bash,sql,yaml,markdown,java,c,cpp,csharp,gorule,rust,php,html,json}

# Phase 6 — WASM 入口
mkdir -p cmd/wasm

# 测试文件放在对应包目录下（Go 惯例）
# core/core_test.go
# language/language_test.go
# 根目录/api_test.go, bench_test.go
```

## 3. WASM 构建

### 构建命令

```bash
# 标准构建
GOOS=js GOARCH=wasm go build -o coder-mate.wasm ./cmd/wasm/

# 优化构建（减小体积，推荐生产环境）
GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o coder-mate.wasm ./cmd/wasm/
```

- `-s`：去掉符号表
- `-w`：去掉 DWARF 调试信息
- 预期产出：3-5MB（gzip 后 < 1MB）

### wasm_exec.js

Go 官方提供的 JS runtime glue，每个 Go 版本对应固定版本：

```bash
# 从 Go 安装目录复制
cp $(go env GOROOT)/misc/wasm/wasm_exec.js ./web/

# 前端加载
# <script src="wasm_exec.js"></script>
```

### 构建约束

`cmd/wasm/main.go` 必须包含构建标签：

```go
//go:build js && wasm

package main
```

## 4. 代码规范

### 格式化

使用 `gofmt`（Go 标准），无需额外配置：

```bash
gofmt -w .          # 格式化所有 .go 文件
gofmt -s -w .       # 格式化 + 简化（推荐）
```

### Lint（可选）

项目零依赖原则，lint 工具作为开发辅助，不纳入构建流程：

```bash
# 如果安装了 golangci-lint
golangci-lint run ./...

# 核心检查项：
# - gofmt 格式
# - go vet 静态分析
# - unused 变量/函数
# - ineffassign 无效赋值
```

### 命名规范

- 包名：小写单词，无下划线（`core`, `language`, `theme`）
- 导出函数/类型：PascalCase（`Parse`, `TokenStream`, `HighlightTheme`）
- 非导出：camelCase（`matchToken`, `parseInlineStyle`）
- 接口：动词或名词（`LanguageAdapter`, `Highlighter`）
- 常量：PascalCase（`FallbackScope`, `InitialState`）

## 5. Git 规范

### 分支策略

- `main`：稳定分支，可编译可测试
- `feat/phase-N`：每个 Phase 一个功能分支
- 合并后删除功能分支

### Commit Message

```
<type>(<scope>): <description>

# 示例
feat(core): implement FSM tokenizer engine
feat(language): add JavaScript language adapter
feat(theme): add Dark+ base theme
test(core): add tokenizer unit tests
```

Type：`feat`, `fix`, `test`, `refactor`, `docs`, `chore`

### 每个 Phase 完成后

1. 确保 `go build ./...` 通过
2. 确保 `go test ./...` 通过
3. 合并到 `main`
4. 打 tag：`v0.1.0-phase1`, `v0.2.0-phase2`, ...

## 6. 开发环境检查清单

开始编码前确认：

- [ ] Go 1.21+ 已安装（`go version`）
- [ ] `go mod init code-mate-core` 已执行
- [ ] `gofmt` 可用（`gofmt -h`）
- [ ] `go build ./...` 通过（空项目）
- [ ] `go test ./...` 通过（空项目）

---

*本文件是工程基建规范，Phase 0 完成后进入 Phase 1 编码。*
