# flush.go 代码阅读文档

## 1. 全局总结
该文件实现了性能指标的刷新和清理逻辑。包含一个后台循环，定期将内存中的指标桶刷新到数据库，并清理过期数据。同时提供 Redis 计数器解析功能。

## 2. 依赖关系
- **fmt**: 格式化错误信息。
- **strconv**: 字符串转整数。
- **time**: 时间处理。
- **github.com/QuantumNous/new-api/common**: 系统错误日志。
- **github.com/QuantumNous/new-api/model**: 数据库操作。
- **github.com/QuantumNous/new-api/setting/perf_metrics_setting**: 性能指标配置。

## 3. 类型定义
无新的类型定义。

## 4. 函数详解
### flushLoop
后台循环，按配置间隔（分钟）执行刷新和清理。检查设置是否启用。

### flushCompletedBuckets
刷新已完成的桶：遍历热桶，跳过当前桶，排空旧桶并写入数据库。写入失败时恢复计数器。

### deleteOldEmptyBucket
删除超过 24 小时的空桶。

### cleanupExpiredMetrics
清理过期指标：根据保留天数删除数据库中的旧记录。

### redisCounters
从 Redis 哈希值构建 counters 结构体。

### parseRedisInt
解析 Redis 整数值，空值返回 0。

## 5. 关键逻辑分析
- **刷新策略**: 定期刷新已完成的桶（非当前时间桶），避免数据丢失。
- **错误恢复**: 写入数据库失败时，将计数器恢复到桶中，等待下次刷新。
- **清理机制**: 删除 24 小时前的空桶，清理过期数据库记录。
- **Redis 集成**: 支持从 Redis 读取计数器数据。

## 6. 关联文件
- **metrics.go**: 使用 hotBuckets 全局变量和 bucketKey、atomicBucket 类型。
- **types.go**: 定义 counters、atomicBucket 等类型。
- **model**: 提供 UpsertPerfMetric 和 DeletePerfMetricsBefore 数据库操作。
- **perf_metrics_setting**: 提供配置（刷新间隔、保留天数等）。