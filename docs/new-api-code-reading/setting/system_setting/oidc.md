# oidc.go 代码阅读文档

## 1. 全局总结

该文件定义 OIDC（OpenID Connect）登录的配置。

## 2. 依赖关系

- `github.com/QuantumNous/new-api/setting/config` — 全局配置管理器

## 3. 类型定义

| 结构体 | 字段 | 类型 | 说明 |
|--------|------|------|------|
| `OIDCSettings` | `Enabled` | `bool` | 是否启用 OIDC 登录 |
| | `ClientId` | `string` | Client ID |
| | `ClientSecret` | `string` | Client Secret |
| | `WellKnown` | `string` | Well-Known 端点 URL |
| | `AuthorizationEndpoint` | `string` | 授权端点 |
| | `TokenEndpoint` | `string` | Token 端点 |
| | `UserInfoEndpoint` | `string` | UserInfo 端点 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `GetOIDCSettings` | `func GetOIDCSettings() *OIDCSettings` | 获取 OIDC 配置 |

## 5. 关键逻辑分析

- 默认配置为空，需管理员手动配置
- 支持自定义端点或从 Well-Known 自动发现

## 6. 关联文件

- `oauth/oidc.go` — OIDC OAuth 实现
- `controller/auth.go` — 认证接口
