# 常量层详细设计 (`constant/`)

## 1. 概述
常量层包含整个应用程序中使用的所有枚举、配置键和全局常量。这些常量定义了系统的各种类型和配置选项。

## 2. 文件详细说明

### 2.1 API 与渠道类型
- **`api_type.go`** -- 基于 `iota` 的 API 类型枚举：
  - `APITypeOpenAI` (0), `APITypeAnthropic` (1), `APITypePaLM` (2), `APITypeBaidu` (3), `APITypeZhipu` (4), `APITypeAli` (5), `APITypeXunfei` (6), `APITypeTencent` (7), `APITypeAzure` (8), `APITypeOpenAIWS` (9), `APITypeMoonshot` (10), `APITypeDeepSeek` (11), `APITypeMiniMax` (12), `APITypeOllama` (13), `APITypeAWS` (14), `APITypeVolcEngine` (15), `APITypeGemini` (16), `APITypeClaude` (17), `APITypeCloudflare` (18), `APITypeMistral` (19), `APITypeOpenRouter` (20), `APITypeXinference` (21), `APITypeCodex` (22), `APITypeCohere` (23), `APITypeDeepSeekJina` (24), `APITypeMidjourney` (25), `APITypeSuno` (26), `APITypeVertex` (27), `APITypeBaiduV2` (28), `APITypeXAI` (29), `APITypeSiliconFlow` (30), `APITypeAliOpen` (31), `APITypePerplexity` (32), `APITypeDify` (33), `APITypeJina` (34), `APITypeLingYiWanWu` (35), `APITypeReplicate` (36), `APITypeJimeng` (37), `APITypeMokaAI` (38), `APIType360` (39)
- **`channel.go`** -- 渠道类型常量：100+ `ChannelType*` 常量（OpenAI, Azure, Anthropic, Gemini, AWS 等）和渠道状态常量（`ChannelStatusEnabled/Disabled`）。
- **`endpoint_type.go`** -- `EndpointType` 字符串常量：`openai`, `openai-response`, `openai-response-compact`, `anthropic`, `gemini`, `jina-rerank`, `image-generation`, `embeddings`, `openai-video`。

### 2.2 上下文与缓存
- **`context_key.go`** -- 用于通过 Go 上下文传递数据的类型化 `ContextKey` 常量：令牌信息、渠道信息、请求元数据、中转信息、计费数据和用户信息（70+ 键）。
- **`cache_key.go`** -- Redis 缓存键格式字符串：`UserGroupKeyFmt`, `UserQuotaKeyFmt`, `UserEnabledKeyFmt`, `UserUsernameKeyFmt` 和令牌字段名常量。

### 2.3 环境
- **`env.go`** -- 全局环境变量持有者：`StreamingTimeout`, `DifyDebug`, `MaxFileDownloadMB`, `ForceStreamOption`, `CountToken`, `MaxRequestBodyMB`, `AzureDefaultAPIVersion`, `TrustedRedirectDomains` 等。

### 2.4 中转
- **`relay/constant/relay_mode.go`** -- `RelayMode*` 常量：所有中转操作模式（聊天补全、嵌入、图像、Midjourney 操作、音频 TTS/STT、重排序、视频、Suno、实时、响应等）以及每种模式的 HTTP 方法/路径辅助。

### 2.5 功能
- **`midjourney.go`** -- Midjourney 操作常量（`IMAGINE`, `DESCRIBE`, `BLEND`, `UPSCALE` 等）和错误码，以及模型到操作的映射表。
- **`task.go`** -- 任务平台常量（`suno`, `mj`）、任务操作常量（`generate`, `textGenerate`, `remixGenerate` 等）和 Suno 模型到操作的映射。
- **`finish_reason.go`** -- LLM 完成原因字符串：`stop`, `tool_calls`, `length`, `function_call`, `content_filter`。
- **`multi_key_mode.go`** -- `MultiKeyMode` 类型：`random` 和 `polling` 模式，用于具有多个 API 密钥的渠道。

### 2.6 支付
- **`waffo_pay_method.go`** -- `WaffoPayMethod` 结构体和 `DefaultWaffoPayMethods` 列表，用于 Waffo 支付方式显示/API 映射（Card, Apple Pay, Google Pay）。

### 2.7 其他
- **`azure.go`** -- Azure 特定常量：`AzureNoRemoveDotTime` -- 用于 Azure 模型名称点移除行为的截止时间戳。
- **`setup.go`** -- 单个布尔值 `Setup` 标志，指示系统是否已通过设置向导初始化。

---

## 关联文件列表

### 常量层核心文件
- `constant/api_type.go` - API 类型枚举
- `constant/channel.go` - 渠道类型常量
- `constant/endpoint_type.go` - 端点类型常量
- `constant/context_key.go` - 上下文键常量
- `constant/cache_key.go` - 缓存键常量
- `constant/env.go` - 环境变量常量
- `constant/midjourney.go` - Midjourney 常量
- `constant/task.go` - 任务常量
- `constant/finish_reason.go` - 完成原因常量
- `constant/multi_key_mode.go` - 多密钥模式
- `constant/waffo_pay_method.go` - Waffo 支付方式
- `constant/azure.go` - Azure 常量
- `constant/setup.go` - 设置状态

### 中转常量
- `relay/constant/relay_mode.go` - 中转模式常量

### 依赖的类型文件
- `types/relay_format.go` - 中转格式常量
- `types/error.go` - 错误类型
