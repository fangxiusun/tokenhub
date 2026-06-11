# topup_creem.go 代码阅读文档

## 1. 全局总结

该文件实现了通过 Creem 进行充值的功能，包括创建支付会话、处理 Webhook 回调和金额查询。

## 2. 依赖关系

- `common` — 通用工具
- `logger` — 日志
- `model` — 充值订单模型
- `setting` — Creem 配置
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `CreemAdaptor` | Creem 支付适配器 |

## 4. 函数详解

### `RequestCreemPay(c *gin.Context)`
创建 Creem 支付会话进行充值。

### `RequestCreemAmount(c *gin.Context)`
查询指定金额的 Creem 实际支付价格。

### Creem Webhook
使用 HMAC-SHA256 验证签名，处理支付成功回调。

## 5. 关键逻辑分析

- 签名验证使用 HMAC-SHA256
- 签名头：`creem-signature`

## 6. 关联文件

- `controller/topup.go` — 充值管理
