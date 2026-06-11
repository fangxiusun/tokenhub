# text.go 代码阅读文档

## 1. 全局总结
本文件实现了 OpenAI 请求到 Mistral 格式的转换逻辑，主要处理 tool call ID 的格式化。Mistral 要求 tool call ID 为 9 位字母数字字符串，而 OpenAI 使用 `call_xxx` 格式，因此需要进行 ID 转换。

## 2. 依赖关系
- **标准库**: `regexp`
- **项目内部**:
  - `github.com/QuantumNous/new-api/common` — 随机字符生成
  - `github.com/QuantumNous/new-api/dto` — 数据传输对象

## 3. 类型定义

### `mistralToolCallIdRegexp` 变量
```go
var mistralToolCallIdRegexp = regexp.MustCompile("^[a-zA-Z0-9]{9}$")
```
Mistral tool call ID 的正则表达式：恰好 9 位字母数字字符。

## 4. 函数详解

### `requestOpenAI2Mistral(request) *dto.GeneralOpenAIRequest`
将 OpenAI 请求转换为 Mistral 格式：

1. **Tool Call ID 格式化**:
   - 遍历所有消息中的 tool calls
   - 如果 ID 不符合 Mistral 格式（9 位字母数字），生成新的随机 ID
   - 使用 `idMap` 映射确保同一 ID 在整个请求中保持一致
   - 同样处理 `tool_call_id` 字段

2. **Content 处理**:
   - 对 assistant 消息，如果只有 tool_calls 没有 content，清空 mediaMessages
   - 将 ImageURL 类型的内容转换为 Mistral 的 `ImageUrl` 字段格式

3. **构建新请求**:
   - 只保留必要字段：model, stream, messages, temperature, topP, tools, toolChoice, maxTokens
   - 处理 `MaxTokens` 和 `MaxCompletionTokens` 的合并逻辑

## 5. 关键逻辑分析

### Tool Call ID 转换
Mistral 对 tool call ID 有严格格式要求（9 位字母数字），而 OpenAI 使用 `call_` 前缀的 UUID 格式。转换策略：
1. 检查 ID 是否已符合格式
2. 不符合则生成新的 9 位随机 ID
3. 使用映射表确保一致性（同一 tool call 在不同消息中使用相同 ID）

### ID 一致性
`idMap` 映射确保了：
- assistant 消息中的 `tool_calls[i].ID`
- 后续 tool 消息中的 `tool_call_id`
- 使用相同的转换后 ID

### Content 清理
当 assistant 消息只有 tool_calls 没有文本内容时，清空 `mediaMessages` 避免发送空内容。

### 请求精简
Mistral 不支持的字段（如 `n`, `frequency_penalty`, `presence_penalty` 等）在转换时被丢弃，只保留 Mistral 支持的字段子集。

## 6. 关联文件
- `adaptor.go` — 在 `ConvertOpenAIRequest` 中调用 `requestOpenAI2Mistral`
