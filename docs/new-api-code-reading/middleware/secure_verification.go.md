# secure_verification.go 代码阅读文档

## 1. 全局总结

该文件实现了安全验证状态检查中间件，包含两个中间件：`SecureVerificationRequired`（强制验证）和 `OptionalSecureVerification`（可选验证）。用于敏感操作前检查用户是否在 5 分钟内通过了安全验证（如密码验证、邮箱验证等）。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `net/http` | HTTP 状态码 |
| `time` | 验证过期时间计算 |
| `github.com/gin-contrib/sessions` | Session 管理 |
| `github.com/gin-gonic/gin` | Gin Web 框架 |

## 3. 类型定义

无自定义类型定义。

## 4. 常量定义

| 常量名 | 值 | 说明 |
|--------|-----|------|
| `SecureVerificationSessionKey` | `"secure_verified_at"` | 验证时间戳的 session key |
| `secureVerificationMethodSessionKey` | `"secure_verified_method"` | 验证方式的 session key |
| `SecureVerificationTimeout` | `300` | 验证有效期（5 分钟） |

## 5. 函数详解

### `SecureVerificationRequired() gin.HandlerFunc`

- **功能**：强制安全验证中间件，未验证或验证过期时返回错误。
- **检查流程**：
  1. 检查用户是否已登录（`id > 0`）
  2. 从 session 获取验证时间戳
  3. 验证时间戳不存在 → 返回 403 `VERIFICATION_REQUIRED`
  4. 时间戳格式错误 → 清除 session，返回 403 `VERIFICATION_INVALID`
  5. 验证已过期（> 5分钟） → 清除 session，返回 403 `VERIFICATION_EXPIRED`
  6. 验证有效 → 放行

### `OptionalSecureVerification() gin.HandlerFunc`

- **功能**：可选安全验证中间件，设置 context 标记但不阻止请求。
- **逻辑**：无论验证状态如何，都继续处理请求，在 context 中设置 `secure_verified` 标记。

### `ClearSecureVerification(c *gin.Context)`

- **功能**：清除安全验证状态，用于登出或强制重新验证场景。

### `clearSecureVerificationSession(session sessions.Session)`

- **功能**：内部函数，清除 session 中的验证相关数据。

## 6. 关键逻辑分析

- **Session 存储**：验证状态存储在服务端 session 中，比 JWT 更安全（可主动失效）。
- **5 分钟有效期**：平衡安全性和用户体验，敏感操作需要定期重新验证。
- **错误码设计**：使用结构化错误码（`VERIFICATION_REQUIRED`/`VERIFICATION_INVALID`/`VERIFICATION_EXPIRED`），便于前端区分处理。
- **可选模式**：`OptionalSecureVerification` 用于需要区分验证状态但不强制阻止的场景（如显示不同的 UI）。

## 7. 关联文件

- `controller/auth.go` — 设置验证状态的控制器逻辑
- `router/router.go` — 中间件注册位置
