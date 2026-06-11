# 应用入口详细设计 (`main.go`)

## 1. 概述
`main.go` 是应用程序的入口点。它初始化所有组件，设置 HTTP 服务器，并启动后台任务。

## 2. 文件详细说明

### 2.1 `main()` - 应用入口
**职责**: 应用程序入口点。

**初始化流程**:
1. **加载环境变量**: 从 `.env` 文件和系统环境变量加载配置
2. **初始化日志**: 设置日志级别和输出目标
3. **初始化数据库**: 建立 SQLite/MySQL/PostgreSQL 连接
4. **初始化 Redis**: 建立可选的 Redis 连接
5. **初始化 i18n**: 加载翻译文件，设置语言支持
6. **初始化 OAuth**: 注册所有 OAuth 提供商
7. **设置 HTTP 服务器**: 配置 Gin 路由和中间件
8. **嵌入前端资源**: 通过 Go embed 指令打包前端构建产物
9. **启动后台任务**: 启动维护任务的后台协程
10. **启动 HTTP 服务器**: 监听指定端口

### 2.2 `setupHttpServer()` - HTTP 服务器配置
**职责**: 配置 Gin HTTP 服务器。

**配置内容**:
- 创建 Gin 引擎实例
- 应用中间件链（CORS、日志、恢复、请求ID 等）
- 设置路由（API、中转、仪表盘、Web）
- 配置静态文件服务（嵌入式前端）
- 设置超时和限制

### 2.3 `startBackgroundTasks()` - 后台任务
**职责**: 启动维护任务的后台协程。

**后台任务列表**:
- **订阅过期检查**: 定期检查并处理过期订阅
- **渠道余额更新**: 定期更新所有渠道的余额
- **渠道模型同步**: 定期从上游同步模型信息
- **日志清理**: 定期清理旧日志文件
- **性能指标收集**: 定期收集系统性能数据
- **Codex 凭证刷新**: 定期刷新 Codex 渠道凭证

### 2.4 嵌入式前端
**职责**: 通过 Go embed 指令将前端构建产物打包进二进制文件。

**实现方式**:
```go
//go:embed web/default/dist/*
var defaultFrontend embed.FS

//go:embed web/classic/dist/*
var classicFrontend embed.FS
```

**主题支持**:
- 默认主题: `web/default/` - React 19, Rsbuild, Base UI
- 经典主题: `web/classic/` - React 18, Vite, Semi Design

## 3. 环境变量配置
| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `PORT` | HTTP 服务器端口 | 3000 |
| `SQL_DSN` | 数据库连接字符串 | - |
| `REDIS_CONN_STRING` | Redis 连接字符串 | - |
| `SESSION_SECRET` | 会话密钥 | - |
| `CRYPTO_SECRET` | 加密密钥 | - |
| `TZ` | 时区 | Asia/Shanghai |
| `LOG_LEVEL` | 日志级别 | info |

## 4. 启动流程图
```
main()
  │
  ├── LoadEnvironment()
  ├── InitLogger()
  ├── InitDatabase()
  ├── InitRedis()
  ├── InitI18n()
  ├── InitOAuth()
  │
  ├── setupHttpServer()
  │   ├── CreateGinEngine()
  │   ├── ApplyMiddleware()
  │   ├── SetupRoutes()
  │   └── EmbedFrontend()
  │
  ├── startBackgroundTasks()
  │   ├── SubscriptionResetTask()
  │   ├── ChannelBalanceTask()
  │   ├── ChannelModelSyncTask()
  │   ├── LogCleanupTask()
  │   └── PerfMetricsTask()
  │
  └── http.ListenAndServe()
```

---

## 关联文件列表

### 应用入口核心文件
- `main.go` - 应用入口

### 依赖的初始化文件
- `common/init.go` - 应用初始化
- `common/redis.go` - Redis 初始化
- `common/database.go` - 数据库初始化
- `i18n/i18n.go` - 国际化初始化
- `oauth/registry.go` - OAuth 注册表
- `router/main.go` - 路由初始化
- `model/main.go` - 模型层初始化

### 依赖的配置文件
- `.env.example` - 环境变量示例
- `go.mod` - Go 模块依赖
- `Dockerfile` - Docker 构建配置
- `docker-compose.yml` - Docker Compose 配置

### 依赖的公共工具文件
- `common/constants.go` - 常量定义
- `common/env.go` - 环境变量工具
- `common/logger.go` - 日志工具
- `common/system_monitor.go` - 系统监控
