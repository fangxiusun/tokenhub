# topup.go 代码阅读文档

## 1. 全局总结

`topup.go` 是充值（TopUp）模块的核心文件，负责管理用户充值订单的完整生命周期。包括订单创建、状态更新、多支付渠道充值回调处理（Stripe、Creem、Waffo 等）、管理员补单、充值记录查询与搜索。所有充值操作均采用数据库事务和行级锁（`FOR UPDATE`）保证并发安全，防止重复充值。

## 2. 依赖关系

| 依赖 | 用途 |
|------|------|
| `common` | 数据库标志（`UsingPostgreSQL`）、时间戳、系统常量（`TopUpStatusPending`、`QuotaPerUnit`）、系统日志 |
| `logger` | 额度格式化输出（`FormatQuota`） |
| `shopspring/decimal` | 高精度金额计算，避免浮点误差 |
| `gorm.io/gorm` | ORM 框架，提供事务、行锁、查询能力 |

## 3. 类型定义

### TopUp（充值订单）

```go
type TopUp struct {
    Id              int     // 主键
    UserId          int     // 用户ID（索引）
    Amount          int64   // 支付金额（分/整数单位）
    Money           float64 // 充值金额（美元）
    TradeNo         string  // 交易单号（唯一索引）
    PaymentMethod   string  // 支付方式
    PaymentProvider string  // 支付供应商
    CreateTime      int64   // 创建时间
    CompleteTime    int64   // 完成时间
    Status          string  // 订单状态
}
```

### 支付方式常量

- `PaymentMethodStripe` / `PaymentMethodCreem` / `PaymentMethodWaffo` / `PaymentMethodWaffoPancake` / `PaymentMethodBalance`

### 支付供应商常量

- `PaymentProviderEpay` / `PaymentProviderStripe` / `PaymentProviderCreem` / `PaymentProviderWaffo` / `PaymentProviderWaffoPancake` / `PaymentProviderBalance`

### 错误变量

| 错误 | 含义 |
|------|------|
| `ErrPaymentMethodMismatch` | 支付方式不匹配 |
| `ErrTopUpNotFound` | 充值订单不存在 |
| `ErrTopUpStatusInvalid` | 充值订单状态无效 |

## 4. 函数详解

### 基础 CRUD

| 函数 | 说明 |
|------|------|
| `Insert()` | 插入新充值订单 |
| `Update()` | 更新充值订单 |
| `GetTopUpById(id)` | 按 ID 查询充值订单 |
| `GetTopUpByTradeNo(tradeNo)` | 按交易单号查询充值订单 |

### 状态更新

| 函数 | 说明 |
|------|------|
| `UpdatePendingTopUpStatus(tradeNo, expectedProvider, targetStatus)` | 在事务中更新待处理订单状态，使用行锁防并发 |
| `ManualCompleteTopUp(tradeNo, callerIp)` | 管理员手动补单，幂等处理，使用 `decimal` 精确计算额度 |

### 充值回调处理

| 函数 | 说明 |
|------|------|
| `Recharge(referenceId, customerId, callerIp)` | Stripe 充值回调：更新订单状态 + 给用户增加额度 |
| `RechargeCreem(referenceId, email, name, callerIp)` | Creem 充值回调：同时尝试回填用户邮箱 |
| `RechargeWaffo(tradeNo, callerIp)` | Waffo 充值回调 |
| `RechargeWaffoPancake(tradeNo)` | Waffo Pancake 充值回调 |

### 查询与搜索

| 函数 | 说明 |
|------|------|
| `GetUserTopUps(userId, pageInfo)` | 分页查询用户充值记录（限 30 天窗口） |
| `GetAllTopUps(pageInfo)` | 管理员查询全平台充值记录（无时间限制） |
| `SearchUserTopUps(userId, keyword, pageInfo)` | 用户按订单号搜索（30 天窗口，防注入 LIKE） |
| `SearchAllTopUps(keyword, pageInfo)` | 管理员按订单号搜索全平台记录 |

### 辅助函数

| 函数 | 说明 |
|------|------|
| `topUpQueryCutoff()` | 计算可查询的最早时间戳（当前时间 - 30 天） |

## 5. 关键逻辑分析

### 并发安全
所有充值操作使用 `DB.Transaction` + `FOR UPDATE` 行级锁，确保同一订单不会被并发处理。

### 多数据库兼容
通过 `common.UsingPostgreSQL` 判断，对 `trade_no` 列使用不同的引用方式（PostgreSQL 用双引号，MySQL/SQLite 用反引号）。

### 幂等性处理
`RechargeWaffo`、`RechargeWaffoPancake`、`ManualCompleteTopUp` 在订单已成功时直接返回 `nil`，防止重复充值。

### LIKE 注入防护
搜索函数使用 `sanitizeLikePattern()` 对用户输入的关键词进行转义，并配合 `ESCAPE '!'` 语法，防止 SQL LIKE 注入。

### 额度计算差异
- Stripe 订单：使用 `Money * QuotaPerUnit`（float64）
- 其他订单：使用 `Amount * QuotaPerUnit`（通过 `decimal` 精确计算）

### 搜索安全上限
`searchTopUpCountHardLimit = 10000`，防止超大表无界 COUNT 触发 DoS。

## 6. 关联文件

| 文件 | 关联内容 |
|------|----------|
| `common/main.go` | `QuotaPerUnit`、`TopUpStatusPending`、`TopUpStatusSuccess` 等常量 |
| `common/utils.go` | `GetTimestamp()`、`GetRandomString()` |
| `model/user.go` | 用户额度更新（`quota` 字段） |
| `model/log.go` | `RecordTopupLog()` 充值日志记录 |
| `model/utils.go` | `sanitizeLikePattern()` SQL LIKE 安全处理 |
| `relay/payment/` | Stripe、Creem、Waffo 等支付渠道的回调处理入口 |
