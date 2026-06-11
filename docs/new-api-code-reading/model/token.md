# token.go 代码阅读文档

## 1. 全局总结

`token.go` 实现令牌（Token）系统的完整功能，是 API 访问控制的核心组件。每个令牌对应一个 API Key，用于认证用户请求。文件包含令牌的 CRUD 操作、验证逻辑、额度管理（增加/减少/批量更新）、模型限制、IP 白名单、搜索功能，以及 Redis 缓存集成。令牌使用 HMAC 哈希存储 Key，支持软删除和异步缓存更新。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `errors` | 标准错误类型 |
| `fmt` | 格式化输出 |
| `strings` | 字符串处理（Key 遮蔽、LIKE 模式清洗、IP 解析） |
| `github.com/QuantumNous/new-api/common` | 通用工具（Redis、时间戳、批量更新等） |
| `github.com/QuantumNous/new-api/setting/operation_setting` | 操作配置（最大令牌数限制） |
| `github.com/bytedance/gopkg/util/gopool` | 高性能 goroutine 池（异步缓存操作） |
| `gorm.io/gorm` | ORM 框架 |

## 3. 类型定义

### Token 结构体
API 令牌数据模型：
- `Id`：主键
- `UserId`：所属用户 ID
- `Key`：API Key（varchar(128)，唯一索引）
- `Status`：状态（1=启用，其他=禁用/耗尽/过期）
- `Name`：名称（索引）
- `CreatedTime`/`AccessedTime`/`ExpiredTime`：时间戳（-1 表示永不过期）
- `RemainQuota`：剩余额度
- `UnlimitedQuota`：是否无限额度
- `ModelLimitsEnabled`：是否启用模型限制
- `ModelLimits`：限制的模型列表（逗号分隔）
- `AllowIps`：IP 白名单（换行分隔）
- `UsedQuota`：已用额度
- `Group`：用户分组
- `CrossGroupRetry`：是否启用跨分组重试（仅 auto 分组有效）
- `DeletedAt`：软删除时间

## 4. 函数详解

### Key 处理

#### (token *Token) Clean()
清空 Key 字段，用于缓存存储前移除敏感数据。

#### MaskTokenKey(key string) string
遮蔽 API Key：长度 ≤4 全部遮蔽，≤8 显示首尾各2字符，>8 显示首尾各4字符。

#### (token *Token) GetFullKey() / GetMaskedKey()
获取完整/遮蔽后的 Key。

### IP 白名单

#### (token *Token) GetIpLimits() []string
解析 `AllowIps` 字段为 IP 列表。支持换行和逗号分隔，自动去除空格和空值。

### 查询

#### GetAllUserTokens(userId int, startIdx, num int) ([]*Token, error)
分页获取用户的所有令牌。

#### sanitizeLikePattern(input string) (string, error)
清洗 LIKE 搜索模式，安全规则：
1. 转义 `!` 和 `_`（使用 `!` 作为 ESCAPE 字符）
2. 拒绝连续的 `%`
3. 最多允许 2 个 `%`
4. 模糊搜索时关键词长度 ≥2

#### SearchUserTokens(userId int, keyword, token string, offset, limit int) ([]*Token, int64, error)
搜索令牌，支持名称和 Key 的 LIKE 搜索。超量用户（超过 `maxTokens` 上限）只允许精确搜索。使用 `sanitizeLikePattern` 防止 SQL 注入。

#### ValidateUserToken(key string) (*Token, error)
验证令牌有效性：检查状态、过期时间、剩余额度。无效时自动更新状态（Redis 未启用时）。

#### GetTokenByIds(id, userId int) (*Token, error)
根据 ID 和用户 ID 查询令牌。

#### GetTokenById(id int) (*Token, error)
根据 ID 查询令牌，成功后异步更新 Redis 缓存。

#### GetTokenByKey(key string, fromDB bool) (*Token, error)
根据 Key 查询令牌。优先从 Redis 缓存获取，未命中则查数据库并回填缓存。

### CRUD

#### (token *Token) Insert() error
插入新令牌。

#### (token *Token) Update() error
更新令牌的可编辑字段（name、status、expired_time、remain_quota、unlimited_quota、model_limits_enabled、model_limits、allow_ips、group、cross_group_retry），成功后异步更新缓存。

#### (token *Token) SelectUpdate() error
选择性更新 `accessed_time` 和 `status`。

#### (token *Token) Delete() error
软删除令牌，成功后异步删除 Redis 缓存。

#### DeleteTokenById(id, userId int) error
根据 ID 和用户 ID 删除令牌。

#### BatchDeleteTokens(ids []int, userId int) (int, error)
批量删除令牌，使用事务确保原子性，删除后异步清理 Redis 缓存。

### 模型限制

#### (token *Token) IsModelLimitsEnabled() bool
检查是否启用模型限制。

#### (token *Token) GetModelLimits() []string
获取限制的模型列表。

#### (token *Token) GetModelLimitsMap() map[string]bool
获取模型限制的 Map 格式。

#### DisableModelLimits(tokenId int) error
禁用令牌的模型限制。

### 额度管理

#### IncreaseTokenQuota(tokenId int, key string, quota int) error
增加令牌额度。支持 Redis 异步更新和批量更新模式。

#### DecreaseTokenQuota(id int, key string, quota int) error
减少令牌额度。逻辑与 `IncreaseTokenQuota` 对称。

#### CountUserTokens(userId int) (int64, error)
统计用户的令牌总数。

### 缓存管理

#### GetTokenKeysByIds(ids []int, userId int) ([]Token, error)
批量获取令牌的 ID 和 Key。

#### InvalidateUserTokensCache(userId int) error
清除指定用户所有令牌的 Redis 缓存，配合用户禁用/删除使用。

## 5. 关键逻辑分析

**异步缓存更新**：使用 `gopool.Go` 异步执行 Redis 操作，避免阻塞主请求流程。缓存更新失败仅记录日志，不影响业务。

**Key 安全**：存储时使用 HMAC 哈希，缓存中通过 `Clean()` 移除原始 Key。前端展示使用 `MaskTokenKey` 遮蔽。

**搜索防注入**：`sanitizeLikePattern` 对 LIKE 通配符进行严格限制和转义，使用 `!` 作为 ESCAPE 字符避免 MySQL 反斜杠问题。超量用户禁止模糊搜索防止全表扫描。

**令牌验证状态机**：`ValidateUserToken` 按优先级检查：数据库错误 → 不存在 → 状态无效 → 过期 → 额度耗尽。Redis 未启用时自动更新状态到数据库。

**批量更新优化**：额度变更支持批量更新模式（`BatchUpdateEnabled`），通过 `addNewRecord` 收集变更记录，由后台 worker 批量写入数据库。

## 6. 关联文件

- `model/token_cache.go`：Redis 缓存操作函数
- `model/user.go`：用户额度和状态管理
- `model/batch_update.go`：批量更新机制
- `common/redis.go`：Redis 工具函数
- `setting/operation_setting/`：系统配置（最大令牌数等）
- `common/constants.go`：令牌状态常量
