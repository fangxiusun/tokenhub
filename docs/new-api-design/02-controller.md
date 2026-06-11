# 控制器层详细设计 (`controller/`)

## 1. 概述
控制器层处理 HTTP 请求/响应逻辑。每个文件对应一个业务领域。控制器解析请求，调用服务/模型层，并格式化 JSON 响应。**不包含业务逻辑**。

## 2. 文件详细说明（按领域分类）

### 2.1 系统与状态管理
- **`setup.go`** -- `GetSetup`, `PostSetup`: 初始系统设置向导（创建 Root 用户，设置运营模式）。
- **`misc.go`** -- `GetStatus`, `GetNotice`, `GetAbout`, `GetUserAgreement`, `GetPrivacyPolicy`, `GetHomePageContent`, `TestStatus`, `SendEmailVerification`, `SendPasswordResetEmail`, `ResetPassword`: 系统状态/配置端点，邮件验证，密码重置。

### 2.2 用户管理
- **`user.go`** (1296行) -- 完整用户生命周期管理：
  - `Login`, `Register`, `Logout` - 认证
  - `GetAllUsers`, `SearchUsers`, `GetUser`, `CreateUser`, `UpdateUser`, `DeleteUser`, `ManageUser` - CRUD
  - `GetSelf`, `UpdateSelf`, `DeleteSelf` - 自服务
  - `GetUserModels`, `GetAffCode`, `TransferAffQuota` - 用户功能
  - `TopUp`, `EmailBind`, `UpdateUserSetting`, `GenerateAccessToken`, `AdminClearUserBinding` - 其他
- **`twofa.go`** -- 双因素认证：`Setup2FA`, `Enable2FA`, `Disable2FA`, `Get2FAStatus`, `RegenerateBackupCodes`, `Verify2FALogin`, `Admin2FAStats`, `AdminDisable2FA`。
- **`passkey.go`** -- WebAuthn/Passkey：`PasskeyStatus`, `PasskeyLoginBegin`, `PasskeyLoginFinish`, `PasskeyRegisterBegin`, `PasskeyRegisterFinish`, `PasskeyDelete`, `AdminResetPasskey`。

### 2.3 OAuth 提供商
- **`oauth.go`** -- `GenerateOAuthCode`, `HandleOAuth`: 所有标准提供商的统一 OAuth 处理器。
- **`github.go`** -- `GitHubOAuth`, `GitHubBind`: GitHub OAuth 登录/绑定。
- **`discord.go`** -- `DiscordOAuth`, `DiscordBind`: Discord OAuth 登录/绑定。
- **`linuxdo.go`** -- `LinuxdoOAuth`, `LinuxDoBind`: LinuxDO OAuth 登录/绑定。
- **`wechat.go`** -- `WeChatAuth`, `WeChatBind`: 微信 OAuth 登录/绑定。
- **`telegram.go`** -- `TelegramLogin`, `TelegramBind`: Telegram OAuth 登录/绑定。
- **`oidc.go`** -- OIDC 提供商登录/绑定。
- **`custom_oauth.go`** (584行) -- 完整自定义 OAuth 提供商 CRUD + 用户绑定管理。

### 2.4 渠道管理
- **`channel.go`** (1900+行) -- 渠道完整生命周期管理：
  - `GetAllChannels`, `SearchChannels`, `GetChannel`, `GetChannelKey`, `AddChannel`, `UpdateChannel`, `DeleteChannel` - CRUD
  - `FetchModels`, `FetchUpstreamModels`, `FixChannelsAbilities` - 模型管理
  - `DisableTagChannels`, `EnableTagChannels`, `EditTagChannels`, `DeleteChannelBatch`, `BatchSetChannelTag` - 批量操作
  - `CopyChannel`, `ManageMultiKeys` - 复制与多密钥管理
  - `OllamaPullModel`, `OllamaPullModelStream`, `OllamaDeleteModel`, `OllamaVersion` - Ollama 管理
  - `RefreshCodexChannelCredential` - Codex 凭证刷新
- **`channel-billing.go`** (505行) -- 渠道余额管理：`UpdateChannelBalance`, `UpdateAllChannelsBalance`, `AutomaticallyUpdateChannels`。
- **`channel-test.go`** (1010行) -- 渠道健康测试：`TestChannel`, `TestAllChannels`, `AutomaticallyTestChannels`（带自动启用/禁用，响应时间跟踪）。
- **`channel_upstream_update.go`** (999行) -- 上游模型自动同步：`ApplyChannelUpstreamModelUpdates`, `DetectChannelUpstreamModelUpdates`, `ApplyAllChannelUpstreamModelUpdates`, `DetectAllChannelUpstreamModelUpdates`, `StartChannelUpstreamModelUpdateTask`。
- **`channel_affinity_cache.go`** -- 渠道亲和性缓存管理。

### 2.5 中转（AI API 代理）
- **`relay.go`** (653行) -- 核心中转处理器：`Relay`, `RelayMidjourney`, `RelayNotImplemented`, `RelayNotFound`, `RelayTask`, `RelayTaskFetch`。包含重试逻辑、渠道错误处理、WebSocket 支持、计费集成、性能指标。
- **`model.go`** (364行) -- 模型列表：`ListModels`, `RetrieveModel`, `DashboardListModels`, `SyncUpstreamPreview`, `SyncUpstreamModels`。

### 2.6 Token 管理
- **`token.go`** (359行) -- API Token 完整 CRUD：`GetAllTokens`, `SearchTokens`, `GetToken`, `GetTokenKey`, `AddToken`, `UpdateToken`, `DeleteToken`, `DeleteTokenBatch`, `GetTokenKeysBatch`, `GetTokenUsage`。

### 2.7 计费与支付
- **`billing.go`** -- `GetSubscription`, `GetUsage`: OpenAI 兼容计费/订阅响应。
- **`topup.go`** -- Epay 支付集成：`RequestEpay`, `EpayNotify`, `RequestAmount`, `GetTopUpInfo`, `GetUserTopUps`, `AdminCompleteTopUp`, `GetAllTopUps`。
- **`topup_stripe.go`** -- Stripe 支付集成。
- **`topup_creem.go`** -- Creem 支付集成。
- **`topup_waffo.go`** -- Waffo 支付集成。
- **`topup_waffo_pancake.go`** -- Waffo-Pancake 支付 + 订阅管理。
- **`subscription.go`** -- 订阅计费：`GetSubscriptionPlans`, `GetSubscriptionSelf`, `UpdateSubscriptionPreference`, `SubscriptionRequestBalancePay`, `SubscriptionRequestEpay`, `SubscriptionRequestStripePay`, `SubscriptionRequestCreemPay`, `SubscriptionRequestWaffoPancakePay`。
- **`redemption.go`** -- 兑换码管理：`GetAllRedemptions`, `SearchRedemptions`, `GetRedemption`, `AddRedemption`, `UpdateRedemption`, `DeleteRedemption`, `DeleteInvalidRedemption`。

### 2.8 配置管理
- **`option.go`** (344行) -- 系统选项/配置管理：`GetOptions`, `UpdateOption`, `ConfirmPaymentCompliance`, `GetChannelAffinityCacheStats`, `ClearChannelAffinityCache`, `ResetModelRatio`, `MigrateConsoleSetting`。
- **`ratio_config.go`** -- `GetRatioConfig`: 模型比例配置。
- **`ratio_sync.go`** -- `GetSyncableChannels`, `FetchUpstreamRatios`: 上游比例同步。
- **`pricing.go`** -- `GetPricing`: 公开定价页面数据。
- **`prefill_group.go`** -- `GetPrefillGroups`, `CreatePrefillGroup`, `UpdatePrefillGroup`, `DeletePrefillGroup`。
- **`group.go`** -- `GetGroups`, `GetUserGroups`: 用户组列表。

### 2.9 日志与数据
- **`log.go`** -- 日志查询与管理：`GetAllLogs`, `GetUserLogs`, `SearchAllLogs`, `SearchUserLogs`, `GetLogByKey`, `GetLogsStat`, `GetLogsSelfStat`, `DeleteHistoryLogs`, `GetChannelAffinityUsageCacheStats`。
- **`usedata.go`** -- `GetAllQuotaDates`, `GetQuotaDatesByUser`, `GetUserQuotaDates`: 配额使用日期查询。

### 2.10 AI 功能
- **`midjourney.go`** -- Midjourney 任务管理：`UpdateMidjourneyTaskBulk`, `GetAllMidjourney`, `GetUserMidjourney`。
- **`task.go`** -- 通用任务管理：`GetAllTask`, `GetUserTask`。
- **`task_video.go`** -- 视频任务管理。
- **`playground.go`** -- `Playground`: 聊天 Playground 端点。

### 2.11 模型与供应商元数据
- **`model_meta.go`** -- 模型元数据 CRUD：`GetAllModelsMeta`, `SearchModelsMeta`, `GetModelMeta`, `CreateModelMeta`, `UpdateModelMeta`, `DeleteModelMeta`。
- **`model_sync.go`** -- 模型与上游同步。
- **`missing_models.go`** -- `GetMissingModels`: 查找本地数据库中缺失的模型。
- **`vendor_meta.go`** -- 供应商元数据 CRUD。

### 2.12 其他
- **`checkin.go`** -- `GetCheckinStatus`, `DoCheckin`: 每日签到系统。
- **`deployment.go`** (810行) -- io.net GPU 部署管理：`GetModelDeploymentSettings`, `TestIoNetConnection`, `GetAllDeployments`, `SearchDeployments`, `GetDeployment`, `CreateDeployment`, `UpdateDeployment`, `DeleteDeployment`, `GetHardwareTypes`, `GetLocations`, `GetAvailableReplicas`, `GetPriceEstimation`, `GetDeploymentLogs`, `ListDeploymentContainers`。
- **`rankings.go`** -- `GetRankings`: 使用量排名。
- **`perf_metrics.go`** -- `GetPerfMetricsSummary`, `GetPerfMetrics`: 性能指标。
- **`performance.go`** -- `GetPerformanceStats`, `ClearDiskCache`, `ResetPerformanceStats`, `ForceGC`, `GetLogFiles`, `CleanupLogFiles`。

---

## 关联文件列表

### 控制器层核心文件
- `controller/user.go` - 用户管理（登录/注册/CRUD/设置）
- `controller/twofa.go` - 双因素认证
- `controller/passkey.go` - WebAuthn/Passkey
- `controller/oauth.go` - OAuth 统一处理
- `controller/github.go` - GitHub OAuth
- `controller/discord.go` - Discord OAuth
- `controller/linuxdo.go` - LinuxDO OAuth
- `controller/wechat.go` - 微信 OAuth
- `controller/telegram.go` - Telegram OAuth
- `controller/oidc.go` - OIDC OAuth
- `controller/custom_oauth.go` - 自定义 OAuth
- `controller/channel.go` - 渠道管理
- `controller/channel-billing.go` - 渠道余额
- `controller/channel-test.go` - 渠道测试
- `controller/channel_upstream_update.go` - 上游模型同步
- `controller/channel_affinity_cache.go` - 渠道亲和性缓存
- `controller/relay.go` - AI API 中转核心
- `controller/model.go` - 模型管理
- `controller/token.go` - Token 管理
- `controller/billing.go` - 计费
- `controller/topup.go` - Epay 充值
- `controller/topup_stripe.go` - Stripe 充值
- `controller/topup_creem.go` - Creem 充值
- `controller/topup_waffo.go` - Waffo 充值
- `controller/topup_waffo_pancake.go` - Waffo-Pancake 充值
- `controller/subscription.go` - 订阅管理
- `controller/redemption.go` - 兑换码管理
- `controller/option.go` - 系统选项
- `controller/ratio_config.go` - 比例配置
- `controller/ratio_sync.go` - 比例同步
- `controller/pricing.go` - 定价
- `controller/prefill_group.go` - 预填充组
- `controller/group.go` - 用户组
- `controller/log.go` - 日志管理
- `controller/usedata.go` - 使用数据
- `controller/midjourney.go` - Midjourney 任务
- `controller/task.go` - 通用任务
- `controller/task_video.go` - 视频任务
- `controller/playground.go` - Playground
- `controller/model_meta.go` - 模型元数据
- `controller/model_sync.go` - 模型同步
- `controller/missing_models.go` - 缺失模型
- `controller/vendor_meta.go` - 供应商元数据
- `controller/setup.go` - 系统设置向导
- `controller/misc.go` - 系统状态/通知
- `controller/checkin.go` - 每日签到
- `controller/deployment.go` - io.net 部署管理
- `controller/rankings.go` - 使用排名
- `controller/perf_metrics.go` - 性能指标
- `controller/performance.go` - 性能管理

### 依赖的服务层文件
- `service/billing.go` - 计费逻辑
- `service/channel.go` - 渠道缓存管理
- `service/channel_select.go` - 渠道选择算法
- `service/quota.go` - 配额计算
- `service/token_counter.go` - Token 计数
- `service/user_notify.go` - 用户通知

### 依赖的模型层文件
- `model/user.go` - 用户模型
- `model/token.go` - Token 模型
- `model/channel.go` - 渠道模型
- `model/log.go` - 日志模型
- `model/option.go` - 系统选项模型
- `model/pricing.go` - 定价模型
