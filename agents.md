# code-mate-core 项目指南

> 本项目是将 TypeScript 版本的语法高亮引擎（coder-mate-core）复刻为 Go + WASM 实现。
> 让 LLM/Agent 快速理解项目结构和规范，高效协作开发。

---

## 项目概述

**目标**：将 `ts/lib/` 下的语法高亮引擎完整复刻到 Go，编译为 WASM 供前端调用。

**核心约束**：
- 输入/输出行为与 TS 版本一致（相同代码 → 相同 HTML/Token 流）
- Go 1.21+，零外部依赖，仅用标准库
- 目录结构遵循 Go 标准项目布局

**当前状态**：需求规格文档（specs/01-10）已编写完成，待进入实现阶段。

---

## 目录结构说明

```
code-mate-core/
├── .agents/                  # Agent 工作区（LLM 协作所需的规范和文档）
│   ├── skills/               # 编码规范和领域知识（Agent 开发时参考）
│   │   ├── belos-street/     # 个人编码风格和最佳实践
│   │   ├── brainstorming/    # 需求分析和设计讨论流程
│   │   ├── golang-best-practices/  # Go 开发规范（WASM/concurrency/testing）
│   │   ├── vibe-flow/        # Web/SaaS 全流程开发
│   │   └── writing-plans/    # 实现计划编写规范
│   └── docs/                 # 项目文档（需求、设计、计划）
│       ├── go-migration-requirements.md  # 需求索引文件（入口）
│       ├── specs/            # 详细规格文档（01-09，按 Phase 组织）
│       └── plans/            # 实现阶段计划
├── ts/                       # TypeScript 源码（参考实现，不要修改）
│   └── lib/
│       ├── core/             # 语言无关核心引擎（types, tokenizer, registry）
│       ├── language/         # 17 种语言适配器
│       └── themes/           # 主题系统（dark-plus + 9 预设）
└── agents.md                 # 本文件 - 项目入口指南
```

## 目录角色

### `.agents/skills/` — Agent 技能与规范

当 Agent 执行编码任务时，应参考此目录下的规范：

| 技能 | 用途 | 触发条件 |
|------|------|----------|
| `golang-best-practices` | Go 开发全规范 | 编写任何 `.go` 文件前必须加载 |
| `belos-street` | 个人编码风格 | 需要遵循特定代码风格时加载 |
| `brainstorming` | 需求分析与设计 | 讨论新功能/设计决策时使用 |
| `writing-plans` | 实现计划编写 | 编写多步骤实现计划时使用 |

**特别说明**：`golang-best-practices/references/` 下包含 37 个 Go 专题参考（WASM、并发、测试、内存管理等），开发时按需查阅。

### `.agents/docs/` — 项目文档

存放需求文档、设计文档、实现计划等。当前阶段的核心文档：

| 文档 | 说明 |
|------|------|
| [go-migration-requirements.md](.agents/docs/go-migration-requirements.md) | 需求索引文件（入口） |
| [specs/](.agents/docs/specs/) | 详细规格文档（01-10），按 Phase 组织 |
| [plans/](.agents/docs/plans/) | 实现阶段计划 |

### `ts/` — TypeScript 参考实现

**不要修改此目录**。它作为 Go 版本的参考实现和一致性验证基准。

核心入口文件：
- `ts/lib/core/types.ts` — 所有核心类型定义
- `ts/lib/core/tokenizer.ts` — FSM 分词引擎（核心算法）
- `ts/lib/core/registry.ts` — 语言注册表
- `ts/lib/api.ts` — 公共 API（Highlighter）
- `ts/lib/language/manager.ts` — 语言管理器

---

## 开发工作流

1. **理解需求** → 阅读 `.agents/docs/go-migration-requirements.md`
2. **撰写计划** → 使用 `writing-plans` skill 创建详细的实现计划
3. **编码实现** → 参考 `golang-best-practices` skill，对照 `ts/` 源码
4. **验证一致性** → 与 TS 版本对比 token 流和 HTML 输出

---

## 快速上手

**首次进入项目的 Agent 应执行**：
1. 阅读本文件（agents.md）了解项目全貌
2. 阅读 [需求索引](.agents/docs/go-migration-requirements.md)，按当前 Phase 加载对应 spec
3. 浏览 `ts/lib/core/` 理解核心算法
4. 加载 `golang-best-practices` skill 获取 Go 开发规范

---

*本文件是 LLM/Agent 协作的入口指南，保持简洁准确。*
