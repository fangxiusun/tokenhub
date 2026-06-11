# adaptor.go 代码阅读文档

## 1. 全局总结
该文件是 OpenAI 渠道的主适配器实现，是整个 relay/channel 体系中最核心、最复杂的适配器。它不仅服务于 OpenAI 渠道，还作为 360、零一万物、Xinference、OpenRouter 等兼容 OpenAI 格式的渠道的底层实现。负责完整的请求生命周期管理：URL 构建、请求头设置、多种格式的请求转换（OpenAI/Claude/Gemini/音频/图像/嵌入/重排序/Responses）、HTTP 请求执行和响应分发。

## 2. 依赖关系
- `bytes`、`encoding/json`、`errors`、`fmt`、`io`、`mime/multipart`、`net/http`、`net/textproto`、`net/url`、`path/filepath`、`strings` — 标准库
- `github.com/QuantumNous/new-api/common` — 通用工具函数
- `github.com/QuantumNous/new-api/constant` — 渠道类型常量
- `github.com/QuantumNous/new-api/dto` — 数据传输对象
- `github.com/QuantumNous/new-api/logger` — 日志
- `github.com/QuantumNous/new-api/relay/channel` — 渠道公共函数
- `github.com/QuantumNous/new-api/relay/channel/ai360` — 360 渠道
- `github.com/QuantumNous/new-api/relay/channel/lingyiwanwu` — 零一万物渠道
- `github.com/QuantumNous/new-api/relay/channel/openrouter` — OpenRouter 渠道
- `github.com/QuantumNous/new-api/relay/channel/xinference` — Xinference 渠道
- `github.com/QuantumNous/new-api/relay/common` — 中继公共信息
- `github.com/QuantumNous/new-api/relay/common_handler` — 通用处理器
- `github.com/QuantumNous/new-api/relay/constant` — 中继常量
- `github.com/QuantumNous/new-api/service` — 业务逻辑
- `github.com/QuantumNous/new-api/setting/model_setting` — 模型设置
- `github.com/QuantumNous/new-api/setting/reasoning` — 推理设置
- `github.com/QuantumNous/new-api/types` — 错误类型
- `github.com/samber/lo` — 泛型工具库
- `github.com/gin-gonic/gin` — HTTP 框架

## 3. 类型定义

### `Adaptor` (struct)
- **作用**：OpenAI 渠道适配器
- **字段**：
  - `ChannelType int` — 渠道类型（OpenAI、Azure、360 等）
  - `ResponseFormat string` — 响应格式（音频相关）

## 4. 函数详解

### `(a *Adaptor) Init(info *relaycommon.RelayInfo)`
- **作用**：初始化适配器
- **逻辑**：记录 ChannelType，如果启用了 thinking_to_content，初始化 ThinkingContentInfo 状态

### `(a *Adaptor) GetRequestURL(info *relaycommon.RelayInfo) (string, error)`
- **作用**：构建请求 URL
- **关键逻辑**：
  - **Realtime 模式**：http→ws、https→wss 协议转换
  - **Azure 渠道**：构建 Azure OpenAI 格式 URL，包含 api-version 和 deployment
  - **Azure Responses API**：特殊处理 `/openai/v1/responses` 路径
  - **Custom 渠道**：替换 URL 中的 `{model}` 占位符
  - **Claude/Gemini 格式**：统一转换为 `/v1/chat/completions`
  - **默认**：使用 `GetFullRequestURL` 构建标准 URL

### `(a *Adaptor) SetupRequestHeader(c *gin.Context, header *http.Header, info *relaycommon.RelayInfo) error`
- **作用**：设置请求头
- **关键逻辑**：
  - **Azure**：使用 `api-key` 而非 `Authorization: Bearer`
  - **OpenAI**：支持 Organization 头
  - **Realtime**：处理 Sec-WebSocket-Protocol 头
  - **OpenRouter**：设置 HTTP-Referer 和 X-OpenRouter-Title
  - **Header Override**：支持自定义 Authorization 头覆盖

### `(a *Adaptor) ConvertOpenAIRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error)`
- **作用**：转换 OpenAI 请求（最复杂的转换函数）
- **关键逻辑**：
  - 非 OpenAI/Azure 渠道移除 StreamOptions
  - **OpenRouter 适配**：
    - 设置 usage.include
    - 处理 `-thinking` 后缀模型名 → reasoning 格式
    - Claude 模型的 THINKING → reasoning 转换
  - **O/GPT-5 模型适配**：
    - max_tokens → max_completion_tokens 转换
    - 移除 temperature（O 系列）/ temperature+top_p+logProbs（GPT-5 系列）
    - 推理力度后缀解析
    - system → developer 角色转换

### `(a *Adaptor) ConvertClaudeRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.ClaudeRequest) (any, error)`
- **作用**：Claude → OpenAI 格式转换
- **逻辑**：调用 `service.ClaudeToOpenAIRequest` 转换，然后调用 ConvertOpenAIRequest 进一步处理

### `(a *Adaptor) ConvertGeminiRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeminiChatRequest) (any, error)`
- **作用**：Gemini → OpenAI 格式转换
- **逻辑**：调用 `service.GeminiToOpenAIRequest` 转换，然后调用 ConvertOpenAIRequest 进一步处理

### `(a *Adaptor) ConvertAudioRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.AudioRequest) (io.Reader, error)`
- **作用**：转换音频请求
- **逻辑**：
  - TTS（Speech）：JSON 序列化
  - STT（Transcription）：multipart/form-data 构建

### `(a *Adaptor) ConvertImageRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.ImageRequest) (any, error)`
- **作用**：转换图像请求
- **逻辑**：
  - Images Edits 模式：构建 multipart/form-data，支持多图片和 mask
  - 其他模式：直接透传

### `(a *Adaptor) ConvertOpenAIResponsesRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.OpenAIResponsesRequest) (any, error)`
- **作用**：转换 Responses API 请求
- **逻辑**：解析推理力度后缀，设置 reasoning.effort

### `(a *Adaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (any, error)`
- **作用**：执行 HTTP 请求
- **逻辑**：根据模式选择 Form/WSS/API 请求方式

### `(a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage any, err *types.NewAPIError)`
- **作用**：响应分发
- **逻辑**：根据 RelayMode 分发到对应的处理器

### `(a *Adaptor) GetModelList() []string` / `GetChannelName() string`
- **作用**：根据 ChannelType 返回对应的模型列表和渠道名称

### 辅助函数
- `isJSONRequest(c *gin.Context) bool` — 检测是否为 JSON 请求
- `detectImageMimeType(filename string) string` — 根据文件扩展名检测图片 MIME 类型

## 5. 关键逻辑分析
- **多渠道复用**：Adaptor 通过 ChannelType 支持多个兼容 OpenAI 格式的渠道，GetModelList 和 GetChannelName 根据渠道类型返回对应配置
- **推理模型适配**：O 系列和 GPT-5 系列模型有特殊的参数要求（temperature 移除、max_tokens 转换、system→developer），通过 `reasoning.ParseOpenAIReasoningEffortFromModelSuffix` 解析模型名后缀中的推理力度
- **OpenRouter 深度适配**：处理 `-thinking` 后缀、Claude 模型的 THINKING 字段、reasoning 格式转换
- **Multipart 处理**：图像编辑和音频转录使用 multipart/form-data，需要手动构建 boundary 和 Part
- **协议转换**：Realtime 模式自动将 HTTP(S) 转换为 WS(S)
- **Header Override**：支持自定义 Authorization 头，避免默认设置被覆盖

## 6. 关联文件
- `relay/channel/openai/relay-openai.go` — 流式/非流式响应处理器
- `relay/channel/openai/relay_image.go` — 图像响应处理器
- `relay/channel/openai/relay_realtime.go` — Realtime WebSocket 处理器
- `relay/channel/openai/relay_responses.go` — Responses API 处理器
- `relay/channel/openai/audio.go` — 音频响应处理器
- `relay/channel/openai/helper.go` — 流式格式转换辅助
- `relay/channel/openai/usage.go` — Usage 后处理
- `relay/channel/openai/constant.go` — 模型列表
- `relay/channel/ai360/adaptor.go` — 360 渠道（复用 OpenAI Adaptor）
- `relay/channel/openrouter/adaptor.go` — OpenRouter 渠道
- `service/request_convert.go` — ClaudeToOpenAIRequest、GeminiToOpenAIRequest
