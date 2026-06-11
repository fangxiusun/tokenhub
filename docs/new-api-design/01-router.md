# 路由层详细设计 (`router/`)

## 1. 概述
路由层负责 HTTP 路由注册，使用 Gin Web 框架。它将 URL 路径映射到控制器函数，应用中间件（鉴权、限流、gzip、CORS），并将端点组织成逻辑分组。

## 2. 文件详细说明

### 2.1 `main.go` - 路由初始化入口
**职责**: 顶层路由器初始化，调用所有子路由器，处理前端基础 URL 重定向。

**关键函数**:
- **`SetRouter(router *gin.Engine, assets ThemeAssets)`**: 入口函数，注册所有路由组（API、Dashboard、Relay、Video、Web）。如果设置了 `FRONTEND_BASE_URL`（且非主节点），则将未匹配路由重定向到外部，而不是提供嵌入式前端。

### 2.2 `api-router.go` - 管理/用户 API 路由
**职责**: 定义所有 `/api/*` 管理/用户 REST 端点。这是最大的路由文件（401行），组织了整个管理面板和用户管理 API。

**关键函数**:
- **`SetApiRouter(router *gin.Engine)`**: 注册所有 API 路由，按领域分组：
  - **系统管理**: `/setup` (初始化), `/status`, `/notice`, `/pricing`, `/rankings`, `/perf-metrics`
  - **用户认证**: register, login, 2FA, passkey, logout, OAuth state/bind
  - **用户自服务**: profile, models, tokens, passkeys, topup, pay (Stripe/Creem/Waffo/Epay), 2FA, checkin, OAuth bindings
  - **管理员用户管理**: CRUD users, manage users, admin topup, 2FA management
  - **订阅管理**: plans, purchase, admin management (balance/epay/stripe/creem/waffo-pancake)
  - **系统选项 (Root)**: system settings, payment compliance, channel affinity cache, model ratio reset
  - **自定义 OAuth (Root)**: CRUD custom OAuth providers, OIDC discovery
  - **性能管理 (Root)**: stats, disk cache, GC, logs
  - **渠道管理 (Admin)**: CRUD, test, batch ops, Ollama pull/delete, Codex OAuth, tag management
  - **Token 管理 (User)**: CRUD, batch operations, key reveal
  - **使用量日志**: all/user/self logs, stats, channel affinity usage cache
  - **兑换码 (Admin)**: CRUD redemption codes
  - **模型元数据 (Admin)**: CRUD model metadata, sync upstream
  - **供应商元数据 (Admin)**: CRUD vendor metadata

### 2.3 `relay-router.go` - AI API 中转路由
**职责**: AI API 中转/代理路由，支持 OpenAI, Claude, Gemini, Midjourney, Suno 等多种格式。

**关键函数**:
- **`SetRelayRouter(router *gin.Engine)`**: 注册中转路由：
  - **模型列表**: `GET /v1/models`, `GET /v1/models/:model`, `GET /v1beta/models`
  - **Playground**: `POST /pg/chat/completions` (用户鉴权，渠道路由)
  - **OpenAI 中转**: `/v1/completions`, `/v1/chat/completions`, `/v1/responses`, `/v1/images/*`, `/v1/embeddings`, `/v1/audio/*`, `/v1/rerank`, `/v1/moderations`, `/v1/realtime` (WebSocket)
  - **Claude 中转**: `POST /v1/messages`
  - **Gemini 中转**: `POST /v1/models/*path`, `POST /v1beta/models/*path`
  - **Midjourney**: `POST /mj/submit/*`, `GET /mj/task/*`
  - **Suno**: `POST /suno/submit/:action`, `POST /suno/fetch`

### 2.4 `dashboard.go` - OpenAI 兼容仪表盘
**职责**: OpenAI 兼容的计费仪表盘端点（旧版 `v1/dashboard/billing/*` 路径）。

**关键函数**:
- **`SetDashboardRouter(router *gin.Engine)`**: 注册：
  - `GET /dashboard/billing/subscription` - 返回订阅/配额信息
  - `GET /dashboard/billing/usage` - 返回使用量信息
  - 两者均支持 `/v1/` 前缀变体

### 2.5 `video-router.go` - 视频生成路由
**职责**: 视频生成和代理路由，支持 Kling, Jimeng 和 OpenAI 兼容视频 API。

**关键函数**:
- **`SetVideoRouter(router *gin.Engine)`**: 注册：
  - `GET /v1/videos/:task_id/content` - 视频内容代理（Token 或用户鉴权）
  - `POST /v1/video/generations`, `GET /v1/video/generations/:task_id` - OpenAI 风格视频生成
  - `POST /kling/v1/videos/text2video`, `POST /kling/v1/videos/image2video` - Kling 原生 API
  - `POST /jimeng/` - Jimeng 官方 API（带请求转换中间件）

### 2.6 `web-router.go` - 前端静态资源服务
**职责**: 嵌入式前端（React SPA）的静态文件服务，支持双主题（default/classic）。

**关键结构体**:
- **`ThemeAssets`**: 持有两个主题的嵌入式 `embed.FS` 和 index HTML 字节。

**关键函数**:
- **`SetWebRouter(router *gin.Engine, assets ThemeAssets)`**: 提供静态资产服务，带 gzip、web 限流和缓存。SPA 回退：为未匹配路由（排除 `/v1`, `/api`, `/assets`）提供 index.html。

## 3. 中间件应用顺序
路由层在注册路由时会应用以下中间件链：
1. **RequestId**: 生成唯一请求 ID
2. **PoweredBy**: 添加 Powered-By 响应头
3. **I18n**: 国际化处理
4. **Logger**: 请求日志记录
5. **Session**: 会话管理
6. **CORS**: 跨域资源共享
7. **Auth**: 鉴权（UserAuth/AdminAuth/RootAuth/TokenAuth）
8. **Distribute**: 渠道路由分发
9. **RateLimit**: 限流
10. **ModelRequestRateLimit**: 模型请求限流
11. **BodyStorageCleanup**: 请求体存储清理

---

## 关联文件列表

### 路由层核心文件
- `router/main.go` - 路由初始化入口
- `router/api-router.go` - 管理/用户 API 路由
- `router/relay-router.go` - AI API 中转路由
- `router/dashboard.go` - OpenAI 兼容仪表盘路由
- `router/video-router.go` - 视频生成路由
- `router/web-router.go` - 前端静态资源路由

### 依赖的控制器文件
- `controller/user.go` - 用户管理控制器
- `controller/channel.go` - 渠道管理控制器
- `controller/token.go` - Token 管理控制器
- `controller/relay.go` - 中转核心控制器
- `controller/model.go` - 模型管理控制器
- `controller/log.go` - 日志管理控制器
- `controller/option.go` - 系统选项控制器
- `controller/billing.go` - 计费控制器
- `controller/subscription.go` - 订阅控制器
- `controller/topup.go` - 充值控制器
- `controller/redemption.go` - 兑换码控制器

### 依赖的中间件文件
- `middleware/auth.go` - 鉴权中间件
- `middleware/distributor.go` - 渠道路由分发中间件
- `middleware/rate-limit.go` - 限流中间件
- `middleware/cors.go` - CORS 中间件
- `middleware/logger.go` - 日志中间件
- `middleware/i18n.go` - 国际化中间件
- `middleware/request-id.go` - 请求 ID 中间件
- `middleware/gzip.go` - Gzip 压缩中间件
- `middleware/body_cleanup.go` - 请求体清理中间件
