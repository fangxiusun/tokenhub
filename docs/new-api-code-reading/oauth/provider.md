# provider.go 代码阅读文档

## 1. 全局概述

本文件定义了 OAuth 提供者的接口（Provider），是整个 OAuth 模块的核心抽象。所有 OAuth 提供者（GitHub、Discord、OIDC、LinuxDO、Generic）都必须实现此接口。

## 2. 依赖关系

- `context` — Go 标准库上下文包
- `github.com/QuantumNous/new-api/model` — 数据模型
- `github.com/gin-gonic/gin` — Web 框架

## 3. 类型定义

### Provider 接口

```go
type Provider interface {
    GetName() string
    IsEnabled() bool
    ExchangeToken(ctx context.Context, code string, c *gin.Context) (*OAuthToken, error)
    GetUserInfo(ctx context.Context, token *OAuthToken) (*OAuthUser, error)
    IsUserIDTaken(providerUserID string) bool
    FillUserByProviderID(user *model.User, providerUserID string) error
    SetProviderUserID(user *model.User, providerUserID string)
    GetProviderPrefix() string
}
```

## 4. 函数详情

### GetName

```go
func GetName() string
```

返回提供者的显示名称（如 "GitHub"、"Discord"）。

### IsEnabled

```go
func IsEnabled() bool
```

返回此 OAuth 提供者是否已启用。

### ExchangeToken

```go
func ExchangeToken(ctx context.Context, code string, c *gin.Context) (*OAuthToken, error)
```

将授权码交换为访问令牌。`gin.Context` 参数用于需要请求信息的提供者（如构造 redirect_uri）。

### GetUserInfo

```go
func GetUserInfo(ctx context.Context, token *OAuthToken) (*OAuthUser, error)
```

使用访问令牌获取用户信息。

### IsUserIDTaken

```go
func IsUserIDTaken(providerUserID string) bool
```

检查提供者用户 ID 是否已与其他账户关联。

### FillUserByProviderID

```go
func FillUserByProviderID(user *model.User, providerUserID string) error
```

通过提供者用户 ID 填充用户模型。

### SetProviderUserID

```go
func SetProviderUserID(user *model.User, providerUserID string)
```

在用户模型上设置提供者用户 ID。

### GetProviderPrefix

```go
func GetProviderPrefix() string
```

返回自动生成用户名的前缀（如 "github_"、"discord_"）。

## 5. 关键逻辑分析

### 接口设计

Provider 接口遵循了单一职责原则，每个方法对应 OAuth 流程中的一个步骤：

1. **认证流程**：`ExchangeToken` → `GetUserInfo`
2. **用户关联**：`IsUserIDTaken` → `FillUserByProviderID` / `SetProviderUserID`
3. **元数据**：`GetName`、`IsEnabled`、`GetProviderPrefix`

### 实现者

| 提供者 | 文件 | 前缀 |
|--------|------|------|
| GitHub | `github.go` | `github_` |
| Discord | `discord.go` | `discord_` |
| OIDC | `oidc.go` | `oidc_` |
| LinuxDO | `linuxdo.go` | `linuxdo_` |
| Generic | `generic.go` | `{slug}_` |

### gin.Context 的使用

`ExchangeToken` 接收 `gin.Context` 参数，这是因为某些提供者（如 LinuxDO）需要从请求中提取主机名来构造 redirect_uri。

## 6. 相关文件

- `oauth/types.go` — OAuthToken、OAuthUser 类型定义
- `oauth/registry.go` — 提供者注册表
- `oauth/github.go` — GitHub 实现
- `oauth/discord.go` — Discord 实现
- `oauth/oidc.go` — OIDC 实现
- `oauth/linuxdo.go` — LinuxDO 实现
- `oauth/generic.go` — 通用 OAuth 实现
