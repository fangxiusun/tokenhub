# adaptor.go 代码阅读文档

## 1. 全局总结
本文件实现了 Perplexity AI 渠道的适配器。`Adaptor` 结构体实现了 `channel.Adaptor` 接口，负责请求格式转换和响应处理。Perplexity 兼容 OpenAI 格式，因此大部分功能委托给 OpenAI 适配器处理，仅在请求转换时做特定调整（如 TopP 限制）。

## 2. 依赖关系
| 包名 | 用途 |
|------|------|
| `errors` | 错误创建 |
| `fmt` | 格式化字符串 |
| `io` | IO 操作 |
| `net/http` | HTTP 请求/响应处理 |
| `github.com/QuantumNous/new-api/dto` | 数据传输对象 |
| `github.com/QuantumNous/new-api/relay/channel` | 渠道公共功能 |
| `github.com/QuantumNous/new-api/relay/channel/openai` | OpenAI 适配器（用于委托） |
| `github.com/QuantumNous/new-api/relay/common` | 中继公共类型 |
| `github.com/QuantumNous/new-api/relay/constant` | 中继常量（RelayModeResponses） |
| `github.com/QuantumNous/new-api/types` | 类型定义 |
| `github.com/samber/lo` | 辅助函数库 |
| `github.com/gin-gonic/gin` | Gin HTTP 框架 |

## 3. 类型定义

### `Adaptor`
Perplexity 渠道适配器结构体，无字段。

## 4. 函数详解

### `ConvertClaudeRequest`
```go
func (a *Adaptor) ConvertClaudeRequest(c *gin.Context, info *relaycommon.RelayInfo, req *dto.ClaudeRequest) (any, error)
```
- **说明**：委托给 `openai.Adaptor.ConvertClaudeRequest` 处理
- **原因**：Perplexity 兼容 OpenAI 格式，Claude 请求可直接转换

### `GetRequestURL`
```go
func (a *Adaptor) GetRequestURL(info *relaycommon.RelayInfo) (string, error)
```
- **说明**：根据中继模式返回不同的端点 URL
  - `RelayModeResponses`：返回 `/v1/responses`
  - 其他模式：返回 `/chat/completions`

### `SetupRequestHeader`
```go
func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error
```
- **说明**：设置标准 `Authorization: Bearer` 请求头

### `ConvertOpenAIRequest`
```go
func (a *Adaptor) ConvertOpenAIRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error)
```
- **说明**：转换 OpenAI 请求为 Perplexity 格式
- **关键逻辑**：
  - 检查请求非 nil
  - 如果 `TopP >= 1`，强制设置为 `0.99`（Perplexity 不支持 TopP=1）
  - 调用 `requestOpenAI2Perplexity` 进行格式转换

### `ConvertOpenAIResponsesRequest`
```go
func (a *Adaptor) ConvertOpenAIResponsesRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.OpenAIResponsesRequest) (any, error)
```
- **说明**：直接透传 OpenAI Responses 请求

### `DoRequest`
```go
func (a *Adaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (any, error)
```
- **说明**：委托给 `channel.DoApiRequest` 发送请求

### `DoResponse`
```go
func (a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage any, err *types.NewAPIError)
```
- **说明**：完全委托给 `openai.Adaptor.DoResponse` 处理
- **原因**：Perplexity 的响应格式与 OpenAI 完全兼容

### 未实现的方法
以下方法返回 "not implemented" 错误：
- `ConvertGeminiRequest` — 不支持 Gemini 格式
- `ConvertAudioRequest` — 不支持音频请求
- `ConvertImageRequest` — 不支持图像请求
- `ConvertEmbeddingRequest` — 不支持嵌入请求

### 其他方法
- `Init` — 空实现
- `ConvertRerankRequest` — 返回 nil, nil
- `GetModelList` — 返回 `ModelList`
- `GetChannelName` — 返回 `ChannelName`

## 5. 关键逻辑分析
- **TopP 限制**：Perplexity 不支持 `TopP=1`，适配器在转换时自动将其限制为 `0.99`。
- **OpenAI 委托模式**：大量功能（Claude 请求转换、响应处理）直接委托给 OpenAI 适配器，减少了代码重复。
- **多端点支持**：根据中继模式（`RelayModeResponses`）支持不同的 API 端点。
- **透传设计**：对于 OpenAI Responses 请求直接透传，不做格式转换。

## 6. 关联文件
- `relay/channel/perplexity/constants.go` — 模型列表和渠道名称
- `relay/channel/perplexity/relay-perplexity.go` — 请求格式转换逻辑
- `relay/channel/openai/adaptor.go` — 被委托的 OpenAI 适配器
- `relay/channel/adapter.go` — `Adaptor` 接口定义
