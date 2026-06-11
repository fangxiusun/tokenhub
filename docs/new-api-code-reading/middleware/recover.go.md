# recover.go 代码阅读文档

## 1. 全局总结

该文件实现了 Relay 请求的 panic 恢复中间件 `RelayPanicRecover`，捕获处理链中的 panic 错误，记录日志和堆栈信息，返回 500 错误响应，防止进程崩溃。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `fmt` | 格式化错误消息 |
| `net/http` | HTTP 状态码 |
| `runtime/debug` | 获取 panic 堆栈信息 |
| `github.com/QuantumNous/new-api/common` | 系统日志记录 |
| `github.com/gin-gonic/gin` | Gin Web 框架 |

## 3. 类型定义

无自定义类型定义。

## 4. 函数详解

### `RelayPanicRecover() gin.HandlerFunc`

- **功能**：创建 panic 恢复中间件。
- **执行流程**：
  1. 使用 `defer` + `recover()` 捕获 panic
  2. 记录 panic 错误信息到系统日志
  3. 记录完整的堆栈跟踪信息
  4. 返回 500 状态码和错误 JSON 响应
  5. 调用 `c.Abort()` 终止后续处理

## 5. 关键逻辑分析

- **错误格式**：返回 OpenAI 兼容的错误格式，包含 `message` 和 `type`（`new_api_panic`）字段。
- **堆栈记录**：使用 `debug.Stack()` 获取完整的 goroutine 堆栈信息，便于调试。
- **错误消息**：包含 GitHub issue 链接，引导用户报告问题。
- **日志级别**：使用 `SysLog` 记录 panic 事件，便于运维监控。

## 6. 关联文件

- `common/log.go` — `SysLog` 系统日志函数
- `router/router.go` — 该中间件在路由中的注册位置
