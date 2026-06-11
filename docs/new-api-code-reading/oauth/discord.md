# discord.go 代码阅读文档

## 1. 全局概述

本文件实现了 Discord OAuth 提供者，支持通过 Discord 账号登录和绑定。使用 Discord API v10 的 OAuth 2.0 授权码流程。

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

### DiscordProvider 结构体

```go
type DiscordProvider struct{}
```

Discord OAuth 提供者实现。

### discordOAuthResponse 结构体

```go
type discordOAuthResponse struct {
    AccessToken  string `json:"access_token"`
    IDToken      string `json:"id_token"`
    RefreshToken string `json:"refresh_token"`
    TokenType    string `json:"token_type"`
    ExpiresIn    int    `json:"expires_in"`
    Scope        string `json:"scope"`
}
```

### discordUser 结构体

```go
type discordUser struct {
    UID  string `json:"id"`
    ID   string `json:"username"`
    Name string `json:"global_name"`
}
```

## 4. 函数详情

### init

```go
func init() {
    Register("discord", &DiscordProvider{})
}
```

### GetName

```go
func (p *DiscordProvider) GetName() string { return "Discord" }
```

### IsEnabled

```go
func (p *DiscordProvider) IsEnabled() bool {
    return system_setting.GetDiscordSettings().Enabled
}
```

### ExchangeToken

```go
func (p *DiscordProvider) ExchangeToken(ctx context.Context, code string, c *gin.Context) (*OAuthToken, error)
```

向 `https://discord.com/api/v10/oauth2/token` 发送 POST 请求交换令牌。

### GetUserInfo

```go
func (p *DiscordProvider) GetUserInfo(ctx context.Context, token *OAuthToken) (*OAuthUser, error)
```

向 `https://discord.com/api/v10/users/@me` 获取用户信息。

### IsUserIDTaken / FillUserByProviderID / SetProviderUserID / GetProviderPrefix

Discord 特定的用户关联方法。

## 5. 关键逻辑分析

### Discord API 特点

- 使用 API v10 版本
- 用户标识使用 Discord 的 Snowflake ID（字符串格式）
- 用户名使用 `username` 字段，显示名使用 `global_name` 字段

### 配置来源

通过 `system_setting.GetDiscordSettings()` 获取：
- `Enabled` — 是否启用
- `ClientId` — Client ID
- `ClientSecret` — Client Secret

### HTTP 超时

- 所有 HTTP 请求设置 5 秒超时（比 GitHub 短）

### 错误处理

错误处理模式与 GitHub 提供者一致，使用 i18n 消息键。

## 6. 相关文件

- `oauth/provider.go` — Provider 接口
- `oauth/types.go` — OAuthToken、OAuthUser 类型
- `setting/system_setting/discord.go` — Discord 设置
- `model/user.go` — `FillUserByDiscordId` 方法
