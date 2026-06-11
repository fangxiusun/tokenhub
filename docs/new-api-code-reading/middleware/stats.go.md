# stats.go 代码阅读文档

## 1. 全局总结

该文件实现了 HTTP 连接统计中间件 `StatsMiddleware`，使用原子操作跟踪当前活跃连接数，并提供 `GetStats` 函数供外部查询统计信息。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `sync/atomic` | 原子操作保证并发安全 |
| `github.com/gin-gonic/gin` | Gin Web 框架 |

## 3. 类型定义

### `HTTPStats`

```go
type HTTPStats struct {
    activeConnections int64  // 当前活跃连接数（未导出）
}
```

### `StatsInfo`

```go
type StatsInfo struct {
    ActiveConnections int64 `json:"active_connections"`  // JSON 序列化的活跃连接数
}
```

## 4. 变量定义

### `globalStats`

```go
var globalStats = &HTTPStats{}
```

- 全局单例，存储 HTTP 统计信息。

## 5. 函数详解

### `StatsMiddleware() gin.HandlerFunc`

- **功能**：创建连接统计中间件。
- **逻辑**：
  1. 请求进入时：`atomic.AddInt64(&globalStats.activeConnections, 1)`
  2. 请求结束时（通过 `defer`）：`atomic.AddInt64(&globalStats.activeConnections, -1)`

### `GetStats() StatsInfo`

- **功能**：获取当前统计信息。
- **实现**：使用 `atomic.LoadInt64` 读取活跃连接数。

## 6. 关键逻辑分析

- **线程安全**：使用 `sync/atomic` 操作保证并发安全，无需加锁。
- **精确计数**：通过 `defer` 确保即使 handler panic 也能正确减少计数。
- **轻量级实现**：仅统计活跃连接数，不记录请求延迟、状态码等信息，适合高并发场景。
- **全局单例**：使用包级变量存储统计信息，所有请求共享同一个计数器。

## 7. 关联文件

- `router/router.go` — 中间件注册位置
- `controller/dashboard.go` — 可能使用 `GetStats` 展示系统状态
