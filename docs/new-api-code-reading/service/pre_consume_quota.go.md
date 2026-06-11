# pre_consume_quota.go 代码阅读文档

## 1. 全局总结

该文件提供旧版的预扣费逻辑，包括预扣费和退还预扣费。是新版 BillingSession 的兼容替代。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | 信任额度 |
| `logger` | 日志 |
| `model` | 额度操作 |
| `relaycommon` | RelayInfo |
| `types` | 错误类型 |

## 3. 函数详解

### `ReturnPreConsumedQuota(c, relayInfo)`
退还预扣费额度（异步执行）

### `PreConsumeQuota(c, preConsumedQuota, relayInfo) *types.NewAPIError`
旧版预扣费逻辑：
1. 查询用户额度
2. 信任额度检查（用户额度+令牌额度都充足时跳过预扣）
3. 预扣令牌额度
4. 扣减用户额度

## 4. 关键逻辑分析

1. **信任旁路**：用户和令牌额度都超过阈值时跳过预扣
2. **异步退还**：使用 gopool.Go 异步执行退款

## 5. 关联文件

- `billing_session.go` — 新版计费会话
- `quota.go` — PostConsumeQuota
