# email-verification-rate-limit.go 代码阅读文档

## 1. 全局总结

该文件实现了邮件验证码发送频率限制中间件 `EmailVerificationRateLimit`，支持 Redis 和内存两种限流后端。限制策略为 30 秒内最多 2 次请求，超出限制返回 429 状态码。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `context` | Redis 操作的上下文 |
| `fmt` | 格式化错误消息 |
| `net/http` | HTTP 状态码常量 |
| `time` | 时间窗口计算 |
| `github.com/QuantumNous/new-api/common` | Redis 客户端、内存限流器、配置项 |
| `github.com/gin-gonic/gin` | Gin Web 框架 |

## 3. 类型定义

无自定义类型定义。

## 4. 常量定义

| 常量名 | 值 | 说明 |
|--------|-----|------|
| `EmailVerificationRateLimitMark` | `"EV"` | 限流键标识前缀 |
| `EmailVerificationMaxRequests` | `2` | 时间窗口内最大请求数 |
| `EmailVerificationDuration` | `30` | 时间窗口（秒） |

## 5. 函数详解

### `EmailVerificationRateLimit() gin.HandlerFunc`

- **功能**：入口函数，根据 Redis 是否可用选择限流后端。
- **逻辑**：Redis 可用时使用 `redisEmailVerificationRateLimiter`，否则使用 `memoryEmailVerificationRateLimiter`。

### `redisEmailVerificationRateLimiter(c *gin.Context)`

- **功能**：基于 Redis 的频率限制。
- **逻辑**：
  1. 使用 `INCR` 命令对键 `emailVerification:EV:{clientIP}` 进行自增
  2. 首次请求（count == 1）时设置 30 秒过期时间
  3. 超出限制（count > 2）时获取 TTL 并返回剩余等待时间
  4. Redis 操作失败时回退到内存限流器

### `memoryEmailVerificationRateLimiter(c *gin.Context)`

- **功能**：基于内存的频率限制。
- **逻辑**：调用 `inMemoryRateLimiter.Request` 进行限流判断。

## 6. 关键逻辑分析

- **双后端策略**：Redis 不可用时无缝回退到内存限流，保证服务可用性。
- **Redis 键结构**：`emailVerification:EV:{clientIP}`，以客户端 IP 为粒度限流。
- **滑动窗口**：使用 Redis `INCR` + `EXPIRE` 实现固定时间窗口限流。
- **用户体验**：超出限制时返回剩余等待时间，而非固定提示。

## 7. 关联文件

- `common/redis.go` — Redis 客户端 `RDB` 定义
- `common/rate_limiter.go` — 内存限流器 `inMemoryRateLimiter` 实现
