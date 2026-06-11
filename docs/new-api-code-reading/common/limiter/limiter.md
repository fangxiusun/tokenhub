# limiter.go 代码阅读文档

## 1. 全局概述

本文件实现了基于 Redis 的令牌桶限流器（RedisLimiter），使用 Lua 脚本保证原子性操作。支持自定义容量、速率和每次请求的令牌消耗量。

## 2. 依赖关系

- `context` — 上下文
- `embed` — 嵌入 Lua 脚本
- `fmt` — 格式化输出
- `sync` — 同步原语
- `github.com/QuantumNous/new-api/common` — 通用工具
- `github.com/go-redis/redis/v8` — Redis 客户端

## 3. 类型定义

### RedisLimiter 结构体

```go
type RedisLimiter struct {
    client         *redis.Client
    limitScriptSHA string
}
```

### Config 结构体

```go
type Config struct {
    Capacity  int64
    Rate      int64
    Requested int64
}
```

### Option 类型

```go
type Option func(*Config)
```

函数选项模式。

## 4. 函数详情

### New

```go
func New(ctx context.Context, r *redis.Client) *RedisLimiter
```

创建 RedisLimiter 单例：
1. 使用 `sync.Once` 确保只初始化一次
2. 预加载 Lua 限流脚本到 Redis
3. 保存脚本 SHA 用于后续 EVALSHA 调用

### Allow

```go
func (rl *RedisLimiter) Allow(ctx context.Context, key string, opts ...Option) (bool, error)
```

检查是否允许请求：
1. 应用默认配置（容量 10，速率 1，请求 1）
2. 应用选项函数覆盖默认值
3. 执行 Lua 限流脚本
4. 返回是否允许（true/false）

### 选项函数

- `WithCapacity(c int64) Option` — 设置桶容量
- `WithRate(r int64) Option` — 设置令牌生成速率
- `WithRequested(n int64) Option` — 设置每次请求消耗的令牌数

## 5. 关键逻辑分析

### 令牌桶算法

令牌桶算法的核心参数：

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `Capacity` | 10 | 桶的最大容量 |
| `Rate` | 1 | 每秒生成的令牌数 |
| `Requested` | 1 | 每次请求消耗的令牌数 |

### Lua 脚本原子性

限流逻辑通过嵌入的 Lua 脚本（`lua/rate_limit.lua`）执行：
- Lua 脚本在 Redis 中原子执行，避免竞态条件
- 使用 `ScriptLoad` 预加载脚本，返回 SHA 哈希
- 后续使用 `EvalSha` 执行，减少网络传输

### 单例模式

```go
var (
    instance *RedisLimiter
    once     sync.Once
)
```

使用 `sync.Once` 实现线程安全的单例模式，确保全局只有一个 RedisLimiter 实例。

### 错误处理

- Lua 脚本执行失败时返回 `false` 和包装的错误
- Redis 连接问题会通过 `EvalSha` 的错误返回

### 使用场景

- API 请求限流
- 用户级别限流
- 渠道级别限流
- 令牌使用限流

## 6. 相关文件

- `common/limiter/lua/rate_limit.lua` — Lua 限流脚本
- `middleware/rate_limit.go` — 限流中间件
- `common/redis.go` — Redis 客户端管理
