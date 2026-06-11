# user_oauth_binding.go 代码阅读文档

## 1. 全局总结

`user_oauth_binding.go` 管理用户与自定义 OAuth 提供商之间的绑定关系。支持创建、更新、删除绑定，提供按用户/提供商查询、查找绑定对应用户、检查占用、统计绑定数等功能。通过唯一索引保证每个用户对每个提供商只有一条绑定，每个 OAuth 账号只绑定一个用户。

## 2. 依赖关系

| 依赖 | 用途 |
|------|------|
| `gorm.io/gorm` | ORM 框架，提供事务支持和查询 |
| `time` | 记录绑定创建时间 |

## 3. 类型定义

### UserOAuthBinding（OAuth 绑定关系）

```go
type UserOAuthBinding struct {
    Id             int       // 主键
    UserId         int       // 用户ID（唯一索引：user_id + provider_id）
    ProviderId     int       // 提供商ID（唯一索引：user_id + provider_id 和 provider_id + provider_user_id）
    ProviderUserId string    // 提供商用户ID（唯一索引：provider_id + provider_user_id）
    CreatedAt      time.Time // 创建时间
}
```

**表名**：`user_oauth_bindings`

**索引设计**：
- `ux_user_provider`：`(UserId, ProviderId)` — 保证每个用户对每个提供商只有一条绑定
- `ux_provider_userid`：`(ProviderId, ProviderUserId)` — 保证每个 OAuth 账号只绑定一个用户

## 4. 函数详解

### 查询

| 函数 | 说明 |
|------|------|
| `GetUserOAuthBindingsByUserId(userId)` | 获取用户的所有 OAuth 绑定 |
| `GetUserOAuthBinding(userId, providerId)` | 获取用户对指定提供商的绑定 |
| `GetUserByOAuthBinding(providerId, providerUserId)` | 通过提供商和 OAuth 用户ID查找对应用户 |
| `GetBindingCountByProviderId(providerId)` | 获取指定提供商的绑定数量 |

### 创建

| 函数 | 说明 |
|------|------|
| `CreateUserOAuthBinding(binding)` | 创建绑定：验证必填字段，检查 OAuth 账号是否已被占用 |
| `CreateUserOAuthBindingWithTx(tx, binding)` | 事务内创建绑定（用于 OAuth 注册的原子操作） |

### 更新

| 函数 | 说明 |
|------|------|
| `UpdateUserOAuthBinding(userId, providerId, newProviderUserId)` | 更新绑定：检查新 OAuth 账号是否被占用，不存在则创建，存在则更新 |

### 删除

| 函数 | 说明 |
|------|------|
| `DeleteUserOAuthBinding(userId, providerId)` | 删除指定绑定 |
| `DeleteUserOAuthBindingsByUserId(userId)` | 删除用户的所有 OAuth 绑定 |

### 检查

| 函数 | 说明 |
|------|------|
| `IsProviderUserIdTaken(providerId, providerUserId)` | 检查 OAuth 账号是否已被其他用户占用 |

## 5. 关键逻辑分析

### 唯一性约束
通过两个复合唯一索引实现：
1. **用户维度**：`(UserId, ProviderId)` — 一个用户对一个提供商只能有一个绑定
2. **OAuth 账号维度**：`(ProviderId, ProviderUserId)` — 一个 OAuth 账号只能绑定一个用户

### 事务安全
`CreateUserOAuthBindingWithTx()` 在事务内检查和创建绑定，用于 OAuth 注册时与用户创建保持原子性。检查占用时也使用事务内的 `tx` 查询，避免竞态条件。

### 更新语义
`UpdateUserOAuthBinding()` 实现了 upsert 逻辑：
1. 检查新 OAuth 账号是否被其他用户占用
2. 检查当前用户是否已有该提供商的绑定
3. 有绑定则更新，无绑定则创建

## 6. 关联文件

| 文件 | 关联内容 |
|------|----------|
| `model/user.go` | 用户模型，OAuth 注册流程中调用绑定函数 |
| `oauth/` | OAuth 提供商实现，登录回调时调用绑定查询/创建 |
| `controller/auth.go` | OAuth 登录/注册 API 入口 |
