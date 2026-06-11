# usedata_rankings.go 代码阅读文档

## 1. 全局总结

`usedata_rankings.go` 实现了基于 `quota_data` 表的模型使用量排名统计功能。提供 Token 总量排名和按时间桶（Bucket）分组的使用量排行，支持自定义时间范围和桶大小，兼容 MySQL 和 PostgreSQL 的聚合表达式差异。

## 2. 依赖关系

| 依赖 | 用途 |
|------|------|
| `common` | 数据库类型判断（`UsingMySQL`），用于区分聚合表达式语法 |
| `gorm.io/gorm` | ORM 框架，构建动态查询 |
| `fmt` | 格式化 SQL 聚合表达式 |

## 3. 类型定义

### RankingQuotaTotal（Token 总量排名）

```go
type RankingQuotaTotal struct {
    ModelName   string // 模型名称
    TotalTokens int64  // Token 总用量
}
```

### RankingQuotaBucket（时间桶排名）

```go
type RankingQuotaBucket struct {
    ModelName string // 模型名称
    Bucket    int64  // 时间桶（向下取整的时间戳）
    Tokens    int64  // 该桶内的 Token 用量
}
```

## 4. 函数详解

| 函数 | 说明 |
|------|------|
| `GetRankingQuotaTotals(startTime, endTime)` | 按模型名分组，统计 Token 总用量，按总量降序排列 |
| `GetRankingQuotaBuckets(startTime, endTime, bucketSize)` | 按模型名 + 时间桶分组，统计各时间段内 Token 用量，桶大小可自定义 |
| `rankingBucketExpr(bucketSize)` | 生成时间桶的 SQL 表达式：MySQL 使用 `FLOOR()`，PostgreSQL 使用整除 |
| `applyRankingQuotaTimeRange(query, startTime, endTime)` | 向 GORM 查询追加时间范围过滤条件 |

## 5. 关键逻辑分析

### 时间桶计算
- MySQL: `FLOOR(created_at / bucketSize) * bucketSize`
- PostgreSQL/SQLite: `(created_at / bucketSize) * bucketSize`（Go 整除天然向下取整）
- 默认桶大小为 3600 秒（1小时），可通过参数自定义

### 查询过滤
- 排除空模型名（`model_name <> ''`）
- 只统计有实际使用量的记录（`sum(token_used) > 0`）
- 支持可选的时间范围过滤（`startTime > 0` 和 `endTime > 0` 时生效）

### 跨数据库兼容
通过 `common.UsingMySQL` 判断数据库类型，生成兼容的 SQL 表达式，确保三种数据库（MySQL、PostgreSQL、SQLite）均可正常运行。

## 6. 关联文件

| 文件 | 关联内容 |
|------|----------|
| `model/usedata.go` | `QuotaData` 模型定义，`quota_data` 表结构 |
| `controller/api.go` | 排名统计 API 入口 |
| `common/main.go` | `UsingMySQL`、`UsingPostgreSQL` 数据库类型标志 |
