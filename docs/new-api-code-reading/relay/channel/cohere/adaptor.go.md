# adaptor.go 代码阅读文档

## 1. 全局总结

本文件实现了 Cohere 频道的适配器，负责请求 URL 构建、请求头设置、请求格式转换（OpenAI → Cohere）和响应处理。支持 Chat 和 Rerank 两种中继模式，将 Cohere 原生 API 格式转换为 OpenAI 兼容格式。

## 2. 依赖关系

| 包 | 用途 |
|---|------|
| `errors` | 错误创建 |
| `fmt` | 字符串格式化 |
| `io` | I/O 操作 |
| `net/http` | HTTP 请求/响应 |
| `dto` | 数据传输对象 |
| `relay/channel` | 通用频道工具 |
| `relay/common` | RelayInfo 上下文 |
| `relay/constant` | 中继模式常量 |
| `types` | 错误类型 |
| `gin` | Web 框架 |

## 3. 类型定义

### `Adaptor` 结构体

```go
type Adaptor struct{}
```

无状态适配器，实现 `channel.Adaptor` 接口。

## 4. 函数详解

### `(a *Adaptor) Init(info *relaycommon.RelayInfo)`
空实现。

### `(a *Adaptor) GetRequestURL(info *relaycommon.RelayInfo) (string, error)`
根据 relay 模式构建 Cohere API URL：
- **Rerank 模式**: `{base}/v1/rerank`
- **其他模式**: `{base}/v1/chat`

### `(a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error`
设置 `Authorization: Bearer {apiKey}` 请求头。

### `(a *Adaptor) ConvertOpenAIRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error)`
将 OpenAI 请求转换为 Cohere 格式，调用 `requestOpenAI2Cohere` 函数。

### `(a *Adaptor) ConvertRerankRequest(c *gin.Context, relayMode int, request dto.RerankRequest) (any, error)`
将 Rerank 请求转换为 Cohere 格式，调用 `requestConvertRerank2Cohere` 函数。

### `(a *Adaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (any, error)`
委托给 `channel.DoApiRequest` 执行 HTTP 请求。

### `(a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage any, err *types.NewAPIError)`
响应处理分发：
- **Rerank 模式**: 调用 `cohereRerankHandler`
- **流式**: 调用 `cohereStreamHandler`
- **非流式**: 调用 `cohereHandler`

### 未实现的方法
- `ConvertGeminiRequest` - 返回 "not implemented"
- `ConvertClaudeRequest` - panic
- `ConvertAudioRequest` - 返回 "not implemented"
- `ConvertImageRequest` - 返回 "not implemented"
- `ConvertEmbeddingRequest` - 返回 "not implemented"
- `ConvertOpenAIResponsesRequest` - 返回 "not implemented"

### `(a *Adaptor) GetModelList() []string`
返回 Cohere 模型列表。

### `(a *Adaptor) GetChannelName() string`
返回 `"cohere"`。

## 5. 关键逻辑分析

### 双模式支持
Cohere 适配器同时支持 Chat（`/v1/chat`）和 Rerank（`/v1/rerank`）两种 API 端点，通过 `GetRequestURL` 和 `DoResponse` 中的 relay 模式判断进行路由。

### 请求转换委托
实际的请求转换逻辑封装在 `relay-cohere.go` 中的 `requestOpenAI2Cohere` 和 `requestConvertRerank2Cohere` 函数中，Adaptor 只负责调度。

### 不支持 Responses API
与 Cloudflare 不同，Cohere 适配器不支持 OpenAI Responses API 格式。

## 6. 关联文件

- `relay/channel/cohere/relay-cohere.go` - 请求转换和响应处理的具体实现
- `relay/channel/cohere/dto.go` - Cohere 特有的 DTO 定义
- `relay/channel/cohere/constant.go` - 模型列表
