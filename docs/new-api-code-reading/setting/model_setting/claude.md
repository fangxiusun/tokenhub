# claude.go 代码阅读文档

## 1. 全局总结

该文件定义 Claude 模型的配置，包括自定义 HTTP 头、默认最大 token 数、Thinking 适配器等配置项。

## 2. 依赖关系

- `net/http` — HTTP 头操作
- `strings` — 字符串处理
- `github.com/QuantumNous/new-api/setting/config` — 全局配置管理器

## 3. 类型定义

| 结构体 | 字段 | 类型 | 说明 |
|--------|------|------|------|
| `ClaudeSettings` | `HeadersSettings` | `map[string]map[string][]string` | 按模型的自定义 HTTP 头 |
| | `DefaultMaxTokens` | `map[string]int` | 按模型的默认最大 token 数 |
| | `ThinkingAdapterEnabled` | `bool` | 是否启用 Thinking 适配器 |
| | `ThinkingAdapterBudgetTokensPercentage` | `float64` | Thinking 预算 token 百分比 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `GetClaudeSettings` | `func GetClaudeSettings() *ClaudeSettings` | 获取 Claude 配置 |
| `WriteHeaders` | `func (c *ClaudeSettings) WriteHeaders(originModel string, httpHeader *http.Header)` | 将配置的头写入 HTTP 请求 |
| `normalizeHeaderListValues` | `func normalizeHeaderListValues(values []string) []string` | 规范化并去重头值 |
| `GetDefaultMaxTokens` | `func (c *ClaudeSettings) GetDefaultMaxTokens(model string) int` | 获取模型的默认最大 token 数 |

## 5. 关键逻辑分析

- `WriteHeaders` 合并配置的头和现有头，使用逗号分隔
- `normalizeHeaderListValues` 处理逗号分隔的值并去重
- 默认最大 token 数为 8192，可通过 "default" 键配置

## 6. 关联文件

- `relay/claude/` — Claude 中继适配器
- `setting/model_setting/` — 其他模型配置
