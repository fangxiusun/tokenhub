# adaptor.go 代码阅读文档

## 1. 全局总结
本文件实现了 SiliconFlow（硅基流动）渠道的适配器。SiliconFlow 是一个提供多种 AI 模型服务的平台，本适配器支持聊天、图像生成、Rerank、嵌入等多种功能。适配器大量委托给 OpenAI 适配器处理，仅在特定场景下做自定义处理（如 FIM 请求、图像请求、Rerank 响应）。

## 2. 依赖关系
| 包名 | 用途 |
|------|------|
| `errors` | 错误创建 |
| `fmt` | 格式化字符串 |
| `io` | IO 操作 |
| `net/http` | HTTP 请求/响应处理 |
| `github.com/QuantumNous/new-api/common` | JSON 封装 |
| `github.com/QuantumNous/new-api/dto` | 数据传输对象 |
| `github.com/QuantumNous/new-api/relay/channel` | 渠道公共功能 |
| `github.com/QuantumNous/new-api/relay/channel/openai` | OpenAI 适配器（委托目标） |
| `github.com/QuantumNous/new-api/relay/common` | 中继公共类型 |
| `github.com/QuantumNous/new-api/relay/constant` | 中继常量（RelayModeRerank） |
| `github.com/QuantumNous/new-api/types` | 类型定义 |
| `github.com/gin-gonic/gin` | Gin HTTP 框架 |
| `github.com/samber/lo` | 辅助函数库 |

## 3. 类型定义

### `Adaptor`
SiliconFlow 渠道适配器结构体，无字段。

## 4. 函数详解

### `ConvertClaudeRequest`
```go
func (a *Adaptor) ConvertClaudeRequest(c *gin.Context, info *relaycommon.RelayInfo, req *dto.ClaudeRequest) (any, error)
```
- 委托给 `openai.Adaptor.ConvertClaudeRequest`

### `ConvertAudioRequest`
```go
func (a *Adaptor) ConvertAudioRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.AudioRequest) (io.Reader, error)
```
- 委托给 `openai.Adaptor.ConvertAudioRequest`

### `ConvertImageRequest`
```go
func (a *Adaptor) ConvertImageRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.ImageRequest) (any, error)
```
- **说明**：将 OpenAI 图像请求转换为 SiliconFlow 图像请求格式
- **逻辑**：
  1. 从 `request.Extra` 解析 `SFImageRequest` 结构
  2. 设置 `Model` 和 `Prompt`
  3. 如果 `SFImageRequest.ImageSize` 为空，使用 OpenAI 的 `Size`
  4. 如果 `SFImageRequest.BatchSize` 为 0，使用 OpenAI 的 `N`

### `GetRequestURL`
```go
func (a *Adaptor) GetRequestURL(info *relaycommon.RelayInfo) (string, error)
```
- **逻辑**：
  - `RelayModeRerank`：返回 `/v1/rerank`
  - 其他模式：拼接完整 URL

### `SetupRequestHeader`
```go
func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error
```
- 设置标准 `Authorization: Bearer` 请求头

### `ConvertOpenAIRequest`
```go
func (a *Adaptor) ConvertOpenAIRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error)
```
- **说明**：处理 FIM（Fill-In-the-Middle）请求的特殊需求
- **逻辑**：如果请求包含 `Prefix` 或 `Suffix` 但没有 `Messages`，添加一个空的 user 消息（SiliconFlow 要求）

### `ConvertRerankRequest`
```go
func (a *Adaptor) ConvertRerankRequest(c *gin.Context, relayMode int, request dto.RerankRequest) (any, error)
```
- 直接透传请求

### `ConvertEmbeddingRequest`
```go
func (a *Adaptor) ConvertEmbeddingRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.EmbeddingRequest) (any, error)
```
- 直接透传请求

### `DoRequest`
```go
func (a *Adaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (any, error)
```
- 委托给 `openai.Adaptor.DoRequest`

### `DoResponse`
```go
func (a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage any, err *types.NewAPIError)
```
- **逻辑**：
  - `RelayModeRerank`：调用 `siliconflowRerankHandler` 处理 Rerank 响应
  - 其他模式：委托给 `openai.Adaptor.DoResponse`

### 未实现的方法
- `ConvertGeminiRequest` — 返回 "not implemented"
- `ConvertOpenAIResponsesRequest` — 返回 "not implemented"

### 其他方法
- `Init` — 空实现
- `GetModelList` — 返回 `ModelList`
- `GetChannelName` — 返回 `ChannelName`

## 5. 关键逻辑分析
- **OpenAI 委托模式**：大部分功能（Claude 请求、音频请求、请求发送、默认响应处理）都委托给 OpenAI 适配器，减少代码重复。
- **FIM 特殊处理**：SiliconFlow 的 FIM 端点要求 `Messages` 数组非空，即使客户端没有发送。
- **图像请求转换**：支持 SiliconFlow 特有的图像参数（`image_size`、`batch_size` 等），通过 `Extra` 字段传递。
- **Rerank 支持**：SiliconFlow 支持 Rerank 功能，有独立的响应处理逻辑。
- **URL 路由**：根据中继模式动态选择 API 端点。

## 6. 关联文件
- `relay/channel/siliconflow/constant.go` — 模型列表和渠道名称
- `relay/channel/siliconflow/dto.go` — SiliconFlow 专用数据结构
- `relay/channel/siliconflow/relay-siliconflow.go` — Rerank 响应处理
- `relay/channel/openai/adaptor.go` — 被委托的 OpenAI 适配器
- `relay/channel/adapter.go` — `Adaptor` 接口定义
