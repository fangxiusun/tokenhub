# subscription_reset_task.go 代码阅读文档

## 1. 全局总结

该文件实现订阅额度重置定时任务，每分钟检查并重置到期的订阅额度，同时清理过期的预扣记录。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | 主节点判断 |
| `logger` | 日志 |
| `model` | 订阅操作 |
| `gopool` | 异步任务 |

## 3. 类型定义

常量：
- `subscriptionResetTickInterval = 1 * time.Minute` — 检查间隔
- `subscriptionResetBatchSize = 300` — 批量大小
- `subscriptionCleanupInterval = 30 * time.Minute` — 清理间隔

## 4. 函数详解

### `StartSubscriptionQuotaResetTask()`
启动订阅额度重置任务（仅主节点）

### `runSubscriptionQuotaResetOnce()`
执行一次维护：
1. 过期到期的订阅
2. 重置需要重置的订阅额度
3. 清理 7 天前的预扣记录（每 30 分钟）

## 5. 关键逻辑分析

1. **CAS 防重入**：防止并发执行
2. **批量处理**：每次最多处理 300 条
3. **定期清理**：预扣记录每 30 分钟清理一次

## 6. 关联文件

- `model/subscription.go` — 订阅数据库操作
