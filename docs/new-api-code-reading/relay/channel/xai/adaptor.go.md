# adaptor.go 代码阅读文档

## 1. 全局总结
xAI（Grok）渠道的适配器实现文件。支持聊天补全、图片生成、响应式 API，以及搜索增强和推理努力级别等特殊功能。

## 2. 依赖关系
- **标准库**: errors, io, net/http, strings
- **内部包**:
  - `github.com/QuantumNous/new-api/dto` — 数据传输对象
  - `github.com/QuantumNous/new-api/relay/channel` — 渠道通用工具
  - `github.com/QuantumNous/new-api/relay/channel/openai` — OpenAI 渠道适配器
  - `github.com/QuantumNous/new-api/relay/common` — RelayInfo
  - `github.com/QuantumNous/new-api/relay/constant` — 中继模式常量
  - `github.com/QuantumNous/new-api/types` — 错误类型
- **外部依赖**:
  - `github.com/gin-gonic/gin` — HTTP 框架
  - `github.com/samber/lo` — 泛型工具

## 3. 类型定义

### Adaptor 结构体
```go
type Adaptor struct {}
```
空结构体。

## 4. 函数详解

### ConvertOpenAIRequest(...)
核心转换逻辑，处理两种特殊场景：

1. **搜索增强模式**（模型名以 `-search` 结尾）：
   - 移除 `-search` 后缀
   - 在请求中添加 `search_parameters: {"mode": "on"}`

2. **Grok-3-mini 推理模式**：
   - 处理 `MaxCompletionTokens` 和 `MaxTokens` 的兼容性
   - 支持 `-high` 和 `-low` 后缀设置推理努力级别

### ConvertImageRequest(...)
将 OpenAI 图片请求转换为 xAI 格式，支持 `model`、`prompt`、`n`、`response_format`。

### GetRequestURL(...)
使用 `relaycommon.GetFullRequestURL` 构建请求 URL。

### SetupRequestHeader(...)
设置标准 Bearer Token 认证头。

### DoResponse(...)
根据中继模式分发：
- 图片生成/编辑: `openai.OpenaiImageHandler`
- 响应式: `openai.OaiResponsesStreamHandler` 或 `openai.OaiResponsesHandler`
- 其他: `xAIStreamHandler` 或 `xAIHandler`

### ConvertOpenAIResponsesRequest(...)
响应式 API 请求转换，自动设置模型名称。

### 未实现的方法
- `ConvertGeminiRequest` / `ConvertClaudeRequest` / `ConvertAudioRequest` / `ConvertEmbeddingRequest`

## 5. 关键逻辑分析

1. **搜索增强**: xAI 的搜索功能通过在模型名后添加 `-search` 后缀触发，转换时自动处理。

2. **推理努力级别**: Grok-3-mini 支持 `-high` 和 `-low` 后缀控制推理深度，通过 `ReasoningEffort` 字段传递。

3. **MaxTokens 兼容性**: Grok-3-mini 使用 `MaxCompletionTokens` 而非 `MaxTokens`，适配器自动处理两者转换。

4. **响应式 API**: xAI 支持 OpenAI 的 Responses API 格式，通过 OpenAI 适配器处理。

## 6. 关联文件
- `xai/constants.go` — 模型列表
- `xai/dto.go` — 请求/响应数据结构
- `xai/text.go` — 文本响应处理器
- `relay/channel/openai/` — OpenAI 格式处理
