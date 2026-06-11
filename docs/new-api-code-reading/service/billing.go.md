# billing.go 代码阅读文档

## 1. 全局总结

该文件是计费系统的顶层入口，提供了预扣费（PreConsume）和后结算（Settle）的统一接口。它封装了 `BillingSession` 的创建和使用，是整个计费生命周期的协调者。所有 AI API 请求的计费流程都通过此文件的函数启动。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `fmt` | 格式化日志消息 |
| `logger` | 结构化日志记录 |
| `relaycommon` | RelayInfo 中继信息传递 |
| `types` | 错误类型定义 |
| `gin` | HTTP 上下文 |

## 3. 类型定义

### 常量
- `BillingSourceWallet = "wallet"` — 钱包计费来源标识
- `BillingSourceSubscription = "subscription"` — 订阅计费来源标识

## 4. 函数详解

### `PreConsumeBilling(c, preConsumedQuota, relayInfo) *types.NewAPIError`
- **作用**：根据用户计费偏好创建 BillingSession 并执行预扣费
- **流程**：调用 `NewBillingSession` 创建会话 → 将 session 存储到 `relayInfo.Billing`
- **返回**：成功返回 nil，失败返回 NewAPIError

### `SettleBilling(ctx, relayInfo, actualQuota) error`
- **作用**：执行计费结算，根据实际消耗与预扣费的差额进行调整
- **逻辑**：
  - 如果有 BillingSession：计算 delta（actualQuota - preConsumed）→ 调用 session.Settle → 发送额度通知
  - 如果无 BillingSession：回退到旧的 PostConsumeQuota 路径（兼容按次计费）
- **通知逻辑**：订阅计费使用 `checkAndSendSubscriptionQuotaNotify`，钱包计费使用 `checkAndSendQuotaNotify`

## 5. 关键逻辑分析

1. **双路径设计**：SettleBilling 同时支持新的 BillingSession 路径和旧的 PostConsumeQuota 路径，确保向后兼容
2. **差额日志**：根据 delta 的正负分别记录"补扣费"、"返还扣费"或"无需调整"日志
3. **通知机制**：结算完成后异步发送额度不足通知

## 6. 关联文件

- `billing_session.go` — BillingSession 结构体和预扣费实现
- `funding_source.go` — 资金来源接口（钱包/订阅）
- `pre_consume_quota.go` — 旧版预扣费逻辑
- `quota.go` — PostConsumeQuota 旧版后结算
