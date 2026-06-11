# oidc.go 代码阅读文档

## 1. 全局概述

本文件实现了 OpenID Connect (OIDC) OAuth 提供者，支持通过任意兼容 OIDC 协议的身份提供商登录。OIDC 是基于 OAuth 2.0 的身份认证层，提供标准化的用户信息获取方式。

## 2. 依赖关系

- `context` — 上下文
- `encoding/json` — JSON 编码
- `fmt` — 格式化输出
- `net/http` — HTTP 客户端
- `net/url` — URL 编码
- `strings` — 字符串操作
- `time` — 超时控制
- `github.com/QuantumNous/new-api/i18n` — 国际化
- `github.com/QuantumNous/new-api/logger` — 日志
- `github.com/QuantumNous/new-api/model` — 数据模型
- `github.com/QuantumNous/new-api/setting/system_setting` — 系统设置
- `github.com/gin-gonic/gin` — Web 框架

## 3. 类型定义

### OIDCProvider 结构体

```go
type OIDCProvider struct{}
```

### oidcOAuthResponse 结构体

```go
type oidcOAuthResponse struct {
    AccessToken  string `json:"access_token"`
    IDToken      string `json:"id_token"`
    RefreshToken string `json:"refresh_token"`
    TokenType    string `json:"token_type"`
    ExpiresIn    int    `json:"expires_in"`
    Scope        string `json:"scope"`
}
```

### oidcUser 结构体

```go
type oidcUser struct {
    OpenID            string `json:"sub"`
    Email             string `json:"email"`
    Name              string `json:"name"`
    PreferredUsername string `json:"preferred_username"`
    Picture           string `json:"picture"`
}
```

## 4. 函数详情

### init

```go
func init() {
    Register("oidc", &OIDCProvider{})
}
```

### GetName / IsEnabled

标准接口实现。

### ExchangeToken

向配置的 `TokenEndpoint` 发送 POST 请求交换令牌。使用 `application/x-www-form-urlencoded` 格式。

### GetUserInfo

向配置的 `UserInfoEndpoint` 获取用户信息。验证 `sub` 和 `email` 字段不为空。

### IsUserIDTaken / FillUserByProviderID / SetProviderUserID / GetProviderPrefix

OIDC 特定的用户关联方法。

## 5. 关键逻辑分析

### OIDC 标准字段

| 字段 | OIDC 标准名 | 说明 |
|------|------------|------|
| `OpenID` | `sub` | 用户唯一标识 |
| `Email` | `email` | 邮箱地址 |
| `Name` | `name` | 显示名称 |
| `PreferredUsername` | `preferred_username` | 首选用户名 |
| `Picture` | `picture` | 头像 URL |

### 配置来源

通过 `system_setting.GetOIDCSettings()` 获取：
- `Enabled` — 是否启用
- `ClientId` — Client ID
- `ClientSecret` — Client Secret
- `TokenEndpoint` — 令牌端点 URL
- `UserInfoEndpoint` — 用户信息端点 URL

### redirect_uri 构造

```go
redirectUri := fmt.Sprintf("%s/oauth/oidc", system_setting.ServerAddress)
```

使用系统配置的服务器地址构造回调 URL。

### 用户信息验证

- 必须提供 `sub`（OpenID）和 `email` 字段
- 缺少任一字段返回 `MsgOAuthUserInfoEmpty` 错误

## 6. 相关文件

- `oauth/provider.go` — Provider 接口
- `oauth/types.go` — OAuthToken、OAuthUser 类型
- `setting/system_setting/oidc.go` — OIDC 设置
- `model/user.go` — `FillUserByOidcId` 方法
