# 模型层详细设计 (`model/`)

## 1. 概述
模型层使用 GORM 处理数据库交互。包含 Redis 缓存层、批量更新支持和数据验证。所有数据库操作都通过此层进行，确保数据一致性和高性能。

## 2. 文件详细说明

### 2.1 核心模型
- **`user.go`** (30932字节) -- 用户模型：
  - `User` 结构体：Id, Username, Password, Role, Status, Email, Quota, UsedQuota, RequestCount, Group, AffCode, Setting
  - CRUD 操作：创建、查询、更新、删除用户
  - 密码处理：bcrypt 哈希、验证
  - 配额管理：余额查询、扣减、补充
  - OAuth 绑定：GitHub, Discord, OIDC, 微信, Telegram, LinuxDO
  - 用户设置：JSON 序列化的用户偏好
- **`token.go`** (14447字节) -- API Token 模型：
  - `Token` 结构体：Id, UserId, Key, Status, Name, ExpiredTime, RemainQuota, UnlimitedQuota, ModelLimitsEnabled, ModelLimits, AllowIps, UsedQuota, Group, CrossGroupRetry
  - 密钥生成：sk- 前缀 + 随机字符串
  - 配额管理：余额查询、扣减
  - 模型限制：逗号分隔的模型列表
  - IP 限制：换行分隔的 IP/CIDR 列表
- **`channel.go`** (32897字节) -- 渠道模型：
  - `Channel` 结构体：Id, Type, Key, BaseURL, Other, Status, Name, Weight, CreatedTime, Balance, BalanceUpdatedTime, Models, ModelMapping, groups, tags, AutoBan, ModelRatio, ModelPrice, BatchMaxInputTokens, BatchMaxOutputTokens, BatchMaxNum, BatchMode
  - CRUD 操作：创建、查询、更新、删除渠道
  - 密钥管理：支持多密钥（逗号分隔）
  - 余额检查：自动检测渠道余额
  - 多密钥支持：随机和轮询模式
- **`ability.go`** (9299字节) -- 渠道能力模型：
  - `Ability` 结构体：Group, Model, ChannelId, Enabled, Priority, Weight
  - CRUD 操作：创建、查询、更新、删除能力
  - 优先级支持：数值越高越优先
  - 权重支持：同优先级内的随机权重
- **`log.go`** (17450字节) -- 日志模型：
  - `Log` 结构体：Id, UserId, CreatedAt, Type, Content, ModelName, TokenName, Quota, TokenId, Username, ChannelId, Channel, PromptTokens, CompletionTokens, Token, QuotaId, IPAddress, GeoIP
  - CRUD 操作：创建、查询、更新、删除日志
  - 使用量统计：按时间、用户、渠道统计
  - 渠道亲和性跟踪：记录用户与渠道的关联
- **`option.go`** (26023字节) -- 系统选项模型：
  - `Option` 结构体：Key, Value
  - 键值配置存储：JSON 序列化的配置项
  - 线程安全：读写锁保护
  - 热更新：运行时配置变更
- **`pricing.go`** (11709字节) -- 定价模型：
  - `Pricing` 结构体：ModelName, QuotaType, ModelRatio, ModelPrice, CompletionRatio, CacheRatio, CreateCacheRatio, AudioRatio, BillingMode, BillingExpr
  - 模型比例管理：输入/输出/缓存 Token 的计费比例
  - 配额计算：基于比例和价格的配额计算

### 2.2 任务管理
- **`task.go`** (16443字节) -- 任务模型：
  - `Task` 结构体：Id, UserId, Action, Model, Status, Progress, SubmitTime, UpdateTime, FailReason, MidjourneyId, VideoId, SunoId, Quota, Type
  - 状态管理：待处理、处理中、已完成、失败
  - 结果存储：JSON 序列化的任务结果
  - 计费集成：任务提交时预扣配额
- **`midjourney.go`** (6319字节) -- Midjourney 任务模型：
  - `Midjourney` 结构体：Id, UserId, Action, Prompt, Parameters, Type, Status, Progress, SubmitTime, ChannelId, Quota, DiscordId, Description, PromptEn, ButtonMessageId, MessageId, MessageContent, RemainingTime, ImageFailReason, ActionStatus
  - 状态跟踪：从提交到完成的完整生命周期
  - 图片存储：URL 和本地路径

### 2.3 支付与计费
- **`topup.go`** (16234字节) -- 充值模型：
  - `TopUp` 结构体：Id, UserId, Username, Amount, Quota, Status, TradeNo, UserBonus, BonusRatio, PaymentId, PaymentMethod
  - 支付记录：交易号、金额、配额
  - 状态管理：待支付、已完成、已取消
- **`redemption.go`** (5500字节) -- 兑换码模型：
  - `Redemption` 结构体：Id, Key, Status, Name, Quota, UsedTimes, UserId, UsedBy, RedeemedAt, CreatedAt
  - 密钥生成：随机字符串
  - 使用跟踪：使用次数、使用者
- **`subscription.go`** (39973字节) -- 订阅模型：
  - `Subscription` 结构体：Id, UserId, PlanId, Status, CurrentPeriodStart, CurrentPeriodEnd, CancelAtPeriodEnd, Quota, UsedQuota
  - 计划管理：订阅计划、计费周期
  - 配额管理：订阅配额、已用配额

### 2.4 用户功能
- **`twofa.go`** (8160字节) -- 双因素认证模型：
  - `TwoFa` 结构体：Id, UserId, Secret, BackupCodes, Enabled
  - 密钥管理：TOTP 密钥生成和存储
  - 备份码：一次性备份码生成和验证
- **`passkey.go`** (8000字节) -- WebAuthn/Passkey 模型：
  - `Passkey` 结构体：Id, UserId, Name, CredentialID, PublicKey, SignCount, Transports, AAGUID, AttestationObject, Authenticator
  - 凭证管理：WebAuthn 凭证存储和验证
- **`user_oauth_binding.go`** (5483字节) -- 用户 OAuth 绑定模型：
  - `UserOAuthBinding` 结构体：Id, UserId, Provider, ProviderUserID
  - 提供商管理：支持多个 OAuth 提供商绑定
- **`checkin.go`** (5798字节) -- 每日签到模型：
  - `CheckIn` 结构体：Id, UserId, ConsecutiveDays, LastCheckInTime
  - 奖励管理：连续签到奖励计算

### 2.5 元数据
- **`model_meta.go`** (6440字节) -- 模型元数据模型：
  - `ModelMeta` 结构体：Id, ModelName, ModelType, DisplayName, Icon, Description, OpenAICompatible, Pricing
  - 上游同步：从上游提供商同步模型信息
- **`vendor_meta.go`** (2610字节) -- 供应商元数据模型：
  - `VendorMeta` 结构体：Id, VendorName, DisplayName, Icon, Description, URL
  - 供应商信息：提供商详细信息
- **`prefill_group.go`** (3579字节) -- 预填充组模型：
  - `PrefillGroup` 结构体：Id, Name, Description, Models, Channels
  - 模板管理：预配置的组模板

### 2.6 缓存
- **`channel_cache.go`** (8082字节) -- 渠道缓存：
  - Redis 基于渠道缓存：批量更新支持
  - 缓存键：`channel:{id}`, `channels:list`
  - 缓存策略：写时失效，定时刷新
- **`user_cache.go`** (6164字节) -- 用户缓存：
  - Redis 基于用户缓存：配额跟踪
  - 缓存键：`user:{id}`, `user:quota:{id}`
  - 缓存策略：读时加载，异步写回
- **`token_cache.go`** (1586字节) -- Token 缓存：
  - Redis 基于 Token 缓存：密钥验证
  - 缓存键：`token:{key}`
  - 缓存策略：写时失效

### 2.7 数据库操作
- **`main.go`** (20776字节) -- 数据库初始化：
  - 数据库连接建立：SQLite, MySQL, PostgreSQL
  - 数据库迁移：自动表结构迁移
  - 批量更新支持：异步批量写入减少数据库压力
  - Redis 连接建立：可选 Redis 缓存层
- **`errors.go`** (529字节) -- 模型错误定义。
- **`db_time.go`** (578字节) -- 数据库时间工具。
- **`utils.go`** (3081字节) -- 模型工具函数。

### 2.8 使用数据
- **`usedata.go`** (5454字节) -- 使用数据模型：
  - `UseData` 结构体：Id, UserId, ModelName, Date, TokenCount, Quota
  - 配额跟踪：按日期和模型统计使用量
- **`usedata_rankings.go`** (1878字节) -- 使用排名。

### 2.9 性能
- **`perf_metric.go`** (3708字节) -- 性能指标模型：
  - `PerfMetric` 结构体：Id, MetricName, MetricValue, Timestamp
  - 系统监控：CPU、内存、磁盘使用率

### 2.10 其他
- **`missing_models.go`** (766字节) -- 查找本地数据库中缺失的模型。
- **`setup.go`** (359字节) -- 设置状态跟踪。
- **`model_extra.go`** (632字节) -- 额外模型信息。

---

## 关联文件列表

### 模型层核心文件
- `model/user.go` - 用户模型
- `model/token.go` - Token 模型
- `model/channel.go` - 渠道模型
- `model/ability.go` - 渠道能力模型
- `model/log.go` - 日志模型
- `model/option.go` - 系统选项模型
- `model/pricing.go` - 定价模型
- `model/task.go` - 任务模型
- `model/midjourney.go` - Midjourney 任务模型
- `model/topup.go` - 充值模型
- `model/redemption.go` - 兑换码模型
- `model/subscription.go` - 订阅模型
- `model/twofa.go` - 双因素认证模型
- `model/passkey.go` - Passkey 模型
- `model/user_oauth_binding.go` - 用户 OAuth 绑定模型
- `model/checkin.go` - 每日签到模型
- `model/model_meta.go` - 模型元数据模型
- `model/vendor_meta.go` - 供应商元数据模型
- `model/prefill_group.go` - 预填充组模型
- `model/channel_cache.go` - 渠道缓存
- `model/user_cache.go` - 用户缓存
- `model/token_cache.go` - Token 缓存
- `model/main.go` - 数据库初始化
- `model/errors.go` - 模型错误
- `model/db_time.go` - 数据库时间工具
- `model/utils.go` - 模型工具
- `model/usedata.go` - 使用数据模型
- `model/usedata_rankings.go` - 使用排名
- `model/perf_metric.go` - 性能指标模型
- `model/missing_models.go` - 缺失模型
- `model/setup.go` - 设置状态
- `model/model_extra.go` - 额外模型信息

### 依赖的公共工具文件
- `common/database.go` - 数据库工具
- `common/redis.go` - Redis 工具
- `common/json.go` - JSON 序列化
- `common/crypto.go` - 加密工具
- `common/constants.go` - 常量定义

### 依赖的常量文件
- `constant/channel.go` - 渠道常量
- `constant/api_type.go` - API 类型常量
- `constant/cache_key.go` - 缓存键常量
