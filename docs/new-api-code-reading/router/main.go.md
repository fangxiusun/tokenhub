# main.go 代码阅读文档

## 1. 全局总结

`router/main.go` 是路由模块的入口文件，负责将所有子路由器（API、Dashboard、Relay、Video、Web）注册到 Gin 引擎上。它通过 `FRONTEND_BASE_URL` 环境变量决定是提供内置前端还是将请求重定向到外部前端地址。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `fmt` | 字符串格式化 |
| `net/http` | HTTP 状态码常量 |
| `os` | 读取环境变量 |
| `strings` | 字符串处理（去除尾部斜杠） |
| `github.com/QuantumNous/new-api/common` | 通用工具（主节点判断、日志） |
| `github.com/QuantumNous/new-api/middleware` | 中间件常量（RouteTagKey） |
| `github.com/gin-gonic/gin` | HTTP 框架 |

## 3. 类型定义

### ThemeAssets（在 web-router.go 中定义，本文件引用）

前端主题资源容器，包含默认主题和经典主题的嵌入式文件系统及首页内容。

## 4. 函数详解

### SetRouter(router *gin.Engine, assets ThemeAssets)

**功能**：总路由注册入口，协调所有子路由的注册。

**执行流程**：
1. 调用 `SetApiRouter` 注册 API 路由
2. 调用 `SetDashboardRouter` 注册 Dashboard 路由
3. 调用 `SetRelayRouter` 注册 Relay（AI API 代理）路由
4. 调用 `SetVideoRouter` 注册视频相关路由
5. 读取 `FRONTEND_BASE_URL` 环境变量：
   - 如果当前是**主节点**且环境变量不为空，忽略该变量（打印日志）
   - 如果环境变量为空，调用 `SetWebRouter` 提供内置前端
   - 如果环境变量不为空，注册 `NoRoute` 处理器，将所有未匹配的请求 **301 重定向**到外部前端

## 5. 关键逻辑分析

### 前端路由选择逻辑

```
FRONTEND_BASE_URL 存在且非空？
  ├─ 是 → 主节点？→ 忽略，使用内置前端
  └─ 否 → 使用内置前端（SetWebRouter）
  └─ 是且非主节点 → 301 重定向到外部前端
```

- **主节点特殊处理**：主节点始终使用内置前端，`FRONTEND_BASE_URL` 被忽略并记录日志
- **URL 清理**：外部前端 URL 会去除尾部 `/` 以避免双斜杠问题
- **重定向机制**：通过 `router.NoRoute` 捕获所有未匹配路由，使用 `http.StatusMovedPermanently`（301）重定向

### 路由注册顺序

按顺序注册：API → Dashboard → Relay → Video → Web，确保各子路由的优先级正确。

## 6. 关联文件

| 文件 | 关联说明 |
|------|---------|
| `router/api-router.go` | API 路由注册 |
| `router/dashboard.go` | Dashboard 兼容路由注册 |
| `router/relay-router.go` | AI API 代理路由注册 |
| `router/video-router.go` | 视频代理路由注册 |
| `router/web-router.go` | 前端静态资源路由注册 |
| `middleware/` | 中间件（RouteTagKey 常量） |
| `common/` | 通用工具（IsMasterNode、SysLog） |
