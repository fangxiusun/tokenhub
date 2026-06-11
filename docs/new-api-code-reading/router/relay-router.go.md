# relay-router.go 代码阅读文档

## 1. 全局总结

`router/relay-router.go` 是 AI API 代理路由的核心文件，定义了所有 AI 模型调用的路由端点。支持 OpenAI、Claude、Gemini、Midjourney、Suno、Kling 等多种 AI 服务的统一代理接口。约 224 行代码，实现了完整的 AI API 路由体系。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `github.com/QuantumNous/new-api/constant` | 渠道类型常量 |
| `github.com/QuantumNous/new-api/controller` | 请求处理器（Relay、ListModels 等） |
| `github.com/QuantumNous/new-api/middleware` | 认证、限流、分发、性能检查中间件 |
| `github.com/QuantumNous/new-api/relay` | Midjourney 图片代理 |
| `github.com/QuantumNous/new-api/types` | 中继格式类型定义 |
| `github.com/gin-gonic/gin` | HTTP 框架 |

## 3. 类型定义

本文件无自定义类型，使用 `types.RelayFormat*` 常量指定中继格式。

## 4. 函数详解

### SetRelayRouter(router *gin.Engine)

**功能**：注册所有 AI API 中继路由。

**全局中间件**（应用于根路径）：
- `middleware.CORS()` — 跨域支持
- `middleware.DecompressRequestMiddleware()` — 请求解压
- `middleware.BodyStorageCleanup()` — 请求体存储清理
- `middleware.StatsMiddleware()` — 统计中间件

### registerMjRouterGroup(relayMjRouter *gin.RouterGroup)

**功能**：注册 Midjourney 相关的路由组，被复用于 `/mj` 和 `/:mode/mj` 两个路径。

## 5. 关键逻辑分析

### 模型列表路由 `/v1/models`

根据请求头自动识别 AI 提供商：

```
请求头判断逻辑：
├─ x-api-key + anthropic-version → Claude 模型列表
├─ x-goog-api-key 或 key 参数 → Gemini 模型检索
└─ 默认 → OpenAI 模型列表
```

**路由端点**：
- `GET /v1/models` — 根据请求头返回对应提供商的模型列表
- `GET /v1/models/:model` — 获取单个模型信息

### Gemini 兼容路由

| 路由 | 功能 |
|------|------|
| `GET /v1beta/models` | Gemini 原生模型列表 |
| `GET /v1beta/openai/models` | Gemini OpenAI 兼容模型列表 |

### Playground 路由

| 路由 | 功能 |
|------|------|
| `POST /pg/chat/completions` | Playground 聊天补全（需 UserAuth + Distribute） |

### 核心中继路由 `/v1`

**中间件链**：
1. `middleware.RouteTag("relay")` — 路由标签
2. `middleware.SystemPerformanceCheck()` — 系统性能检查
3. `middleware.TokenAuth()` — Token 认证
4. `middleware.ModelRequestRateLimit()` — 模型请求限流

#### WebSocket 路由

| 路由 | 格式 | 功能 |
|------|------|------|
| `GET /v1/realtime` | `RelayFormatOpenAIRealtime` | OpenAI Realtime WebSocket |

#### HTTP 路由

| 路由 | 格式 | 功能 |
|------|------|------|
| `POST /v1/messages` | `RelayFormatClaude` | Claude 消息接口 |
| `POST /v1/completions` | `RelayFormatOpenAI` | OpenAI 补全接口 |
| `POST /v1/chat/completions` | `RelayFormatOpenAI` | OpenAI 聊天补全 |
| `POST /v1/responses` | `RelayFormatOpenAIResponses` | OpenAI Responses 接口 |
| `POST /v1/responses/compact` | `RelayFormatOpenAIResponsesCompaction` | Responses 紧凑格式 |
| `POST /v1/edits` | `RelayFormatOpenAIImage` | 图片编辑 |
| `POST /v1/images/generations` | `RelayFormatOpenAIImage` | 图片生成 |
| `POST /v1/images/edits` | `RelayFormatOpenAIImage` | 图片编辑 |
| `POST /v1/embeddings` | `RelayFormatEmbedding` | 嵌入向量 |
| `POST /v1/audio/transcriptions` | `RelayFormatOpenAIAudio` | 音频转录 |
| `POST /v1/audio/translations` | `RelayFormatOpenAIAudio` | 音频翻译 |
| `POST /v1/audio/speech` | `RelayFormatOpenAIAudio` | 文本转语音 |
| `POST /v1/rerank` | `RelayFormatRerank` | 重排序 |
| `POST /v1/engines/:model/embeddings` | `RelayFormatGemini` | Gemini 嵌入 |
| `POST /v1/models/*path` | `RelayFormatGemini` | Gemini 模型调用 |
| `POST /v1/moderations` | `RelayFormatOpenAI` | 内容审核 |

#### 未实现路由

以下路由返回 `RelayNotImplemented`：
- 图片变体 (`/images/variations`)
- 文件管理 (`/files` 系列)
- 微调 (`/fine-tunes` 系列)
- 模型删除 (`DELETE /models/:model`)

### Midjourney 路由

**路径**：`/mj` 和 `/:mode/mj`

**中间件**：`RouteTag("relay")` + `SystemPerformanceCheck()`

**端点**（通过 `registerMjRouterGroup` 注册）：

| 路由 | 功能 |
|------|------|
| `GET /image/:id` | 获取图片（无需认证） |
| `POST /submit/action` | 提交操作 |
| `POST /submit/shorten` | 缩短提示词 |
| `POST /submit/modal` | 提交模态 |
| `POST /submit/imagine` | 图片生成 |
| `POST /submit/change` | 修改操作 |
| `POST /submit/simple-change` | 简单修改 |
| `POST /submit/describe` | 图片描述 |
| `POST /submit/blend` | 图片融合 |
| `POST /submit/edits` | 图片编辑 |
| `POST /submit/video` | 视频生成 |
| `GET /task/:id/fetch` | 获取任务 |
| `GET /task/:id/image-seed` | 获取图片种子 |
| `POST /task/list-by-condition` | 条件查询任务 |
| `POST /insight-face/swap` | 人脸替换 |
| `POST /submit/upload-discord-images` | 上传 Discord 图片 |

### Suno 路由

**路径**：`/suno`

| 路由 | 功能 |
|------|------|
| `POST /submit/:action` | 提交音乐生成任务 |
| `POST /fetch` | 获取任务状态 |
| `GET /fetch/:id` | 获取任务状态（GET 方式） |

### Gemini 原生路由

**路径**：`/v1beta`

| 路由 | 功能 |
|------|------|
| `POST /models/*path` | Gemini 模型调用 |

**中间件链**：
1. `RouteTag("relay")` — 路由标签
2. `SystemPerformanceCheck()` — 性能检查
3. `TokenAuth()` — Token 认证
4. `ModelRequestRateLimit()` — 模型限流
5. `Distribute()` — 负载分发

### 请求分发模式

所有中继路由都使用 `middleware.Distribute()` 中间件进行负载分发，支持：
- 渠道选择
- 负载均衡
- 失败重试

## 6. 关联文件

| 文件 | 关联说明 |
|------|---------|
| `router/main.go` | 调用 SetRelayRouter 注册路由 |
| `controller/relay.go` | Relay 处理器（核心中继逻辑） |
| `relay/` | AI API 代理实现（渠道适配器） |
| `types/relay.go` | RelayFormat* 类型定义 |
| `constant/channel.go` | ChannelType* 常量定义 |
| `middleware/distribute.go` | 负载分发中间件 |
| `middleware/performance.go` | 系统性能检查中间件 |
| `middleware/rate_limit.go` | 模型请求限流中间件 |
