# dto.go 代码阅读文档

## 1. 全局总结
本文件定义了 Google PaLM API 专用的数据传输对象，包括请求和响应的结构体。这些结构体用于 PaLM v1beta2 generateMessage API 的序列化和反序列化。

## 2. 依赖关系
| 包名 | 用途 |
|------|------|
| `github.com/QuantumNous/new-api/dto` | 引用 `dto.Message` 类型用于响应中的消息字段 |

## 3. 类型定义

### `PaLMChatMessage`
PaLM 聊天消息结构体。

| 字段 | 类型 | JSON 标签 | 说明 |
|------|------|-----------|------|
| `Author` | `string` | `json:"author"` | 消息作者（如 "0"、"1" 或 "user"、"assistant"） |
| `Content` | `string` | `json:"content"` | 消息内容 |

### `PaLMFilter`
PaLM 内容过滤器结构体。

| 字段 | 类型 | JSON 标签 | 说明 |
|------|------|-----------|------|
| `Reason` | `string` | `json:"reason"` | 过滤原因 |
| `Message` | `string` | `json:"message"` | 过滤消息 |

### `PaLMPrompt`
PaLM 提示结构体。

| 字段 | 类型 | JSON 标签 | 说明 |
|------|------|-----------|------|
| `Messages` | `[]PaLMChatMessage` | `json:"messages"` | 聊天消息列表 |

### `PaLMChatRequest`
PaLM 聊天请求结构体。

| 字段 | 类型 | JSON 标签 | 说明 |
|------|------|-----------|------|
| `Prompt` | `PaLMPrompt` | `json:"prompt"` | 提示内容 |
| `Temperature` | `*float64` | `json:"temperature,omitempty"` | 温度参数，控制随机性 |
| `CandidateCount` | `int` | `json:"candidateCount,omitempty"` | 生成的候选回复数量 |
| `TopP` | `float64` | `json:"topP,omitempty"` | Top-P 采样参数 |
| `TopK` | `uint` | `json:"topK,omitempty"` | Top-K 采样参数 |

### `PaLMError`
PaLM API 错误结构体。

| 字段 | 类型 | JSON 标签 | 说明 |
|------|------|-----------|------|
| `Code` | `int` | `json:"code"` | 错误码 |
| `Message` | `string` | `json:"message"` | 错误消息 |
| `Status` | `string` | `json:"status"` | 错误状态 |

### `PaLMChatResponse`
PaLM 聊天响应结构体。

| 字段 | 类型 | JSON 标签 | 说明 |
|------|------|-----------|------|
| `Candidates` | `[]PaLMChatMessage` | `json:"candidates"` | 候选回复列表 |
| `Messages` | `[]dto.Message` | `json:"messages"` | 消息列表（使用标准 dto.Message 类型） |
| `Filters` | `[]PaLMFilter` | `json:"filters"` | 内容过滤器列表 |
| `Error` | `PaLMError` | `json:"error"` | 错误信息 |

## 4. 函数详解
本文件无函数定义。

## 5. 关键逻辑分析
- **请求结构**：PaLM 的请求格式与 OpenAI 不同，使用 `prompt.messages` 包裹消息，而非直接的 `messages` 字段。
- **温度指针类型**：`Temperature` 使用 `*float64` 指针类型，允许区分"未设置"和"设置为 0"。
- **错误处理**：PaLM 的错误信息包含在响应体中（而非 HTTP 状态码），需要在响应处理时检查 `Error` 字段。
- **候选回复**：PaLM 支持返回多个候选回复（`Candidates`），这与 OpenAI 的单回复模式不同。

## 6. 关联文件
- `relay/channel/palm/adaptor.go` — 使用这些结构体进行请求/响应处理
- `relay/channel/palm/relay-palm.go` — 将 PaLM 响应转换为 OpenAI 格式
- `relay/channel/palm/constants.go` — 渠道常量定义
