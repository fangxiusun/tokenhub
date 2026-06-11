# adaptor.go 代码阅读文档

## 1. 全局总结
本文件是 Mistral 渠道的适配器入口，实现了 `channel.Adaptor` 接口。Mistral 是欧洲的 AI 公司，提供多种规模的开源和闭源语言模型。该适配器主要支持文本生成，使用标准 Bearer Token 认证，请求格式需要特殊的 tool call ID 格式化处理。

## 2. 依赖关系
- **标准库**: `errors`, `io`, `net/http`
- **项目内部**:
  - `github.com/QuantumNous/new-api/dto` — 数据传输对象
  - `github.com/QuantumNous/new-api/relay/channel` — 通用渠道工具
  - `github.com/QuantumNous/new-api/relay/channel/openai` — OpenAI 处理器
  - `github.com/QuantumNous/new-api/relay/common` — RelayInfo
  - `github.com/QuantumNous/new-api/types` — 错误类型
- **外部依赖**: `github.com/gin-gonic/gin`

## 3. 类型定义

### `Adaptor` 结构体
```go
type Adaptor struct {}
```
Mistral 渠道适配器。

## 4. 函数详解

### `ConvertGeminiRequest` / `ConvertClaudeRequest` / `ConvertAudioRequest` / `ConvertImageRequest` / `ConvertEmbeddingRequest` / `ConvertOpenAIResponsesRequest`
均返回 `"not implemented"` 错误或 panic。Mistral 渠道不支持这些格式转换。

### `ConvertClaudeRequest`
注意：使用 `panic("implement me")` 而非返回错误。

### `Init(info)`
空实现。

### `GetRequestURL(info) (string, error)`
使用 `relaycommon.GetFullRequestURL` 构建完整请求 URL，直接拼接基础 URL 和请求路径。

### `SetupRequestHeader(c, req, info) error`
设置标准 API 请求头，并添加 `Authorization: Bearer {apiKey}`。

### `ConvertOpenAIRequest(c, info, request) (any, error)`
调用 `requestOpenAI2Mistral` 对 OpenAI 请求进行格式化处理（主要是 tool call ID 格式化）。

### `ConvertRerankRequest`
返回 nil（不支持 rerank）。

### `DoRequest(c, info, requestBody) (any, error)`
调用 `channel.DoApiRequest` 发送请求。

### `DoResponse(c, resp, info) (usage, err)`
根据流式标志选择处理器：
- 流式 → `openai.OaiStreamHandler`
- 非流式 → `openai.OpenaiHandler`

### `GetModelList() []string`
返回 `ModelList`。

### `GetChannelName() string`
返回 `ChannelName`（`"mistral"`）。

## 5. 关键逻辑分析

### 请求 URL 构建
与大多数渠道不同，Mistral 使用 `GetFullRequestURL` 直接拼接 URL，而非手动构建端点路径。这说明 Mistral 使用标准的 OpenAI 兼容 API 路径。

### Tool Call ID 格式化
`ConvertOpenAIRequest` 调用 `requestOpenAI2Mistral` 进行格式化，Mistral 要求 tool call ID 为 9 位字母数字字符串，这与 OpenAI 的 `call_xxx` 格式不同。

### 简单的响应处理
Mistral 的响应处理直接委托给 OpenAI 处理器，说明 Mistral 的 API 响应格式与 OpenAI 完全兼容。

## 6. 关联文件
- `text.go` — `requestOpenAI2Mistral` 函数实现
- `constants.go` — 模型列表和渠道名称
