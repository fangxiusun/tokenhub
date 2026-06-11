# topup_stripe.go 代码阅读文档

## 1. 全局总结

该文件实现了通过 Stripe 进行充值的功能，包括创建 Checkout Session、处理 Webhook 回调和金额查询。

## 2. 依赖关系

- `common` — 通用工具
- `logger` — 日志
- `model` — 充值订单模型
- `setting` — Stripe 配置
- `setting/operation_setting` — 运营设置
- `gin-gonic/gin` — HTTP 框架
- `stripe-go/v81` — Stripe SDK

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `StripePayRequest` | Stripe 充值请求 |
| `StripeAdaptor` | Stripe 支付适配器 |

## 4. 函数详解

### `RequestStripePay(c *gin.Context)`
创建 Stripe Checkout Session 进行充值。

### `RequestStripeAmount(c *gin.Context)`
查询指定金额的 Stripe 实际支付价格。

### Stripe Webhook
处理支付成功的 Webhook，完成充值订单。

## 5. 关键逻辑分析

- 使用 Stripe Checkout Session 模式
- 支持 Webhook 签名验证
- 订单号使用随机字符串生成

## 6. 关联文件

- `controller/topup.go` — 充值管理
- `setting/` — Stripe 配置
