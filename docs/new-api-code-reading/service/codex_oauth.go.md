# codex_oauth.go 代码阅读文档

## 1. 全局总结

该文件实现了 Codex (OpenAI) 的 OAuth 2.0 认证流程，包括授权码交换、令牌刷新、PKCE 支持、JWT 解析等功能。是 Codex 通道凭证管理的底层实现。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | JSON 解码 |
| `crypto/rand` | 安全随机数生成 |
| `crypto/sha256` | PKCE challenge 生成 |
| `encoding/base64` | Base64 编解码 |
| `net/http` | HTTP 客户端 |
| `net/url` | URL 构建 |

## 3. 类型定义

### `CodexOAuthTokenResult`
令牌刷新结果：`AccessToken`、`RefreshToken`、`ExpiresAt`

### `CodexOAuthAuthorizationFlow`
授权流程参数：`State`、`Verifier`、`Challenge`、`AuthorizeURL`

## 4. 函数详解

### `RefreshCodexOAuthToken(ctx, refreshToken)` / `RefreshCodexOAuthTokenWithProxy`
- 使用 refresh_token 获取新的 access_token
- 支持代理配置

### `ExchangeCodexAuthorizationCode(ctx, code, verifier)` / WithProxy
- 使用授权码交换令牌
- 支持 PKCE 验证

### `CreateCodexOAuthAuthorizationFlow()`
- 生成 state、PKCE verifier/challenge
- 构建授权 URL

### `ExtractCodexAccountIDFromJWT(token)` / `ExtractEmailFromJWT(token)`
- 从 JWT 中提取 account_id 和 email
- 仅解码不验证签名

### `decodeJWTClaims(token)`
- 解析 JWT payload 部分

## 5. 关键逻辑分析

1. **PKCE 支持**：使用 S256 方法生成 code_challenge
2. **固定 Client ID**：使用硬编码的 `app_EMoamEEZ73f0CkXaXp7hrann`
3. **JWT 仅解码**：不验证签名，仅提取 claims
4. **20秒超时**：HTTP 客户端默认超时

## 6. 关联文件

- `codex_credential_refresh.go` — 凭证刷新
- `http_client.go` — 代理客户端
