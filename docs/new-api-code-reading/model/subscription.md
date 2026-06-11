# subscription.go 代码阅读文档

## 1. 全局总结

`subscription.go` 实现了完整的订阅（Subscription）系统，是整个文件中最大的文件（1321行）。它包含订阅套餐（Plan）、订阅订单（Order）、用户订阅实例（UserSubscription）三个核心数据模型，以及套餐购买、订单完成、额度预扣/退款、配额重置、过期处理、分组升降级等完整的业务逻辑。该系统支持多种支付方式（外部支付、余额支付、管理员绑定），并实现了幂等性保证和 Redis+内存混合缓存。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `errors` | 标准错误类型 |
| `fmt` | 格式化输出 |
| `strconv` | 整数转字符串（缓存 key 生成） |
| `strings` | 字符串修剪和比较 |
| `sync` | `sync.Once` 用于缓存单例初始化 |
| `time` | 时间计算（订阅时长、重置周期） |
| `github.com/QuantumNous/new-api/common` | 通用工具（时间戳、随机字符串、环境变量、Redis 等） |
| `github.com/QuantumNous/new-api/pkg/cachex` | Redis+内存混合缓存框架 |
| `github.com/samber/hot` | 高性能热缓存库（LRU） |
| `github.com/shopspring/decimal` | 精确十进制运算（金额计算） |
| `gorm.io/gorm` | ORM 框架 |

## 3. 类型定义

### 常量

**时长单位**：`year`、`month`、`day`、`hour`、`custom`

**重置周期**：`never`、`daily`、`weekly`、`monthly`、`custom`

### SubscriptionPlan 结构体
订阅套餐定义：
- 基础信息：`Title`、`Subtitle`、`PriceAmount`、`Currency`
- 时长配置：`DurationUnit`、`DurationValue`、`CustomSeconds`
- 显示控制：`Enabled`、`SortOrder`
- 支付配置：`AllowBalancePay`、`StripePriceId`、`CreemProductId`、`WaffoPancakeProductId`
- 限制：`MaxPurchasePerUser`
- 分组升级：`UpgradeGroup`
- 额度配置：`TotalAmount`、`QuotaResetPeriod`、`QuotaResetCustomSeconds`

### SubscriptionOrder 结构体
订阅订单：
- `UserId`、`PlanId`、`Money`
- `TradeNo`：交易号（唯一索引）
- `PaymentMethod`、`PaymentProvider`
- `Status`：pending/success/expired
- `ProviderPayload`：支付提供商原始数据

### UserSubscription 结构体
用户订阅实例：
- `AmountTotal`/`AmountUsed`：总额度/已用额度
- `StartTime`/`EndTime`：起止时间
- `Status`：active/expired/cancelled
- `Source`：来源（order/admin）
- `LastResetTime`/`NextResetTime`：配额重置时间
- `UpgradeGroup`/`PrevUserGroup`：分组升降级记录

### SubscriptionPreConsumeRecord 结构体
预扣幂等记录，基于 `RequestId` 去重。

### SubscriptionPreConsumeResult 结构体
预扣操作返回结果，包含订阅 ID、预扣额度、总额度等。

### SubscriptionSummary 结构体
订阅摘要，包含 `UserSubscription` 指针。

### SubscriptionPlanInfo 结构体
套餐简要信息（ID 和标题），用于缓存查询。

### SubscriptionDuration* 常量
时长单位常量。

### SubscriptionReset* 常量
重置周期常量。

## 4. 函数详解

### 缓存相关

#### subscriptionPlanCacheTTL() / subscriptionPlanInfoCacheTTL()
从环境变量读取缓存 TTL，默认分别为 300 秒和 120 秒。

#### subscriptionPlanCacheCapacity() / subscriptionPlanInfoCacheCapacity()
从环境变量读取缓存容量，默认分别为 5000 和 10000。

#### getSubscriptionPlanCache() / getSubscriptionPlanInfoCache()
使用 `sync.Once` 单例初始化混合缓存（Redis + LRU 内存）。

#### InvalidateSubscriptionPlanCache(planId int)
清除指定套餐缓存和所有套餐信息缓存。

### 套餐与订单

#### GetSubscriptionPlanById(id int) (*SubscriptionPlan, error)
带缓存的套餐查询。先查缓存，未命中则查数据库并回填缓存。

#### calcPlanEndTime(start time.Time, plan *SubscriptionPlan) (int64, error)
根据套餐时长配置计算结束时间。支持年、月、日、小时、自定义秒数。

#### NormalizeResetPeriod(period string) string
标准化重置周期字符串，无效值返回 `never`。

#### calcNextResetTime(base time.Time, plan *SubscriptionPlan, endUnix int64) int64
计算下次配额重置时间。支持每日、每周（对齐到下周一）、每月（对齐到下月1号）、自定义周期。

#### CreateUserSubscriptionFromPlanTx(tx *gorm.DB, ...) (*UserSubscription, error)
事务中创建用户订阅。检查购买上限，计算结束时间和重置时间，处理分组升级。

#### CompleteSubscriptionOrder(tradeNo string, ...) error
完成订阅订单（幂等）。使用行锁查询订单，验证状态后创建订阅实例，更新 TopUp 记录，记录日志。

#### ExpireSubscriptionOrder(tradeNo string, ...) error
过期订阅订单，将 pending 状态改为 expired。

#### AdminBindSubscription(userId, planId int, sourceNote string) (string, error)
管理员绑定订阅，无需支付。创建订阅实例并可能触发分组升级。

#### PurchaseSubscriptionWithBalance(userId, planId int) error
使用余额购买订阅。事务中检查套餐有效性、余额充足性，扣减额度，创建订阅和订单。

### 查询

#### GetAllActiveUserSubscriptions(userId int) ([]SubscriptionSummary, error)
获取用户所有活跃订阅。

#### HasActiveUserSubscription(userId int) (bool, error)
轻量级检查用户是否有活跃订阅。

#### GetAllUserSubscriptions(userId int) ([]SubscriptionSummary, error)
获取用户所有订阅（包括过期的）。

#### CountUserSubscriptionsByPlan(userId, planId int) (int64, error)
统计用户对某套餐的购买次数（用于限制购买上限）。

### 管理

#### AdminInvalidateUserSubscription(userSubscriptionId int) (string, error)
管理员取消订阅，立即结束并可能回退分组。

#### AdminDeleteUserSubscription(userSubscriptionId int) (string, error)
管理员硬删除订阅。

#### ExpireDueSubscriptions(limit int) (int, error)
批量过期到期的订阅，处理分组回退。

#### ResetDueSubscriptions(limit int) (int, error)
批量重置到期的订阅配额。

### 预扣/退款

#### PreConsumeUserSubscription(requestId string, userId int, ...) (*SubscriptionPreConsumeResult, error)
预扣订阅额度（幂等）。遍历活跃订阅找到有足够余额的，创建预扣记录并增加已用额度。

#### RefundSubscriptionPreConsume(requestId string) error
退款预扣额度（幂等）。调用 `PostConsumeUserSubscriptionDelta` 减少已用额度。

#### PostConsumeUserSubscriptionDelta(userSubscriptionId int, delta int64) error
更新订阅已用额度增量（正数消费，负数退款）。使用行锁确保安全。

### 清理

#### CleanupSubscriptionPreConsumeRecords(olderThanSeconds int64) (int64, error)
清理旧的预扣幂等记录。

#### GetSubscriptionPlanInfoByUserSubscriptionId(userSubscriptionId int) (*SubscriptionPlanInfo, error)
带缓存的套餐信息查询。

## 5. 关键逻辑分析

**幂等性保证**：`PreConsumeUserSubscription` 和 `RefundSubscriptionPreConsume` 都基于 `RequestId` 去重，即使重复调用也不会产生副作用。

**混合缓存策略**：使用 `cachex.HybridCache` 实现 Redis + LRU 内存的二级缓存，Redis 不可用时自动降级到纯内存缓存。

**分组升降级**：购买订阅时升级用户分组（`UpgradeGroup`），取消/过期时回退到之前的分组（`PrevUserGroup`）。回退前会检查是否有其他活跃订阅仍需要高级分组。

**配额重置机制**：支持每日/每周/每月/自定义周期的配额自动重置。`maybeResetUserSubscriptionWithPlanTx` 处理跨周期的多次重置（如长期未访问后一次性重置多个周期）。

**精确金额计算**：使用 `shopspring/decimal` 库避免浮点数精度问题，确保余额扣减的准确性。

**数据库兼容性**：`trade_no` 是保留字，通过 `common.UsingPostgreSQL` 判断使用不同的引号方式。

## 6. 关联文件

- `model/user.go`：用户额度和分组更新
- `model/topup.go`：`TopUp` 结构体和充值记录
- `model/log.go`：`RecordLog()` 日志记录
- `model/batch_update.go`：`cacheDecrUserQuota()` 缓存扣减
- `model/user_group_cache.go`：`UpdateUserGroupCache()` 分组缓存更新
- `pkg/cachex/`：混合缓存框架
- `common/redis.go`：Redis 工具函数
- `setting/operation_setting/`：系统配置
