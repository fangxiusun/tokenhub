# config.go 代码阅读文档

## 1. 全局总结

该文件定义性能指标（Performance Metrics）设置，控制指标采集的开关、刷新间隔、时间桶大小和数据保留天数。

## 2. 依赖关系

- `github.com/QuantumNous/new-api/setting/config` — 全局配置管理器

## 3. 类型定义

| 结构体 | 字段 | 类型 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `PerfMetricsSetting` | `Enabled` | `bool` | `true` | 是否启用性能指标 |
| | `FlushInterval` | `int` | `5` | 刷新间隔（分钟） |
| | `BucketTime` | `string` | `"hour"` | 时间桶大小 |
| | `RetentionDays` | `int` | `0` | 数据保留天数（0 表示永久） |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `GetSetting` | `func GetSetting() PerfMetricsSetting` | 获取性能指标设置 |
| `GetBucketSeconds` | `func GetBucketSeconds() int64` | 获取时间桶大小（秒） |
| `GetFlushIntervalMinutes` | `func GetFlushIntervalMinutes() int` | 获取刷新间隔（分钟） |

## 5. 关键逻辑分析

- 时间桶支持：`"minute"` (60s)、`"5min"` (300s)、`"hour"` (3600s)
- 刷新间隔最小值为 1 分钟
- `RetentionDays` 为 0 表示数据永久保留

## 6. 关联文件

- `model/perf_metric.go` — 性能指标模型
- `middleware/performance.go` — 性能监控中间件
