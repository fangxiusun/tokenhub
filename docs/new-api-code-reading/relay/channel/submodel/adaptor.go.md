# adaptor.go 代码阅读文档

## 1. 全局总结
本文件实现了 SubModel 渠道的适配器。SubModel 是一个提供多种开源大模型服务的平台，适配器主要支持聊天功能，将请求和响应处理委托给 OpenAI 适配器。大部分非聊天功能（音频、图像、嵌入、Rerank 等）不被支持。

## 2. 依赖关系
| 包名 | 用途 |
|------|------|
| `errors` | 错误创建 |
| `io` | IO 操作 |
| `net/http` | HTTP 请求/响应处理 |
| `github.com/QuantumNous/new-api/dto` | 数据传输对象 |
| `github.com/QuantumNous/new-api/relay/channel` | 渠道公共功能 |
| `github.com/QuantumNous/new-api/relay/channel/openai` | OpenAI 适配器（委托目标） |
| `github.com/QuantumNous/new-api/relay/common` | 中继公共类型 |
| `github.com/QuantumNous/new-api/types` | 类型定义 |
| `github.com/gin-gonic/gin` | Gin HTTP 框架 |

## 3. 类型定义

### `Adaptor`
SubModel 渠道适配器结构体，无字段。

## 4. 函数详解

### 未实现的方法（返回 "endpoint not supported" 错误）
| 方法 | 说明 |
|------|------|
| `ConvertGeminiRequest` | 不支持 Gemini 格式 |
| `ConvertClaudeRequest` | 不支持 Claude 格式 |
| `ConvertAudioRequest` | 不支持音频请求 |
| `ConvertImageRequest` | 不支持图像请求 |
| `ConvertRerankRequest` | 不支持 Rerank 请求 |
| `ConvertEmbeddingRequest` | 不支持嵌入请求 |
| `ConvertOpenAIResponsesRequest` | 不支持 Responses 请求 |

### `Init`
```go
func (a *Adaptor) Init(info *relaycommon.RelayInfo)
```
- 空实现

### `GetRequestURL`
```go
func (a *Adaptor) GetRequestURL(info *relaycommon.RelayInfo) (string, error)
```
- 拼接完整请求 URL

### `SetupRequestHeader`
```go
func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error
```
- 设置标准 `Authorization: Bearer` 请求头

### `ConvertOpenAIRequest`
```go
func (a *Adaptor) ConvertOpenAIRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error)
```
- 直接透传请求（仅检查 nil）

### `DoRequest`
```go
func (a *Adaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (any, error)
```
- 委托给 `channel.DoApiRequest`

### `DoResponse`
```go
func (a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage any, err *types.NewAPIError)
```
- **逻辑**：
  - **流式**：调用 `openai.OaiStreamHandler`
  - **非流式**：调用 `openai.OpenaiHandler`
- 委托给 OpenAI 适配器的公开处理函数

### `GetModelList` / `GetChannelName`
```go
func (a *Adaptor) GetModelList() []string
func (a *Adaptor) GetChannelName() string
```
- 返回渠道常量中定义的模型列表和名称

## 5. 关键逻辑分析
- **最小化实现**：SubModel 适配器只实现了聊天功能，其他功能全部返回不支持错误。
- **OpenAI 兼容**：由于 SubModel 兼容 OpenAI 格式，大部分逻辑直接委托给 OpenAI 适配器。
- **流式/非流式分流**：`DoResponse` 根据 `IsStream` 标志分别调用不同的 OpenAI 处理函数。
- **请求透传**：`ConvertOpenAIRequest` 直接返回原始请求，不做任何格式转换。

## 6. 关联文件
- `relay/channel/submodel/constants.go` — 模型列表和渠道名称
- `relay/channel/openai/adaptor.go` — 被委托的 OpenAI 适配器
- `relay/channel/openai/handler.go` — `OaiStreamHandler` 和 `OpenaiHandler` 函数
- `relay/channel/adapter.go` — `Adaptor` 接口定义
- `relay/channel/channel.go` — 公共请求发送逻辑
