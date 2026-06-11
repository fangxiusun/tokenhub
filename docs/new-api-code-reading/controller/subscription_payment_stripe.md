# subscription_payment_stripe.go 代码阅读文档

## 1. 全局总结

该文件实现了通过 Stripe 支付订阅计划的功能，包括创建 Checkout Session 和处理 Stripe Webhook 回调。

## 2. 依赖关系

- `common` — 通用工具
- `logger` — 日志
- `model` — 订阅模型
- `setting` — Stripe 配置
- `gin-gonic/gin` — HTTP 框架
- `stripe-go/v81` — Stripe SDK

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `SubscriptionStripePayRequest` | Stripe 订阅支付请求 |

## 4. 函数详解

### `SubscriptionRequestStripePay(c *gin.Context)`
创建 Stripe Checkout Session。验证计划存在且已启用 → 创建 Stripe Session → 返回支付 URL。

### Stripe Webhook 处理
处理 Stripe 支付成功的 Webhook 事件，完成订阅激活。

## 5. 关键逻辑分析

- 使用 Stripe Checkout Session 模式
- 需要支付合规确认
- 支持 Webhook 签名验证

## 6. 关联文件

- `controller/subscription.go` — 订阅管理
- `setting/` — Stripe 配置
