# metrics.go 代码阅读文档

## 1. 全局总结
该文件实现了性能指标的核心逻辑，包括记录样本、查询指标、汇总统计等。使用热桶（内存）和 Redis 双写机制，支持分组查询和时间序列数据。

## 2. 依赖关系
- **context**: Redis 操作超时。
- **fmt**: 格式化。
- **math**: 数学运算。
- **sort**: 排序。
- **sync**: 并发控制。
- **time**: 时间处理。
- **github.com/QuantumNous/new-api/common**: Redis 客户端和工具。
- **github.com/QuantumNous/new-api/model**: 数据库操作。
- **github.com/QuantumNous/new-api/relay/common**: RelayInfo 类型。
- **github.com/QuantumNous/new-api/setting/perf_metrics_setting**: 配置。

## 3. 类型定义
无新的类型定义，但使用了 types.go 中的 Sample、QueryParams、QueryResult 等类型。

## 4. 函数详解
### Init
初始化性能指标系统，启动后台刷新循环。

### RecordRelaySample
从 RelayInfo 计算指标并记录：延迟、TTFT（首 token 时间）、生成时间等。

### Record
记录样本：验证设置，构建桶键，存储到热桶和 Redis。

### Query
查询指标：合并数据库和热桶数据，构建时间序列结果。

### QuerySummaryAll
查询所有模型的汇总统计：计算平均延迟、成功率、平均 TPS 等。

### allowedGroupSet
构建允许的分组集合。

### bucketStart
计算桶的开始时间戳（按桶秒数对齐）。

### mergeCounters
合并计数器到映射中。

### buildQueryResult
构建查询结果：按分组组织时间序列数据。

### bucketPoint
构建桶点数据。

### avg
计算平均值。

### successRate
计算成功率。

### avgTps
计算平均 TPS（每秒 token 数）。

### recordRedis
记录样本到 Redis：使用哈希结构存储计数器，设置 1 小时过期。

### mergeRedisActiveBuckets
合并 Redis 中的活动桶数据到查询结果。

### redisBucketKey
构建 Redis 桶键。

## 5. 关键逻辑分析
- **双写机制**: 同时写入内存热桶和 Redis，确保数据一致性和查询性能。
- **时间序列**: 按时间桶组织数据，支持分组和时间范围查询。
- **聚合计算**: 在查询时合并数据库和内存数据，实时计算汇总统计。
- **并发安全**: 使用 atomic.Int64 确保计数器原子操作。

## 6. 关联文件
- **flush.go**: 使用 hotBuckets 进行刷新。
- **types.go**: 定义类型。
- **model**: 数据库操作。
- **relay/common**: RelayInfo 类型。
- **perf_metrics_setting**: 配置。