# adaptor.go 代码阅读文档

## 1. 全局总结
本文件是 Gemini 渠道的适配器入口，实现了 `channel.Adaptor` 接口，负责在统一的 OpenAI 格式与 Gemini 原生格式之间进行请求/响应的转换。它是 Gemini 渠道的核心调度层，决定了如何构建请求 URL、设置认证头、转换各种格式的请求体（OpenAI/Claude/Gemini/Image/Embedding），以及如何处理不同类型的响应（流式/非流式、聊天/嵌入/图片生成）。

## 2. 依赖关系
- **标准库**: `errors`, `fmt`, `io`, `net/http`, `strings`
- **项目内部**:
  - `github.com/QuantumNous/new-api/dto` — 数据传输对象定义
  - `github.com/QuantumNous/new-api/relay/channel` — 通用渠道工具（请求构造、头设置）
  - `github.com/QuantumNous/new-api/relay/channel/openai` — OpenAI 渠道适配器（用于 Claude→OpenAI 转换）
  - `github.com/QuantumNous/new-api/relay/common` — 中继通用信息（RelayInfo）
  - `github.com/QuantumNous/new-api/relay/constant` — 中继模式常量
  - `github.com/QuantumNous/new-api/setting/model_setting` — 模型设置（Gemini 版本、安全设置等）
  - `github.com/QuantumNous/new-api/setting/reasoning` — 推理努力级别处理
  - `github.com/QuantumNous/new-api/types` — 错误类型定义
- **外部依赖**:
  - `github.com/gin-gonic/gin` — Web 框架
  - `github.com/samber/lo` — 泛型工具库

## 3. 类型定义

### `Adaptor` 结构体
```go
type Adaptor struct {}
```
Gemini 渠道适配器，空结构体，通过方法实现 `channel.Adaptor` 接口。

## 4. 函数详解

### `ConvertGeminiRequest(c, info, request) (any, error)`
将 Gemini 原生请求进行预处理：确保第一个 content 的 role 默认为 `"user"`，并为 YouTube 视频链接自动补充 MIME 类型 `video/webm`。

### `ConvertClaudeRequest(c, info, req) (any, error)`
将 Claude 格式请求先通过 OpenAI 适配器转换为 OpenAI 格式，再调用 `ConvertOpenAIRequest` 转换为 Gemini 格式。

### `ConvertAudioRequest(c, info, request) (io.Reader, error)`
未实现，返回 `"not implemented"` 错误。

### `ConvertImageRequest(c, info, request) (any, error)`
将 OpenAI ImageRequest 转换为 Gemini Imagen 请求格式。仅支持 `imagen` 前缀的模型。将 `size` 转换为 `aspectRatio`（支持 `1:1`, `3:2`, `2:3`, `9:16`, `16:9` 等），支持用户直接指定宽高比。根据 `Quality` 参数映射 `imageSize`（`1K`/`2K`）。

### `Init(info)`
空实现，无需初始化。

### `GetRequestURL(info) (string, error)`
构建 Gemini API 请求 URL：
- 首先处理 thinking 适配器相关的模型名后缀剥离（`-thinking-<budget>`, `-thinking`, `-nothinking`, effort 后缀）
- 获取 Gemini API 版本号
- 根据模型类型和中继模式选择不同端点：
  - `imagen` 模型 → `models/{model}:predict`
  - 嵌入模型 → `models/{model}:embedContent` 或 `:batchEmbedContents`
  - 流式 → `models/{model}:streamGenerateContent?alt=sse`
  - 非流式 → `models/{model}:generateContent`

### `SetupRequestHeader(c, req, info) error`
设置请求头：调用通用渠道头设置函数，并添加 `x-goog-api-key` 认证头。

### `ConvertOpenAIRequest(c, info, request) (any, error)`
调用 `CovertOpenAI2Gemini` 将 OpenAI 请求转换为 Gemini 格式。

### `ConvertRerankRequest(c, relayMode, request) (any, error)`
返回 nil，Gemini 不支持 rerank。

### `ConvertEmbeddingRequest(c, info, request) (any, error)`
将 OpenAI 嵌入请求转换为 Gemini 批量嵌入格式。为每个输入创建独立的嵌入请求对象，包含 `model` 和 `content`。对较新的嵌入模型（`text-embedding-004` 等）支持 `outputDimensionality` 参数。始终使用批量嵌入端点。

### `ConvertOpenAIResponsesRequest(c, info, request) (any, error)`
未实现。

### `DoRequest(c, info, requestBody) (any, error)`
调用 `channel.DoApiRequest` 发送请求。

### `DoResponse(c, resp, info) (usage, err)`
响应处理调度器，根据中继模式和模型类型选择不同处理器：
- `RelayModeGemini` + 嵌入 → `NativeGeminiEmbeddingHandler`
- `RelayModeGemini` + 流式 → `GeminiTextGenerationStreamHandler`
- `RelayModeGemini` + 非流式 → `GeminiTextGenerationHandler`
- `imagen` 模型 → `GeminiImageHandler`
- 嵌入模型 → `GeminiEmbeddingHandler`
- 流式 → `GeminiChatStreamHandler`
- 非流式 → `GeminiChatHandler`

### `GetModelList() []string`
返回 Gemini 支持的模型列表。

### `GetChannelName() string`
返回渠道名称 `"google gemini"`。

## 5. 关键逻辑分析

### 请求 URL 构建策略
`GetRequestURL` 实现了多级路由逻辑：先处理 thinking 适配器的模型名清理，再根据模型前缀（`imagen`、`text-embedding`、`embedding`、`gemini-embedding`）和流式标志选择不同的 API 端点。流式请求使用 `?alt=sse` 查询参数。

### 图片请求转换
`ConvertImageRequest` 将标准的 `size` 字符串（如 `"1024x1024"`）映射为 Gemini 的 `aspectRatio` 格式。支持用户直接传入宽高比字符串。质量参数映射到 `imageSize`，但仅 Standard 和 Ultra 模型支持。

### 嵌入请求批量构建
始终构建批量嵌入请求格式（`requests` 数组），即使只有单个输入。这确保了请求格式与端点的一致性。

### 响应路由
`DoResponse` 根据多种条件（中继模式、模型前缀、流式标志）进行多级分发，是整个适配器的核心调度逻辑。

## 6. 关联文件
- `relay-gemini.go` — 核心转换和响应处理函数实现
- `relay-gemini-native.go` — Gemini 原生模式的响应处理
- `constant.go` — 模型列表和渠道常量
- `relay/channel/openai/adaptor.go` — OpenAI 适配器，用于 Claude 格式转换
- `relay/channel/channel.go` — 通用渠道工具函数
