# adaptor.go 代码阅读文档

## 1. 全局总结
本文件是即梦（Jimeng）渠道的适配器入口，实现了 `channel.Adaptor` 接口。即梦是字节跳动的 AI 图像生成平台，该适配器主要支持图像生成（CVProcess）功能，同时通过 OpenAI 兼容接口支持文本对话。使用自定义的 HMAC-SHA256 签名认证机制。

## 2. 依赖关系
- **标准库**: `encoding/json`, `errors`, `fmt`, `io`, `net/http`
- **项目内部**:
  - `github.com/QuantumNous/new-api/dto` — 数据传输对象
  - `github.com/QuantumNous/new-api/relay/channel` — 通用渠道工具
  - `github.com/QuantumNous/new-api/relay/channel/openai` — OpenAI 处理器（用于文本对话响应）
  - `github.com/QuantumNous/new-api/relay/common` — RelayInfo
  - `github.com/QuantumNous/new-api/relay/constant` — 中继模式常量
  - `github.com/QuantumNous/new-api/types` — 错误类型
- **外部依赖**: `github.com/gin-gonic/gin`

## 3. 类型定义

### `Adaptor` 结构体
```go
type Adaptor struct {}
```
即梦渠道适配器。

### `LogoInfo` 结构体
```go
type LogoInfo struct {
    AddLogo         bool    `json:"add_logo,omitempty"`
    Position        int     `json:"position,omitempty"`
    Language        int     `json:"language,omitempty"`
    Opacity         float64 `json:"opacity,omitempty"`
    LogoTextContent string  `json:"logo_text_content,omitempty"`
}
```
水印配置信息，控制图片水印的开关、位置、语言、透明度和文字内容。

### `imageRequestPayload` 结构体
```go
type imageRequestPayload struct {
    ReqKey     string   `json:"req_key"`
    Prompt     string   `json:"prompt"`
    Seed       int64    `json:"seed,omitempty"`
    Width      int      `json:"width,omitempty"`
    Height     int      `json:"height,omitempty"`
    UsePreLLM  bool     `json:"use_pre_llm,omitempty"`
    UseSR      bool     `json:"use_sr,omitempty"`
    ReturnURL  bool     `json:"return_url,omitempty"`
    LogoInfo   LogoInfo `json:"logo_info,omitempty"`
    ImageUrls  []string `json:"image_urls,omitempty"`
    BinaryData []string `json:"binary_data_base64,omitempty"`
}
```
即梦图像生成请求载荷，包含服务标识、提示词、随机种子、图片尺寸、超分辨率开关、返回 URL 开关、水印配置、输入图片等参数。

## 4. 函数详解

### `ConvertGeminiRequest` / `ConvertClaudeRequest` / `ConvertAudioRequest` / `ConvertEmbeddingRequest` / `ConvertRerankRequest` / `ConvertOpenAIResponsesRequest`
均返回 `"not implemented"` 错误，即梦渠道不支持这些转换。

### `Init(info)`
空实现。

### `GetRequestURL(info) (string, error)`
构建请求 URL：`{baseUrl}/?Action=CVProcess&Version=2022-08-31`，使用即梦的视觉处理 API。

### `SetupRequestHeader(c, header, info) error`
返回 `"not implemented"` 错误。即梦使用自定义签名而非标准 header 设置，在 `DoRequest` 中通过 `Sign` 函数处理。

### `ConvertOpenAIRequest(c, info, request) (any, error)`
直接透传 OpenAI 请求，不做格式转换。

### `ConvertImageRequest(c, info, request) (any, error)`
将 OpenAI ImageRequest 转换为即梦 imageRequestPayload：
- 使用 `request.Model` 作为 `req_key`
- 默认返回 URL 格式
- 通过 `ExtraFields` 支持即梦特有的参数扩展

### `DoRequest(c, info, requestBody) (any, error)`
自定义请求发送流程：
1. 获取请求 URL
2. 创建 http.Request
3. 调用 `Sign` 函数进行 HMAC-SHA256 签名
4. 通过 `channel.DoRequest` 发送请求

### `DoResponse(c, resp, info) (usage, err)`
响应处理：
- 图像生成模式 → `jimengImageHandler`
- 流式 → `openai.OaiStreamHandler`
- 非流式 → `openai.OpenaiHandler`

### `GetModelList() []string`
返回 `ModelList`。

### `GetChannelName() string`
返回 `ChannelName`（`"jimeng"`）。

## 5. 关键逻辑分析

### 认证机制
即梦使用 HMAC-SHA256 签名认证，而非标准的 Bearer Token。API Key 格式为 `ak|sk`（accessKey|secretKey），签名过程在 `sign.go` 中实现。这与大多数 AI API 的认证方式不同。

### 图像生成载荷
`imageRequestPayload` 使用即梦特有的 `req_key` 字段标识服务（如 `jimeng_high_aes_general_v21_L`），支持中文/英文提示词、超分辨率、水印等特有功能。

### 混合响应模式
文本对话使用 OpenAI 格式的处理器，图像生成使用自定义处理器，实现了同一渠道对不同能力的混合支持。

## 6. 关联文件
- `sign.go` — HMAC-SHA256 签名实现
- `image.go` — 图像生成响应处理
- `constants.go` — 模型列表和渠道名称
- `relay/channel/openai/` — OpenAI 兼容的文本对话处理
