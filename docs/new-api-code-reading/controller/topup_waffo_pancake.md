# topup_waffo_pancake.go 代码阅读文档

## 1. 全局总结

该文件实现了通过 Waffo Pancake 进行充值的功能，包括创建支付会话、处理 Webhook 回调和金额查询。

## 2. 依赖关系

- `common` — 通用工具
- `logger` — 日志
- `model` — 充值订单模型
- `service` — 支付服务
- `setting` — Waffo Pancake 配置
- `setting/operation_setting` — 运营设置
- `gin-gonic/gin` — HTTP 框架
- `shopspring/decimal` — 精确计算

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `WaffoPancakePayRequest` | Waffo Pancake 充值请求 |

## 4. 函数详解

### `RequestWaffoPancakeAmount(c *gin.Context)`
查询指定金额的 Waffo Pancake 实际支付价格。

### `RequestWaffoPancakePay(c *gin.Context)`
创建 Waffo Pancake 支付订单进行充值。

### Waffo Pancake Webhook
处理支付成功的 Webhook 回调。

## 5. 关键逻辑分析

- 使用 Waffo Pancake SDK
- 支持 Webhook 签名验证

## 6. 关联文件

- `controller/topup.go` — 充值管理
