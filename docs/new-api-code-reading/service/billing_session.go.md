# billing_session.go 代码阅读文档

## 1. 全局总结

该文件实现了统一计费会话 `BillingSession`，封装了单次请求的预扣费/结算/退款生命周期。支持钱包和订阅两种资金来源，实现了信任额度旁路（Trust Quota）优化，以及订阅预扣费的额外预留机制。是整个计费系统的核心组件。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | 信任额度配置、日志 |
| `logger` | 结构化日志 |
| `model` | 数据库操作（令牌额度、用户额度、订阅） |
| `relaycommon` | RelayInfo 中继信息 |
| `types` | 错误类型 |
| `gopool` | 异步任务池 |
| `gin` | HTTP 上下文 |
| `sync` | 互斥锁 |

## 3. 类型定义

### `BillingSession`
核心计费会话结构体，封装：
- `relayInfo` — 中继信息
- `funding` — 资金来源（FundingSource 接口）
- `preConsumedQuota` — 实际预扣额度
- `tokenConsumed` — 令牌额度实际扣减量
- `extraReserved` — 发送前补充预扣的额度
- `trusted` — 是否命中信任额度旁路
- `fundingSettled` / `settled` / `refunded` — 状态标记
- `mu` — 并发安全互斥锁

## 4. 核心方法

### `Settle(actualQuota int) error`
- 根据实际消耗进行差额结算
- 两步提交：资金来源 → 令牌额度
- 资金来源已提交但令牌调整失败时，标记 `fundingSettled` 防止 Refund 误退

### `Refund(c *gin.Context)`
- 幂等安全的退款操作
- 异步执行（gopool.Go）
- 退还资金来源 + 令牌额度 + 订阅额外预留

### `Reserve(targetQuota int) error`
- 发送前补充预扣额度
- 先预留资金来源 → 再预留令牌额度
- 失败时原子回滚

### `preConsume(c, quota) *types.NewAPIError`
- 统一预扣费入口
- 信任额度旁路检查 → 令牌预扣 → 资金来源预扣
- 任一步骤失败时原子回滚

### `shouldTrust(c) bool`
- 统一信任额度检查
- 异步任务（ForcePreConsume=true）不允许信任旁路
- 订阅不能启用信任旁路（原因：PreConsumeUserSubscription 要求 amount>0）

## 5. 关键逻辑分析

1. **信任额度旁路**：用户额度和令牌额度都超过阈值时，跳过预扣费直接放行
2. **订阅不信任**：订阅计费必须预扣（需要创建预扣记录并锁定订阅）
3. **原子性保证**：预扣费失败时回滚已完成的步骤
4. **fundingSettled 保护**：防止对已提交的资金来源执行退款
5. **syncRelayInfo**：将 BillingSession 状态同步到 RelayInfo 兼容字段

## 6. 关联文件

- `billing.go` — 顶层计费入口
- `funding_source.go` — 资金来源接口实现
- `pre_consume_quota.go` — 旧版预扣费
- `quota.go` — PostConsumeQuota
