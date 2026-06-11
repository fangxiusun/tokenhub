# funding_source.go 代码阅读文档

## 1. 全局总结

该文件定义了计费资金来源接口 `FundingSource` 及其两个实现：钱包（WalletFunding）和订阅（SubscriptionFunding）。是计费系统的核心抽象层。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `model` | 额度扣减/增加操作 |
| `time` | 重试间隔 |

## 3. 类型定义

### `FundingSource` 接口
- `Source() string` — 返回资金来源标识
- `PreConsume(amount int) error` — 预扣额度
- `Settle(delta int) error` — 差额结算
- `Refund() error` — 退还预扣费

### `WalletFunding`
- `userId` — 用户 ID
- `consumed` — 已预扣额度

### `SubscriptionFunding`
- `requestId` — 请求 ID（幂等键）
- `userId` / `modelName` — 用户和模型信息
- `amount` — 预扣额度
- `subscriptionId` — 订阅 ID
- `preConsumed` — 实际预扣量
- `AmountTotal` / `AmountUsedAfter` / `PlanId` / `PlanTitle` — 订阅详情

## 4. 函数详解

### WalletFunding
- `PreConsume`：调用 `model.DecreaseUserQuota` 扣减用户额度
- `Settle`：正数补扣，负数退还
- `Refund`：增加用户额度（非幂等，不能重试）

### SubscriptionFunding
- `PreConsume`：调用 `model.PreConsumeUserSubscription` 创建预扣记录
- `Settle`：调用 `model.PostConsumeUserSubscriptionDelta` 调整差额
- `Refund`：调用 `model.RefundSubscriptionPreConsume`（幂等，可重试）

### `refundWithRetry(fn) error`
- 通用退款重试函数
- 最多 3 次尝试，递增延迟（200ms/400ms/600ms）
- 仅用于基于事务的退款函数

## 5. 关键逻辑分析

1. **幂等性差异**：Wallet 退款非幂等不能重试，Subscription 退款幂等可重试
2. **订阅预扣**：amount 参数被忽略，使用构造时的值
3. **订阅计划信息**：PreConsume 后获取计划详情

## 6. 关联文件

- `billing_session.go` — 使用 FundingSource 的计费会话
- `billing.go` — 计费入口
