# token_cache.go 代码阅读文档

## 1. 全局总结

`token_cache.go` 实现令牌（Token）的 Redis 缓存操作层。所有缓存操作使用 HMAC 哈希后的 Key 作为 Redis 的 Hash Key，存储完整的 Token 对象。支持缓存写入、删除、额度原子增减、单字段更新和缓存读取。该文件是 `token.go` 中缓存相关操作的底层实现。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `fmt` | 格式化 Redis Key |
| `time` | 缓存过期时间 |
| `github.com/QuantumNous/new-api/common` | Redis 操作函数和 HMAC 生成 |
| `github.com/QuantumNous/new-api/constant` | Redis 字段名常量 |

## 3. 类型定义

无类型定义。

## 4. 函数详解

### cacheSetToken(token Token) error
将 Token 对象写入 Redis Hash。流程：
1. 使用 `common.GenerateHMAC(token.Key)` 生成哈希 Key
2. 调用 `token.Clean()` 清除原始 Key（安全考虑）
3. 使用 `RedisHSetObj` 以 Hash 结构存储，设置过期时间（`RedisKeyCacheSeconds` 秒）

### cacheDeleteToken(key string) error
从 Redis 中删除指定令牌的缓存。使用 HMAC 哈希 Key 构造 Redis Key 后调用 `RedisDelKey`。

### cacheIncrTokenQuota(key string, increment int64) error
原子增加令牌的 `remain_quota` 字段。使用 `RedisHIncrBy` 对 Hash 的指定字段做原子递增，避免并发下的竞态条件。

### cacheDecrTokenQuota(key string, decrement int64) error
原子减少令牌的 `remain_quota` 字段。内部调用 `cacheIncrTokenQuota` 并取反。

### cacheSetTokenField(key string, field string, value string) error
更新令牌缓存中的单个字段。使用 `RedisHSetField` 操作 Hash 的指定字段。

### cacheGetTokenByKey(key string) (*Token, error)
从 Redis 缓存中读取 Token 对象：
1. 使用 HMAC 哈希 Key
2. 检查 Redis 是否启用
3. 使用 `RedisHGetObj` 读取 Hash 并反序列化为 Token
4. 读取成功后恢复原始 Key（因为写入时已清空）

## 5. 关键逻辑分析

**HMAC 哈希 Key**：Redis 中存储的 Key 不是原始 API Key，而是其 HMAC 哈希值。这确保即使 Redis 被访问，原始 API Key 也不会泄露。

**Hash 结构存储**：使用 Redis Hash 而非 String 存储 Token，使得可以对单个字段（如 `remain_quota`）做原子操作，而不需要读取-修改-写入整个对象。

**原子额度操作**：`cacheIncrTokenQuota` 使用 `HINCRBY` 命令，这是 Redis 的原子操作，即使多个进程同时修改也不会产生竞态条件。

**Key 恢复**：`cacheGetTokenByKey` 在读取后将原始 Key 写回 Token 对象，因为写入缓存时已通过 `Clean()` 清空了 Key。

**过期时间管理**：缓存写入时设置 TTL，由 Redis 自动过期。过期后下次请求会从数据库重新加载。

## 6. 关联文件

- `model/token.go`：调用这些缓存函数的上层逻辑（`GetTokenByKey`、`Update`、`Delete` 等）
- `common/redis.go`：`RedisHSetObj`、`RedisHGetObj`、`RedisHIncrBy`、`RedisDelKey` 等 Redis 操作函数
- `common/utils.go`：`GenerateHMAC()` 哈希函数
- `constant/token.go`：`TokenFiledRemainQuota` 等字段常量
