# payment_compliance.go 代码阅读文档

## 1. 全局总结

该文件实现了支付合规确认功能。管理员在启用支付相关功能前需要确认合规声明，记录确认时间、用户和 IP。

## 2. 依赖关系

- `common` — 通用工具
- `i18n` — 国际化错误消息
- `logger` — 日志
- `model` — Option 模型
- `setting/operation_setting` — 合规版本常量
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `PaymentComplianceRequest` | 合规确认请求 |

## 4. 函数详解

### `requirePaymentCompliance(c *gin.Context) bool`
检查支付合规是否已确认，未确认则返回 i18n 错误消息。

### `ConfirmPaymentCompliance(c *gin.Context)`
确认支付合规。记录确认状态、条款版本、时间、用户 ID 和客户端 IP。仅允许 dashboard session 认证（不允许 API access token）。

## 5. 关键逻辑分析

- 合规确认后不可通过通用设置接口修改（`option.go` 中有拦截）
- 记录完整的审计信息（版本、时间、用户、IP）
- 仅允许 session 认证，防止 API token 绕过

## 6. 关联文件

- `setting/operation_setting/` — 合规版本常量
- `controller/option.go` — 合规字段保护
