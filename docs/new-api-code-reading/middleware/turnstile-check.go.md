# turnstile-check.go 代码阅读文档

## 1. 全局总结

该文件实现了 Cloudflare Turnstile 人机验证中间件 `TurnstileCheck`，在启用时检查用户是否通过了 Turnstile 验证。支持 session 缓存验证结果，避免重复验证。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `encoding/json` | 解析 Turnstile 验证响应 |
| `net/http` | HTTP 请求和状态码 |
| `net/url` | URL 编码表单数据 |
| `github.com/QuantumNous/new-api/common` | 配置项和日志 |
| `github.com/gin-contrib/sessions` | Session 管理 |
| `github.com/gin-gonic/gin` | Gin Web 框架 |

## 3. 类型定义

### `turnstileCheckResponse`

```go
type turnstileCheckResponse struct {
    Success bool `json:"success"`
}
```

## 4. 函数详解

### `TurnstileCheck() gin.HandlerFunc`

- **功能**：创建 Turnstile 验证中间件。
- **执行流程**：
  1. 检查 `TurnstileCheckEnabled` 配置是否启用
  2. 检查 session 中是否有验证通过标记
  3. 从 URL query 获取 `turnstile` token
  4. 向 Cloudflare 验证端点发送 POST 请求
  5. 验证成功 → 保存 session 标记，放行
  6. 验证失败 → 返回错误信息

## 5. 关键逻辑分析

- **Session 缓存**：验证通过后在 session 中设置 `turnstile: true`，后续请求无需重复验证。
- **Cloudflare 集成**：使用 `https://challenges.cloudflare.com/turnstile/v0/siteverify` 端点进行服务端验证。
- **验证参数**：
  - `secret`：Turnstile 密钥
  - `response`：客户端提交的 token
  - `remoteip`：客户端 IP 地址
- **配置开关**：通过 `common.TurnstileCheckEnabled` 控制是否启用，未启用时直接放行。
- **错误处理**：网络错误和验证失败都返回 200 状态码 + `success: false`，而非 HTTP 错误码。

## 6. 关联文件

- `common/config.go` — `TurnstileCheckEnabled` 和 `TurnstileSecretKey` 配置
- `router/router.go` — 中间件注册位置
