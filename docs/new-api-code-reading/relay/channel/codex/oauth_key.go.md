# oauth_key.go 代码阅读文档

## 1. 全局总结

本文件定义了 Codex 频道使用的 OAuth 密钥数据结构和解析逻辑。Codex 频道的 API Key 不是简单的 Bearer Token，而是一个包含多种 OAuth 令牌的 JSON 对象。

## 2. 依赖关系

| 包 | 用途 |
|---|------|
| `errors` | 错误创建 |
| `common` | JSON 反序列化工具（Unmarshal） |

## 3. 类型定义

### `OAuthKey` 结构体

```go
type OAuthKey struct {
    IDToken      string `json:"id_token,omitempty"`
    AccessToken  string `json:"access_token,omitempty"`
    RefreshToken string `json:"refresh_token,omitempty"`
    AccountID    string `json:"account_id,omitempty"`
    LastRefresh  string `json:"last_refresh,omitempty"`
    Email        string `json:"email,omitempty"`
    Type         string `json:"type,omitempty"`
    Expired      string `json:"expired,omitempty"`
}
```

Codex OAuth 密钥结构：
- `IDToken`: ID 令牌（JWT 格式）
- `AccessToken`: 访问令牌（用于 API 认证）
- `RefreshToken`: 刷新令牌（用于获取新的访问令牌）
- `AccountID`: ChatGPT 账户 ID（用于 `chatgpt-account-id` 请求头）
- `LastRefresh`: 上次刷新时间
- `Email`: 关联的邮箱地址
- `Type`: 令牌类型
- `Expired`: 过期时间

## 4. 函数详解

### `ParseOAuthKey(raw string) (*OAuthKey, error)`
解析 JSON 格式的 OAuth 密钥字符串：
1. 检查输入是否为空
2. 使用 `common.Unmarshal` 将 JSON 字符串反序列化为 `OAuthKey` 结构体
3. 返回解析后的密钥或错误信息

## 5. 关键逻辑分析

### 密钥格式要求
Codex 频道要求 API Key 以 `{` 开头（JSON 对象），这在 `adaptor.go` 的 `SetupRequestHeader` 中进行了前置校验。只有通过此格式校验后才会调用 `ParseOAuthKey`。

### 使用 `common.Unmarshal` 而非 `encoding/json`
遵循项目规则 1，使用 `common` 包的 JSON 工具函数而非直接使用 `encoding/json`。

### 两个必需字段
`access_token` 和 `account_id` 是必需的，在 `SetupRequestHeader` 中会进行空值检查。

## 6. 关联文件

- `relay/channel/codex/adaptor.go` - 在 `SetupRequestHeader` 中调用 `ParseOAuthKey` 解析密钥
- `common/json.go` - 提供 JSON 反序列化工具
