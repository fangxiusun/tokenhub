# rate-limit.go 代码阅读文档

## 1. 全局总结

本文件实现了一个基于内存的速率限制器（In-Memory Rate Limiter），用于在不依赖外部存储（如 Redis）的情况下对请求进行频率控制。采用滑动窗口算法，通过记录每个 key 的请求时间戳队列来判断是否超过限制。包含自动清理过期条目的后台协程机制。

## 2. 依赖关系

| 包名 | 用途 |
|------|------|
| `sync` | 提供互斥锁（Mutex），保证并发安全 |
| `time` | 提供时间戳获取和定时清理功能 |

本文件不依赖项目内其他模块，是独立的工具组件。

## 3. 类型定义

### InMemoryRateLimiter

```go
type InMemoryRateLimiter struct {
    store              map[string]*[]int64  // 请求时间戳存储，key 为限流标识
    mutex              sync.Mutex           // 互斥锁，保护 store 的并发访问
    expirationDuration time.Duration        // 过期时长，用于清理和判断
}
```

**字段说明：**
- `store`：以字符串 key 为索引，值为指向 `[]int64` 的指针。每个 key 对应一个时间戳队列，记录该 key 的最近请求时间（Unix 时间戳，秒）。
- `mutex`：保护 `store` 的读写操作，确保并发安全。
- `expirationDuration`：配置的过期时长，同时用于清理周期和判断条目是否过期。

## 4. 函数详解

### Init(expirationDuration time.Duration)

```go
func (l *InMemoryRateLimiter) Init(expirationDuration time.Duration)
```

初始化速率限制器。采用双重检查锁定（Double-Checked Locking）模式确保只初始化一次：
1. 先检查 `store` 是否为 nil（无锁快路径）
2. 加锁后再次检查（防止并发重复初始化）
3. 创建 map，设置过期时长
4. 如果过期时长大于 0，启动后台清理协程

### clearExpiredItems()

```go
func (l *InMemoryRateLimiter) clearExpiredItems()
```

后台清理协程，周期性删除过期的 key：
- 每隔 `expirationDuration` 执行一次清理
- 遍历所有 key，判断最后一个时间戳是否已过期
- 过期条件：队列为空，或当前时间与队列最后一个元素的差值超过过期时长（秒）
- 注意：此协程永远不会退出（无限循环）

### Request(key string, maxRequestNum int, duration int64) bool

```go
func (l *InMemoryRateLimiter) Request(key string, maxRequestNum int, duration int64) bool
```

核心限流判断方法，基于滑动窗口算法：

**参数：**
- `key`：限流标识（如用户 ID、IP 地址等）
- `maxRequestNum`：时间窗口内允许的最大请求数
- `duration`：时间窗口大小（单位：秒）

**逻辑流程：**
1. 加锁并记录当前时间戳
2. 如果 key 已存在：
   - 若队列长度 < 最大请求数：直接添加时间戳，允许请求
   - 若队列已满：检查队列最早时间戳是否已过窗口期
     - 已过期：移除最早时间戳，添加新时间戳，允许请求
     - 未过期：拒绝请求
3. 如果 key 不存在：创建新队列，添加时间戳，允许请求
4. 解锁

**返回值：** `true` 表示允许请求，`false` 表示请求被限流。

## 5. 关键逻辑分析

### 5.1 滑动窗口算法

该实现采用的是**固定窗口 + 队列**的变体滑动窗口算法：
- 队列存储每个请求的精确时间戳
- 当队列满时，检查最旧请求是否已超出窗口
- 这比简单固定窗口更精确，能避免窗口边界处的突发流量

### 5.2 并发安全策略

- 使用 `sync.Mutex` 保护所有对 `store` 的读写操作
- `Init` 方法使用双重检查锁定模式避免重复初始化
- `clearExpiredItems` 在清理时也持有锁

### 5.3 内存管理

- 使用 `map[string]*[]int64` 而非 `map[string][]int64`，通过指针减少 map 值的复制开销
- 清理协程周期性运行，防止内存无限增长
- 清理条件基于最后一个时间戳而非第一个，可能在某些情况下保留更多数据

### 5.4 潜在问题

1. **清理协程永不退出**：没有停止机制，如果速率限制器被销毁，协程会继续运行
2. **时间戳精度**：使用 Unix 秒级时间戳，对于高并发场景精度可能不足
3. **清理条件**：只检查最后一个时间戳，如果中间有大量过期数据但最后一个未过期，不会被清理

## 6. 关联文件

- `new-api/middleware/ratelimiter/` — 速率限制中间件，可能使用本组件
- `new-api/common/redis.go` — Redis 版本的速率限制实现
- `new-api/setting/` — 配置管理，可能包含速率限制相关配置
