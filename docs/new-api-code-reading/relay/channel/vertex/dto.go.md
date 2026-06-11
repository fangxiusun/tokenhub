# dto.go 代码阅读文档

## 1. 全局总结
Vertex AI 渠道的数据传输对象定义文件。定义了 Vertex AI Claude 请求结构体和请求复制函数。

## 2. 依赖关系
- **标准库**: encoding/json
- **内部包**: `github.com/QuantumNous/new-api/dto`

## 3. 类型定义

### VertexAIClaudeRequest
```go
type VertexAIClaudeRequest struct {
    AnthropicVersion string              `json:"anthropic_version"`
    Messages         []dto.ClaudeMessage `json:"messages"`
    System           any                 `json:"system,omitempty"`
    MaxTokens        *uint               `json:"max_tokens,omitempty"`
    StopSequences    []string            `json:"stop_sequences,omitempty"`
    Stream           *bool               `json:"stream,omitempty"`
    Temperature      *float64            `json:"temperature,omitempty"`
    TopP             *float64            `json:"top_p,omitempty"`
    TopK             *int                `json:"top_k,omitempty"`
    Tools            any                 `json:"tools,omitempty"`
    ToolChoice       any                 `json:"tool_choice,omitempty"`
    Thinking         *dto.Thinking       `json:"thinking,omitempty"`
    OutputConfig     json.RawMessage     `json:"output_config,omitempty"`
}
```
Vertex AI 上 Claude 请求的特化结构体。与标准 `dto.ClaudeRequest` 的区别：
- 增加 `AnthropicVersion` 字段（必需）
- `OutputConfig` 使用 `json.RawMessage` 类型以支持灵活的 JSON 结构

## 4. 函数详解

### copyRequest(req *dto.ClaudeRequest, version string) *VertexAIClaudeRequest
将标准 Claude 请求复制为 Vertex AI Claude 请求格式，设置 `AnthropicVersion` 为指定版本。

## 5. 关键逻辑分析
- Vertex AI 的 Claude API 与原生 Claude API 的主要区别是需要额外的 `anthropic_version` 字段
- 使用 `any` 类型的 `System` 和 `Tools` 字段以支持灵活的 JSON 结构
- `OutputConfig` 使用 `json.RawMessage` 以透传任意 JSON 配置

## 6. 关联文件
- `vertex/adaptor.go` — 在 `ConvertClaudeRequest` 中使用 `copyRequest`
- `relay/channel/dto/claude.go` — 标准 Claude 请求结构体
