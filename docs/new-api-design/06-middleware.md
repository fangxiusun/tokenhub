# 中间件层详细设计 (`middleware/`)

## 1. 概述
中间件层包含 HTTP 中间件函数，用于鉴权、限流、CORS、日志记录和请求处理。这些中间件在请求到达控制器之前执行，负责横切关注点。

## 2. 文件详细说明

### 2.1 鉴权
- **`auth.go`** (12742字节) -- 核心鉴权中间件：
  - `UserAuth()`: 要求 `Role >= RoleCommonUser`（普通用户）
  - `AdminAuth()`: 要求 `Role >= RoleAdminUser`（管理员）
  - `RootAuth()`: 要求 `Role >= RoleRootUser`（超级管理员）
  - `TokenAuth()`: 处理 API Token 鉴权
  - **鉴权流程**:
    1. 检查会话中的 `username`, `role`, `id`, `status`
    2. 如果无会话，回退到 `Authorization` 头（访问令牌）
    3. 验证用户状态（启用/禁用）和角色
    4. 设置上下文值：`username`, `role`, `id`, `group`, `use_access_token`
  - **Token 鉴权流程**:
    1. 从 `Authorization` 头、`Sec-WebSocket-Protocol`、`x-api-key` 或查询参数提取密钥
    2. 通过 `model.ValidateUserToken(key)` 验证令牌
    3. 检查状态、过期时间、剩余配额
    4. 检查 IP 限制（如果配置了）
    5. 验证用户状态和组权限
    6. 设置令牌上下文：`token_id`, `token_key`, `token_quota`, `token_model_limit` 等
- **`secure_verification.go`** (3349字节) -- `SecureVerification()`: 敏感操作的 2FA 验证。

### 2.2 限流
- **`rate-limit.go`** (6074字节) -- 全局限流：
  - 支持 Redis 和内存两种后端
  - 滑动窗口算法
  - 配置：`GlobalAPIRateLimit`（每 IP 全局 API 限流）
  - 内存限流器：时间戳基础的滑动窗口，可配置过期和后台清理
- **`model-rate-limit.go`** (6046字节) -- 每模型限流：
  - 按模型名称的请求限流
  - 支持按用户组配置
- **`email-verification-rate-limit.go`** (1893字节) -- 邮件验证限流：
  - `EmailVerificationRateLimit`: 每 60 秒 10 次
  - 用于密码重置和邮箱验证

### 2.3 请求处理
- **`distributor.go`** (19672字节) -- 渠道路由分发中间件：
  - 从请求体提取模型名称
  - 解析用户组（令牌覆盖或用户默认）
  - 支持 "auto" 组（按顺序尝试用户可访问的组）
  - 调用 `service.CacheGetRandomSatisfiedChannel()` 查询能力缓存
  - 按优先级（降序）和权重随机选择渠道
  - 成功后记录渠道亲和性
- **`request-id.go`** (708字节) -- 生成唯一请求 ID。
- **`request_body_limit.go`** (1077字节) -- 限制请求体大小。
- **`body_cleanup.go`** (514字节) -- 处理后清理请求体。
- **`gzip.go`** (1729字节) -- Gzip 压缩中间件。
- **`cors.go`** (527字节) -- CORS 中间件。
- **`cache.go`** (389字节) -- 缓存中间件。
- **`disable-cache.go`** (287字节) -- 禁用浏览器缓存。

### 2.4 日志与监控
- **`logger.go`** (820字节) -- 请求日志中间件。
- **`stats.go`** (791字节) -- 统计信息收集中间件。
- **`performance.go`** (2047字节) -- 性能监控中间件。

### 2.5 国际化
- **`i18n.go`** (1442字节) -- 国际化中间件：
  - 从 `Accept-Language` 头提取语言
  - 设置上下文语言

### 2.6 安全
- **`recover.go`** (700字节) -- Panic 恢复中间件。
- **`turnstile-check.go`** (1775字节) -- Cloudflare Turnstile CAPTCHA 验证。

### 2.7 提供商适配器
- **`jimeng_adapter.go`** (1848字节) -- 即梦请求适配器中间件。
- **`kling_adapter.go`** (1192字节) -- Kling 请求适配器中间件。

### 2.8 导航
- **`header_nav.go`** (2646字节) -- 导航头中间件。

### 2.9 工具
- **`utils.go`** (955字节) -- 中间件工具函数。

---

## 关联文件列表

### 中间件层核心文件
- `middleware/auth.go` - 核心鉴权中间件
- `middleware/secure_verification.go` - 敏感操作验证
- `middleware/rate-limit.go` - 全局限流
- `middleware/model-rate-limit.go` - 模型限流
- `middleware/email-verification-rate-limit.go` - 邮件验证限流
- `middleware/distributor.go` - 渠道路由分发
- `middleware/request-id.go` - 请求 ID 生成
- `middleware/request_body_limit.go` - 请求体限制
- `middleware/body_cleanup.go` - 请求体清理
- `middleware/gzip.go` - Gzip 压缩
- `middleware/cors.go` - CORS 处理
- `middleware/cache.go` - 缓存处理
- `middleware/disable-cache.go` - 禁用缓存
- `middleware/logger.go` - 请求日志
- `middleware/stats.go` - 统计收集
- `middleware/performance.go` - 性能监控
- `middleware/i18n.go` - 国际化
- `middleware/recover.go` - Panic 恢复
- `middleware/turnstile-check.go` - Turnstile 验证
- `middleware/jimeng_adapter.go` - 即梦适配器
- `middleware/kling_adapter.go` - Kling 适配器
- `middleware/header_nav.go` - 导航头
- `middleware/utils.go` - 工具函数

### 依赖的服务层文件
- `service/channel.go` - 渠道缓存管理
- `service/channel_select.go` - 渠道选择算法
- `service/channel_affinity.go` - 渠道亲和性
- `service/quota.go` - 配额计算

### 依赖的模型层文件
- `model/user.go` - 用户模型
- `model/token.go` - Token 模型
- `model/channel.go` - 渠道模型
- `model/ability.go` - 能力模型

### 依赖的公共工具文件
- `common/redis.go` - Redis 工具
- `common/rate-limit.go` - 内存限流器
- `common/constants.go` - 常量定义
- `common/gin.go` - Gin 工具
