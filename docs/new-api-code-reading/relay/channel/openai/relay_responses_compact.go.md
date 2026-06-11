# relay_responses_compact.go 代码阅读文档

## 1. 全局总结
该文件实现了 OpenAI Responses API 的 Compact（压缩）模式响应处理。Compact 模式是 Responses API 的一种简化响应格式，返回压缩后的响应体。此处理器解析 Compact 响应并提取 usage 信息。

## 2. 依赖关系
- `io`、`net/http` — 标准库
- `github.com/QuantumNous/new-api/common` — 通用工具函数
- `github.com/QuantumNous/new-api/dto` — 数据传输对象
- `github.com/QuantumNous/new-api/service` — 业务逻辑
- `github.com/QuantumNous/new-api/types` — 错误类型
- `github.com/gin-gonic/gin` — HTTP 框架

## 3. 类型定义
无自定义类型定义。

## 4. 函数详解

### `OaiResponsesCompactionHandler(c *gin.Context, resp *http.Response) (*dto.Usage, *types.NewAPIError)`
- **作用**：处理 Responses API 的 Compact 模式响应
- **参数**：Gin 上下文、HTTP 响应
- **返回**：使用量信息、错误
- **关键逻辑**：
  1. 读取完整响应体
  2. 反序列化为 `OpenAIResponsesCompactionResponse`
  3. 检测 OpenAI 错误
  4. 将原始响应体写入客户端
  5. 从 Compact 响应中提取 usage：
    - InputTokens → PromptTokens
    - OutputTokens → CompletionTokens
    - TotalTokens → TotalTokens
    - InputTokensDetails.CachedTokens → PromptTokensDetails.CachedTokens

## 5. 关键逻辑分析
- **Compact 模式**：Responses API 的 Compact 端点（`/responses/compact`）返回压缩格式的响应，与标准 Responses API 的 usage 字段结构一致
- **Usage 映射**：将 Responses API 的 input/output tokens 映射到标准的 prompt/completion tokens
- **缓存 Token**：支持从 InputTokensDetails 中提取 cached_tokens 用于精确计费
- **透传响应**：原始响应体直接写入客户端，不做格式转换

## 6. 关联文件
- `relay/channel/openai/adaptor.go` — DoResponse 中 RelayModeResponsesCompact 分发到此处理器
- `relay/channel/openai/relay_responses.go` — 标准 Responses API 处理器
- `dto/responses.go` — OpenAIResponsesCompactionResponse DTO 定义
