# secure_verification.go 代码阅读文档

## 1. 全局总结

该文件实现了通用安全验证框架，支持 2FA 和 Passkey 两种验证方式。验证成功后在 session 中记录时间戳，供后续敏感操作（如查看渠道密钥）使用。

## 2. 依赖关系

- `common` — 通用工具
- `model` — 用户/2FA/Passkey 模型
- `gin-contrib/sessions` — 会话管理
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `UniversalVerifyRequest` | 通用验证请求 |
| `VerificationStatusResponse` | 验证状态响应 |

## 4. 函数详解

### `UniversalVerify(c *gin.Context)`
通用验证入口。根据 `method` 参数选择 2FA 或 Passkey 验证方式。

### `setSecureVerificationSession(c, method) (int64, error)`
验证成功后设置 session 标记。

### `consumePasskeyReady(c) (bool, error)`
消费 Passkey Ready 标记（一次性使用）。

## 5. 关键逻辑分析

- 验证有效期：5 分钟（`SecureVerificationTimeout`）
- Passkey Ready 标记有效期：60 秒
- Passkey 验证分两步：`PasskeyVerifyFinish` 设置 Ready 标记 → `UniversalVerify` 消费标记
- session 中存储：验证时间戳、验证方法、Passkey Ready 时间戳

## 6. 关联文件

- `controller/passkey.go` — Passkey 验证流程
- `controller/channel.go` — `validateTwoFactorAuth` 2FA 验证
