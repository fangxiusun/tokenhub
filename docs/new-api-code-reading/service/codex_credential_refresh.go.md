# codex_credential_refresh.go 代码阅读文档

## 1. 全局总结

该文件实现了 Codex 通道凭证的刷新功能。当 Codex OAuth 令牌过期或即将过期时，通过 refresh_token 获取新的 access_token，并更新到数据库中。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | JSON 序列化 |
| `constant` | 通道类型常量 |
| `model` | 通道数据更新 |

## 3. 类型定义

### `CodexCredentialRefreshOptions`
刷新选项：`ResetCaches bool` — 是否重置缓存

### `CodexOAuthKey`
Codex OAuth 密钥结构体：
- `IDToken` / `AccessToken` / `RefreshToken` — OAuth 令牌
- `AccountID` / `Email` — 账户信息
- `LastRefresh` / `Expired` — 时间戳
- `Type` — 密钥类型

## 4. 函数详解

### `parseCodexOAuthKey(raw string) (*CodexOAuthKey, error)`
解析 JSON 格式的 Codex OAuth 密钥

### `RefreshCodexChannelCredential(ctx, channelID, opts) (*CodexOAuthKey, *model.Channel, error)`
凭证刷新主流程：
1. 查询通道信息并验证类型
2. 解析当前 OAuth 密钥
3. 调用 `RefreshCodexOAuthTokenWithProxy` 刷新令牌
4. 更新 access_token、refresh_token、过期时间
5. 从 JWT 中提取 account_id 和 email
6. 更新数据库
7. 可选重置缓存

## 5. 关键逻辑分析

1. **10秒超时**：刷新操作有严格的超时限制
2. **JWT 解析**：从新的 access_token 中提取账户信息
3. **缓存重置**：刷新后可选择重置通道缓存和代理客户端缓存

## 6. 关联文件

- `codex_oauth.go` — OAuth 令牌刷新实现
- `codex_credential_refresh_task.go` — 自动刷新定时任务
- `http_client.go` — 代理客户端
