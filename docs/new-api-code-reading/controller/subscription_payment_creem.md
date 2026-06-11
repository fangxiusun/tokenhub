# subscription_payment_creem.go 代码阅读文档

## 1. 全局总结

该文件实现了通过 Creem 支付订阅计划的功能，包括创建支付会话和处理 Webhook 回调。

## 2. 依赖关系

- `common` — 通用工具
- `logger` — 日志
- `model` — 订阅模型
- `setting` — Creem 配置
- `setting/operation_setting` — 运营设置
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `SubscriptionCreemPayRequest` | Creem 订阅支付请求 |

## 4. 函数详解

### `SubscriptionRequestCreemPay(c *gin.Context)`
创建 Creem 支付会话。验证计划 → 创建订单 → 调用 Creem API → 返回支付 URL。

### Creem Webhook 处理
处理 Creem 支付成功的 Webhook，使用 HMAC-SHA256 签名验证。

## 5. 关键逻辑分析

- 使用 HMAC-SHA256 进行 Webhook 签名验证
- 签名头：`creem-signature`

## 6. 关联文件

- `controller/subscription.go` — 订阅管理
