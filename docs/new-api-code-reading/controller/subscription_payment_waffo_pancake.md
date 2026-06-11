# subscription_payment_waffo_pancake.go 代码阅读文档

## 1. 全局总结

该文件实现了通过 Waffo Pancake 支付订阅计划的功能，包括创建支付会话和处理 Webhook 回调。

## 2. 依赖关系

- `common` — 通用工具
- `logger` — 日志
- `model` — 订阅模型
- `service` — 支付服务
- `setting` — Waffo Pancake 配置
- `gin-gonic/gin` — HTTP 框架
- `shopspring/decimal` — 精确计算

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `SubscriptionWaffoPancakePayRequest` | Waffo Pancake 订阅支付请求 |

## 4. 函数详解

### `SubscriptionRequestWaffoPancakePay(c *gin.Context)`
创建 Waffo Pancake 支付会话。

### Waffo Pancake Webhook 处理
处理支付成功的 Webhook 回调。

## 5. 关键逻辑分析

- 需要支付合规确认
- 支持 Webhook 签名验证

## 6. 关联文件

- `controller/subscription.go` — 订阅管理
