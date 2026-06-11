# adaptor.go 代码阅读文档

## 1. 全局总结

本文件实现了 Cloudflare Workers AI 频道的适配器（Adaptor），是该频道与 relay 框架之间的桥梁。Adaptor 实现了 `channel.Adaptor` 接口，负责请求 URL 构建、请求头设置、请求格式转换以及响应处理。支持 Chat Completions、Embeddings、Responses、Audio（STT）和 Completions（旧版）等多种 relay 模式。

## 2. 依赖关系

| 包 | 用途 |
|---|------|
| `bytes` | 构建内存缓冲区，用于音频文件上传 |
| `errors` | 错误创建 |
| `fmt` | 字符串格式化（URL 拼接） |
| `io` | I/O 操作（Reader 接口） |
| `net/http` | HTTP 请求/响应处理 |
| `dto` | 数据传输对象（OpenAI 请求/响应格式） |
| `relay/channel` | 通用频道工具函数（DoApiRequest、SetupApiRequestHeader） |
| `relay/channel/openai` | OpenAI 频道的 Responses 处理器复用 |
| `relay/common` | RelayInfo 等中继上下文信息 |
| `relay/constant` | 中继模式常量（RelayModeChatCompletions 等） |
| `types` | 错误类型（NewAPIError） |
| `gin` | Web 框架上下文 |

## 3. 类型定义

### `Adaptor` 结构体

```go
type Adaptor struct{}
```

无状态的适配器结构体，所有方法通过值接收器调用。实现了 `channel.Adaptor` 接口的全部方法。

## 4. 函数详解

### `(a *Adaptor) Init(info *relaycommon.RelayInfo)`
初始化适配器，当前为空实现（无需初始化状态）。

### `(a *Adaptor) GetRequestURL(info *relaycommon.RelayInfo) (string, error)`
根据 relay 模式构建 Cloudflare Workers AI 的 API URL：
- **ChatCompletions**: `{base}/client/v4/accounts/{accountId}/ai/v1/chat/completions`
- **Embeddings**: `{base}/client/v4/accounts/{accountId}/ai/v1/embeddings`
- **Responses**: `{base}/client/v4/accounts/{accountId}/ai/v1/responses`
- **其他**: `{base}/client/v4/accounts/{accountId}/ai/run/{modelName}`（通用推理端点）

### `(a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error`
设置请求头：调用通用工具设置基础头，然后添加 `Authorization: Bearer {apiKey}`。

### `(a *Adaptor) ConvertOpenAIRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error)`
将 OpenAI 格式请求转换为 Cloudflare 格式：
- **Completions 模式**: 调用 `convertCf2CompletionsRequest` 转换为 `CfRequest`
- **其他模式**: 直接透传 `GeneralOpenAIRequest`

### `(a *Adaptor) ConvertOpenAIResponsesRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.OpenAIResponsesRequest) (any, error)`
Responses 格式请求直接透传，不做转换。

### `(a *Adaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (any, error)`
委托给 `channel.DoApiRequest` 执行实际的 HTTP 请求。

### `(a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage any, err *types.NewAPIError)`
响应处理的核心分发逻辑，根据 relay 模式和流式标志选择处理器：
- **Embeddings / ChatCompletions**: 流式调用 `cfStreamHandler`，非流式调用 `cfHandler`
- **Responses**: 流式复用 `openai.OaiResponsesStreamHandler`，非流式复用 `openai.OaiResponsesHandler`
- **AudioTranslation / AudioTranscription**: 调用 `cfSTTHandler`

### `(a *Adaptor) ConvertAudioRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.AudioRequest) (io.Reader, error)`
从 multipart 表单中读取上传的音频文件，将其内容写入内存缓冲区并返回。

### 其他未实现的转换方法
- `ConvertGeminiRequest` - 返回 "not implemented"
- `ConvertClaudeRequest` - panic（不应被调用）
- `ConvertImageRequest` - 返回 "not implemented"
- `ConvertRerankRequest` - 直接透传
- `ConvertEmbeddingRequest` - 直接透传

### `(a *Adaptor) GetModelList() []string`
返回支持的模型列表（定义在 `constant.go`）。

### `(a *Adaptor) GetChannelName() string`
返回频道名称 `"cloudflare"`。

## 5. 关键逻辑分析

### URL 路径构建模式
Cloudflare 使用 `/client/v4/accounts/{accountId}/ai/` 作为统一前缀，不同模式对应不同子路径。默认推理端点使用 `/ai/run/{modelName}` 路径。

### 响应处理的流式/非流式分支
响应处理采用 `switch + if` 的双层分发：外层按 relay 模式分发，内层按流式标志分发。对于 Responses 模式直接复用 OpenAI 的处理器，减少了代码重复。

### 音频请求处理
音频请求采用从 multipart form 中直接提取文件内容的方式，将二进制数据包装在 `bytes.Buffer` 中作为 `io.Reader` 返回，供下游使用。

## 6. 关联文件

- `relay/cloudflare/relay_cloudflare.go` - 包含 `cfHandler`、`cfStreamHandler`、`cfSTTHandler` 等响应处理函数
- `relay/cloudflare/dto.go` - Cloudflare 特有的请求/响应 DTO 定义
- `relay/cloudflare/constant.go` - 模型列表和频道名称常量
- `relay/channel/openai/` - OpenAI 适配器，提供 Responses 模式的复用处理器
- `relay/channel/` - 通用频道工具（DoApiRequest、SetupApiRequestHeader）
