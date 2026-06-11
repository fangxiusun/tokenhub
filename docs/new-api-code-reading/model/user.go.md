# user.go 代码阅读文档

## 1. 全局总结

`user.go` 是用户模块的核心文件，定义了 `User` 结构体及完整的用户管理功能。涵盖用户 CRUD、多平台 OAuth 绑定（GitHub/Discord/OIDC/WeChat/Telegram/LinuxDO）、额度管理（含批量更新）、邀请机制、登录验证、Redis 缓存集成、边栏配置初始化等。是整个系统中最大、最复杂的 model 文件。

## 2. 依赖关系

| 依赖 | 用途 |
|------|------|
| `common` | 密码哈希、JWT token、UUID、Redis 操作、系统配置常量、用户状态/角色定义 |
| `dto` | `UserSetting` 用户设置 DTO |
| `logger` | 额度日志格式化 |
| `setting/operation_setting` | 支付合规确认检查 |
| `gorm.io/gorm` | ORM 框架 |
| `bytedance/gopkg/util/gopool` | 异步 goroutine 池，用于缓存更新 |

## 3. 类型定义

### User（用户）

```go
type User struct {
    Id               int            // 用户ID
    Username         string         // 用户名（唯一索引，最大20字符）
    Password         string         // 密码哈希
    OriginalPassword string         // 原始密码（仅用于修改验证，不存库）
    DisplayName      string         // 显示名
    Role             int            // 角色（admin/common/root）
    Status           int            // 状态（enabled/disabled）
    Email            string         // 邮箱（索引）
    GitHubId         string         // GitHub ID（索引）
    DiscordId        string         // Discord ID（索引）
    OidcId           string         // OIDC ID（索引）
    WeChatId         string         // 微信 ID（索引）
    TelegramId       string         // Telegram ID（索引）
    LinuxDOId        string         // Linux DO ID（索引）
    AccessToken      *string        // 系统管理 Token（唯一索引）
    Quota            int            // 当前额度
    UsedQuota        int            // 已用额度
    RequestCount     int            // 请求次数
    Group            string         // 用户组（默认 "default"）
    AffCode          string         // 邀请码（唯一索引）
    AffCount         int            // 邀请人数
    AffQuota         int            // 邀请剩余额度
    AffHistoryQuota  int            // 邀请历史额度
    InviterId        int            // 邀请人ID（索引）
    Setting          string         // 用户设置（JSON，text 类型）
    Remark           string         // 备注
    StripeCustomer   string         // Stripe 客户ID（索引）
    CreatedAt        int64          // 创建时间（自动）
    LastLoginAt      int64          // 最后登录时间
    DeletedAt        gorm.DeletedAt // 软删除标记
}
```

## 4. 函数详解

### 用户创建

| 函数 | 说明 |
|------|------|
| `Insert(inviterId)` | 创建新用户：密码哈希、初始额度、邀请码、边栏配置初始化、邀请奖励 |
| `InsertWithTx(tx, inviterId)` | 事务内创建用户（OAuth 注册用） |
| `FinalizeOAuthUserCreation(inviterId)` | OAuth 用户创建后的事务外处理（边栏配置、日志、邀请奖励） |

### 用户查询

| 函数 | 说明 |
|------|------|
| `GetUserById(id, selectAll)` | 按 ID 查询，`selectAll=false` 时排除密码 |
| `GetAllUsers(pageInfo)` | 分页查询所有用户（含软删除，排除密码） |
| `SearchUsers(keyword, group, role, status, startIdx, num)` | 多条件搜索用户 |
| `GetUserById(id)` | 按 ID 获取用户 |
| `GetUserIdByAffCode(affCode)` | 按邀请码获取用户ID |
| `GetUsernameById(id, fromDB)` | 带 Redis 缓存的用户名查询 |
| `GetRootUser()` | 获取超级管理员 |
| `RootUserExists()` | 检查超级管理员是否存在 |

### 用户更新

| 函数 | 说明 |
|------|------|
| `Update(updatePassword)` | 更新用户信息并刷新缓存 |
| `Edit(updatePassword)` | 精确更新指定字段（username/display_name/group/remark） |
| `ClearBinding(bindingType)` | 清除指定类型的 OAuth 绑定 |

### 用户删除

| 函数 | 说明 |
|------|------|
| `Delete()` | 软删除用户并清除缓存 |
| `HardDelete()` | 硬删除用户 |
| `DeleteUserById(id)` | 按 ID 软删除 |
| `HardDeleteUserById(id)` | 按 ID 硬删除 |

### OAuth 绑定查询

| 函数 | 说明 |
|------|------|
| `FillUserByGitHubId()` | 按 GitHub ID 填充用户 |
| `FillUserByDiscordId()` | 按 Discord ID 填充用户 |
| `FillUserByOidcId()` | 按 OIDC ID 填充用户 |
| `FillUserByWeChatId()` | 按微信 ID 填充用户 |
| `FillUserByTelegramId()` | 按 Telegram ID 填充用户 |
| `FillUserByLinuxDOId()` | 按 Linux DO ID 填充用户 |
| `Is*AlreadyTaken()` | 检查各平台 ID 是否已被占用 |

### 额度管理

| 函数 | 说明 |
|------|------|
| `GetUserQuota(id, fromDB)` | 带 Redis 缓存的额度查询 |
| `IncreaseUserQuota(id, quota, db)` | 增加用户额度（支持批量更新模式） |
| `DecreaseUserQuota(id, quota, db)` | 减少用户额度（支持批量更新模式） |
| `DeltaUpdateUserQuota(id, delta)` | 根据正负值增减额度 |
| `TransferAffQuotaToQuota(quota)` | 将邀请额度转为可用额度（事务+行锁） |

### 使用量更新

| 函数 | 说明 |
|------|------|
| `UpdateUserUsedQuotaAndRequestCount(id, quota)` | 更新已用额度和请求次数（支持批量更新） |
| `UpdateUserLastLoginAt(id)` | 更新最后登录时间 |

### 认证相关

| 函数 | 说明 |
|------|------|
| `ValidateAndFill()` | 验证用户名/密码并填充用户信息 |
| `ValidateAccessToken(token)` | 验证系统管理 Token |
| `ResetUserPasswordByEmail(email, password)` | 通过邮箱重置密码 |

### 用户设置

| 函数 | 说明 |
|------|------|
| `GetSetting()` | 获取用户设置（JSON 反序列化） |
| `SetSetting(setting)` | 设置用户设置（JSON 序列化） |
| `GetUserSetting(id, fromDB)` | 带 Redis 缓存的用户设置查询 |
| `GetUserGroup(id, fromDB)` | 带 Redis 缓存的用户组查询 |

### 辅助函数

| 函数 | 说明 |
|------|------|
| `ToBaseUser()` | 转换为 `UserBase` 缓存结构 |
| `generateDefaultSidebarConfigForRole(role)` | 根据角色生成默认边栏配置 |
| `IsAdmin(userId)` | 检查用户是否为管理员 |
| `IsEmailAlreadyTaken(email)` | 检查邮箱是否已被占用 |

## 5. 关键逻辑分析

### Redis 缓存策略
- `GetUserQuota`、`GetUserGroup`、`GetUsernameById`、`GetUserSetting` 均采用"Redis 优先 + DB 回退"模式
- DB 查询成功后异步（`gopool.Go`）更新 Redis 缓存
- `fromDB` 参数控制是否跳过 Redis 直接查 DB

### 批量更新机制
- `IncreaseUserQuota`、`DecreaseUserQuota`、`UpdateUserUsedQuotaAndRequestCount` 在 `BatchUpdateEnabled` 时将变更暂存内存，由 `utils.go` 的 `batchUpdate()` 定时批量刷库

### OAuth 注册流程
- `Insert()` 用于普通注册，`InsertWithTx()` + `FinalizeOAuthUserCreation()` 用于 OAuth 注册
- 分离事务内和事务外操作，确保用户创建和 OAuth 绑定的原子性

### 边栏配置初始化
- `generateDefaultSidebarConfigForRole()` 根据角色（root/admin/common）生成不同权限的边栏模块配置
- 超级管理员可访问所有功能，管理员不能访问系统设置，普通用户无管理区域

### 密码安全
- 密码使用 `common.Password2Hash()` 哈希存储
- 查询时可通过 `Omit("password")` 排除密码字段
- `OriginalPassword` 字段标记为 `gorm:"-:all"`，不存库

## 6. 关联文件

| 文件 | 关联内容 |
|------|----------|
| `model/user_cache.go` | `UserBase` 缓存结构、Redis 缓存操作 |
| `model/utils.go` | 批量更新机制、`shouldUpdateRedis()` |
| `model/user_oauth_binding.go` | OAuth 绑定关系管理 |
| `dto/setting.go` | `UserSetting` DTO 定义 |
| `common/auth.go` | 密码哈希与验证 |
| `common/redis.go` | Redis 操作封装 |
| `controller/user.go` | 用户管理 API 入口 |
| `middleware/auth.go` | JWT 验证中间件 |
