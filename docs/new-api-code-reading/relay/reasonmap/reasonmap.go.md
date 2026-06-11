# reasonmap.go 代码阅读文档

## 1. 全局总结

本文件提供了 Claude 和 OpenAI 之间停止原因（stop reason / finish reason）的双向映射函数。

## 2. 依赖关系

- `constant`: FinishReasonContentFilter 常量

## 3. 类型定义

本文件无自定义类型定义。

## 4. 函数详解

### `ClaudeStopReasonToOpenAIFinishReason(stopReason) string`
- Claude → OpenAI 映射：
  - `stop_sequence` / `end_turn` → `stop`
  - `max_tokens` → `length`
  - `tool_use` → `tool_calls`
  - `refusal` → `content_filter`

### `OpenAIFinishReasonToClaudeStopReason(finishReason) string`
- OpenAI → Claude 映射：
  - `stop` → `end_turn`
  - `stop_sequence` → `stop_sequence`
  - `length` / `max_tokens` → `max_tokens`
  - `content_filter` → `refusal`
  - `tool_calls` → `tool_use`

## 5. 关键逻辑分析

1. **双向映射**: 两个函数互为逆映射
2. **默认透传**: 未匹配的值直接返回原值
3. **大小写不敏感**: 使用 `strings.ToLower` 进行匹配

## 6. 关联文件

- `relay/channel/claude/`: Claude 适配器使用此映射
- `constant/finish_reason.go`: FinishReasonContentFilter 定义
