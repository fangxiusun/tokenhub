# http_client.go 代码阅读文档

## 1. 全局总结

该文件管理全局 HTTP 客户端，支持代理配置（HTTP/HTTPS/SOCKS5），实现连接池管理和 SSRF 防护。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | TLS 配置、连接池参数 |
| `system_setting` | 获取设置 |
| `proxy` | SOCKS5 代理 |

## 3. 类型定义

全局变量：
- `httpClient` — 默认 HTTP 客户端
- `proxyClients` — 代理客户端缓存（map[proxyURL]*http.Client）

## 4. 函数详解

### `InitHttpClient()`
初始化默认 HTTP 客户端，配置：
- 连接池参数
- HTTP/2 支持
- 环境变量代理
- TLS 配置
- 超时设置
- 重定向检查（SSRF 防护）

### `GetHttpClient() *http.Client`
获取默认客户端

### `GetHttpClientWithProxy(proxyURL) (*http.Client, error)`
获取支持代理的客户端

### `NewProxyHttpClient(proxyURL) (*http.Client, error)`
创建代理客户端：
- HTTP/HTTPS：标准 HTTP 代理
- SOCKS5：SOCKS5 代理拨号器
- 缓存已创建的客户端

### `ResetProxyClientCache()`
清空代理客户端缓存

### `checkRedirect(req, via) error`
重定向检查：
- SSRF 防护验证
- 最多 10 次重定向

## 5. 关键逻辑分析

1. **客户端缓存**：代理客户端按 URL 缓存，避免重复创建
2. **连接池复用**：全局共享连接池配置
3. **重定向防护**：每次重定向都进行 SSRF 验证

## 6. 关联文件

- `download.go` — 使用 HTTP 客户端下载
- `common/redis.go` — Redis 客户端
