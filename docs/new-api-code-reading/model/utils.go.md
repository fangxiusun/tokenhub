# utils.go 代码阅读文档

## 1. 全局总结

`utils.go` 提供 model 层的通用工具函数，核心功能是**批量更新机制**（Batch Update）。通过内存缓冲 + 定时刷库的方式，将高频的小额度变更（用户额度、Token 额度、已用额度、渠道额度、请求次数）聚合后批量写入数据库，显著降低数据库写入压力。

## 2. 依赖关系

| 依赖 | 用途 |
|------|------|
| `common` | 批量更新间隔配置（`BatchUpdateInterval`）、Redis 开关、系统日志 |
| `gorm.io/gorm` | `gorm.ErrRecordNotFound` 错误判断 |
| `bytedance/gopkg/util/gopool` | 异步 goroutine 池，启动批量更新定时任务 |
| `sync` | 互斥锁保护各类型的批量更新缓冲区 |
| `time` | 定时器控制刷库间隔 |

## 3. 类型定义

### 批量更新类型常量

```go
const (
    BatchUpdateTypeUserQuota        = 0 // 用户额度
    BatchUpdateTypeTokenQuota       = 1 // Token 额度
    BatchUpdateTypeUsedQuota        = 2 // 已用额度
    BatchUpdateTypeChannelUsedQuota = 3 // 渠道已用额度
    BatchUpdateTypeRequestCount     = 4 // 请求次数
    BatchUpdateTypeCount            = 5 // 类型总数
)
```

### 全局变量

| 变量 | 类型 | 说明 |
|------|------|------|
| `batchUpdateStores` | `[]map[int]int` | 每种类型的内存缓冲区，key 为 ID，value 为累加值 |
| `batchUpdateLocks` | `[]sync.Mutex` | 每种类型对应的互斥锁 |

## 4. 函数详解

| 函数 | 说明 |
|------|------|
| `init()` | 初始化所有类型的缓冲区和锁 |
| `InitBatchUpdater()` | 启动后台 goroutine，按 `BatchUpdateInterval` 秒间隔执行批量更新 |
| `addNewRecord(type_, id, value)` | 向指定类型的缓冲区追加记录（累加同 ID 的值） |
| `batchUpdate()` | 执行批量更新：快照并清空缓冲区，然后分类处理 |
| `RecordExist(err)` | 判断 GORM 查询结果是否为"记录不存在" |
| `shouldUpdateRedis(fromDB, err)` | 判断是否需要更新 Redis 缓存 |

## 5. 关键逻辑分析

### 批量更新流程

```
请求 → addNewRecord() → 内存缓冲区 → batchUpdate() → 数据库
                                      ↑
                              定时器触发（每 N 秒）
```

1. **写入阶段**：高频请求调用 `addNewRecord()` 将变更追加到内存 map 中，相同 ID 的值会累加
2. **快照阶段**：`batchUpdate()` 执行时，先加锁获取当前所有缓冲区的快照，然后立即清空原缓冲区
3. **更新阶段**：
   - Token 额度和渠道额度：逐条更新
   - 用户额度、已用额度、请求次数：合并到同一用户的多条记录后批量更新

### 合并优化
`batchUpdate()` 将 `UserQuota`、`UsedQuota`、`RequestCount` 三种类型的变更按用户 ID 合并，一次 SQL 更新用户的所有字段，减少数据库查询次数。

### 线程安全
每种类型有独立的互斥锁，不同类型之间不互相阻塞。`batchUpdate()` 中快照-清空操作在锁内完成，更新操作在锁外执行。

### 辅助函数

- `RecordExist(err)`: 标准化 GORM 记录不存在的判断
- `shouldUpdateRedis(fromDB, err)`: 仅在从 DB 读取成功且 Redis 已启用时返回 `true`，用于触发缓存异步刷新

## 6. 关联文件

| 文件 | 关联内容 |
|------|----------|
| `model/user.go` | `IncreaseUserQuota()`、`DecreaseUserQuota()`、`UpdateUserUsedQuotaAndRequestCount()` 调用 `addNewRecord()` |
| `model/token.go` | `increaseTokenQuota()`、`updateChannelUsedQuota()` 被批量更新调用 |
| `model/main.go` | `InitBatchUpdater()` 在应用启动时调用 |
| `common/constants.go` | `BatchUpdateInterval`、`BatchUpdateEnabled` 配置 |
