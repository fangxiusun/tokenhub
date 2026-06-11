# request_body_limit.go 代码阅读文档

## 1. 全局总结

该文件实现了匿名用户请求体大小限制中间件 `AnonymousRequestBodyLimit`，限制未登录用户请求体的最大字节数，防止超大请求导致的资源耗尽攻击。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `bytes` | 构建请求体缓冲区 |
| `io` | 请求体读取和限制 |
| `net/http` | HTTP 状态码 |
| `github.com/QuantumNous/new-api/common` | 限制配置、错误判断 |
| `github.com/gin-gonic/gin` | Gin Web 框架 |

## 3. 类型定义

无自定义类型定义。

## 4. 函数详解

### `AnonymousRequestBodyLimit() gin.HandlerFunc`

- **功能**：创建请求体大小限制中间件。
- **执行流程**：
  1. 获取匿名请求体大小限制配置
  2. 如果限制 <= 0 或请求体为空，跳过检查
  3. 读取请求体并检查大小
  4. 超出限制返回 413 Request Entity Too Large
  5. 读取错误返回 400 Bad Request
  6. 替换请求体并更新 `ContentLength`

### `readAnonymousRequestBody(body io.Reader, maxBytes int64) ([]byte, error)`

- **功能**：读取请求体并限制大小。
- **实现**：使用 `io.LimitReader` 读取 `maxBytes+1` 字节，如果实际读取超过 `maxBytes` 则返回 `ErrRequestBodyTooLarge` 错误。

## 5. 关键逻辑分析

- **精确限制**：读取 `maxBytes+1` 字节来检测是否超出限制，避免读取过多数据。
- **请求体重写**：读取后用 `bytes.NewReader` 重建请求体，确保后续 handler 可以正常读取。
- **ContentLength 更新**：替换请求体后更新 `ContentLength`，保持一致性。
- **错误区分**：区分请求体过大（413）和其他读取错误（400），返回不同的状态码。
- **仅限匿名用户**：从命名和设计来看，该中间件仅用于未登录用户的请求，已登录用户可能有独立的限制机制。

## 6. 关联文件

- `common/config.go` — `GetAnonymousRequestBodyLimitBytes` 配置获取函数
- `common/error.go` — `IsRequestBodyTooLargeError` 和 `ErrRequestBodyTooLarge` 定义
