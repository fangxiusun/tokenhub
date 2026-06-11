# adaptor.go 代码阅读文档

## 1. 全局总结
本文件是 Jina 渠道的适配器入口，实现了 `channel.Adaptor` 接口。Jina 是一个专注于嵌入（Embedding）和重排序（Rerank）的 AI 服务提供商。该适配器支持 rerank 和 embeddings 两种中继模式，使用标准 Bearer Token 认证。

## 2. 依赖关系
- **标准库**: `errors`, `fmt`, `io`, `net/http`
- **项目内部**:
  - `github.com/QuantumNous/new-api/dto` — 数据传输对象
  - `github.com/QuantumNous/new-api/relay/channel` — 通用渠道工具
  - `github.com/QuantumNous/new-api/relay/channel/openai` — OpenAI 处理器
  - `github.com/QuantumNous/new-api/relay/common` — RelayInfo
  - `github.com/QuantumNous/new-api/relay/common_handler` — 通用 rerank 处理器
  - `github.com/QuantumNous/new-api/relay/constant` — 中继模式常量
  - `github.com/QuantumNous/new-api/types` — 错误类型
- **外部依赖**: `github.com/gin-gonic/gin`

## 3. 类型定义

### `Adaptor` 结构体
```go
type Adaptor struct {}
```
Jina 渠道适配器。

## 4. 函数详解

### `ConvertGeminiRequest` / `ConvertClaudeRequest` / `ConvertAudioRequest` / `ConvertImageRequest` / `ConvertOpenAIResponsesRequest`
均返回 `"not implemented"` 错误或 panic。Jina 渠道不支持这些格式转换。

### `ConvertClaudeRequest`
注意：此函数使用 `panic("implement me")` 而非返回错误，这是一个潜在的代码质量问题。

### `Init(info)`
空实现。

### `GetRequestURL(info) (string, error)`
根据中继模式构建请求 URL：
- `RelayModeRerank` → `{baseUrl}/v1/rerank`
- `RelayModeEmbeddings` → `{baseUrl}/v1/embeddings`
- 其他模式返回错误

### `SetupRequestHeader(c, req, info) error`
设置标准 API 请求头，并添加 `Authorization: Bearer {apiKey}`。

### `ConvertOpenAIRequest(c, info, request) (any, error)`
直接透传 OpenAI 请求。

### `ConvertRerankRequest(c, relayMode, request) (any, error)`
直接透传 Rerank 请求。

### `ConvertEmbeddingRequest(c, info, request) (any, error)`
透传嵌入请求，但清空 `EncodingFormat` 字段（Jina 不使用此参数）。

### `DoRequest(c, info, requestBody) (any, error)`
调用 `channel.DoApiRequest` 发送请求。

### `DoResponse(c, resp, info) (usage, err)`
响应处理：
- `RelayModeRerank` → `common_handler.RerankHandler`
- `RelayModeEmbeddings` → `openai.OpenaiHandler`

### `GetModelList() []string`
返回 `ModelList`。

### `GetChannelName() string`
返回 `ChannelName`（`"jina"`）。

## 5. 关键逻辑分析

### 功能定位
Jina 渠道专注于两个 AI 能力：嵌入（Embedding）和重排序（Rerank），不支持文本生成、图像生成等功能。这反映了 Jina 作为搜索和信息检索基础设施提供商的定位。

### 复用通用处理器
Rerank 响应直接使用 `common_handler.RerankHandler`，嵌入使用 `openai.OpenaiHandler`，说明 Jina 的 API 响应格式与 OpenAI 兼容。

### 代码质量注意
`ConvertClaudeRequest` 使用 `panic` 而非返回错误，与其他函数的处理方式不一致。

## 6. 关联文件
- `relay-jina.go` — 空文件（仅包声明）
- `constant.go` — 模型列表和渠道名称
- `relay/common_handler/` — 通用 rerank 处理器
