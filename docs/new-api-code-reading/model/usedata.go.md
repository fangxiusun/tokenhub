# usedata.go 代码阅读文档

## 1. 全局总结

`usedata.go` 负责用户使用数据（QuotaData）的采集、缓存与持久化。采用内存缓存 + 定时写入数据库的策略，按小时粒度聚合用户模型使用数据（Token 数、配额、请求数），用于数据看板的图表展示。

## 2. 依赖关系

| 依赖 | 用途 |
|------|------|
| `common` | 数据导出开关（`DataExportEnabled`）、导出间隔（`DataExportInterval`）、系统日志 |
| `gorm.io/gorm` | ORM 框架，数据库查询与更新 |
| `sync` | 互斥锁保护缓存并发访问 |
| `time` | 定时器控制数据刷新间隔 |
| `fmt` | 缓存 key 格式化 |

## 3. 类型定义

### QuotaData（配额数据）

```go
type QuotaData struct {
    Id        int    // 主键
    UserID    int    // 用户ID（索引）
    Username  string // 用户名（复合索引）
    ModelName string // 模型名（复合索引）
    CreatedAt int64  // 创建时间（小时粒度，复合索引）
    TokenUsed int    // 使用的 Token 数
    Count     int    // 请求次数
    Quota     int    // 消耗的配额
}
```

## 4. 函数详解

### 数据采集

| 函数 | 说明 |
|------|------|
| `LogQuotaData(userId, username, modelName, quota, createdAt, tokenUsed)` | 记录使用数据，时间精确到小时（向下取整到 3600 秒边界），线程安全 |
| `logQuotaDataCache(...)` | 内部函数，将数据追加到内存缓存中，按 `userId-username-modelName-timestamp` 作为 key |

### 数据持久化

| 函数 | 说明 |
|------|------|
| `SaveQuotaDataCache()` | 将内存缓存刷入数据库：已有记录则累加更新，否则插入新记录 |
| `increaseQuotaData(...)` | 使用 GORM 表达式原子累加 `count`、`quota`、`token_used` |
| `UpdateQuotaData()` | 后台 goroutine，按 `DataExportInterval` 分钟间隔调用 `SaveQuotaDataCache` |

### 数据查询

| 函数 | 说明 |
|------|------|
| `GetQuotaDataByUsername(username, startTime, endTime)` | 按用户名查询时间段内的使用数据 |
| `GetQuotaDataByUserId(userId, startTime, endTime)` | 按用户ID查询时间段内的使用数据 |
| `GetQuotaDataGroupByUser(startTime, endTime)` | 按用户名和时间分组聚合 |
| `GetAllQuotaDates(startTime, endTime, username)` | 全局查询：指定用户名则按用户查，否则按模型名+时间分组聚合 |

## 5. 关键逻辑分析

### 缓存策略
- 使用 `map[string]*QuotaData` 作为内存缓存，key 由 `userId-username-modelName-timestamp` 组成
- 时间戳向下取整到小时（`createdAt - createdAt % 3600`），确保同小时数据聚合

### 定时刷新
- `UpdateQuotaData()` 在独立 goroutine 中运行，无限循环
- 每隔 `DataExportInterval` 分钟调用一次 `SaveQuotaDataCache()`
- 仅在 `DataExportEnabled = true` 时执行

### 并发安全
- 所有缓存读写操作通过 `CacheQuotaDataLock` 互斥锁保护
- `SaveQuotaDataCache()` 在持锁状态下遍历缓存并写库，写完后清空缓存

### 数据合并
- 写入数据库时先查询是否存在同 key 记录
- 存在则使用 `gorm.Expr` 原子累加，不存在则创建新记录

## 6. 关联文件

| 文件 | 关联内容 |
|------|----------|
| `common/constants.go` | `DataExportEnabled`、`DataExportInterval` 配置 |
| `controller/api.go` | 数据看板 API 入口，调用查询函数 |
| `middleware/relay.go` | 请求处理完成后调用 `LogQuotaData()` 记录使用数据 |
| `model/main.go` | 数据库初始化，`quota_data` 表结构 |
