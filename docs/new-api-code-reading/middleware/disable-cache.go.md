# disable-cache.go 代码阅读文档

## 1. 全局总结

该文件实现了一个完全禁用 HTTP 缓存的中间件 `DisableCache`，通过设置多个 HTTP 头来确保浏览器和代理服务器不会缓存响应内容。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `github.com/gin-gonic/gin` | Gin Web 框架 |

## 3. 类型定义

无自定义类型定义。

## 4. 函数详解

### `DisableCache() gin.HandlerFunc`

- **功能**：创建一个禁用缓存的中间件。
- **设置的响应头**：
  - `Cache-Control: no-store, no-cache, must-revalidate, private, max-age=0`
  - `Pragma: no-cache`
  - `Expires: 0`
- **返回值**：`gin.HandlerFunc`

## 5. 关键逻辑分析

- **三重禁用策略**：
  - `Cache-Control`：综合设置 `no-store`（不存储）、`no-cache`（必须重新验证）、`must-revalidate`（过期后必须重新验证）、`private`（仅限私有缓存）、`max-age=0`（立即过期）
  - `Pragma: no-cache`：兼容 HTTP/1.0 客户端
  - `Expires: 0`：设置过期时间为 0（已过期）
- **使用场景**：适用于敏感数据接口（如用户信息、账单数据等），确保每次请求都从服务器获取最新数据。

## 6. 关联文件

- `middleware/cache.go` — 有缓存控制的中间件，与本中间件策略相反
