# chat_via_responses.go 代码阅读文档

## 1. 全局总结
该文件实现了将 OpenAI Responses API 的响应转换为 Chat Completions API 格式的功能，包括非流式和流式两种模式。这是为了让使用 Chat Completions 格式的客户端能够透明地使用 Responses API 后端，同时保持 OpenAI、Claude 和 Gemini 三种输出格式的兼容。

## 2. 依赖关系
- `fmt`、`io`、`net/http`、`strings`、`time` — 标准库
- `github.com/QuantumNous/new-api/common` — 通用工具函数
- `github.com/QuantumNous/new-api/dto` — 数据传输对象
- `github.com/QuantumNous/new-api/logger` — 日志
- `github.com/QuantumNous/new-api/relay/common` — 中继公共信息
- `github.com/QuantumNous/new-api/relay/helper` — 流式处理辅助
- `github.com/QuantumNous/new-api/service` — 业务逻辑（格式转换）
- `github.com/QuantumNous/new-api/types` — 错误类型
- `github.com/gin-gonic/gin` — HTTP 框架

## 3. 类型定义
无自定义类型定义。

## 4. 函数详解

### `responsesStreamIndexKey(itemID string, idx *int) string`
- **作用**：生成 Responses 流中用于跟踪索引的 key

### `stringDeltaFromPrefix(prev string, next string) string`
- **作用**：计算两个累积字符串之间的增量（delta）
- **逻辑**：如果 next 以 prev 开头，返回差异部分；否则返回整个 next

### `OaiResponsesToChatHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)`
- **作用**：非流式处理 — 将 Responses 响应转换为 Chat Completions 响应
- **关键逻辑**：
  1. 读取并解析 Responses 响应
  2. 调用 `service.ResponsesResponseToChatCompletionsResponse` 进行格式转换
  3. 如果 usage 为空或为零，通过输出文本估算 token
  4. 根据 RelayFormat 转换为 Claude/Gemini/OpenAI 格式

### `OaiResponsesToChatStreamHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)`
- **作用**：流式处理 — 将 Responses 流事件转换为 Chat Completions 流
- **关键逻辑**：
  1. 使用 StreamScannerHandler 处理流式事件
  2. 事件类型分发：
     - `response.created` — 提取模型和创建时间
     - `response.reasoning_summary_text.delta` — 推理摘要增量 → reasoning_content
     - `response.output_text.delta` — 输出文本增量 → content
     - `response.output_item.added/done` — 函数调用处理
     - `response.function_call_arguments.delta` — 函数参数增量
     - `response.completed` — 使用量统计，发送 stop 信号
     - `response.error/failed` — 错误处理
  3. 支持 tool_calls 的增量流式传输（按 ID 跟踪索引和名称）
  4. 最终发送 usage 和 [DONE]

### `sendToolCallDelta(callID string, name string, argsDelta string) bool` (闭包)
- **作用**：发送工具调用增量到客户端
- **逻辑**：维护 toolCallIndexByID 映射，确保每个工具调用有正确的索引和名称

### `sendReasoningSummaryDelta(delta string) bool` (闭包)
- **作用**：发送推理摘要增量，转换为 reasoning_content 格式
- **逻辑**：处理摘要之间的分隔符（双换行）

## 5. 关键逻辑分析
- **双格式桥接**：将 Responses API 的事件驱动格式转换为 Chat Completions 的 delta 流格式
- **工具调用状态机**：通过多个 map（toolCallIndexByID、toolCallNameByID、toolCallArgsByID）跟踪每个工具调用的状态，支持增量参数传输
- **推理摘要处理**：Responses API 的 `reasoning_summary_text` 事件转换为 Chat Completions 的 `reasoning_content` 字段，支持分隔符插入
- **Fallback 估算**：当流中未包含 usage 时，通过累积的输出文本进行 token 估算
- **多格式输出**：支持 OpenAI/Claude/Gemini 三种输出格式，Claude 格式需要额外初始化 ClaudeConvertInfo
- **Finish Reason**：如果有工具调用但无文本输出，finish_reason 设为 `tool_calls`

## 6. 关联文件
- `relay/channel/openai/adaptor.go` — DoResponse 中分发到此处理器
- `relay/channel/openai/helper.go` — HandleStreamFormat 用于格式转换
- `relay/channel/openai/relay_responses.go` — 原生 Responses API 处理
- `service/response_convert.go` — ResponsesResponseToChatCompletionsResponse 转换
- `dto/responses.go` — Responses API 的 DTO 定义
