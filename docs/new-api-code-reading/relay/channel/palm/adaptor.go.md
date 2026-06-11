# adaptor.go 代码阅读文档

## 1. 全局总结
本文件实现了 Google PaLM（Pathways Language Model）渠道的适配器。`Adaptor` 结构体实现了 `channel.Adaptor` 接口，负责将统一的 OpenAI 格式请求转换为 PaLM API 格式，并处理 PaLM 的响应。该适配器支持流式和非流式两种响应模式。

## 2. 依赖关系
| 包名 | 用途 |
|------|------|
| `errors` | 错误创建 |
| `fmt` | 格式化字符串 |
| `io` | IO 操作 |
| `net/http` | HTTP 请求/响应处理 |
| `github.com/QuantumNous/new-api/dto` | 数据传输对象 |
| `github.com/QuantumNous/new-api/relay/channel` | 渠道公共功能（API 请求发送、请求头设置） |
| `github.com/QuantumNous/new-api/relay/common` | 中继公共类型（RelayInfo） |
| `github.com/QuantumNous/new-api/service` | 服务层（用量计算、响应处理） |
| `github.com/QuantumNous/new-api/types` | 类型定义（错误类型） |
| `github.com/gin-gonic/gin` | Gin HTTP 框架 |

## 3. 类型定义

### `Adaptor`
PaLM 渠道适配器结构体，无字段，通过方法实现接口。

## 4. 函数详解

### `ConvertGeminiRequest`
```go
func (a *Adaptor) ConvertGeminiRequest(*gin.Context, *relaycommon.RelayInfo, *dto.GeminiChatRequest) (any, error)
```
- **状态**：未实现，返回 "not implemented" 错误
- **说明**：PaLM 不支持 Gemini 格式的请求转换

### `ConvertClaudeRequest`
```go
func (a *Adaptor) ConvertClaudeRequest(*gin.Context, *relaycommon.RelayInfo, *dto.ClaudeRequest) (any, error)
```
- **状态**：未实现，调用 `panic("implement me")`
- **说明**：PaLM 不支持 Claude 格式的请求转换

### `ConvertAudioRequest`
```go
func (a *Adaptor) ConvertAudioRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.AudioRequest) (io.Reader, error)
```
- **状态**：未实现，返回 "not implemented" 错误

### `ConvertImageRequest`
```go
func (a *Adaptor) ConvertImageRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.ImageRequest) (any, error)
```
- **状态**：未实现，返回 "not implemented" 错误

### `Init`
```go
func (a *Adaptor) Init(info *relaycommon.RelayInfo)
```
- **说明**：初始化适配器，当前为空实现

### `GetRequestURL`
```go
func (a *Adaptor) GetRequestURL(info *relaycommon.RelayInfo) (string, error)
```
- **返回**：PaLM v1beta2 的 chat-bison-001 模型的 generateMessage 端点 URL
- **格式**：`{base_url}/v1beta2/models/chat-bison-001:generateMessage`

### `SetupRequestHeader`
```go
func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error
```
- **说明**：设置请求头，使用 `x-goog-api-key` 而非标准的 `Authorization` 头
- **逻辑**：先调用公共方法设置通用头，再添加 Google API 密钥

### `ConvertOpenAIRequest`
```go
func (a *Adaptor) ConvertOpenAIRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error)
```
- **说明**：直接透传请求对象（不做格式转换），仅检查 nil

### `DoRequest`
```go
func (a *Adaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (any, error)
```
- **说明**：委托给 `channel.DoApiRequest` 发送 HTTP 请求

### `DoResponse`
```go
func (a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage any, err *types.NewAPIError)
```
- **说明**：根据 `info.IsStream` 分流处理：
  - **流式**：调用 `palmStreamHandler`，通过 `service.ResponseText2Usage` 计算用量
  - **非流式**：调用 `palmHandler`，返回标准 OpenAI 格式的用量信息

### `GetModelList` / `GetChannelName`
```go
func (a *Adaptor) GetModelList() []string
func (a *Adaptor) GetChannelName() string
```
- **说明**：返回渠道常量中定义的模型列表和渠道名称

## 5. 关键逻辑分析
- **API Key 认证方式**：PaLM 使用 `x-goog-api-key` 请求头进行认证，而非标准的 `Authorization: Bearer` 方式。
- **固定模型端点**：请求 URL 固定为 `chat-bison-001` 模型，不支持动态切换模型。
- **请求透传**：`ConvertOpenAIRequest` 直接返回原始请求对象，表明 PaLM 的 API 格式与 OpenAI 格式相似。
- **未实现的功能**：音频、图像、Gemini、Claude、嵌入、Rerank 等功能均未实现，说明该渠道专注于文本聊天。

## 6. 关联文件
- `relay/channel/palm/constants.go` — 模型列表和渠道名称定义
- `relay/channel/palm/dto.go` — PaLM 专用数据结构
- `relay/channel/palm/relay-palm.go` — PaLM 响应处理逻辑
- `relay/channel/adapter.go` — `Adaptor` 接口定义
- `relay/channel/channel.go` — `DoApiRequest` 和 `SetupApiRequestHeader` 公共函数
