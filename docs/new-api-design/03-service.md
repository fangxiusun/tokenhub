# 服务层详细设计 (`service/`)

## 1. 概述
服务层包含核心业务逻辑、工具函数和外部集成。它被控制器层和中转系统调用。这是系统中最复杂的层次，包含了计费、渠道路由、Token 计数等关键功能。

## 2. 文件详细说明（按领域分类）

### 2.1 渠道选择与亲和性
- **`channel.go`** -- 渠道缓存管理、随机渠道选择、渠道禁用/启用逻辑。
- **`channel_select.go`** -- 渠道选择算法：优先级、权重、基于组的路由。
- **`channel_affinity.go`** -- 渠道亲和性规则引擎：按用户/密钥的粘性路由。

### 2.2 计费与配额
- **`billing.go`** -- 核心计费逻辑：预扣费、结算、退款配额操作。
- **`billing_session.go`** -- 计费会话管理（事务性计费上下文）。
- **`quota.go`** -- 配额计算、扣减和补充。
- **`text_quota.go`** -- 基于文本的配额计算（模型比例、完成比例、缓存比例）。
- **`pre_consume_quota.go`** -- API 调用前预扣配额（失败时退款）。
- **`tiered_settle.go`** -- 使用表达式系统的分层/动态计费结算。
- **`tool_billing.go`** -- 工具调用计费计算。
- **`violation_fee.go`** -- 违规费用计算（CSAM 检测、安全策略违规）。
- **`funding_source.go`** -- 资金源管理（钱包 vs 订阅）。

### 2.3 Token 计数与估算
- **`token_counter.go`** -- `EstimateRequestToken`, `CountTokenRealtime`, `CountTokenInput`, `CountAudioTokenInput/Output`, `CountTextToken`: 用于计费的 Token 计数。
- **`token_estimator.go`** -- `EstimateToken`, `EstimateTokenByModel`: 无需完整分词器的快速 Token 估算。
- **`tokenizer.go`** -- `InitTokenEncoders`, `getTokenEncoder`, `getTokenNum`: 基于 tiktoken 的 Token 编码。

### 2.4 用户与鉴权
- **`group.go`** -- `GetUserUsableGroups`, `GetUserGroupRatio`: 用户组解析和比例查找。
- **`passkey/service.go`**, **`passkey/session.go`**, **`passkey/user.go`** -- WebAuthn/Passkey 服务实现。

### 2.5 通知
- **`user_notify.go`** -- `NotifyRootUser`, `NotifyUpstreamModelUpdateWatchers`, `NotifyUser`: 多渠道通知（邮件、Bark、Gotify、Telegram、Discord、Webhook）。
- **`notify-limit.go`** -- 通知限流。
- **`webhook.go`** -- `SendWebhookNotify`: 带 HMAC 签名的 Webhook 通知。

### 2.6 HTTP 与网络
- **`http_client.go`** -- `NewProxyHttpClient`: 代理感知的 HTTP 客户端创建。
- **`http.go`** -- HTTP 工具函数。
- **`download.go`** -- 文件下载工具。

### 2.7 文件处理
- **`file_service.go`** -- 文件服务抽象。
- **`file_decoder.go`** -- 文件解码工具。

### 2.8 内容与安全
- **`sensitive.go`** -- `CheckSensitiveText`: 敏感内容过滤。
- **`convert.go`** -- 数据转换工具。
- **`str.go`** -- 字符串工具。

### 2.9 任务管理
- **`task.go`** -- 异步 AI 任务（视频生成等）的提交和结果处理。
- **`task_billing.go`** -- 任务特定计费逻辑。
- **`task_polling.go`** -- 任务状态轮询。

### 2.10 Midjourney
- **`midjourney.go`** -- Midjourney 特定业务逻辑。

### 2.11 AI 兼容性
- **`openai_chat_responses_compat.go`** -- OpenAI Chat-to-Responses 格式兼容性。
- **`openai_chat_responses_mode.go`** -- OpenAI Responses 模式处理。

### 2.12 OpenAI 兼容子包 (`openaicompat/`)
- **`chat_to_responses.go`** -- 将 Chat Completions 请求转换为 Responses API 格式。
- **`responses_to_chat.go`** -- 将 Responses API 请求转换为 Chat Completions 格式。
- **`policy.go`** -- 转换策略。
- **`regex.go`** -- 格式检测的正则表达式。

### 2.13 排名与分析
- **`rankings.go`** -- 使用量排名计算。
- **`log_info_generate.go`** -- 日志信息生成辅助函数。
- **`usage_helpr.go`** -- `ResponseText2Usage`, `ValidUsage`: 从响应中提取使用量数据。

### 2.14 支付集成
- **`waffo_pancake.go`** -- `CreateWaffoPancakeCheckoutSession`, `VerifyConfiguredWaffoPancakeWebhook`, `ResolveWaffoPancakeTradeNo`, `CreateWaffoPancakePrimaryStore`, `SaveWaffoPancakeConfig`, `ListWaffoPancakeCatalog`: Waffo-Pancake 支付服务。
- **`epay.go`** -- Epay 支付服务。

### 2.15 Codex 集成
- **`codex_oauth.go`** -- Codex OAuth 令牌交换和 JWT 解析。
- **`codex_credential_refresh.go`** -- `RefreshCodexChannelCredential`: Codex 凭证刷新。
- **`codex_credential_refresh_task.go`** -- 后台 Codex 凭证刷新任务。
- **`codex_wham_usage.go`** -- `FetchCodexWhamUsage`: Codex 使用量 API 客户端。

### 2.16 其他
- **`return_path.go`** -- 返回路径处理。
- **`subscription_reset_task.go`** -- 后台订阅过期任务。
- **`error.go`** -- `RelayErrorHandler`: 上游中继错误解析和标准化。

---

## 关联文件列表

### 服务层核心文件
- `service/channel.go` - 渠道缓存管理
- `service/channel_select.go` - 渠道选择算法
- `service/channel_affinity.go` - 渠道亲和性
- `service/billing.go` - 核心计费逻辑
- `service/billing_session.go` - 计费会话管理
- `service/quota.go` - 配额计算
- `service/text_quota.go` - 文本配额计算
- `service/pre_consume_quota.go` - 预扣配额
- `service/tiered_settle.go` - 分层结算
- `service/tool_billing.go` - 工具调用计费
- `service/violation_fee.go` - 违规费用
- `service/funding_source.go` - 资金源管理
- `service/token_counter.go` - Token 计数
- `service/token_estimator.go` - Token 估算
- `service/tokenizer.go` - Token 编码器
- `service/group.go` - 用户组管理
- `service/passkey/service.go` - Passkey 服务
- `service/passkey/session.go` - Passkey 会话
- `service/passkey/user.go` - Passkey 用户
- `service/user_notify.go` - 用户通知
- `service/notify-limit.go` - 通知限流
- `service/webhook.go` - Webhook 通知
- `service/http_client.go` - HTTP 客户端
- `service/http.go` - HTTP 工具
- `service/download.go` - 文件下载
- `service/file_service.go` - 文件服务
- `service/file_decoder.go` - 文件解码
- `service/sensitive.go` - 敏感词检查
- `service/convert.go` - 数据转换
- `service/str.go` - 字符串工具
- `service/task.go` - 任务管理
- `service/task_billing.go` - 任务计费
- `service/task_polling.go` - 任务轮询
- `service/midjourney.go` - Midjourney 逻辑
- `service/openai_chat_responses_compat.go` - Chat-Responses 兼容
- `service/openai_chat_responses_mode.go` - Responses 模式
- `service/openaicompat/chat_to_responses.go` - Chat 转 Responses
- `service/openaicompat/responses_to_chat.go` - Responses 转 Chat
- `service/openaicompat/policy.go` - 转换策略
- `service/openaicompat/regex.go` - 正则表达式
- `service/rankings.go` - 使用排名
- `service/log_info_generate.go` - 日志信息生成
- `service/usage_helpr.go` - 使用量帮助
- `service/waffo_pancake.go` - Waffo-Pancake 支付
- `service/epay.go` - Epay 支付
- `service/codex_oauth.go` - Codex OAuth
- `service/codex_credential_refresh.go` - Codex 凭证刷新
- `service/codex_credential_refresh_task.go` - Codex 凭证刷新任务
- `service/codex_wham_usage.go` - Codex 使用量
- `service/return_path.go` - 返回路径
- `service/subscription_reset_task.go` - 订阅重置任务
- `service/error.go` - 中继错误处理

### 依赖的模型层文件
- `model/channel.go` - 渠道模型
- `model/ability.go` - 能力模型
- `model/user.go` - 用户模型
- `model/token.go` - Token 模型
- `model/pricing.go` - 定价模型
- `model/log.go` - 日志模型

### 依赖的常量文件
- `constant/channel.go` - 渠道常量
- `constant/api_type.go` - API 类型常量
- `constant/context_key.go` - 上下文键常量
