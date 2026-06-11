# task_billing.go 代码阅读文档

## 1. 全局总结

该文件实现异步任务的计费逻辑，包括任务消费日志记录、令牌额度调整、资金来源调整、任务退款、以及差额结算。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | 字符串匹配 |
| `constant` | 任务常量 |
| `logger` | 日志 |
| `model` | 额度操作 |
| `relaycommon` | RelayInfo |
| `ratio_setting` | 倍率配置 |
| `gin` | HTTP 上下文 |

## 3. 函数详解

### `LogTaskConsumption(c, info)`
记录任务消费日志（仅记录，不涉及实际扣费）

### `resolveTokenKey(ctx, tokenId, taskID) string`
通过 TokenId 运行时获取令牌 Key

### `taskIsSubscription(task) bool`
判断任务是否通过订阅计费

### `taskAdjustFunding(task, delta) error`
调整任务资金来源（钱包/订阅）

### `taskAdjustTokenQuota(ctx, task, delta)`
调整任务令牌额度

### `RefundTaskQuota(ctx, task, reason)`
统一任务失败退款：
1. 退还资金来源
2. 退还令牌额度
3. 记录退款日志

### `RecalculateTaskQuota(ctx, task, actualQuota, reason)`
通用差额结算：
1. 计算 delta
2. 调整资金来源
3. 调整令牌额度
4. 记录日志

### `RecalculateTaskQuotaByTokens(ctx, task, totalTokens)`
根据实际 token 消耗重新计费

## 4. 关键逻辑分析

1. **按次计费**：支持固定价格的按次计费模式
2. **差额结算**：预扣费与实际消耗的差额处理
3. **令牌 Key 获取**：运行时从数据库获取，不从 PrivateData 读取

## 5. 关联文件

- `task_polling.go` — 任务轮询
- `quota.go` — PostConsumeQuota
