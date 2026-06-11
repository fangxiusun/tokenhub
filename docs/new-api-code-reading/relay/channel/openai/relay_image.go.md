# relay_image.go 代码阅读文档

## 1. 全局总结
该文件负责处理 OpenAI 的图像生成和编辑 API 响应，支持非流式、SSE 流式和 JSON-as-Stream 三种响应模式。包含完整的图像流式响应重建逻辑，能够将上游响应转换为标准的 SSE 格式输出。

## 2. 依赖关系
- `encoding/json`、`fmt`、`io`、`net/http`、`strings`、`time` — 标准库
- `github.com/QuantumNous/new-api/common` — 通用工具函数
- `github.com/QuantumNous/new-api/dto` — 数据传输对象
- `github.com/QuantumNous/new-api/logger` — 日志
- `github.com/QuantumNous/new-api/relay/common` — 中继公共信息
- `github.com/QuantumNous/new-api/relay/helper` — 流式处理辅助
- `github.com/QuantumNous/new-api/service` — 业务逻辑
- `github.com/QuantumNous/new-api/types` — 错误类型
- `github.com/gin-gonic/gin` — HTTP 框架

## 3. 类型定义
无自定义类型定义。

## 4. 函数详解

### `OpenaiImageHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)`
- **作用**：处理非流式图像响应（generations/edits）
- **逻辑**：读取响应 → 解析 usage → normalizeOpenAIUsage → applyUsagePostProcessing → 写入客户端

### `normalizeOpenAIUsage(usage *dto.Usage)`
- **作用**：将 OpenAI Images API 的 usage 格式（input_tokens/output_tokens）映射到标准字段（prompt_tokens/completion_tokens）
- **关键字段映射**：InputTokens→PromptTokens、OutputTokens→CompletionTokens、InputTokensDetails→PromptTokensDetails

### `OpenaiImageStreamHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)`
- **作用**：处理图像流式响应
- **逻辑**：
  1. 检查 Content-Type，区分 SSE 流和 JSON-as-Stream
  2. 错误状态码回退到非流式处理
  3. SSE 流：使用 StreamScannerHandler 逐帧处理
  4. 检测错误事件（error/upstream_error）
  5. 提取 usage 信息
  6. 重建 SSE 帧（event + data）

### `writeOpenaiImageStreamChunk(c *gin.Context, data []byte)`
- **作用**：重建 SSE 帧，从 JSON 的 type 字段提取 event 名称

### `isOpenAIImageStreamErrorEvent(data []byte) bool`
- **作用**：检测图像流中的错误事件
- **逻辑**：检查 type 为 "error"/"upstream_error" 或存在 error 字段

### `extractOpenAIImageStreamErrorMessage(data []byte) string`
- **作用**：提取错误事件的错误消息
- **逻辑**：依次尝试 message 字段、error.message 字段、error 原始内容

### `OpenaiImageJSONAsStreamHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)`
- **作用**：将 JSON 格式的图像响应转换为 SSE 流格式
- **逻辑**：
  1. 读取并解析 JSON 响应
  2. 设置 SSE 头
  3. 遍历每张图片，发送 `image_generation.completed` 事件
  4. 发送 `[DONE]` 终止信号

### `writeOpenaiImageStreamPayload(c *gin.Context, eventName string, payload any) error`
- **作用**：写入 SSE 事件帧（event + data）

### `writeOpenaiImageStreamDone(c *gin.Context) error`
- **作用**：写入 `[DONE]` 终止信号

## 5. 关键逻辑分析
- **三种响应模式**：非流式（JSON）、SSE 流式（text/event-stream）、JSON-as-Stream（JSON 但需要转为 SSE）
- **Usage 标准化**：OpenAI 图像 API 使用 input_tokens/output_tokens 而非 prompt_tokens/completion_tokens，需要归一化处理
- **错误检测**：通过 JSON 内容（而非 SSE event 行）检测错误事件，因为 StreamScanner 只传递 data 行
- **SSE 重建**：从 JSON 的 type 字段重建 event 行，保持与 OpenAI 官方格式一致
- **客户端断开检测**：通过 StreamStatus 的 EndReason 跟踪客户端连接状态

## 6. 关联文件
- `relay/channel/openai/adaptor.go` — DoResponse 中分发到图像处理器
- `relay/channel/openai/usage.go` — applyUsagePostProcessing usage 后处理
- `relay/helper/stream_scanner.go` — StreamScannerHandler 流式扫描器
- `dto/image.go` — ImageResponse、SimpleResponse 等 DTO 定义
