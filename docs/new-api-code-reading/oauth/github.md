# github.go 代码阅读文档

## 1. 全局概述

本文件实现了 GitHub OAuth 提供者，支持通过 GitHub 账号登录和绑定。使用 GitHub 的 OAuth 2.0 授权码流程。

## 2. 依赖关系

- `bytes` — 字节缓冲区
- `context` — 上下文
- `encoding/json` — JSON 编码
- `fmt` — 格式化输出
- `io` — I/O 操作
- `net/http` — HTTP 客户端
- `strconv` — 字符串转换
- `time` — 超时控制
- `github.com/QuantumNous/new-api/common` — 通用工具
- `github.com/QuantumNous/new-api/i18n` — 国际化
- `github.com/QuantumNous/new-api/logger` — 日志
- `github.com/QuantumNous/new-api/model` — 数据模型
- `github.com/gin-gonic/gin` — Web 框架

## 3. 类型定义

### GitHubProvider 结构体

```go
type GitHubProvider struct{}
```

GitHub OAuth 提供者实现。

### gitHubOAuthResponse 结构体

```go
type gitHubOAuthResponse struct {
    AccessToken string `json:"access_token"`
    Scope       string `json:"scope"`
    TokenType   string `json:"token_type"`
}
```

GitHub OAuth 令牌响应。

### gitHubUser 结构体

```go
type gitHubUser struct {
    Id    int64  `json:"id"`
    Login string `json:"login"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

GitHub 用户信息。

## 4. 函数详情

### init

```go
func init() {
    Register("github", &GitHubProvider{})
}
```

在包初始化时自动注册 GitHub 提供者。

### GetName

```go
func (p *GitHubProvider) GetName() string { return "GitHub" }
```

### IsEnabled

```go
func (p *GitHubProvider) IsEnabled() bool {
    return common.GitHubOAuthEnabled
}
```

### ExchangeToken

```go
func (p *GitHubProvider) ExchangeToken(ctx context.Context, code string, c *gin.Context) (*OAuthToken, error)
```

将授权码交换为访问令牌：
1. 向 `https://github.com/login/oauth/access_token` 发送 POST 请求
2. 使用 JSON 格式传递 client_id、client_secret、code
3. 解析响应获取 access_token

### GetUserInfo

```go
func (p *GitHubProvider) GetUserInfo(ctx context.Context, token *OAuthToken) (*OAuthUser, error)
```

获取 GitHub 用户信息：
1. 向 `https://api.github.com/user` 发送 GET 请求
2. 使用 Bearer Token 认证
3. 解析用户 ID、Login、Name、Email

### IsUserIDTaken

```go
func (p *GitHubProvider) IsUserIDTaken(providerUserID string) bool {
    return model.IsGitHubIdAlreadyTaken(providerUserID)
}
```

### FillUserByProviderID

```go
func (p *GitHubProvider) FillUserByProviderID(user *model.User, providerUserID string) error {
    user.GitHubId = providerUserID
    return user.FillUserByGitHubId()
}
```

### SetProviderUserID

```go
func (p *GitHubProvider) SetProviderUserID(user *model.User, providerUserID string) {
    user.GitHubId = providerUserID
}
```

### GetProviderPrefix

```go
func (p *GitHubProvider) GetProviderPrefix() string { return "github_" }
```

## 5. 关键逻辑分析

### 用户标识策略

- **主标识**：使用 GitHub 数字 ID（`int64`），永不改变
- **兼容标识**：在 `Extra` 中存储 `legacy_id`（Login），用于从旧版本迁移

### HTTP 超时

- `ExchangeToken` 和 `GetUserInfo` 都设置 20 秒超时
- 使用 `http.NewRequestWithContext` 支持请求取消

### 错误处理

- 空授权码 → `MsgOAuthInvalidCode`
- 空访问令牌 → `MsgOAuthTokenFailed`
- 非 200 状态码 → `MsgOAuthGetUserErr`（附带状态码信息）
- 空 ID 或 Login → `MsgOAuthUserInfoEmpty`

### 配置来源

- `common.GitHubOAuthEnabled` — 是否启用
- `common.GitHubClientId` — Client ID
- `common.GitHubClientSecret` — Client Secret

## 6. 相关文件

- `oauth/provider.go` — Provider 接口
- `oauth/types.go` — OAuthToken、OAuthUser 类型
- `oauth/registry.go` — 注册表
- `model/user.go` — `FillUserByGitHubId` 方法
- `common/global.go` — GitHub 配置变量
