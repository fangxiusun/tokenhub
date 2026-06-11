# quota.go 代码阅读文档

## 1. 全局总结

该文件是计费系统的核心计算层，实现了音频额度计算、WebSocket 实时流计费、音频计费、令牌额度预扣/后扣、以及额度通知等功能。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | 额度单位、信任额度 |
| `constant` | 上下文键名 |
| `dto` | Usage 结构体 |
| `model` | 额度操作 |
| `billingexpr` | 分级计费 |
| `relaycommon` | RelayInfo |
| `ratio_setting` | 倍率配置 |
| `decimal` | 精确计算 |
| `gopool` | 异步任务 |

## 3. 类型定义

### `TokenDetails`
Token 详情：`TextTokens`、`AudioTokens`

### `QuotaInfo`
额度计算信息：Token 详情、模型/分组倍率、是否使用固定价格

## 4. 函数详解

### `calculateAudioQuota(info QuotaInfo) int`
音频额度计算：
- 固定价格模式：`price * quotaPerUnit * groupRatio`
- 倍率模式：`(textInput + textOutput*completionRatio + audioInput*audioRatio + audioOutput*audioRatio*audioCompletionRatio) * ratio`

### `PreWssConsumeQuota(ctx, relayInfo, usage) error`
WebSocket 实时流预扣费

### `PostWssConsumeQuota(ctx, relayInfo, modelName, usage, extraContent)`
WebSocket 实时流后结算

### `PostAudioConsumeQuota(ctx, relayInfo, usage, extraContent)`
音频请求后结算

### `PreConsumeTokenQuota(relayInfo, quota) error`
令牌额度预扣

### `PostConsumeQuota(relayInfo, quota, preConsumedQuota, sendEmail) error`
统一后结算：
1. 扣减资金来源（钱包/订阅）
2. 扣减令牌额度
3. 可选发送额度通知

### `CalcOpenRouterCacheCreateTokens(usage, priceData) int`
从 OpenRouter 成本反算缓存创建 token 数

## 5. 关键逻辑分析

1. **精确计算**：使用 decimal 库避免浮点精度问题
2. **分级计费**：支持 tiered_expr 模式的差额结算
3. **额度通知**：余额低于阈值时异步发送通知
4. **订阅通知**：订阅额度单独的通知逻辑

## 6. 关联文件

- `text_quota.go` — 文本额度计算
- `billing.go` — 计费入口
- `user_notify.go` — 通知发送
