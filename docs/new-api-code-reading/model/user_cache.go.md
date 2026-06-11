# user_cache.go 代码阅读文档

## 1. 全局总结

`user_cache.go` 实现了用户数据的 Redis 缓存层。定义了轻量级的 `UserBase` 结构体作为缓存数据模型，提供完整的缓存读写、失效、原子操作（额度增减）以及单字段更新功能。采用 Redis Hash 存储，支持缓存回源和异步刷新。

## 2. 依赖关系

| 依赖 | 用途 |
|------|------|
| `common` | Redis 操作封装（`RedisHSetObj`、`RedisHGetObj`、`RedisHIncrBy`、`RedisHSetField`、`RedisDelKey`）、系统日志、Redis 开关 |
| `constant` | Context Key 常量（`ContextKeyUserGroup` 等） |
| `dto` | `UserSetting` 用户设置 DTO |
| `gin` | HTTP 上下文（`gin.Context`） |
| `bytedance/gopkg/util/gopool` | 异步 goroutine 池 |

## 3. 类型定义

### UserBase（用户缓存数据）

```go
type UserBase struct {
    Id       int    // 用户ID
    Group    string // 用户组
    Email    string // 邮箱
    Quota    int    // 额度
    Status   int    // 状态
    Username string // 用户名
    Setting  string // 用户设置（JSON 字符串）
}
```

## 4. 函数详解

### 缓存读取

| 函数 | 说明 |
|------|------|
| `GetUserCache(userId)` | 获取完整用户缓存：Redis 优先，失败回退 DB，成功后异步刷新 Redis |
| `cacheGetUserBase(userId)` | 从 Redis Hash 获取 `UserBase` |
| `getUserGroupCache(userId)` | 获取缓存中的用户组 |
| `getUserQuotaCache(userId)` | 获取缓存中的额度 |
| `getUserStatusCache(userId)` | 获取缓存中的状态 |
| `getUserNameCache(userId)` | 获取缓存中的用户名 |
| `getUserSettingCache(userId)` | 获取缓存中的用户设置 |

### 缓存写入

| 函数 | 说明 |
|------|------|
| `updateUserCache(user)` | 将完整 `UserBase` 写入 Redis Hash，设置过期时间 |
| `updateUserStatusCache(userId, status)` | 更新缓存中的状态字段 |
| `updateUserQuotaCache(userId, quota)` | 更新缓存中的额度字段 |
| `updateUserGroupCache(userId, group)` | 更新缓存中的用户组字段 |
| `UpdateUserGroupCache(userId, group)` | 导出版本的用户组缓存更新 |
| `updateUserNameCache(userId, username)` | 更新缓存中的用户名字段 |
| `updateUserSettingCache(userId, setting)` | 更新缓存中的设置字段 |

### 缓存失效

| 函数 | 说明 |
|------|------|
| `invalidateUserCache(userId)` | 删除用户缓存（内部版本） |
| `InvalidateUserCache(userId)` | 删除用户缓存（导出版本，供上层调用） |

### 原子操作

| 函数 | 说明 |
|------|------|
| `cacheIncrUserQuota(userId, delta)` | Redis Hash 原子增加额度 |
| `cacheDecrUserQuota(userId, delta)` | Redis Hash 原子减少额度 |

### 辅助方法

| 函数 | 说明 |
|------|------|
| `WriteContext(c)` | 将 `UserBase` 数据写入 Gin Context |
| `GetSetting()` | 反序列化 `Setting` JSON 字符串为 `UserSetting` |
| `GetUserLanguage(userId)` | 获取用户的语言偏好设置 |

## 5. 关键逻辑分析

### 缓存架构
- 使用 Redis Hash 存储，key 格式为 `user:{userId}`
- Hash 字段对应 `UserBase` 的各属性，支持单字段读写和原子操作
- 缓存过期时间由 `common.RedisKeyCacheSeconds()` 控制

### 缓存回源策略
- `GetUserCache()` 先尝试 Redis，失败则回退 DB
- DB 查询成功后通过 `gopool.Go()` 异步更新 Redis，避免阻塞主流程
- `shouldUpdateRedis(fromDB, err)` 判断是否需要更新缓存：仅在 `fromDB=true` 且无错误时更新

### Context 注入
`WriteContext()` 将用户信息写入 Gin Context，供后续中间件和 handler 使用：
- `ContextKeyUserGroup` → 用户组
- `ContextKeyUserQuota` → 额度
- `ContextKeyUserStatus` → 状态
- `ContextKeyUserEmail` → 邮箱
- `ContextKeyUserName` → 用户名
- `ContextKeyUserSetting` → 用户设置

### 原子额度操作
使用 `RedisHIncrBy` 实现原子增减，避免并发场景下的额度不一致问题。

## 6. 关联文件

| 文件 | 关联内容 |
|------|----------|
| `model/user.go` | `User.ToBaseUser()` 转换方法，各 `Get*()` 函数的缓存回退逻辑 |
| `model/utils.go` | `shouldUpdateRedis()` 缓存更新判断 |
| `common/redis.go` | Redis 操作封装函数 |
| `constant/context.go` | Context Key 常量定义 |
| `dto/setting.go` | `UserSetting` DTO 定义 |
| `middleware/auth.go` | 从 Context 读取用户信息 |
