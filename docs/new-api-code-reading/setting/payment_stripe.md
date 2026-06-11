# payment_stripe.go 代码阅读文档

## 1. 全局总结

该文件定义 Stripe 支付网关的配置变量。

## 2. 依赖关系

无外部依赖。

## 3. 类型定义

| 变量名 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `StripeApiSecret` | `string` | `""` | Stripe API 密钥 |
| `StripeWebhookSecret` | `string` | `""` | Stripe Webhook 签名密钥 |
| `StripePriceId` | `string` | `""` | Stripe 价格 ID |
| `StripeUnitPrice` | `float64` | `8.0` | Stripe 单价（人民币） |
| `StripeMinTopUp` | `int` | `1` | 最小充值金额 |
| `StripePromotionCodesEnabled` | `bool` | `false` | 是否启用促销码 |

## 4. 函数详解

无函数定义，仅变量声明。

## 5. 关键逻辑分析

- `StripeUnitPrice` 默认 8.0 表示 1 USD = 8 RMB 的换算汇率
- 所有变量为包级全局变量

## 6. 关联文件

- `controller/option.go` — 管理界面配置接口
- `service/stripe.go` — Stripe 支付服务实现
