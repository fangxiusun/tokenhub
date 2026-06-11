# oidc.go 代码阅读文档

## 1. 全局总结

该文件实现了 OIDC（OpenID Connect）OAuth 登录和绑定功能。支持通过 OIDC 兼容的身份提供商登录/注册和账号绑定。

## 2. 依赖关系

- `common` — 通用工具
- `model` — 用户模型
- `setting/system_setting` — OIDC 配置
- `gin-contrib/sessions` — 会话管理
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `OidcResponse` | OIDC Token 响应 |
| `OidcUser` | OIDC 用户信息（sub、email、name 等） |

## 4. 函数详解

### `getOidcUserInfoByCode(code string) (*OidcUser, error)`
通过授权码获取 OIDC 用户信息。流程：交换 Token → 获取用户信息。

### `OidcAuth(c *gin.Context)`
OIDC 登录/注册处理器。验证 state → 获取用户信息 → 创建或登录用户。

### `OidcBind(c *gin.Context)`
OIDC 账号绑定处理器。

## 5. 关键逻辑分析

- OIDC 端点从 `system_setting.GetOIDCSettings()` 配置读取
- 使用 preferred_username 作为用户名，回退为 `oidc_{maxUserId+1}`
- 新用户邮箱直接从 OIDC 响应中获取
- state 验证用于 CSRF 防护

## 6. 关联文件

- `controller/user.go` — `setupLogin` 登录设置
- `setting/system_setting/` — OIDC 配置
