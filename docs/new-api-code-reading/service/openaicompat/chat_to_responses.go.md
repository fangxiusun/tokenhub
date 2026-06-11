# chat_to_responses.go 代码阅读文档

## 1. 全局总结

本文件实现了 OpenAI Chat Completions API 请求格式到 Responses API 请求格式的转换逻辑。核心函数 `ChatCompletionsRequestToResponsesRequest` 将 `GeneralOpenAIRequest`（Chat 格式）转换为 `OpenAIResponsesRequest`（Responses 格式），包括消息内容的重新映射、工具调用的格式转换、响应格式的适配等。这是 API 网关兼容层的重要组成部分，使上游仅支持 Responses API 的提供商能够透明地服务 Chat Completions 请求。

## 2. 依赖关系

| 包 | 用途 |
|---|---|
| `encoding/json` | JSON 序列化（`json.RawMessage`） |
| `errors` | 错误创建 |
| `fmt` | 格式化字符串 |
| `strings` | 字符串处理（修剪、拼接等） |
| `github.com/QuantumNous/new-api/common` | JSON 封装（`Marshal`/`Unmarshal`）、类型转换（`Interface2String`）、指针工具（`GetPointer`） |
| `github.com/QuantumNous/new-api/dto` | 请求/响应 DTO 定义（`GeneralOpenAIRequest`、`OpenAIResponsesRequest`、`MessageImageUrl` 等） |
| `github.com/samber/lo` | `FromPtrOr`、`ToPtr` 等辅助函数 |

## 3. 类型定义

本文件未定义新类型，仅使用 `dto` 包中的既有类型。

## 4. 函数详解

### `normalizeChatImageURLToString(v any) any`

**签名**: `func normalizeChatImageURLToString(v any) any`

**功能**: 将图片 URL 的多种表示形式（`string`、`map[string]any`、`dto.MessageImageUrl`、`*dto.MessageImageUrl`）统一归一化为字符串。

**逻辑**: 使用类型 switch 分别处理：
- `string` → 直接返回
- `map[string]any` → 提取 `url` 字段
- `dto.MessageImageUrl` → 提取 `Url` 字段
- `*dto.MessageImageUrl` → 非空时提取 `Url` 字段
- 其他类型 → 原样返回

### `convertChatResponseFormatToResponsesText(reqFormat *dto.ResponseFormat) json.RawMessage`

**签名**: `func convertChatResponseFormatToResponsesText(reqFormat *dto.ResponseFormat) json.RawMessage`

**功能**: 将 Chat Completions 格式的 `response_format` 转换为 Responses API 的 `text.format` 结构。

**逻辑**:
- 若 `reqFormat` 为空或类型为空，返回 `nil`
- 构建 `{"type": <type>}` 基础映射
- 对 `json_schema` 类型做特殊处理：解析 schema 中的字段并展平到 format 映射中，若存在嵌套的 `json_schema` 字段则将其内容提升到顶层
- 最终包装为 `{"format": {...}}` 格式的 `json.RawMessage`

### `ChatCompletionsRequestToResponsesRequest(req *dto.GeneralOpenAIRequest) (*dto.OpenAIResponsesRequest, error)`

**签名**: `func ChatCompletionsRequestToResponsesRequest(req *dto.GeneralOpenAIRequest) (*dto.OpenAIResponsesRequest, error)`

**功能**: 核心转换函数，将 Chat Completions 请求完整转换为 Responses 请求。

**参数校验**:
- `req` 不能为 `nil`
- `req.Model` 不能为空
- `req.N` 不能大于 1（Responses API 不支持多生成）

**消息转换逻辑**:

| Chat 角色 | Responses 映射 |
|---|---|
| `system` / `developer` | 合并为 `instructions`（多个系统消息用 `\n\n` 拼接） |
| `tool` / `function` | 转换为 `function_call_output` 类型项，包含 `call_id` 和 `output`；若缺少 `call_id` 则降级为普通 user 消息 |
| `assistant` | 普通消息 + 追加 `function_call` 类型项（从 `ToolCalls` 中提取） |
| `user` | 普通消息项 |

**内容格式转换**:

| Chat 内容类型 | Responses 内容类型 |
|---|---|
| 文本（user） | `input_text` |
| 文本（assistant） | `output_text` |
| 图片 URL | `input_image` |
| 音频输入 | `input_audio` |
| 文件 | `input_file` |
| 视频 URL | `input_video` |

**工具转换**:
- `function` 类型工具直接映射字段（`name`、`description`、`parameters`）
- 未知类型工具尽量保留原始结构

**ToolChoice 转换**:
- 字符串值直接传递
- `{"type":"function","function":{"name":"..."}}` 格式转换为 `{"type":"function","name":"..."}`

**其他字段映射**:
- `max_tokens` 和 `max_completion_tokens` 取较大值作为 `max_output_tokens`
- `top_p` 保持指针语义
- `reasoning_effort` 映射为 `reasoning` 结构（`summary` 固定为 `"detailed"`）
- `parallel_tool_calls` 直接传递

## 5. 关键逻辑分析

1. **指令合并策略**: 多条 system/developer 消息被合并为单一 `instructions` 字段，使用双换行分隔。这保留了原始提示词的段落结构。

2. **工具调用的降级处理**: 当 `tool`/`function` 响应缺少 `call_id` 时，系统将其降级为普通 user 消息并附带 `[tool_output_missing_call_id]` 前缀，确保对话不会因格式不匹配而中断。

3. **图片 URL 归一化**: 通过 `normalizeChatImageURLToString` 处理多种图片 URL 表示形式，兼容上游可能传递的不同格式。

4. **n>1 拒绝**: 直接拒绝 `n>1` 的请求，因为 Responses API 不支持批量生成。

## 6. 关联文件

- `new-api/dto/request.go` — `GeneralOpenAIRequest` 定义
- `new-api/dto/response.go` — `OpenAIResponsesRequest` 定义
- `new-api/service/openaicompat/responses_to_chat.go` — 反向转换（Responses → Chat）
- `new-api/service/openaicompat/policy.go` — 策略判断是否启用转换
- `new-api/common/json.go` — JSON 工具函数
