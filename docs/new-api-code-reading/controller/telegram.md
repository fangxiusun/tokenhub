# telegram.go 代码阅读文档

## 1. 全局总结

该文件实现了 Telegram OAuth 登录和绑定功能。使用 Telegram Bot 进行用户认证。

## 2. 依赖关系

- `common` — 通用工具、Telegram 配置
- `model` — 用户模型
- `gin-contrib/sessions` — 会话管理
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

无。

## 4. 函数详解

### `TelegramBind(c *gin.Context)`
Telegram 账号绑定处理器。验证 Telegram 授权参数 → 获取用户信息 → 绑定到当前用户。

### `TelegramOAuth(c *gin.Context)`
Telegram OAuth 登录/注册处理器。验证 state → 获取用户信息 → 创建或登录用户。

### `checkTelegramAuthorization(params, botToken) bool`
验证 Telegram 授权参数的 HMAC-SHA256 签名。

## 5. 关键逻辑分析

- 使用 HMAC-SHA256 验证 Telegram 授权数据的完整性
- 签名验证：按参数名排序 → 拼接 → 计算 HMAC
- 新用户注册时用户名格式为 `telegram_{maxUserId+1}`
- 支持邀请码机制

## 6. 关联文件

- `controller/user.go` — `setupLogin` 登录设置
