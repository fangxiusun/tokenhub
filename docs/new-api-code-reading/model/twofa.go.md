# twofa.go 代码阅读文档

## 1. 全局总结

`twofa.go` 实现了双因素认证（2FA/TOTP）的完整功能，包括 2FA 设置的创建、启用、禁用、TOTP 验证、备用码管理以及账户锁定机制。支持失败次数限制和自动锁定，为用户账户提供额外的安全保障。

## 2. 依赖关系

| 依赖 | 用途 |
|------|------|
| `common` | TOTP 验证（`ValidateTOTPCode`）、备用码哈希/验证（`HashBackupCode`、`ValidateBackupCode`、`NormalizeBackupCode`、`ValidatePasswordAndHash`）、失败次数/锁定时长配置（`MaxFailAttempts`、`LockoutDuration`）、系统日志 |
| `gorm.io/gorm` | ORM 框架，提供事务、软/硬删除、错误处理 |
| `time` | 时间计算（锁定时间、使用时间） |
| `fmt` | 错误信息格式化 |

## 3. 类型定义

### TwoFA（2FA 设置）

```go
type TwoFA struct {
    Id             int            // 主键
    UserId         int            // 用户ID（唯一索引）
    Secret         string         // TOTP 密钥（不返回前端）
    IsEnabled      bool           // 是否启用
    FailedAttempts int            // 失败尝试次数（默认0）
    LockedUntil    *time.Time     // 锁定截止时间
    LastUsedAt     *time.Time     // 最后使用时间
    CreatedAt      time.Time      // 创建时间
    UpdatedAt      time.Time      // 更新时间
    DeletedAt      gorm.DeletedAt // 软删除标记
}
```

### TwoFABackupCode（备用码）

```go
type TwoFABackupCode struct {
    Id        int            // 主键
    UserId    int            // 用户ID（索引）
    CodeHash  string         // 备用码哈希（不返回前端）
    IsUsed    bool           // 是否已使用
    UsedAt    *time.Time     // 使用时间
    CreatedAt time.Time      // 创建时间
    DeletedAt gorm.DeletedAt // 软删除标记
}
```

## 4. 函数详解

### 2FA 查询与状态

| 函数 | 说明 |
|------|------|
| `GetTwoFAByUserId(userId)` | 根据用户ID获取 2FA 设置，未设置返回 `nil` |
| `IsTwoFAEnabled(userId)` | 检查用户是否启用了 2FA |
| `GetTwoFAStats()` | 获取 2FA 统计信息（总用户数、启用数、启用率），供管理员使用 |

### 2FA 生命周期管理

| 函数 | 说明 |
|------|------|
| `Create()` | 创建 2FA 设置，检查用户是否已存在设置及用户是否存在 |
| `Update()` | 更新 2FA 设置 |
| `Delete()` | 硬删除 2FA 及相关备用码（事务操作） |
| `Enable()` | 启用 2FA，重置失败次数 |
| `DisableTwoFA(userId)` | 禁用用户 2FA（删除设置和备用码） |

### 验证逻辑

| 函数 | 说明 |
|------|------|
| `ValidateTOTPAndUpdateUsage(code)` | 验证 TOTP 码，成功重置失败次数并更新使用时间，失败增加计数 |
| `ValidateBackupCodeAndUpdateUsage(code)` | 验证备用码，逻辑同 TOTP |
| `ValidateBackupCode(userId, code)` | 验证并标记备用码为已使用 |

### 备用码管理

| 函数 | 说明 |
|------|------|
| `CreateBackupCodes(userId, codes)` | 事务中删除旧备用码并创建新备用码（哈希存储） |
| `GetUnusedBackupCodeCount(userId)` | 获取未使用的备用码数量 |

### 账户锁定

| 函数 | 说明 |
|------|------|
| `IsLocked()` | 检查账户是否被锁定（当前时间 < LockedUntil） |
| `IncrementFailedAttempts()` | 增加失败次数，达到阈值自动锁定 |
| `ResetFailedAttempts()` | 重置失败次数并解除锁定 |

## 5. 关键逻辑分析

### 安全设计
- TOTP 密钥（`Secret`）和备用码哈希（`CodeHash`）均标记为 `json:"-"`，不会返回给前端
- 备用码使用 `common.HashBackupCode()` 进行哈希存储，验证时使用 `ValidatePasswordAndHash()` 比对

### 锁定机制
- 失败次数达到 `common.MaxFailAttempts` 时自动锁定
- 锁定时长为 `common.LockoutDuration` 秒
- 验证成功自动重置失败次数和锁定状态

### 备用码生命周期
- 创建新备用码时先删除旧码（软删除）
- 验证成功后标记 `IsUsed = true` 并记录 `UsedAt`
- 删除 2FA 时硬删除所有相关备用码（`Unscoped()`）

### 统计功能
`GetTwoFAStats()` 返回总用户数、启用用户数和启用率，供管理后台展示。

## 6. 关联文件

| 文件 | 关联内容 |
|------|----------|
| `common/auth.go` | `ValidateTOTPCode()`、`HashBackupCode()`、`ValidateBackupCode()`、`NormalizeBackupCode()` |
| `common/constants.go` | `MaxFailAttempts`、`LockoutDuration` 配置 |
| `model/user.go` | 用户存在性验证（`First` 查询） |
| `controller/auth.go` | 2FA 设置/验证/禁用的 API 入口 |
