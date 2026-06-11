# dto.go 代码阅读文档

## 1. 全局总结

本文件定义了 Cohere 频道特有的数据传输对象，包括聊天请求/响应结构、Rerank 请求/响应结构以及元数据结构。

## 2. 依赖关系

| 包 | 用途 |
|---|------|
| `dto` | 通用 DTO（RerankResponseResult） |

## 3. 类型定义

### `CohereRequest` 结构体

```go
type CohereRequest struct {
    Model       string        `json:"model"`
    ChatHistory []ChatHistory `json:"chat_history"`
    Message     string        `json:"message"`
    Stream      bool          `json:"stream"`
    MaxTokens   uint          `json:"max_tokens"`
    SafetyMode  string        `json:"safety_mode,omitempty"`
}
```

Cohere Chat API 请求格式：
- `Model`: 模型标识
- `ChatHistory`: 历史对话记录
- `Message`: 当前用户消息
- `Stream`: 是否流式输出
- `MaxTokens`: 最大生成 token 数
- `SafetyMode`: 安全模式（可选）

### `ChatHistory` 结构体

```go
type ChatHistory struct {
    Role    string `json:"role"`
    Message string `json:"message"`
}
```

聊天历史条目，角色值为 `CHATBOT`、`SYSTEM`、`USER`（注意不是标准的 assistant/system/user）。

### `CohereResponse` 结构体

```go
type CohereResponse struct {
    IsFinished   bool                  `json:"is_finished"`
    EventType    string                `json:"event_type"`
    Text         string                `json:"text,omitempty"`
    FinishReason string                `json:"finish_reason,omitempty"`
    Response     *CohereResponseResult `json:"response"`
}
```

Cohere 流式响应格式：
- `IsFinished`: 是否已完成
- `EventType`: 事件类型
- `Text`: 生成的文本片段
- `FinishReason`: 停止原因
- `Response`: 完整响应结果（仅在完成时存在）

### `CohereResponseResult` 结构体

```go
type CohereResponseResult struct {
    ResponseId   string     `json:"response_id"`
    FinishReason string     `json:"finish_reason,omitempty"`
    Text         string     `json:"text"`
    Meta         CohereMeta `json:"meta"`
}
```

Cohere 非流式/完成响应格式。

### `CohereRerankRequest` 结构体

```go
type CohereRerankRequest struct {
    Documents       []any  `json:"documents"`
    Query           string `json:"query"`
    Model           string `json:"model"`
    TopN            int    `json:"top_n"`
    ReturnDocuments bool   `json:"return_documents"`
}
```

Cohere Rerank API 请求格式。

### `CohereRerankResponseResult` 结构体

```go
type CohereRerankResponseResult struct {
    Results []dto.RerankResponseResult `json:"results"`
    Meta    CohereMeta                 `json:"meta"`
}
```

Cohere Rerank 响应格式。

### `CohereMeta` 结构体

```go
type CohereMeta struct {
    BilledUnits CohereBilledUnits `json:"billed_units"`
}
```

元数据，包含计费单元信息。

### `CohereBilledUnits` 结构体

```go
type CohereBilledUnits struct {
    InputTokens  int `json:"input_tokens"`
    OutputTokens int `json:"output_tokens"`
}
```

计费 token 数量。

### `CohereTokens` 结构体（未使用）

```go
type CohereTokens struct {
    InputTokens  int `json:"input_tokens"`
    OutputTokens int `json:"output_tokens"`
}
```

定义但未在当前代码中使用（`CohereMeta` 中的 `Tokens` 字段被注释掉了）。

## 4. 函数详解

无函数定义。

## 5. 关键逻辑分析

### 角色名称差异
Cohere 使用独特的角色命名：`CHATBOT`（而非 assistant）、`SYSTEM`、`USER`。这在 `relay-cohere.go` 的转换函数中处理。

### 流式/非流式响应的双重结构
`CohereResponse` 用于流式响应，包含 `IsFinished` 标志和增量 `Text`。`CohereResponseResult` 用于非流式响应或流式完成时的最终结果。

### 计费信息来源
Cohere API 在响应的 `meta.billed_units` 中直接返回 token 用量，这比 Cloudflare 的文本估算方式更精确。

## 6. 关联文件

- `relay/channel/cohere/relay-cohere.go` - 使用这些 DTO 进行请求转换和响应解析
- `relay/channel/cohere/adaptor.go` - 调度使用
