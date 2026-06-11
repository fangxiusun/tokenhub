# relay_info.go 代码阅读文档

## 1. 全局总结

本文件定义了 Relay 模块的核心上下文结构 `RelayInfo`，以及相关的辅助类型和工厂函数。`RelayInfo` 承载了请求处理过程中的所有状态信息，包括用户信息、渠道信息、计费数据、流式状态等。

## 2. 依赖关系

- `dto`: 请求/响应 DTO
- `types`: RelayFormat、PriceData
- `pkg/billingexpr`: 分层计费快照
- `constant`: 上下文键名
- `gin`: HTTP 框架
- `gorilla/websocket`: WebSocket

## 3. 类型定义

### `RelayInfo`
核心上下文结构，包含：
- **用户信息**: UserId, UserGroup, UserQuota, UserEmail
- **令牌信息**: TokenId, TokenKey, TokenGroup, TokenUnlimited
- **请求信息**: RelayMode, RelayFormat, OriginModelName, RequestURLPath, IsStream
- **渠道信息**: ChannelMeta（嵌入）
- **计费信息**: PriceData, Billing, BillingSource, SubscriptionId
- **流式状态**: StreamStatus, ThinkingContentInfo
- **任务信息**: TaskRelayInfo（嵌入）

### `ChannelMeta`
渠道元数据：
- ChannelType, ChannelId, ApiType, ApiVersion, ApiKey
- ParamOverride, HeadersOverride
- ChannelSetting, ChannelOtherSettings
- UpstreamModelName, IsModelMapped, SupportStreamOptions

### `TaskRelayInfo`
任务上下文：Action, OriginTaskID, PublicTaskID, LockedChannel

### `TaskSubmitReq`
任务提交请求：Prompt, Model, Mode, Image, Images, Size, Duration, Metadata

### `TaskInfo`
任务结果信息：Code, TaskID, Status, Url, Progress, CompletionTokens

### 其他辅助类型
- `ThinkingContentInfo`: 思考内容状态
- `ClaudeConvertInfo`: Claude 转换信息
- `RerankerInfo`: Rerank 信息
- `ResponsesUsageInfo`: Responses 使用量信息
- `TokenCountMeta`: Token 计数元数据

## 4. 工厂函数

- `GenRelayInfo(c, relayFormat, request, ws)`: 统一入口
- `GenRelayInfoOpenAI(c, request)`: OpenAI 格式
- `GenRelayInfoClaude(c, request)`: Claude 格式
- `GenRelayInfoGemini(c, request)`: Gemini 格式
- `GenRelayInfoWs(c, ws)`: WebSocket 格式
- `GenRelayInfoResponses(c, request)`: Responses 格式
- `GenRelayInfoEmbedding(c, request)`: Embedding 格式
- `GenRelayInfoRerank(c, request)`: Rerank 格式
- `GenRelayInfoImage(c, request)`: Image 格式
- `GenRelayInfoOpenAIAudio(c, request)`: Audio 格式

## 5. 关键逻辑分析

1. **InitChannelMeta**: 从 gin.Context 初始化渠道元数据，设置 API 类型、密钥、覆盖规则等
2. **RequestConversionChain**: 记录请求格式转换链（如 openai → openai_responses）
3. **RemoveDisabledFields**: 移除渠道设置中禁用的字段（service_tier, inference_geo, speed, store 等）
4. **RemoveGeminiDisabledFields**: 移除 Gemini 特定的禁用字段（functionResponse.id）
5. **流式支持渠道**: 定义了支持 StreamOptions 的渠道类型列表

## 6. 关联文件

- `relay/channel/adapter.go`: Adaptor 接口
- `relay/common/billing.go`: BillingSettler 接口
- `types/relay_format.go`: RelayFormat 类型
