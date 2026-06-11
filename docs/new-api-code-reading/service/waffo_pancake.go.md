# waffo_pancake.go 代码阅读文档

## 1. 全局总结

该文件实现 Waffo Pancake 支付集成，包括结账会话创建、Webhook 验证、产品/商店管理等功能。是支付系统的核心组件。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `model` | 订单查询 |
| `setting` | Pancake 配置 |
| `pancake` | Pancake SDK |

## 3. 类型定义

### `WaffoPancakePriceSnapshot`
每会话价格覆盖：`Amount`、`TaxCategory`

### `WaffoPancakeCreateSessionParams`
结账会话创建参数

### `WaffoPancakeCheckoutSession`
结账会话响应：SessionID、CheckoutURL、Token 等

### `WaffoPancakeWebhookEvent` / `WaffoPancakeWebhookData`
Webhook 事件结构体

### `WaffoPancakeCatalog` / `WaffoPancakeCatalogStore` / `WaffoPancakeCatalogProduct`
产品目录结构体

### `WaffoPancakePairResult`
商店+产品创建结果

## 4. 核心函数

### `CreateWaffoPancakeCheckoutSession(ctx, params)`
创建认证模式结账会话

### `VerifyConfiguredWaffoPancakeWebhook(payload, signatureHeader)`
验证 Webhook 签名

### `ResolveWaffoPancakeTradeNo(event)`
解析 Webhook 事件为本地订单号

### `CreateWaffoPancakePrimaryStore/Product/Pair`
创建商店/产品/商店+产品对

### `ListWaffoPancakeCatalog(ctx, merchantID, privateKey)`
查询 Pancake 产品目录

### `SaveWaffoPancakeConfig(ctx, ...)`
持久化 Pancake 配置

## 5. 关键逻辑分析

1. **认证模式**：BuyerIdentity 绑定用户，确保订单归属
2. **幂等结账**：OrderMerchantExternalID 用于去重
3. **身份验证**：Webhook 处理时验证 buyer identity 匹配
4. **目录过滤**：仅显示 active 状态的产品

## 6. 关联文件

- `setting/waffo_pancake.go` — Pancake 配置
- `model/topup.go` — 充值订单
