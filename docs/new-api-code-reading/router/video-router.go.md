# video-router.go 代码阅读文档

## 1. 全局总结

`router/video-router.go` 定义了视频生成相关的路由端点，支持 OpenAI 兼容视频 API、Kling 视频 API 和即梦（Jimeng）视频 API。约 52 行代码，提供视频生成、查询和代理功能。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `github.com/QuantumNous/new-api/controller` | 请求处理器（RelayTask、VideoProxy） |
| `github.com/QuantumNous/new-api/middleware` | 认证、分发、请求转换中间件 |
| `github.com/gin-gonic/gin` | HTTP 框架 |

## 3. 类型定义

本文件无自定义类型。

## 4. 函数详解

### SetVideoRouter(router *gin.Engine)

**功能**：注册所有视频相关的路由端点。

## 5. 关键逻辑分析

### 视频代理路由 `/v1/videos`

| 路由 | 中间件 | 功能 |
|------|--------|------|
| `GET /v1/videos/:task_id/content` | `TokenOrUserAuth()` | 视频内容代理（支持 Session 或 Token 认证） |

**特点**：使用 `TokenOrUserAuth()` 中间件，同时支持 Dashboard Session 认证和 API Token 认证，方便前端和 API 客户端都能访问视频内容。

### OpenAI 兼容视频路由 `/v1`

| 路由 | 功能 |
|------|------|
| `POST /v1/video/generations` | 提交视频生成任务 |
| `GET /v1/video/generations/:task_id` | 查询视频生成状态 |
| `POST /v1/videos/:video_id/remix` | 视频混音/二次编辑 |
| `POST /v1/videos` | 提交视频生成任务（OpenAI 兼容格式） |
| `GET /v1/videos/:task_id` | 查询视频生成状态 |

**中间件**：`TokenAuth()` + `Distribute()`（Token 认证 + 负载分发）

### Kling 视频路由 `/kling/v1`

| 路由 | 功能 |
|------|------|
| `POST /kling/v1/videos/text2video` | 文本生成视频 |
| `POST /kling/v1/videos/image2video` | 图片生成视频 |
| `GET /kling/v1/videos/text2video/:task_id` | 查询文本生视频状态 |
| `GET /kling/v1/videos/image2video/:task_id` | 查询图片生视频状态 |

**中间件**：`KlingRequestConvert()` + `TokenAuth()` + `Distribute()`

`KlingRequestConvert()` 中间件负责将 Kling 原生请求格式转换为系统内部格式。

### 即梦（Jimeng）视频路由 `/jimeng`

| 路由 | 功能 |
|------|------|
| `POST /jimeng/` | 提交即梦视频任务 |

**中间件**：`JimengRequestConvert()` + `TokenAuth()` + `Distribute()`

`JimengRequestConvert()` 中间件负责将即梦官方 API 格式（`?Action=CVSync2AsyncSubmitTask&Version=2022-08-31`）转换为系统内部格式。

### 架构设计特点

1. **统一任务模型**：所有视频生成都使用 `RelayTask` / `RelayTaskFetch` 处理器，通过任务 ID 追踪异步任务
2. **请求格式转换**：通过中间件（`KlingRequestConvert`、`JimengRequestConvert`）适配不同提供商的 API 格式
3. **负载分发**：所有视频路由都使用 `Distribute()` 中间件进行负载分发
4. **灵活认证**：视频代理支持 Session 和 Token 双认证模式

## 6. 关联文件

| 文件 | 关联说明 |
|------|---------|
| `router/main.go` | 调用 SetVideoRouter 注册路由 |
| `controller/task.go` | RelayTask / RelayTaskFetch 处理器实现 |
| `controller/video.go` | VideoProxy 视频代理处理器 |
| `middleware/auth.go` | TokenOrUserAuth 双认证中间件 |
| `middleware/distribute.go` | 负载分发中间件 |
| `middleware/request_convert.go` | Kling/Jimeng 请求转换中间件 |
