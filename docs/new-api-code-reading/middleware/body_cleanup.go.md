# body_cleanup.go 代码阅读文档

## 1. 全局总结

该文件实现了一个请求体存储清理中间件 `BodyStorageCleanup`，在每个请求处理完成后自动清理请求过程中产生的磁盘/内存缓存（包括请求体存储和文件来源缓存），防止资源泄漏。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `github.com/QuantumNous/new-api/common` | 提供 `CleanupBodyStorage` 函数，清理请求体存储 |
| `github.com/QuantumNous/new-api/service` | 提供 `CleanupFileSources` 函数，清理文件来源缓存 |
| `github.com/gin-gonic/gin` | Gin Web 框架 |

## 3. 类型定义

无自定义类型定义。

## 4. 函数详解

### `BodyStorageCleanup() gin.HandlerFunc`

- **功能**：创建一个 Gin 中间件，在请求处理完成后执行清理操作。
- **执行流程**：
  1. 调用 `c.Next()` 处理请求
  2. 请求结束后调用 `common.CleanupBodyStorage(c)` 清理请求体存储
  3. 调用 `service.CleanupFileSources(c)` 清理文件来源缓存（如 URL 下载的文件）
- **返回值**：`gin.HandlerFunc`

## 5. 关键逻辑分析

- **清理时机**：使用 `c.Next()` 后置清理模式，确保在所有后续中间件和 handler 执行完毕后才进行清理。
- **双重清理**：分别清理请求体存储和文件来源缓存，覆盖两种不同的临时数据类型。
- **资源管理**：该中间件适用于需要缓存请求体（如大文件上传、URL 下载等场景）的路由，避免磁盘或内存资源持续累积。

## 6. 关联文件

- `common/body_storage.go` — `CleanupBodyStorage` 的实现
- `service/file_source.go` — `CleanupFileSources` 的实现
