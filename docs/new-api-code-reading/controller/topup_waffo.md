# topup_waffo.go 代码阅读文档

## 1. 全局总结

该文件实现了通过 Waffo（全球支付网关）进行充值的功能，支持沙箱和生产两种环境。

## 2. 依赖关系

- `common` — 通用工具
- `logger` — 日志
- `model` — 充值订单模型
- `service` — 支付服务
- `setting` — Waffo 配置
- `setting/operation_setting` — 运营设置
- `gin-gonic/gin` — HTTP 框架
- `waffo-com/waffo-go` — Waffo SDK

## 3. 类型定义

无。

## 4. 函数详解

### `getWaffoSDK() (*waffo.Waffo, error)`
获取 Waffo SDK 实例，根据沙箱/生产环境选择配置。

### `RequestWaffoPay(c *gin.Context)`
创建 Waffo 支付订单进行充值。

### `RequestWaffoAmount(c *gin.Context)`
查询指定金额的 Waffo 实际支付价格。

### Waffo Webhook
处理支付成功的 Webhook 回调。

## 5. 关键逻辑分析

- 支持沙箱/生产环境切换
- 使用 Waffo Go SDK

## 6. 关联文件

- `controller/topup.go` — 充值管理
