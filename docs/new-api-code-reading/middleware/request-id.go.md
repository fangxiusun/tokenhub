# request-id.go 代码阅读文档

## 1. 全局总结

该文件实现了请求 ID 生成中间件 `RequestId`，为每个请求生成唯一标识符，注入到 Gin context、请求 context 和响应头中，用于请求追踪和日志关联。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `context` | 向请求 context 注入 ID |
| `crypto/sha256` | 计算构建信息哈希 |
| `encoding/hex` | 十六进制编码 |
| `runtime/debug` | 获取构建信息 |
| `github.com/QuantumNous/new-api/common` | 时间字符串、随机字符串、常量 |
| `github.com/gin-gonic/gin` | Gin Web 框架 |

## 3. 类型定义

无自定义类型定义。

## 4. 变量定义

### `_bp` (包级变量)

```go
var _bp = func() string { ... }()
```

- **功能**：在包初始化时计算一次的构建标识符。
- **逻辑**：
  - 从 `debug.ReadBuildInfo()` 获取模块路径
  - 对路径取 SHA-256 哈希的前 4 字节（8 个十六进制字符）
  - 如果无法获取构建信息，使用 8 位随机字符串

## 5. 函数详解

### `RequestId() func(c *gin.Context)`

- **功能**：创建请求 ID 生成中间件。
- **ID 格式**：`{时间字符串}{构建标识}{8位随机字符串}`
- **注入位置**：
  1. `c.Set(common.RequestIdKey, id)` — Gin context
  2. `context.WithValue(c.Request.Context(), common.RequestIdKey, id)` — Go request context
  3. `c.Header(common.RequestIdKey, id)` — HTTP 响应头

## 6. 关键逻辑分析

- **ID 唯一性**：结合时间戳、构建标识和随机字符串，确保全局唯一。
- **构建标识**：通过模块路径的哈希值区分不同构建版本，便于排查问题。
- **三重注入**：确保在 Gin 中间件链、Go 原生 context 和 HTTP 响应中都能获取请求 ID。
- **初始化时机**：`_bp` 在包加载时计算一次，避免每次请求重复计算。

## 7. 关联文件

- `common/constant.go` — `RequestIdKey` 常量定义
- `common/util.go` — `GetTimeString` 和 `GetRandomString` 函数
