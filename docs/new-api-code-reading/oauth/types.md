# types.go 代码阅读文档

## 1. 全局概述

本文件定义了 OAuth 模块的核心数据类型，包括令牌、用户信息和错误类型。这些类型在所有 OAuth 提供者之间共享。

## 2. 依赖关系

本文件无外部依赖。

## 3. 类型定义

### OAuthToken 结构体

```go
type OAuthToken struct {
    AccessToken  string `json:"access_token"`
    TokenType    string `json:"token_type"`
    RefreshToken string `json:"refresh_token,omitempty"`
    ExpiresIn    int    `json:"expires_in,omitempty"`
    Scope        string `json:"scope,omitempty"`
    IDToken      string `json:"id_token,omitempty"`
}
```

OAuth 提供者返回的令牌信息。

### OAuthUser 结构体

```go
type OAuthUser struct {
    ProviderUserID string
    Username       string
    DisplayName    string
    Email          string
    Extra          map[string]any
}
```

从 OAuth 提供者获取的用户信息。

### OAuthError 结构体

```go
type OAuthError struct {
    MsgKey   string
    Params   map[string]any
    RawError string
}
```

可翻译的 OAuth 错误类型。

### AccessDeniedError 结构体

```go
type AccessDeniedError struct {
    Message string
}
```

直接面向用户的访问拒绝错误。

## 4. 函数详情

### OAuthError 方法

- `Error() string` — 返回错误消息（优先使用 RawError）
- `NewOAuthError(msgKey string, params map[string]any) *OAuthError` — 创建新的 OAuth 错误
- `NewOAuthErrorWithRaw(msgKey string, params map[string]any, rawError string) *OAuthError` — 创建带原始错误的 OAuth 错误

### AccessDeniedError 方法

- `Error() string` — 返回错误消息

## 5. 关键逻辑分析

### OAuthToken 字段说明

| 字段 | 说明 |
|------|------|
| `AccessToken` | 访问令牌（必需） |
| `TokenType` | 令牌类型（如 "Bearer"） |
| `RefreshToken` | 刷新令牌（可选） |
| `ExpiresIn` | 过期时间（秒） |
| `Scope` | 授权范围 |
| `IDToken` | ID 令牌（OIDC） |

### OAuthUser 字段说明

| 字段 | 说明 |
|------|------|
| `ProviderUserID` | 提供者用户唯一标识（必需） |
| `Username` | 用户名（如 GitHub login） |
| `DisplayName` | 显示名称 |
| `Email` | 邮箱 |
| `Extra` | 额外的提供者特定数据 |

### OAuthError 的设计

`OAuthError` 支持 i18n 翻译：
- `MsgKey` — 消息键，用于翻译
- `Params` — 模板参数
- `RawError` — 原始错误（用于日志，不返回给用户）

### AccessDeniedError 的使用

`AccessDeniedError` 用于直接向用户显示拒绝消息（如自定义 OAuth 的访问策略拒绝），不同于 `OAuthError` 的翻译机制。

## 6. 相关文件

- `oauth/provider.go` — Provider 接口使用这些类型
- `oauth/generic.go` — AccessDeniedError 在通用 OAuth 中使用
- `i18n/keys.go` — OAuth 相关的消息键
