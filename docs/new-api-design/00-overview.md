# New-API 系统设计详细参考文档

## 1. 整体架构与目录结构

### 1.1 核心分层架构
项目采用经典的 **分层架构**，数据流向清晰，各层职责明确，便于维护和扩展。

**请求处理流程:**
```
HTTP Request
    │
    ▼
[中间件链] (middleware/)
    │  RequestId → PoweredBy → I18n → Logger → Session
    │  → CORS → TokenAuth / UserAuth / AdminAuth / RootAuth
    │  → Distribute (渠道路由) → RateLimit → ModelRequestRateLimit
    │  → BodyStorageCleanup
    ▼
[路由层]  (router/)
    │  将 URL 路径映射到控制器函数
    │  Relay 路由: /v1/chat/completions, /v1/messages, /v1beta/models/*
    │  API 路由:   /api/user/*, /api/channel/*, /api/token/* 等
    ▼
[控制器层]  (controller/)
    │  解析请求 → 调用服务/模型 → 返回 JSON
    │  对于中转请求: controller.Relay() 编排整个中转流水线
    ▼
[服务层]  (service/)
    │  核心业务逻辑: 计费、配额、渠道路由、Token 计数
    │  对于中转请求: PreConsumeBilling → 重试循环 → SettleBilling → LogTaskConsumption
    ▼
[模型层]  (model/)
    │  GORM 数据库操作 + Redis 缓存层
    │  支持批量更新以应对高并发写入
    ▼
[中转适配层]  (relay/channel/<provider>/)
    │  Init → GetRequestURL → SetupRequestHeader → ConvertRequest → DoRequest → DoResponse
    ▼
上游 AI 提供商 (OpenAI, Claude, Gemini 等)
```

### 1.2 技术栈
| 层级 | 技术选型 | 说明 |
|---|---|---|
| **语言** | Go 1.25+ | 高性能并发处理 |
| **Web 框架** | Gin v1.9 | 轻量级高性能 HTTP 框架 |
| **ORM** | GORM v2 | 支持 MySQL, PostgreSQL, SQLite |
| **缓存** | Redis (go-redis/v8) + 内存缓存 | 分布式缓存与本地缓存结合 |
| **鉴权** | Session (Cookie), JWT, WebAuthn, OAuth | 多种鉴权方式支持 |
| **前端** | React 19, TypeScript, Rsbuild | 现代化前端技术栈 |
| **UI 框架** | Base UI, Tailwind CSS v4 | 原子化 CSS + 无样式组件库 |
| **状态管理** | Zustand | 轻量级状态管理 |
| **表单** | React Hook Form + Zod | 类型安全的表单处理 |
| **国际化** | 后端: go-i18n; 前端: i18next | 多语言支持 |
| **支付** | Stripe, Creem, Waffo, ePay | 多种支付网关集成 |
| **容器化** | Docker + docker-compose | 标准化部署 |

### 1.3 核心目录结构详解
| 目录 | 文件数 | 职责描述 |
|---|---|---|
| `main.go` | 1 | 应用入口，初始化所有资源（DB, Redis, i18n, OAuth），启动后台协程，配置 Gin 服务器 |
| `router/` | 6 | HTTP 路由注册层，按功能拆分：API管理、中转代理、仪表盘、视频、Web静态资源 |
| `controller/` | 73 | 请求处理器（约 559 个导出函数），涵盖用户管理、OAuth、渠道路由、计费、模型等全部业务接口 |
| `service/` | 56+ | 业务逻辑层（约 457 个导出函数），包含计费结算、渠道路由、Token 计数、任务管理、支付集成等 |
| `model/` | 39 | 数据模型与数据库交互层，GORM 模型定义、Redis 缓存、批量更新、数据库迁移 |
| `relay/` | 30+ | AI API 中转核心，包含适配器接口定义、32 个厂商适配器、10 个任务适配器 |
| `middleware/` | 24 | HTTP 中间件：鉴权、限流、CORS、请求ID、日志、渠道路由分发、请求体存储 |
| `dto/` | 29 | 数据传输对象：OpenAI, Claude, Gemini, 音频, 图像, 视频等请求/响应结构体 |
| `constant/` | 14 | 常量定义：渠道类型、API类型、缓存键、上下文键、端点类型等 |
| `types/` | 5+ | 共享类型定义：错误类型、中转格式枚举、价格数据、文件源 |
| `setting/` | 49 | 配置管理子包：计费设置、模型设置、运营设置、系统设置、性能设置等 9 个子包 |
| `common/` | 46 | 公共工具层：JSON封装、加密、Redis客户端、数据库初始化、SSRF防护、限流、分页等 |
| `i18n/` | 2+ | 后端国际化：go-i18n 集成，支持英语、简体中文、繁体中文 |
| `oauth/` | 8 | OAuth 提供商实现：GitHub, Discord, OIDC, LinuxDO, 通用 OAuth 等 |
| `pkg/` | 4子包 | 内部工具包：计费表达式引擎、混合缓存、io.net 部署客户端、性能指标 |
| `logger/` | 1 | 结构化日志：基于 Zap 的日志系统 |
| `web/default/` | - | 默认前端主题（React 19, Rsbuild, TanStack Router） |
| `web/classic/` | - | 经典前端主题（React 18, Vite, Semi Design） |

---

## 关联文件列表

### 根目录文件
- `main.go` - 应用入口
- `go.mod` / `go.sum` - Go 模块依赖
- `Dockerfile` / `docker-compose.yml` - 容器化配置
- `makefile` - 构建脚本
- `.env.example` - 环境变量示例

### 核心目录
- `router/` - 路由层
- `controller/` - 控制器层
- `service/` - 服务层
- `model/` - 模型层
- `relay/` - 中转层
- `middleware/` - 中间件层
- `common/` - 公共工具层
- `dto/` - 数据传输对象
- `constant/` - 常量定义
- `types/` - 类型定义
- `setting/` - 配置管理
- `i18n/` - 国际化
- `oauth/` - OAuth 实现
- `pkg/` - 内部工具包
- `logger/` - 日志系统
- `web/` - 前端代码
