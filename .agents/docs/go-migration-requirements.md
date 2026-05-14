# code-mate-core Go 版本复刻需求文档

> **索引文件** — 每个 Phase 加载对应的 spec 即可，无需加载全部。

---

## Spec 文档索引

| # | 文档 | 内容 | 所属 Phase |
|---|------|------|:---:|
| 01 | [项目概述与技术约束](specs/01-project-overview.md) | 项目目标、质量指标、技术约束 | 全部 |
| 02 | [架构设计](specs/02-architecture.md) | TS 架构、数据流、Go 目录结构、TS→Go 映射 | 全部 |
| 03 | [核心分词引擎](specs/03-core-engine.md) | types.go、tokenizer.go（FSM 7 函数）、registry.go | Phase 1 |
| 04 | [语言适配器系统](specs/04-language-adapters.md) | manager、builtins、17 语言表、文件结构、regex 兼容 | Phase 2, 5 |
| 05 | [主题系统](specs/05-theme-system.md) | HighlightTheme、registry、variant、dark+、9 presets | Phase 3 |
| 06 | [公共 API](specs/06-public-api.md) | Highlighter 接口、CodeToHtml、escapeHtml、缓存、深拷贝 | Phase 4 |
| 07 | [WASM 入口](specs/07-wasm-entry.md) | 构建约束、JS API 导出、前端加载示例 | Phase 6 |
| 08 | [测试策略](specs/08-testing-strategy.md) | 单元测试、一致性验证、基准测试 | Phase 7 |
| 09 | [实现阶段](specs/09-implementation-phases.md) | Phase 1-7 完整文件清单 | 全部 |
| 10 | [工程基建](specs/10-engineering-setup.md) | go.mod、WASM 构建、lint、Git 规范 | Phase 0 |

---

## 快速导航

- **首次了解项目** → 01 + 02
- **开始 Phase 1 编码** → 03
- **开始 Phase 2/5 编码** → 04
- **开始 Phase 3 编码** → 05
- **开始 Phase 4 编码** → 06
- **开始 Phase 6 编码** → 07
- **开始 Phase 7 测试** → 08
- **查看完整文件清单** → 09
- **工程基建准备** → 10

---

## 实现计划

- **详细执行计划**：[2026-05-14-go-migration.md](plans/2026-05-14-go-migration.md) — 按对话拆分的逐步实现指南
- **Phase 文件清单**：[09-实现阶段](specs/09-implementation-phases.md)

---

*本文件是索引文件，详细规格见 `specs/` 目录。*
