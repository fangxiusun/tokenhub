# payment_webhook_availability.go 代码阅读文档

## 1. 全局总结

该文件提供了支付渠道（Stripe、Creem、Waffo、Waffo Pancake、Epay）的可用性检查函数。判断各支付渠道是否已正确配置并启用。

## 2. 依赖关系

- `setting` — 支付渠道配置
- `setting/operation_setting` — 支付合规和支付方式配置

## 3. 类型定义

无。

## 4. 函数详解

每个支付渠道提供三个检查函数：
- `is*TopUpEnabled()` — 是否启用充值
- `is*WebhookConfigured()` — Webhook 是否已配置
- `is*WebhookEnabled()` — Webhook 是否启用

支持的渠道：Stripe、Creem、Waffo（含沙箱模式）、Waffo Pancake、Epay。

## 5. 关键逻辑分析

- 所有渠道都需要支付合规确认（`isPaymentComplianceConfirmed()`）
- Stripe 需要 API Secret + Webhook Secret + Price ID
- Creem 需要 API Key + Products（非空且非 `[]`）
- Waffo 支持沙箱/生产两套配置
- Epay 需要 Webhook 配置 + 支付方式列表非空

## 6. 关联文件

- `setting/` — 各支付渠道配置
- `controller/payment_compliance.go` — 支付合规
