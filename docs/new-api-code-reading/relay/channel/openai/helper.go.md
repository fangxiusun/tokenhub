# helper.go 代码阅读文档

## 1. 全局总结
该文件提供了 OpenAI 渠道的辅助函数，主要负责流式响应的格式处理和 token 数据提取。核心功能包括：将 OpenAI 流式响应转换为 Claude/Gemini 格式、从流式数据中提取 token 使用量、处理最后一个响应帧和最终响应的发送。

## 2. 依赖关系
- `strings` — 标准库
- `github.com/QuantumNous/new-api/common` — 通用工具函数
- `github.com/QuantumNous/new-api/dto` — 数据传输对象
- `github.com/QuantumNous/new-api/logger` — 日志
- `github.com/QuantumNous/new-api/relay/common` — 中继公共信息
- `github.com/QuantumNous/new-api/relay/constant` — 中继常量
- `github.com/QuantumNous/new-api/relay/helper` — 流式处理辅助
- `github.com/QuantumNous/new-api/service` — 业务逻辑（格式转换、token 计数）
- `github.com/QuantumNous/new-api/types` — 错误类型
- `github.com/samber/lo` — 泛型工具库
- `github.com/gin-gonic/gin` — HTTP 框架

## 3. 类型定义
无自定义类型定义。

## 4. 函数详解

### `HandleStreamFormat(c *gin.Context, info *relaycommon.RelayInfo, data string, forceFormat bool, thinkToContent bool) error`
- **作用**：根据目标格式分发流式数据处理
- **逻辑**：根据 `RelayFormat` 选择 OpenAI/Claude/Gemini 处理路径

### `handleClaudeFormat(c *gin.Context, data string, info *relaycommon.RelayInfo) error`
- **作用**：将 OpenAI 流式响应转换为 Claude 格式
- **逻辑**：解析 OpenAI chunk → 调用 `service.StreamResponseOpenAI2Claude` → 发送 Claude 格式数据

### `handleGeminiFormat(c *gin.Context, data string, info *relaycommon.RelayInfo) error`
- **作用**：将 OpenAI 流式响应转换为 Gemini 格式
- **逻辑**：解析 OpenAI chunk → 调用 `service.StreamResponseOpenAI2Gemini` → 发送 Gemini 格式数据

### `ProcessStreamResponse(streamResponse dto.ChatCompletionsStreamResponse, responseTextBuilder *strings.Builder, toolCount *int) error`
- **作用**：从流式响应中提取文本内容和工具调用信息
- **逻辑**：累积 content、reasoning_content 和 tool_calls 的 name+arguments

### `processTokenData(relayMode int, data string, responseTextBuilder *strings.Builder, toolCount *int) error`
- **作用**：根据 RelayMode 处理流式数据的 token 提取
- **逻辑**：聊天补全和文本补全使用不同的解析结构

### `processCompletionsStreamResponse(streamResponse dto.CompletionsStreamResponse, responseTextBuilder *strings.Builder)`
- **作用**：处理文本补全的流式响应，提取 text 字段

### `handleLastResponse(lastStreamData string, ...) error`
- **作用**：解析最后一个流式数据帧，提取 responseId、createAt、systemFingerprint、model 和 usage
- **逻辑**：判断是否包含有效的 stream usage，决定是否需要发送最后一个响应

### `HandleFinalResponse(c *gin.Context, info *relaycommon.RelayInfo, ...)`
- **作用**：发送最终响应（usage 和 [DONE]）
- **逻辑**：根据 RelayFormat 分发：
  - OpenAI：发送 usage 帧（如果需要）+ [DONE]
  - Claude：发送 Claude 格式的终止响应
  - Gemini：发送 Gemini 格式的终止响应

### `sendResponsesStreamData(c *gin.Context, streamResponse dto.ResponsesStreamResponse, data string)`
- **作用**：发送 Responses API 流式数据

## 5. 关键逻辑分析
- **格式转换架构**：通过 `HandleStreamFormat` 统一入口，根据目标格式分发到不同的转换器，实现一个上游响应支持多种客户端格式
- **Token 估算**：当上游未返回 usage 时，通过累积的文本内容进行 token 计数估算
- **Claude 格式特殊性**：Claude 响应需要维护 `ClaudeConvertInfo` 状态，包括 Usage 和 Done 标志
- **Gemini 格式特殊性**：Gemini 使用 `c.Render` 直接写入 SSE 数据，不需要额外的 event/data 分隔
- **Tool Count**：跟踪工具调用数量，用于估算 completion tokens（每个工具调用约 7 tokens）

## 6. 关联文件
- `relay/channel/openai/relay-openai.go` — OaiStreamHandler 调用这些辅助函数
- `relay/channel/openai/relay_responses.go` — Responses API 处理
- `relay/helper/sse.go` — SSE 响应辅助函数（StringData、ObjectData、ClaudeData 等）
- `service/response_convert.go` — 格式转换服务
