# openai_chat_responses_compat.go 代码阅读文档

## 1. 全局总结

该文件提供 OpenAI Chat Completions 和 Responses API 之间的兼容转换层。是 `openaicompat` 包的薄包装。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `dto` | 请求/响应结构体 |
| `openaicompat` | 实际转换逻辑 |

## 3. 函数详解

### `ChatCompletionsRequestToResponsesRequest(req) (*dto.OpenAIResponsesRequest, error)`
Chat Completions 请求转 Responses 请求

### `ResponsesResponseToChatCompletionsResponse(resp, id) (*dto.OpenAITextResponse, *dto.Usage, error)`
Responses 响应转 Chat Completions 响应

### `ExtractOutputTextFromResponses(resp) string`
从 Responses 响应中提取文本输出

## 4. 关联文件

- `service/openaicompat/` — 实际转换实现
