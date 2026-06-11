# billing.go 代码阅读文档

## 1. 全局总结

该文件实现了 OpenAI 兼容的计费（Billing）相关接口，包括获取订阅信息和使用量统计。核心功能是将内部额度（quota）按照不同的展示类型（USD/CNY/TOKENS）转换为前端可读的金额值。

## 2. 依赖关系

- `common` — 配额单位常量、Token 状态开关
- `model` — 用户/Token 额度查询
- `setting/operation_setting` — 额度展示类型、汇率配置
- `types` — OpenAI 错误格式
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

无自定义类型定义。使用了 `model.Token` 等外部类型。

## 4. 函数详解

### `GetSubscription(c *gin.Context)`
获取当前用户的订阅/额度信息。根据 `DisplayTokenStatEnabled` 开关决定从 Token 或 User 维度读取额度。支持三种额度展示模式：USD（除以 QuotaPerUnit）、CNY（先转 USD 再乘汇率）、TOKENS（直接使用 token 数量）。返回 `OpenAISubscriptionResponse` 结构。

### `GetUsage(c *gin.Context)`
获取当前用户的使用量统计。同样支持三种展示模式，返回 `OpenAIUsageResponse`，其中 `TotalUsage` 以百分之一美元为单位（乘以 100）。

## 5. 关键逻辑分析

- 额度计算：`quota = remainQuota + usedQuota`，展示金额根据 `GetQuotaDisplayType()` 进行汇率转换
- 无限额度 Token 的特殊处理：`amount` 直接设为 1 亿
- 过期时间 `expiredTime` 小于等于 0 时归零

## 6. 关联文件

- `controller/channel-billing.go` — 定义了 `OpenAISubscriptionResponse` 和 `OpenAIUsageResponse` 类型
- `model/token.go` — Token 模型和额度查询
- `setting/operation_setting/` — 额度展示配置
