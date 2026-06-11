# service_account.go 代码阅读文档

## 1. 全局总结
Vertex AI 渠道的 Google Cloud 服务账号认证实现。负责将服务账号凭证转换为 JWT，再通过 Google OAuth2 端点换取 Access Token，并使用异步缓存管理 Token 生命周期。

## 2. 依赖关系
- **标准库**: crypto/rsa, crypto/x509, encoding/json, encoding/pem, errors, fmt, net/http, net/url, strings, time
- **内部包**:
  - `github.com/QuantumNous/new-api/relay/common` — RelayInfo
  - `github.com/QuantumNous/new-api/service` — HTTP 客户端（含代理支持）
- **外部依赖**:
  - `github.com/bytedance/gopkg/cache/asynccache` — 异步缓存
  - `github.com/golang-jwt/jwt/v5` — JWT 签发和验证

## 3. 类型定义

### Credentials
```go
type Credentials struct {
    ProjectID    string `json:"project_id"`
    PrivateKeyID string `json:"private_key_id"`
    PrivateKey   string `json:"private_key"`
    ClientEmail  string `json:"client_email"`
    ClientID     string `json:"client_id"`
}
```
Google Cloud 服务账号凭证结构体。

### Cache
```go
var Cache = asynccache.NewAsyncCache(...)
```
异步缓存实例，用于存储 Access Token：
- 刷新周期: 35 分钟
- 过期时间: 30 分钟
- 支持多渠道/多密钥索引的缓存键

## 4. 函数详解

### getAccessToken(a *Adaptor, info *relaycommon.RelayInfo) (string, error)
获取 Access Token（内部方法）：
1. 构建缓存键（支持多密钥模式）
2. 从缓存获取 Token，命中则直接返回
3. 未命中则创建签名 JWT 并换取 Token
4. 存入缓存并返回

### createSignedJWT(email, privateKeyPEM string) (string, error)
创建签名 JWT：
1. 清理 PEM 格式中的换行和头尾标记
2. 解析 PKCS8 私钥，转换为 RSA 私钥
3. 构建 JWT Claims（iss、scope、aud、exp、iat）
4. 使用 RS256 算法签名

### exchangeJwtForAccessToken(signedJWT string, info *relaycommon.RelayInfo) (string, error)
通过 Google OAuth2 端点换取 Access Token：
- 端点: `https://www.googleapis.com/oauth2/v4/token`
- 支持代理配置
- 使用 `jwt-bearer` 授权类型

### AcquireAccessToken(creds Credentials, proxy string) (string, error)
公开的 Access Token 获取方法（供外部调用），支持指定代理。

### exchangeJwtForAccessTokenWithProxy(signedJWT string, proxy string) (string, error)
带代理支持的 JWT 换 Token 实现。

## 5. 关键逻辑分析

1. **JWT 有效期**: Access Token 有效期 30 分钟，缓存刷新周期 35 分钟，确保在 Token 过期前刷新。

2. **多密钥支持**: 缓存键格式为 `access-token-{channelId}` 或 `access-token-{channelId}-{multiKeyIndex}`，支持同一渠道的多个服务账号。

3. **私钥处理**: 从 JSON 配置中提取的私钥可能包含转义的换行符 `\n`，需要统一清理后重新组装为标准 PEM 格式。

4. **代理支持**: 通过 `service.NewProxyHttpClient` 创建代理 HTTP 客户端，支持企业代理环境。

5. **Scope 固定**: JWT 的 scope 固定为 `https://www.googleapis.com/auth/cloud-platform`，授予完整的 Cloud Platform 访问权限。

## 6. 关联文件
- `vertex/adaptor.go` — 在 `SetupRequestHeader` 中调用 `getAccessToken`
- `vertex/dto.go` — Credentials 结构体在此文件中定义
- `service/http_client.go` — HTTP 客户端和代理支持
