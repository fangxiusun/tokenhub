# endpoint_type.go 代码阅读文档

## 1. 全局概述

本文件定义了端点类型（EndpointType）常量，用于标识 AI API 的不同端点格式。与 `api_type.go` 中的 API 类型不同，端点类型更关注请求/响应的协议格式。

## 2. 依赖关系

本文件无外部依赖。

## 3. 类型定义

### EndpointType 类型

```go
type EndpointType string
```

字符串类型的端点标识符，用于区分不同的 API 协议格式。

## 4. 函数详情

本文件无函数定义。

## 5. 关键逻辑分析

### 端点类型常量

| 常量名 | 值 | 说明 |
|--------|-----|------|
| `EndpointTypeOpenAI` | `"openai"` | OpenAI Chat Completions 格式 |
| `EndpointTypeOpenAIResponse` | `"openai-response"` | OpenAI Responses API 格式 |
| `EndpointTypeOpenAIResponseCompact` | `"openai-response-compact"` | OpenAI Responses API 紧凑格式 |
| `EndpointTypeAnthropic` | `"anthropic"` | Anthropic Claude Messages 格式 |
| `EndpointTypeGemini` | `"gemini"` | Google Gemini 格式 |
| `EndpointTypeJinaRerank` | `"jina-rerank"` | Jina Rerank 格式 |
| `EndpointTypeImageGeneration` | `"image-generation"` | 图像生成格式 |
| `EndpointTypeEmbeddings` | `"embeddings"` | Embeddings 格式 |
| `EndpointTypeOpenAIVideo` | `"openai-video"` | OpenAI 视频生成格式 |

### 与 API 类型的区别

- `APIType` 标识的是**服务商**（如 OpenAI、Anthropic、Gemini）
- `EndpointType` 标识的是**协议格式**（如 OpenAI 格式、Anthropic 格式）

一个服务商可能支持多种端点类型。例如，同一个渠道可能同时支持 OpenAI Chat 格式和 OpenAI Responses 格式。

### 被注释的端点类型

文件中注释掉了 Midjourney、Suno、Kling、Jimeng 的端点类型，说明这些可能使用了独立的代理协议而非标准 API 格式。

## 6. 相关文件

- `constant/api_type.go` — API 类型常量
- `relay/relay_info.go` — 根据端点类型分发请求
- `relay/adaptor.go` — 根据端点类型选择适配器
