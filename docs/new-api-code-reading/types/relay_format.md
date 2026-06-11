# relay_format.go 代码阅读文档

## 1. 全局概述

本文件定义了中继格式（RelayFormat）常量，用于标识 AI API 请求的协议格式。这些常量在路由分发和请求转换中使用。

## 2. 依赖关系

本文件无外部依赖。

## 3. 类型定义

### RelayFormat 类型

```go
type RelayFormat string
```

字符串类型的中继格式标识符。

## 4. 函数详情

本文件无函数定义。

## 5. 关键逻辑分析

### 中继格式常量

| 常量名 | 值 | 说明 |
|--------|-----|------|
| `RelayFormatOpenAI` | `"openai"` | OpenAI Chat Completions 格式 |
| `RelayFormatClaude` | `"claude"` | Anthropic Claude 格式 |
| `RelayFormatGemini` | `"gemini"` | Google Gemini 格式 |
| `RelayFormatOpenAIResponses` | `"openai_responses"` | OpenAI Responses API 格式 |
| `RelayFormatOpenAIResponsesCompaction` | `"openai_responses_compaction"` | OpenAI Responses API 紧凑格式 |
| `RelayFormatOpenAIAudio` | `"openai_audio"` | OpenAI Audio 格式 |
| `RelayFormatOpenAIImage` | `"openai_image"` | OpenAI Image 格式 |
| `RelayFormatOpenAIRealtime` | `"openai_realtime"` | OpenAI Realtime 格式 |
| `RelayFormatRerank` | `"rerank"` | Rerank 格式 |
| `RelayFormatEmbedding` | `"embedding"` | Embedding 格式 |
| `RelayFormatTask` | `"task"` | 任务格式 |
| `RelayFormatMjProxy` | `"mj_proxy"` | Midjourney 代理格式 |

### 与 EndpointType 的区别

- `RelayFormat` 关注的是**请求的协议格式**，用于路由分发
- `EndpointType` 关注的是**端点的协议类型**，用于适配器选择

### 使用场景

- 路由层根据 URL 路径确定 `RelayFormat`
- 中继层根据 `RelayFormat` 选择对应的请求转换器
- 不同格式的请求使用不同的解析和序列化逻辑

## 6. 相关文件

- `constant/endpoint_type.go` — 端点类型常量
- `relay/relay_info.go` — 中继信息中使用格式常量
- `router/` — 路由层根据路径设置格式
