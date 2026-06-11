# api_type.go 代码阅读文档

## 1. 全局概述

本文件定义了系统支持的所有 API 类型常量，使用 Go 的 `iota` 枚举机制。这些常量用于标识不同的 AI 服务商接口类型，在路由分发、渠道管理和中继请求处理中广泛使用。

## 2. 依赖关系

本文件无外部依赖，仅使用 Go 内置的 `iota` 常量生成器。

## 3. 类型定义

本文件未定义自定义类型，常量以 `int` 类型通过 `iota` 自动生成。

## 4. 函数详情

本文件无函数定义。

## 5. 关键逻辑分析

### API 类型常量

使用 `iota` 从 0 开始自增，定义了 34 种 API 类型：

| 常量名 | 值 | 对应服务商 |
|--------|-----|-----------|
| `APITypeOpenAI` | 0 | OpenAI |
| `APITypeAnthropic` | 1 | Anthropic (Claude) |
| `APITypePaLM` | 2 | Google PaLM |
| `APITypeBaidu` | 3 | 百度文心一言 |
| `APITypeZhipu` | 4 | 智谱 AI |
| `APITypeAli` | 5 | 阿里通义千问 |
| `APITypeXunfei` | 6 | 讯飞星火 |
| `APITypeAIProxyLibrary` | 7 | AI Proxy Library |
| `APITypeTencent` | 8 | 腾讯混元 |
| `APITypeGemini` | 9 | Google Gemini |
| `APITypeZhipuV4` | 10 | 智谱 V4 |
| `APITypeOllama` | 11 | Ollama (本地) |
| `APITypePerplexity` | 12 | Perplexity |
| `APITypeAws` | 13 | AWS Bedrock |
| `APITypeCohere` | 14 | Cohere |
| `APITypeDify` | 15 | Dify |
| `APITypeJina` | 16 | Jina |
| `APITypeCloudflare` | 17 | Cloudflare Workers AI |
| `APITypeSiliconFlow` | 18 | SiliconFlow |
| `APITypeVertexAi` | 19 | Google Vertex AI |
| `APITypeMistral` | 20 | Mistral AI |
| `APITypeDeepSeek` | 21 | DeepSeek |
| `APITypeMokaAI` | 22 | Moka AI |
| `APITypeVolcEngine` | 23 | 火山引擎 |
| `APITypeBaiduV2` | 24 | 百度 V2 |
| `APITypeOpenRouter` | 25 | OpenRouter |
| `APITypeXinference` | 26 | Xinference |
| `APITypeXai` | 27 | xAI (Grok) |
| `APITypeCoze` | 28 | Coze |
| `APITypeJimeng` | 29 | 即梦 |
| `APITypeMoonshot` | 30 | Moonshot (Kimi) |
| `APITypeSubmodel` | 31 | Submodel |
| `APITypeMiniMax` | 32 | MiniMax |
| `APITypeReplicate` | 33 | Replicate |
| `APITypeCodex` | 34 | Codex |
| `APITypeDummy` | 35 | 占位符（仅用于计数） |

### `APITypeDummy` 的作用

`APITypeDummy` 是哨兵值，仅用于统计 API 类型总数。注释明确指示不要在其后添加新的 API 类型。

## 6. 相关文件

- `constant/channel.go` — 渠道类型常量，与 API 类型一一对应
- `constant/endpoint_type.go` — 端点类型定义
- `relay/channel/` — 各服务商适配器，根据 API 类型分发请求
