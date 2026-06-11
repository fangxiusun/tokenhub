# relay_responses.go 代码阅读文档

## 1. 全局总结
该文件实现了 OpenAI Responses API 的原生响应处理，包括非流式和流式两种模式。非流式处理器解析完整的 Responses 响应并提取 usage；流式处理器处理 SSE 事件流，支持文本增量、工具调用（web_search_call）和完成事件，同时处理图像生成调用的元数据。

## 2. 依赖关系
- `fmt`、`io`、`net/http`、`strings` — 标准库
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

### `OaiResponsesHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)`
- **作用**：处理非流式 Responses 响应
- **关键逻辑**：
  1. 读取并解析响应体
  2. 检测图像生成调用（设置上下文标记）
  3. 写入客户端
  4. 从 usage 中提取 input_tokens/output_tokens
  5. 解析内置工具使用量（web_search 等）

### `OaiResponsesStreamHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)`
- **作用**：处理流式 Responses 响应
- **关键逻辑**：
  1. 使用 StreamScannerHandler 处理 SSE 流
  2. 事件类型分发：
    - `response.completed` — 提取 usage，检测图像生成调用
    - `response.output_text.delta` — 累积输出文本
    - `response.output_item.done` — 处理 web_search_call 工具调用计数
  3. 如果流中未包含 usage，通过累积文本估算 completion tokens
  4. 如果 prompt_tokens 为 0，使用估算值
  5. 发送 Responses 格式的流数据

## 5. 关键逻辑分析
- **图像生成检测**：通过 `HasImageGenerationCall()` 检测响应中是否包含图像生成调用，并记录 quality 和 size 信息
- **内置工具计数**：遍历响应中的 tools 字段，匹配 `BuiltInTools` 配置中的工具类型，累加 CallCount 用于计费
- **Web Search 特殊处理**：在流式模式下，检测 `response.output_item.done` 事件中的 `web_search_call` 类型，单独计数
- **Usage Fallback**：当上游未返回有效 usage 时，通过 `CountTextToken` 本地估算 completion tokens
- **Prompt Token 估算**：如果 completion tokens 有值但 prompt tokens 为 0，使用 `GetEstimatePromptTokens` 估算

## 6. 关联文件
- `relay/channel/openai/adaptor.go` — DoResponse 中分发到 Responses 处理器
- `relay/channel/openai/chat_via_responses.go` — Responses→Chat 格式转换
- `relay/channel/openai/helper.go` — sendResponsesStreamData 辅助函数
- `dto/responses.go` — OpenAIResponsesResponse、ResponsesStreamResponse DTO 定义
- `relay/helper/stream_scanner.go` — StreamScannerHandler 流式扫描器
