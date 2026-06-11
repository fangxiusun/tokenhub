# redis.go 代码阅读文档

## 1. 全局总结

本文件封装了 Redis 客户端的初始化和常用操作函数，提供统一的 Redis 操作接口。支持字符串读写、Hash 结构的序列化/反序列化、原子自增操作等。使用 go-redis v8 客户端库，支持事务管道（TxPipeline）以保证操作原子性。通过环境变量 `REDIS_CONN_STRING` 控制是否启用 Redis。

## 2. 依赖关系

| 包名 | 用途 |
|------|------|
| `context` | 提供超时控制和请求上下文 |
| `errors` | 错误判断（`errors.Is`） |
| `fmt` | 格式化输出 |
| `os` | 读取环境变量 |
| `reflect` | 反射机制，用于结构体字段遍历和类型转换 |
| `strconv` | 字符串与数值类型转换 |
| `time` | 时间处理 |
| `github.com/go-redis/redis/v8` | Redis 客户端库 |
| `gorm.io/gorm` | GORM ORM 框架（用于 DeletedAt 类型） |

**项目内依赖：**
- `common.SyncFrequency` — 同步频率配置（用于 Redis key 缓存过期时间）
- `common.SysLog` / `common.FatalLog` — 日志函数
- `common.DebugEnabled` — 调试模式开关
- `common.GetEnvOrDefault` — 环境变量读取工具

## 3. 类型定义

### 全局变量

```go
var RDB *redis.Client          // Redis 客户端实例
var RedisEnabled = true        // Redis 是否启用标志
```

### RedisKeyCacheSeconds()

```go
func RedisKeyCacheSeconds() int
```

返回 `SyncFrequency` 的值，用于 Redis key 的缓存过期时间。

## 4. 函数详解

### InitRedisClient() error

```go
func InitRedisClient() (err error)
```

初始化 Redis 客户端：
1. 检查 `REDIS_CONN_STRING` 环境变量是否设置
2. 若未设置，禁用 Redis 并返回
3. 设置默认 `SYNC_FREQUENCY` 为 60 秒
4. 解析 Redis 连接字符串
5. 从 `REDIS_POOL_SIZE` 读取连接池大小（默认 10）
6. 创建 Redis 客户端
7. 执行 Ping 测试连接（5 秒超时）
8. 调试模式下输出连接信息

### ParseRedisOption() *redis.Options

```go
func ParseRedisOption() *redis.Options
```

解析 Redis 连接字符串并返回配置选项。失败时调用 `FatalLog` 终止程序。

### RedisSet(key string, value string, expiration time.Duration) error

```go
func RedisSet(key string, value string, expiration time.Duration) error
```

设置 Redis 字符串键值对，支持过期时间。

### RedisGet(key string) (string, error)

```go
func RedisGet(key string) (string, error)
```

获取 Redis 字符串值。

### RedisDel(key string) error

```go
func RedisDel(key string) error
```

删除 Redis 键。

### RedisDelKey(key string) error

```go
func RedisDelKey(key string) error
```

删除 Redis 键（与 `RedisDel` 功能相同，可能是历史遗留）。

### RedisHSetObj(key string, obj interface{}, expiration time.Duration) error

```go
func RedisHSetObj(key string, obj interface{}, expiration time.Duration) error
```

将结构体序列化为 Redis Hash 存储：
1. 通过反射遍历结构体所有字段
2. 跳过 `gorm.DeletedAt` 类型字段
3. 处理指针类型（nil 值存为空字符串）
4. 布尔类型转换为 "true"/"false" 字符串
5. 其他类型使用 `%v` 格式化
6. 使用事务管道执行 HSet + 可选 Expire

### RedisHGetObj(key string, obj interface{}) error

```go
func RedisHGetObj(key string, obj interface{}) error
```

从 Redis Hash 反序列化为结构体：
1. 使用 HGetAll 获取所有字段
2. 通过反射将字段值设置到目标结构体
3. 支持的类型：string、int/int64、bool、struct（仅 gorm.DeletedAt）
4. 指针类型自动初始化并解引用

### RedisIncr(key string, delta int64) error

```go
func RedisIncr(key string, delta int64) error
```

原子自增操作，带 TTL 保护：
1. 先获取 key 的剩余 TTL
2. 如果 key 存在且有 TTL，使用事务管道执行自增并重新设置过期时间
3. 如果 key 不存在或无 TTL，不执行操作

### RedisHIncrBy(key, field string, delta int64) error

```go
func RedisHIncrBy(key, field string, delta int64) error
```

Hash 字段原子自增，带 TTL 保护。逻辑与 `RedisIncr` 类似。

### RedisHSetField(key, field string, value interface{}) error

```go
func RedisHSetField(key, field string, value interface{}) error
```

设置 Hash 单个字段，带 TTL 保护。使用事务管道保证原子性。

## 5. 关键逻辑分析

### 5.1 连接管理

- 通过环境变量 `REDIS_CONN_STRING` 控制是否启用
- 连接池大小可配置（`REDIS_POOL_SIZE`，默认 10）
- 初始化时进行 Ping 测试，失败则 Fatal 退出

### 5.2 结构体序列化机制

`RedisHSetObj` 和 `RedisHGetObj` 使用反射实现结构体与 Hash 的双向转换：
- **序列化**：遍历字段，跳过 `gorm.DeletedAt`，布尔值特殊处理
- **反序列化**：按字段名匹配，支持指针类型自动初始化
- **事务保证**：使用 `TxPipeline` 确保 HSet 和 Expire 的原子性

### 5.3 TTL 保护模式

`RedisIncr`、`RedisHIncrBy`、`RedisHSetField` 三个函数采用统一模式：
1. 先查询 key 的 TTL
2. 如果 key 存在且有 TTL，在事务中执行操作并重新设置相同的过期时间
3. 如果 key 不存在，不执行任何操作（返回 nil）

这种设计确保了自增操作不会意外删除 key 的过期时间。

### 5.4 调试日志

所有 Redis 操作在 `DebugEnabled` 为 true 时输出详细日志，包括 key、value、过期时间等信息。

### 5.5 潜在问题

1. **RedisDel 和 RedisDelKey 重复**：两个函数功能完全相同，可能是历史遗留
2. **反射性能**：结构体序列化/反序列化使用反射，在高频调用场景下可能有性能开销
3. **错误处理**：`InitRedisClient` 中 Redis 连接失败调用 `FatalLog`，可能导致进程退出

## 6. 关联文件

- `new-api/common/env.go` — `GetEnvOrDefault` 函数定义
- `new-api/common/log.go` — `SysLog`、`FatalLog` 函数定义
- `new-api/setting/operation.go` — `SyncFrequency` 配置
- `new-api/model/` — 使用 `gorm.DeletedAt` 的数据模型
- `new-api/middleware/` — 使用 Redis 的中间件（如认证缓存）
