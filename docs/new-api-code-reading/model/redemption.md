# redemption.go 代码阅读文档

## 1. 全局总结

`redemption.go` 实现兑换码（Redemption Code）系统的完整 CRUD 和业务逻辑。支持兑换码的创建、查询、搜索、兑换（充值到用户额度）、删除和批量清理。兑换操作使用数据库事务和行锁（`FOR UPDATE`）确保并发安全，防止同一兑换码被重复使用。还支持兑换码的过期机制。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `errors` | 标准错误类型 |
| `fmt` | 格式化日志消息 |
| `strconv` | 字符串转整数（搜索关键词解析） |
| `github.com/QuantumNous/new-api/common` | 通用工具函数（时间戳、随机延迟、PostgreSQL 判断等） |
| `github.com/QuantumNous/new-api/logger` | 日志工具（额度格式化） |
| `gorm.io/gorm` | ORM 框架，数据库操作 |

## 3. 类型定义

### Redemption 结构体
兑换码数据模型，包含：
- `Id`：主键
- `UserId`：创建者用户 ID
- `Key`：兑换码字符串（32位字符，唯一索引）
- `Status`：状态（1=启用，其他=已使用/已禁用）
- `Name`：名称（索引）
- `Quota`：兑换额度（默认100）
- `CreatedTime`：创建时间
- `RedeemedTime`：兑换时间
- `Count`：仅用于 API 请求的辅助字段（不存入数据库）
- `UsedUserId`：使用该兑换码的用户 ID
- `DeletedAt`：软删除时间
- `ExpiredTime`：过期时间（0 表示不过期）

## 4. 函数详解

### GetAllRedemptions(startIdx, num int) ([]*Redemption, int64, error)
分页获取所有兑换码，按 ID 降序排列。使用事务确保数据一致性，返回总数和分页数据。

### SearchRedemptions(keyword string, startIdx, num int) ([]*Redemption, int64, error)
搜索兑换码。如果关键词可转为整数，则同时匹配 ID 和名称前缀；否则仅匹配名称前缀。

### GetRedemptionById(id int) (*Redemption, error)
根据 ID 获取单个兑换码。

### Redeem(key string, userId int) (int, error)
核心兑换函数。使用数据库事务执行：
1. 加行锁（`FOR UPDATE`）查询兑换码
2. 验证兑换码状态（必须为启用状态）
3. 验证是否过期
4. 增加用户额度
5. 更新兑换码状态为已使用
6. 记录充值日志
7. 使用 `common.RandomSleep()` 做随机延迟防止并发攻击

### (redemption *Redemption) Insert() error
插入新兑换码记录。

### (redemption *Redemption) SelectUpdate() error
选择性更新兑换码的 `redeemed_time` 和 `status` 字段。

### (redemption *Redemption) Update() error
更新兑换码的 `name`、`status`、`quota`、`redeemed_time`、`expired_time` 字段。

### (redemption *Redemption) Delete() error
软删除兑换码。

### DeleteRedemptionById(id int) error
根据 ID 删除兑换码。

### DeleteInvalidRedemptions() (int64, error)
批量清理无效兑换码：已使用、已禁用、或已过期的兑换码。返回删除行数。

## 5. 关键逻辑分析

**并发安全**：`Redeem` 函数使用 `FOR UPDATE` 行锁 + `RandomSleep()` 随机延迟双重机制防止兑换码被并发重复使用。数据库事务确保额度增加和状态更新的原子性。

**跨数据库兼容**：`key` 是 SQL 保留字，MySQL 使用反引号引用，PostgreSQL 使用双引号引用。通过 `common.UsingPostgreSQL` 判断选择正确的引用方式。

**软删除**：`DeletedAt` 字段使用 GORM 的软删除机制，`Delete()` 方法只设置删除时间而非物理删除。

**批量清理**：`DeleteInvalidRedemptions` 一次性清理多种无效状态的兑换码，用于定期维护任务。

## 6. 关联文件

- `common/utils.go`：`GetTimestamp()`、`RandomSleep()`、`GetRandomString()` 等工具函数
- `model/user.go`：用户额度更新
- `model/log.go`：`RecordLog()` 日志记录
- `common/constants.go`：`RedemptionCodeStatusEnabled/Used/Disabled` 状态常量
