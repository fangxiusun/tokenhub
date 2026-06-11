# dashboard.go 代码阅读文档

## 1. 全局总结

`router/dashboard.go` 定义了 Dashboard 兼容路由，提供 OpenAI 风格的计费和用量查询接口。这是一个轻量级路由文件（仅 23 行），主要用于兼容旧版 API 格式（`/dashboard/billing/*` 和 `/v1/dashboard/billing/*`）。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `github.com/QuantumNous/new-api/controller` | 请求处理器（GetSubscription、GetUsage） |
| `github.com/QuantumNous/new-api/middleware` | 路由标签、限流、CORS、Token 认证 |
| `github.com/gin-contrib/gzip` | Gzip 压缩 |
| `github.com/gin-gonic/gin` | HTTP 框架 |

## 3. 类型定义

本文件无自定义类型。

## 4. 函数详解

### SetDashboardRouter(router *gin.Engine)

**功能**：注册 Dashboard 计费相关的兼容路由。

**路由配置**：
- 基路径：`/`（根路径）
- 路由标签：`"old_api"`（标记为旧版 API）
- 中间件链：
  1. `middleware.RouteTag("old_api")` — 路由标签
  2. `gzip.Gzip(gzip.DefaultCompression)` — 响应压缩
  3. `middleware.GlobalAPIRateLimit()` — 全局限流
  4. `middleware.CORS()` — 跨域支持
  5. `middleware.TokenAuth()` — Token 认证

**注册的端点**：

| 路由 | 处理器 | 功能 |
|------|--------|------|
| `GET /dashboard/billing/subscription` | `controller.GetSubscription` | 获取订阅信息 |
| `GET /v1/dashboard/billing/subscription` | `controller.GetSubscription` | 获取订阅信息（v1 前缀） |
| `GET /dashboard/billing/usage` | `controller.GetUsage` | 获取用量信息 |
| `GET /v1/dashboard/billing/usage` | `controller.GetUsage` | 获取用量信息（v1 前缀） |

## 5. 关键逻辑分析

### 双路径兼容设计

每个端点都注册了两个版本：
- `/dashboard/billing/*` — 原始路径格式
- `/v1/dashboard/billing/*` — 带 v1 前缀的格式

这种设计确保了与不同版本客户端的兼容性。

### 认证机制

所有 Dashboard 路由都使用 `TokenAuth()` 中间件进行认证，要求请求携带有效的 API Token。这与标准 OpenAI Dashboard API 的认证方式一致。

### 路由标签

使用 `"old_api"` 作为路由标签，与 `"api"` 和 `"relay"` 标签区分，便于日志追踪和流量分析。

## 6. 关联文件

| 文件 | 关联说明 |
|------|---------|
| `router/main.go` | 调用 SetDashboardRouter 注册路由 |
| `controller/` | GetSubscription、GetUsage 处理器实现 |
| `middleware/token_auth.go` | Token 认证中间件 |
| `middleware/cors.go` | CORS 中间件 |
