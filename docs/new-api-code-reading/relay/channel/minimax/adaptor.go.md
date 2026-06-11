# adaptor.go 代码阅读文档

## 1. 全局总结
本文件是 MiniMax 渠道的适配器入口，实现了 `channel.Adaptor` 接口。MiniMax 是一个提供文本生成、语音合成（TTS）和图像生成能力的 AI 平台。该适配器支持多种中继模式，并能根据 `RelayFormat` 自动选择 Claude 或 OpenAI 格式的处理器。

## 2. 依赖关系
- **标准库**: `bytes`, `encoding/json`, `errors`, `fmt`, `io`, `net/http`
- **项目内部**:
  - `github.com/QuantumNous/new-api/dto` — 数据传输对象
  - `github.com/QuantumNous/new-api/relay/channel` — 通用渠道工具
  - `github.com/QuantumNous/new-api/relay/channel/claude` — Claude 适配器
  - `github.com/QuantumNous/new-api/relay/channel/openai` — OpenAI 适配器
  - `github.com/QuantumNous/new-api/relay/common` — RelayInfo
  - `github.com/QuantumNous/new-api/relay/constant` — 中继模式常量
  - `github.com/QuantumNous/new-api/types` — 错误类型
- **外部依赖**: `github.com/gin-gonic/gin`, `github.com/samber/lo`

## 3. 类型定义

### `Adaptor` 结构体
```go
type Adaptor struct {}
```
MiniMax 渠道适配器。

## 4. 函数详解

### `ConvertGeminiRequest`
返回 `"not implemented"` 错误。

### `ConvertClaudeRequest(c, info, req) (any, error)`
委托给 `claude.Adaptor` 的 `ConvertClaudeRequest` 方法处理。这使得 MiniMax 能够接受 Claude 格式的请求。

### `ConvertAudioRequest(c, info, request) (io.Reader, error)`
仅支持 `RelayModeAudioSpeech` 模式，将 OpenAI AudioRequest 转换为 MiniMax TTS 请求：
- 提取 voice、speed、output_format 参数
- 构建 `MiniMaxTTSRequest` 结构体
- 支持通过 `Metadata` 扩展即梦特有参数
- 将响应格式设置为 `"hex"` 或 `"url"`
- 返回 JSON 序列化后的请求体

### `ConvertImageRequest(c, info, request) (any, error)`
仅支持 `RelayModeImagesGenerations` 模式，调用 `oaiImage2MiniMaxImageRequest` 转换图像请求。

### `Init(info)`
空实现。

### `GetRequestURL(info) (string, error)`
委托给 `relay-minimax.go` 中的 `GetRequestURL` 函数。

### `SetupRequestHeader(c, req, info) error`
设置标准 API 请求头，并添加 `Authorization: Bearer {apiKey}`。

### `ConvertOpenAIRequest(c, info, request) (any, error)`
直接透传 OpenAI 请求，不做格式转换。

### `ConvertRerankRequest`
返回 nil（不支持 rerank）。

### `ConvertEmbeddingRequest`
直接透传嵌入请求。

### `ConvertOpenAIResponsesRequest`
返回 `"not implemented"` 错误。

### `DoRequest(c, info, requestBody) (any, error)`
调用 `channel.DoApiRequest` 发送请求。

### `DoResponse(c, resp, info) (usage, err)`
多模式响应处理：
- `RelayModeAudioSpeech` → `handleTTSResponse`
- `RelayModeImagesGenerations` → `miniMaxImageHandler`
- Claude 格式 → `claude.Adaptor.DoResponse`
- 默认 → `openai.Adaptor.DoResponse`

### `GetModelList() []string`
返回 `ModelList`。

### `GetChannelName() string`
返回 `ChannelName`（`"minimax"`）。

## 5. 关键逻辑分析

### 多能力支持
MiniMax 是少数同时支持文本生成、TTS 和图像生成的渠道之一。适配器通过 `RelayMode` 和 `RelayFormat` 两个维度进行请求/响应分发。

### Claude 兼容性
MiniMax 通过委托给 `claude.Adaptor` 实现了 Claude 格式的支持，这是一种复用策略，避免重复实现 Claude 格式的转换逻辑。

### TTS 请求转换
音频请求转换时，支持通过 `Metadata` 字段传递 MiniMax 特有的参数（如 `pronunciation_dict`、`timbre_weights`、`voice_modify` 等），实现了对 OpenAI AudioRequest 的扩展。

### 响应格式选择
`DoResponse` 根据 `RelayFormat` 选择 Claude 或 OpenAI 处理器，实现了同一渠道对不同客户端格式的兼容。

## 6. 关联文件
- `relay-minimax.go` — 请求 URL 构建
- `tts.go` — TTS 请求/响应类型定义和处理
- `image.go` — 图像生成请求转换和响应处理
- `constants.go` — 模型列表
- `relay/channel/claude/` — Claude 格式处理
- `relay/channel/openai/` — OpenAI 格式处理
