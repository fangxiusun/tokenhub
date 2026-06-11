# finish_reason.go 代码阅读文档

## 1. 全局概述

本文件定义了 AI 模型响应的完成原因（Finish Reason）常量，这些常量与 OpenAI API 规范兼容，用于标识模型响应结束的原因。

## 2. 依赖关系

本文件无外部依赖。

## 3. 类型定义

本文件无自定义类型定义，常量为字符串类型。

## 4. 函数详情

本文件无函数定义。

## 5. 关键逻辑分析

### Finish Reason 常量

| 变量名 | 值 | 说明 |
|--------|-----|------|
| `FinishReasonStop` | `"stop"` | 模型正常完成响应 |
| `FinishReasonToolCalls` | `"tool_calls"` | 模型请求调用工具/函数 |
| `FinishReasonLength` | `"length"` | 达到 `max_tokens` 限制 |
| `FinishReasonFunctionCall` | `"function_call"` | 模型请求调用函数（旧版） |
| `FinishReasonContentFilter` | `"content_filter"` | 响应被内容过滤器拦截 |

### 使用场景

- 在流式传输中，每个 chunk 的 `finish_reason` 字段使用这些常量
- 在非流式响应中，最后一个 choice 的 `finish_reason` 字段使用这些常量
- 系统根据 `finish_reason` 决定是否继续处理（如 `tool_calls` 需要执行工具调用）

### 与 OpenAI 规范的兼容性

这些常量值完全遵循 OpenAI Chat Completions API 的 `finish_reason` 字段规范，确保与上游和下游系统的兼容性。

## 6. 相关文件

- `relay/` — 中继层处理流式和非流式响应时使用这些常量
- `relay/channel/openai/` — OpenAI 适配器直接使用这些常量
- `dto/` — 响应 DTO 中引用这些常量
