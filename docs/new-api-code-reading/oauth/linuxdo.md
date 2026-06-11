# linuxdo.go 代码阅读文档

## 1. 全局概述

本文件实现了 LinuxDO OAuth 提供者，支持通过 LinuxDO 社区账号登录。LinuxDO 是一个技术社区，本实现使用其 OAuth 2.0 接口，并包含信任等级（Trust Level）检查机制。

## 2. 依赖关系

- `context` — 上下文
- `encoding/base64` — Base64 编码
- `encoding/json` — JSON 编码
- `fmt` — 格式化输出
- `net/http` — HTTP 客户端
- `net/url` — URL 编码
- `strconv` — 字符串转换
- `strings` — 字符串操作
- `time` — 超时控制
- `github.com/QuantumNous/new-api/common` — 通用工具
- `github.com/QuantumNous/new-api/i18n` — 国际化
- `github.com/QuantumNous/new-api/logger` — 日志
- `github.com/QuantumNous/new-api/model` — 数据模型
- `github.com/gin-gonic/gin` — Web 框架

## 3. 类型定义

### LinuxDOProvider 结构体

```go
type LinuxDOProvider struct{}
```

### linuxdoUser 结构体

```go
type linuxdoUser struct {
    Id         int    `json:"id"`
    Username   string `json:"username"`
    Name       string `json:"name"`
    Active     bool   `json:"active"`
    TrustLevel int    `json:"trust_level"`
    Silenced   bool   `json:"silenced"`
}
```

### TrustLevelError 结构体

```go
type TrustLevelError struct {
    Required int
    Current  int
}
```

信任等级不足的专用错误类型。

## 4. 函数详情

### init

```go
func init() {
    Register("linuxdo", &LinuxDOProvider{})
}
```

### ExchangeToken

使用 Basic Auth 认证交换令牌：
1. 构造 `Basic base64(client_id:client_secret)` 认证头
2. 从请求中动态构造 redirect_uri（支持 HTTP/HTTPS）
3. 向 token endpoint 发送 POST 请求

### GetUserInfo

获取用户信息并检查信任等级：
1. 向 user endpoint 发送 GET 请求
2. 解析用户信息
3. 检查信任等级是否满足最低要求
4. 不满足时返回 `TrustLevelError`

## 5. 关键逻辑分析

### Basic Auth 认证

LinuxDO 使用 Basic Auth 而非 POST 参数传递客户端凭证：
```go
credentials := client_id + ":" + client_secret
basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials))
```

### redirect_uri 动态构造

从请求中提取协议和主机名：
```go
scheme := "http"
if c.Request.TLS != nil {
    scheme = "https"
}
redirectURI := fmt.Sprintf("%s://%s/api/oauth/linuxdo", scheme, c.Request.Host)
```

### 信任等级检查

- 最低信任等级由 `common.LinuxDOMinimumTrustLevel` 配置
- 低于要求时返回 `TrustLevelError`，包含所需和当前等级
- 信任等级信息存储在 `Extra` 中

### 用户信息字段

| 字段 | 说明 |
|------|------|
| `Id` | LinuxDO 用户 ID（整数） |
| `Username` | 用户名 |
| `Name` | 显示名称 |
| `Active` | 是否活跃 |
| `TrustLevel` | 信任等级 |
| `Silenced` | 是否被静音 |

### 配置来源

- `common.LinuxDOOAuthEnabled` — 是否启用
- `common.LinuxDOClientId` — Client ID
- `common.LinuxDOClientSecret` — Client Secret
- `common.LinuxDOMinimumTrustLevel` — 最低信任等级

### 环境变量覆盖

Token endpoint 和 User endpoint 支持环境变量覆盖：
- `LINUX_DO_TOKEN_ENDPOINT`
- `LINUX_DO_USER_ENDPOINT`

## 6. 相关文件

- `oauth/provider.go` — Provider 接口
- `oauth/types.go` — OAuthToken、OAuthUser 类型
- `model/user.go` — `FillUserByLinuxDOId` 方法
- `common/global.go` — LinuxDO 配置变量
