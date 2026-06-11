# responses_to_chat.go 代码阅读文档

## 1. 全局总结

本文件实现了 OpenAI Responses API 响应格式到 Chat Completions 响应格式的反向转换。核心函数 `ResponsesResponseToChatCompletionsResponse` 将 `OpenAIResponsesResponse` 转换为 `OpenAITextResponse`，包括文本提取、用量统计映射、工具调用重建等。与 `chat_to_responses.go` 构成双向转换对。

## 2. 依赖关系

| 包 | 用途 |
|---|---|
| `errors` | 错误创建 |
| `strings` | 字符串处理（修剪、构建器） |
| `github.com/QuantumNous/new-api/dto` | 响应 DTO 定义（`OpenAIResponsesResponse`、`OpenAITextResponse`、`Usage`、`ToolCallResponse` 等） |

## 3. 类型定义

本文件未定义新类型。

## 4. 函数详解

### `ResponsesResponseToChatCompletionsResponse(resp, id) (*dto.OpenAITextResponse, *dto.Usage, error)`

**签名**:
```go
func ResponsesResponseToChatCompletionsResponse(
    resp *dto.OpenAIResponsesResponse,
    id string,
) (*dto.OpenAITextResponse, *dto.Usage, error)
```

**功能**: 将 Responses API 响应转换为 Chat Completions 标准响应格式。

**逻辑**:
1. **空值校验**: `resp` 不能为 `nil`
2. **文本提取**: 调用 `ExtractOutputTextFromResponses` 提取输出文本
3. **用量映射**:

| Responses 字段 | Chat Completions 字段 |
|---|---|
| `InputTokens` | `PromptTokens` / `InputTokens` |
| `OutputTokens` | `CompletionTokens` / `OutputTokens` |
| `TotalTokens` | `TotalTokens`（若无则自动计算） |
| `InputTokensDetails.CachedTokens` | `PromptTokensDetails.CachedTokens` |
| `InputTokensDetails.ImageTokens` | `PromptTokensDetails.ImageTokens` |
| `InputTokensDetails.AudioTokens` | `PromptTokensDetails.AudioTokens` |
| `CompletionTokenDetails.ReasoningTokens` | `CompletionTokenDetails.ReasoningTokens` |

4. **工具调用重建**: 当无文本输出但有 `function_call` 类型的 output 时，重建 `ToolCallResponse` 列表：
   - 优先使用 `CallId`，若为空则使用 `ID`
   - 提取函数名和参数字符串

5. **完成原因判断**:
   - 有工具调用 → `"tool_calls"`
   - 无工具调用 → `"stop"`

6. **消息构建**: 组装 `dto.Message`，若有工具调用则清空 `Content`

7. **返回**: 构造完整的 `OpenAITextResponse` 并返回

### `ExtractOutputTextFromResponses(resp) string`

**签名**: `func ExtractOutputTextFromResponses(resp *dto.OpenAIResponsesResponse) string`

**功能**: 从 Responses 响应中提取输出文本。

**逻辑**:
1. **优先级策略**: 首先遍历 `output` 数组，筛选 `type == "message"` 且 `role` 为空或 `"assistant"` 的项
2. 从这些项的 `content` 中提取 `type == "output_text"` 的文本内容
3. **降级策略**: 若优先级策略未产出文本，则遍历所有 output 项的所有 content，提取任何非空 `Text`
4. 使用 `strings.Builder` 高效拼接多段文本

## 5. 关键逻辑分析

1. **双向对称**: 本文件与 `chat_to_responses.go` 构成完整的双向转换对。Chat → Responses 在 `ChatCompletionsRequestToResponsesRequest` 中完成，反向在本文件中完成。

2. **文本优先级**: `ExtractOutputTextFromResponses` 实现了明确的优先级策略——优先提取 assistant 消息类型的输出文本，确保在混合输出（文本+工具调用）场景下返回正确的文本内容。

3. **Usage 双向兼容**: 同时设置 `PromptTokens`/`InputTokens` 和 `CompletionTokens`/`OutputTokens`，确保下游无论使用哪种命名约定都能正确读取。

4. **工具调用容错**: 对 `CallId` 和 `ID` 的回退处理，以及对空名称的跳过，体现了对上游响应格式差异的兼容。

## 6. 关联文件

- `new-api/service/openaicompat/chat_to_responses.go` — 反向转换（Chat → Responses）
- `new-api/dto/response.go` — `OpenAIResponsesResponse`、`OpenAITextResponse` 定义
- `new-api/dto/types.go` — `ToolCallResponse`、`Usage` 等类型定义
