# cache.go 代码阅读文档

## 1. 全局总结

该文件实现了一个 HTTP 缓存控制中间件 `Cache`，根据请求路径设置不同的 `Cache-Control` 响应头。根路径 `/` 设置为 `no-cache`（不缓存），其他路径设置为一周的缓存时间。同时添加了 `Cache-Version` 头用于缓存版本控制。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `github.com/gin-gonic/gin` | Gin Web 框架 |

## 3. 类型定义

无自定义类型定义。

## 4. 函数详解

### `Cache() func(c *gin.Context)`

- **功能**：创建一个缓存控制中间件。
- **逻辑**：
  - 请求 URI 为 `/` 时：设置 `Cache-Control: no-cache`
  - 其他请求 URI：设置 `Cache-Control: max-age=604800`（604800 秒 = 7 天）
  - 所有请求都设置 `Cache-Version` 头为固定的 SHA-256 哈希值
- **返回值**：`func(c *gin.Context)`

## 5. 关键逻辑分析

- **差异化缓存策略**：根路径通常是 API 健康检查或入口页面，不应被浏览器缓存；其他静态资源路径（如前端 JS/CSS）允许长期缓存。
- **缓存版本控制**：`Cache-Version` 头使用固定的哈希值，当需要强制刷新客户端缓存时，修改此值即可使所有客户端重新请求资源，无需等待 `max-age` 过期。
- **注意**：该中间件仅设置响应头，不做实际缓存存储，依赖浏览器或 CDN 的缓存机制。

## 6. 关联文件

- `middleware/disable-cache.go` — 完全禁用缓存的中间件，与本中间件互补
