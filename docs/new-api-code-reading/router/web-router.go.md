# web-router.go 代码阅读文档

## 1. 全局总结

`router/web-router.go` 负责前端静态资源的路由注册和主题切换。支持默认主题（Default）和经典主题（Classic）两套前端，通过嵌入式文件系统（embed.FS）提供静态文件服务，并实现了 SPA 路由回退机制。约 46 行代码。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `embed` | Go 嵌入式文件系统 |
| `net/http` | HTTP 状态码和 Content-Type |
| `strings` | 字符串前缀判断 |
| `github.com/QuantumNous/new-api/common` | 嵌入文件夹封装、主题判断 |
| `github.com/QuantumNous/new-api/controller` | RelayNotFound 处理器 |
| `github.com/QuantumNous/new-api/middleware` | 限流、缓存中间件 |
| `github.com/gin-contrib/gzip` | Gzip 压缩 |
| `github.com/gin-contrib/static` | 静态文件服务 |
| `github.com/gin-gonic/gin` | HTTP 框架 |

## 3. 类型定义

### ThemeAssets

```go
type ThemeAssets struct {
    DefaultBuildFS   embed.FS  // 默认主题的嵌入式文件系统
    DefaultIndexPage []byte    // 默认主题的首页 HTML 内容
    ClassicBuildFS   embed.FS  // 经典主题的嵌入式文件系统
    ClassicIndexPage []byte    // 经典主题的首页 HTML 内容
}
```

**功能**：封装两套前端主题的资源，由 `main.go` 通过 `//go:embed` 指令注入。

## 4. 函数详解

### SetWebRouter(router *gin.Engine, assets ThemeAssets)

**功能**：注册前端静态资源路由和 SPA 回退处理。

**执行流程**：
1. 使用 `common.EmbedFolder` 将嵌入式文件系统转换为 Gin 可用的文件系统
2. 使用 `common.NewThemeAwareFS` 创建主题感知的文件系统（根据当前主题选择资源）
3. 注册全局中间件：Gzip 压缩、Web 限流、缓存
4. 使用 `static.Serve` 提供根路径 `/` 的静态文件服务
5. 注册 `NoRoute` 处理器处理 SPA 路由回退

## 5. 关键逻辑分析

### 主题切换机制

```go
defaultFS := common.EmbedFolder(assets.DefaultBuildFS, "web/default/dist")
classicFS := common.EmbedFolder(assets.ClassicBuildFS, "web/classic/dist")
themeFS := common.NewThemeAwareFS(defaultFS, classicFS)
```

- `NewThemeAwareFS` 创建一个主题感知的文件系统包装器
- 运行时通过 `common.GetTheme()` 动态选择主题
- 支持默认主题（React 19 + Rsbuild）和经典主题（React 18 + Vite）

### SPA 路由回退

```go
router.NoRoute(func(c *gin.Context) {
    // 1. API 路径返回 404 JSON
    if strings.HasPrefix(c.Request.RequestURI, "/v1") ||
       strings.HasPrefix(c.Request.RequestURI, "/api") ||
       strings.HasPrefix(c.Request.RequestURI, "/assets") {
        controller.RelayNotFound(c)
        return
    }
    // 2. 其他路径返回对应主题的 index.html
    c.Header("Cache-Control", "no-cache")
    if common.GetTheme() == "classic" {
        c.Data(http.StatusOK, "text/html; charset=utf-8", assets.ClassicIndexPage)
    } else {
        c.Data(http.StatusOK, "text/html; charset=utf-8", assets.DefaultIndexPage)
    }
})
```

**逻辑说明**：
- `/v1`、`/api`、`/assets` 前缀的请求返回 `RelayNotFound`（JSON 404）
- 其他所有路径返回对应主题的 `index.html`，实现前端路由（SPA）
- 设置 `Cache-Control: no-cache` 确保 HTML 内容始终新鲜

### 中间件配置

| 中间件 | 功能 |
|--------|------|
| `gzip.Gzip(gzip.DefaultCompression)` | 响应压缩，减少传输大小 |
| `middleware.GlobalWebRateLimit()` | Web 全局限流，防止滥用 |
| `middleware.Cache()` | 响应缓存，提升性能 |

### 嵌入式资源路径

- 默认主题：`web/default/dist`
- 经典主题：`web/classic/dist`
- 这些路径在 `main.go` 中通过 `//go:embed` 指令嵌入到二进制文件中

## 6. 关联文件

| 文件 | 关联说明 |
|------|---------|
| `router/main.go` | 调用 SetWebRouter，传递 ThemeAssets |
| `main.go` | 通过 `//go:embed` 嵌入前端资源 |
| `common/embed.go` | EmbedFolder、NewThemeAwareFS 实现 |
| `common/theme.go` | GetTheme 主题判断函数 |
| `controller/not_found.go` | RelayNotFound 处理器 |
| `middleware/rate_limit.go` | GlobalWebRateLimit 限流中间件 |
| `middleware/cache.go` | Cache 缓存中间件 |
| `web/default/` | 默认前端资源 |
| `web/classic/` | 经典前端资源 |
