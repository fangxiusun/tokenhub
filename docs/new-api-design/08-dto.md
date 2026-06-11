# 数据传输对象层详细设计 (`dto/`)

## 1. 概述
数据传输对象（DTO）层包含所有 API 端点和提供商集成的请求和响应结构体。这些结构体定义了系统内部和外部 API 之间的数据契约。

## 2. 文件详细说明

### 2.1 OpenAI 兼容
- **`openai_request.go`** -- **核心请求 DTO**：`GeneralOpenAIRequest` -- 统一的聊天补全请求，包含约 100 个字段（messages, tools, response_format, reasoning, audio 等），以及 `RequestMessage`, `ToolCall`, `ToolChoice`, `StreamOptions`, Token 计数元数据和零值保留逻辑。
- **`openai_response.go`** -- 所有 OpenAI 响应类型：`SimpleResponse`, `TextResponse`, `TextStreamResponse`, `ImageResponse`, `Usage`, `OpenAITextResponseChoice`, `OpenAIResponse`（Responses API）, `OpenAIImageResponse`，以及流式增量类型。
- **`openai_image.go`** -- `ImageRequest`/`ImageResponse` 用于 OpenAI 图像生成 API，包括 `ImageEditRequest`, `ImageVariationRequest`, `ImageEditResponse`。
- **`openai_video.go`** -- `OpenAIVideo` 结构体和状态常量用于 OpenAI 视频生成 API（queued/in_progress/completed/failed）。
- **`openai_compaction.go`** -- `OpenAIResponsesCompactionResponse` 用于 OpenAI Responses 压缩端点，带使用量跟踪。
- **`openai_responses_compaction_request.go`** -- `OpenAIResponsesCompactionRequest` 用于 Responses 压缩 API。

### 2.2 Claude/Anthropic
- **`claude.go`** -- 完整的 Claude/Anthropic API 请求/响应类型：消息、媒体、使用量、流式增量、内容块（文本、图像、工具使用、思考、缓存控制）和转换辅助。

### 2.3 Gemini/Google
- **`gemini.go`** -- 完整的 Gemini API 请求/响应类型：聊天内容、安全设置、生成配置、工具配置、函数调用、基础元数据，以及双向 camelCase/snake_case 反序列化。

### 2.4 其他提供商
- **`midjourney.go`** -- Midjourney 代理 API 类型：imagine, describe, blend, upscale, swap-face 操作的请求/响应；任务状态轮询。
- **`suno.go`** -- Suno AI 音乐生成类型：`SunoSubmitReq`, `SunoDataResponse`, `SunoSong`, `SunoLyrics`。
- **`rerank.go`** -- `RerankRequest`/`RerankResponse` 用于 Jina 重排序 API，带 Token 计数元数据提取。

### 2.5 通用类型
- **`request_common.go`** -- `Request` 接口和 `BaseRequest`：所有中转请求类型的通用契约（`GetTokenCountMeta()`, `IsStream()`, `SetModelName()`）。
- **`error.go`** -- 错误 DTO：`OpenAIErrorWithStatusCode`, `GeneralErrorResponse`（带回退字段提取，用于多样化的上游错误格式）和 `GetOpenAIError()` 标准化函数。
- **`values.go`** -- 多态值类型：`StringValue`, `IntValue`, `Float64Value`，带自定义 JSON 反序列化，接受字符串和数字输入（用于上游 API 不一致地类型化数值字段的情况）。

### 2.6 任务类型
- **`task.go`** -- 通用任务响应类型：`TaskError`, `TaskResponse[T]`，用于异步任务平台（Suno, Midjourney 等）。
- **`video.go`** -- `VideoRequest`/`VideoResponse`/`VideoTaskResponse` 用于第三方视频生成 API（Kling 等）。

### 2.7 用户与设置
- **`user_settings.go`** -- `UserSetting` 结构体用于每用户偏好：通知类型、配额警告、webhook/邮件/Bark/Gotify 设置、计费偏好、语言。
- **`channel_settings.go`** -- `ChannelSettings`, `ChannelOtherSettings` 结构体用于每渠道配置：代理、thinking-to-content、Azure responses 版本、Vertex 密钥类型、OpenRouter 企业版、Claude beta 等。
- **`notify.go`** -- `Notify` 结构体用于内部通知事件（配额超限、渠道更新/测试），带模板化内容。

### 2.8 定价与元数据
- **`pricing.go`** -- 模型元数据 DTO：`OpenAIModels`, `AnthropicModel`, `GeminiModel` 结构体用于模型列表/定价信息。
- **`ratio_sync.go`** -- 上游比例同步 DTO：`UpstreamDTO`, `UpstreamRequest`, `TestResult`, `DifferenceItem` 用于比较本地与上游模型比例。

### 2.9 其他
- **`audio.go`** -- `AudioRequest` 和 `AudioResponse` 用于 OpenAI 兼容的 TTS/STT API。
- **`embedding.go`** -- `EmbeddingRequest`/`EmbeddingResponse` 用于文本嵌入 API，带 Ollama 特定扩展。
- **`playground.go`** -- `PlayGroundRequest` 用于 API Playground 端点。
- **`realtime.go`** -- OpenAI Realtime API 类型：`RealtimeEvent`, `RealtimeSession`, `RealtimeItem`, `RealtimeAudio` 用于 WebSocket 实时音频流。
- **`sensitive.go`** -- `SensitiveResponse` 用于敏感词检测结果。

---

## 关联文件列表

### DTO 层核心文件
- `dto/openai_request.go` - OpenAI 请求结构体
- `dto/openai_response.go` - OpenAI 响应结构体
- `dto/openai_image.go` - OpenAI 图像结构体
- `dto/openai_video.go` - OpenAI 视频结构体
- `dto/openai_compaction.go` - OpenAI 压缩响应
- `dto/openai_responses_compaction_request.go` - OpenAI 压缩请求
- `dto/claude.go` - Claude 结构体
- `dto/gemini.go` - Gemini 结构体
- `dto/midjourney.go` - Midjourney 结构体
- `dto/suno.go` - Suno 结构体
- `dto/rerank.go` - Rerank 结构体
- `dto/request_common.go` - 请求通用接口
- `dto/error.go` - 错误结构体
- `dto/values.go` - 多态值类型
- `dto/task.go` - 任务响应类型
- `dto/video.go` - 视频结构体
- `dto/user_settings.go` - 用户设置
- `dto/channel_settings.go` - 渠道设置
- `dto/notify.go` - 通知结构体
- `dto/pricing.go` - 定价结构体
- `dto/ratio_sync.go` - 比例同步结构体
- `dto/audio.go` - 音频结构体
- `dto/embedding.go` - 嵌入结构体
- `dto/playground.go` - Playground 结构体
- `dto/realtime.go` - Realtime API 结构体
- `dto/sensitive.go` - 敏感词结构体

### 依赖的类型文件
- `types/relay_format.go` - 中转格式常量
- `types/error.go` - 错误类型
- `types/price_data.go` - 价格数据

### 依赖的常量文件
- `constant/finish_reason.go` - 完成原因常量
- `constant/task.go` - 任务常量
