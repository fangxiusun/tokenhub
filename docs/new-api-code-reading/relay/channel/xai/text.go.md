# text.go 代码阅读文档

## 1. 全局总结
xAI 渠道的文本响应处理文件，包含流式和非流式响应处理器，以及响应格式转换逻辑。

## 2. 依赖关系
- **标准库**: io, net/http, strings
- **内部包**:
  - `github.com/QuantumNous/new-api/common` — JSON 处理、日志
  - `github.com/QuantumNous/new-api/dto` — 数据传输对象
  - `github.com/QuantumNous/new-api/relay/channel/openai` — OpenAI 流式处理
  - `github.com/QuantumNous/new-api/relay/common` — RelayInfo
  - `github.com/QuantumNous/new-api/relay/helper` — 流式处理工具
  - `github.com/QuantumNous/new-api/service` — 使用量计算、响应体处理
  - `github.com/QuantumNous/new-api/types` — 错误类型
- **外部依赖**:
  - `github.com/gin-gonic/gin` — HTTP 框架

## 3. 类型定义
无新增类型定义。

## 4. 函数详解

### streamResponseXAI2OpenAI(xAIResp *dto.ChatCompletionsStreamResponse, usage *dto.Usage) *dto.ChatCompletionsStreamResponse
将 xAI 流式响应转换为 OpenAI 格式。主要处理 `Usage` 字段中的 `CompletionTokens`。

### xAIStreamHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)
流式响应处理器：
1. 设置 SSE 响应头
2. 使用 `helper.StreamScannerHandler` 扫描流数据
3. 对每个数据块：反序列化 → 提取 Usage → 转换为 OpenAI 格式 → 发送
4. 如果流中没有 Usage，则从响应文本估算使用量
5. 考虑工具调用的 token 开销（每个工具调用约 7 token）

### xAIHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)
非流式响应处理器：
1. 读取完整响应体
2. 反序列化为 `ChatCompletionResponse`
3. 计算 `CompletionTokens`（TotalTokens - PromptTokens）
4. 计算 `TextTokens`（CompletionTokens - ReasoningTokens）
5. 重新编码并写回客户端

## 5. 关键逻辑分析

1. **Usage 计算**: xAI 的 Usage 可能不包含 `CompletionTokens`，需要通过 `TotalTokens - PromptTokens` 计算。

2. **ReasoningTokens 处理**: xAI 的响应包含 `ReasoningTokens`，需要从 `CompletionTokens` 中减去得到实际文本 token 数。

3. **流式 Usage 估算**: 如果流式响应中没有 Usage 信息，使用 `service.ResponseText2Usage` 从文本内容估算，并额外加上工具调用的 token。

4. **工具调用开销**: 每个工具调用约 7 token，在没有流式 Usage 时需要手动添加。

5. **响应透传**: 非流式响应重新编码后直接写回客户端，保持原有状态码。

## 6. 关联文件
- `xai/adaptor.go` — 调用 `xAIStreamHandler` 和 `xAIHandler`
- `xai/dto.go` — `ChatCompletionResponse` 结构体定义
- `relay/channel/openai/` — OpenAI 流式处理工具
- `relay/helper/sse.go` — SSE 流处理工具
- `service/response.go` — 使用量计算
