# perf_metric.go 代码阅读文档

## 1. 全局总结
该文件定义了性能指标（PerfMetric）的数据模型和相关的数据库操作函数。性能指标用于记录 AI 模型中继（relay）操作的性能数据，包括请求计数、成功计数、延迟、首字节时间（TTFT）、输出 tokens 等。该文件提供了性能指标的插入/更新、查询、汇总和清理功能。

## 2. 依赖关系
- **time**: 用于计算时间范围。
- **gorm.io/gorm**: 用于数据库 ORM 操作。
- **gorm.io/gorm/clause**: 用于实现 upsert（插入或更新）操作。

## 3. 类型定义
### PerfMetric 结构体
表示性能指标的数据模型，包含以下字段：
- `Id`: 主键 ID
- `ModelName`: 模型名称（联合唯一索引的一部分）
- `Group`: 用户组（联合唯一索引的一部分）
- `BucketTs`: 时间桶时间戳（联合唯一索引的一部分，建立单独索引）
- `RequestCount`: 请求计数
- `SuccessCount`: 成功计数
- `TotalLatencyMs`: 总延迟（毫秒）
- `TtftSumMs`: 首字节时间总和（毫秒）
- `TtftCount`: 首字节时间计数
- `OutputTokens`: 输出 tokens 数量
- `GenerationMs`: 生成时间（毫秒）

### PerfMetricSummary 结构体
表示性能指标汇总的数据模型：
- `ModelName`: 模型名称
- `RequestCount`: 请求计数
- `SuccessCount`: 成功计数
- `TotalLatencyMs`: 总延迟（毫秒）
- `OutputTokens`: 输出 tokens 数量
- `GenerationMs`: 生成时间（毫秒）

## 4. 函数详解
### 表名方法
- **(PerfMetric) TableName() string**
  - 返回自定义表名 "perf_metrics"。

### 数据操作函数
- **UpsertPerfMetric(metric *PerfMetric) error**
  - 插入或更新性能指标：使用 upsert 操作，如果记录存在则累加各字段值。
  - 如果 metric 为 nil 或 RequestCount 为 0，则跳过操作。

- **GetPerfMetrics(modelName string, group string, startTs int64, endTs int64) ([]PerfMetric, error)**
  - 查询指定模型、组和时间范围内的性能指标。
  - 按时间戳升序排列。

- **GetPerfMetricsSummaryAll(startTs int64, endTs int64, groups []string) ([]PerfMetricSummary, error)**
  - 查询所有模型的性能指标汇总，按模型分组。
  - 支持按用户组过滤。
  - 只返回请求数大于 0 的模型。

- **DeletePerfMetricsBefore(cutoffTs int64) error**
  - 删除指定时间戳之前的性能指标数据。
  - 用于数据清理和保留策略。

- **PerfMetricStartTime(hours int) int64**
  - 计算指定小时数之前的时间戳。
  - 如果 hours <= 0，默认为 24 小时。

## 5. 关键逻辑分析
- **Upset 操作**: 使用 GORM 的 `clause.OnConflict` 实现 upsert，当记录存在时累加各字段值，而不是覆盖。这确保了性能指标的正确聚合。
- **联合唯一索引**: 通过 `model_name`、`group`、`bucket_ts` 三个字段的联合唯一索引，确保同一模型、组、时间桶只有一条记录。
- **时间桶机制**: 性能指标按时间桶（BucketTs）聚合，便于时间范围查询和数据清理。
- **数据清理**: `DeletePerfMetricsBefore` 支持按时间清理旧数据，防止表无限增长。
- **汇总查询**: `GetPerfMetricsSummaryAll` 使用 SQL 聚合函数（SUM）和 GROUP BY 实现高效汇总。

## 6. 关联文件
- **relay/perf.go**: 可能包含性能指标的收集和上报逻辑。
- **controller/perf.go**: 可能包含处理性能指标查询的 HTTP 请求控制器。
- **router/perf.go**: 可能包含性能指标相关的路由定义。
- **model/main.go**: 提供数据库连接和跨数据库兼容性变量（如 commonGroupCol）。