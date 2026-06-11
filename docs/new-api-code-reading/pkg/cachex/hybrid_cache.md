# hybrid_cache.go 代码阅读文档

## 1. 全局总结
该文件实现了一个混合缓存系统，支持 Redis 和内存两种后端。当 Redis 启用时使用 Redis，否则回退到内存热缓存（hot.HotCache）。提供统一的 Get、Set、Keys、Purge、Delete 等操作，并支持命名空间隔离。

## 2. 依赖关系
- **context**: 用于 Redis 操作的超时控制。
- **errors**: 错误检查。
- **strings**: 字符串处理。
- **sync**: 同步原语（sync.Once）。
- **time**: 时间处理。
- **github.com/go-redis/redis/v8**: Redis 客户端。
- **github.com/samber/hot**: 内存热缓存库。

## 3. 类型定义
### HybridCacheConfig[V any]
泛型结构体，配置混合缓存：
- **Namespace**: 命名空间，用于键隔离。
- **Redis**: Redis 客户端。
- **RedisCodec**: Redis 值编解码器。
- **RedisEnabled**: 函数，返回是否启用 Redis。
- **Memory**: 函数，返回内存缓存实例。

### HybridCache[V any]
泛型结构体，混合缓存实现：
- **ns**: 命名空间。
- **redis**: Redis 客户端。
- **redisCodec**: Redis 编解码器。
- **redisEnabled**: Redis 启用函数。
- **memOnce**: 内存缓存初始化同步原语。
- **memInit**: 内存缓存初始化函数。
- **mem**: 内存缓存实例。

## 4. 函数详解
### NewHybridCache
构造函数，根据配置创建 HybridCache 实例。

### FullKey
返回带命名空间的完整键。

### redisOn
判断 Redis 是否启用：Redis 客户端和编解码器非空，且 RedisEnabled 函数返回 true（或未设置）。

### memCache
返回内存缓存实例，使用 sync.Once 确保只初始化一次。

### Get
获取值：若 Redis 启用则从 Redis 获取，否则从内存获取。处理 Redis.Nil 错误（键不存在）。

### SetWithTTL
设置值并指定 TTL：若 Redis 启用则写入 Redis，否则写入内存。

### Keys
返回所有键：若 Redis 启用则扫描匹配模式的键，否则返回内存缓存的所有键。

### scanKeys
内部函数，使用 Redis SCAN 命令扫描匹配模式的键，分批获取。

### Purge
清除所有键：若 Redis 启用则扫描并删除所有匹配键，否则清除内存缓存。

### DeleteByPrefix
按前缀删除键：构建完整前缀，扫描匹配键并删除。

### DeleteMany
批量删除键：接受原始键或完整键，使用 Redis Pipeline 批量删除（使用 UNLINK 非阻塞命令）。

### Capacity
返回缓存容量：若 Redis 启用则返回 (0, 0)，否则返回内存缓存容量。

### Algorithm
返回缓存算法：若 Redis 启用则返回 ("redis", "")，否则返回内存缓存算法。

## 5. 关键逻辑分析
- **后端切换**: 根据 RedisEnabled 动态切换 Redis 和内存后端，提供无缝回退。
- **命名空间隔离**: 使用 Namespace 为键添加前缀，避免不同用例间的键冲突。
- **超时控制**: 所有 Redis 操作设置超时（获取/设置 2s，扫描 30s，删除 10s）。
- **批量删除**: 使用 Redis Pipeline 和 UNLINK 命令，提高批量删除性能。
- **内存缓存**: 使用 hot.HotCache（LRU 算法），支持 TTL 和容量限制。

## 6. 关联文件
- **codec.go**: 提供 ValueCodec 接口，用于 Redis 值编解码。
- **namespace.go**: 提供 Namespace 类型，用于键命名空间管理。