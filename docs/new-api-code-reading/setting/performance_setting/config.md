# config.go 代码阅读文档

## 1. 全局总结

该文件定义性能设置，包括磁盘缓存配置和性能监控阈值，并提供与 `common` 包的同步机制。

## 2. 依赖关系

- `github.com/QuantumNous/new-api/common` — 磁盘缓存和性能监控配置
- `github.com/QuantumNous/new-api/setting/config` — 全局配置管理器

## 3. 类型定义

| 结构体 | 字段 | 类型 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `PerformanceSetting` | `DiskCacheEnabled` | `bool` | `false` | 是否启用磁盘缓存 |
| | `DiskCacheThresholdMB` | `int` | `10` | 触发磁盘缓存的阈值（MB） |
| | `DiskCacheMaxSizeMB` | `int` | `1024` | 磁盘缓存最大大小（MB） |
| | `DiskCachePath` | `string` | `""` | 磁盘缓存目录 |
| | `MonitorEnabled` | `bool` | `true` | 是否启用性能监控 |
| | `MonitorCPUThreshold` | `int` | `90` | CPU 使用率阈值（%） |
| | `MonitorMemoryThreshold` | `int` | `90` | 内存使用率阈值（%） |
| | `MonitorDiskThreshold` | `int` | `95` | 磁盘使用率阈值（%） |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `GetPerformanceSetting` | `func GetPerformanceSetting() *PerformanceSetting` | 获取性能设置 |
| `UpdateAndSync` | `func UpdateAndSync()` | 更新配置并同步到 common 包 |
| `GetCacheStats` | `func GetCacheStats() common.DiskCacheStats` | 获取缓存统计信息 |
| `ResetStats` | `func ResetStats()` | 重置统计信息 |

## 5. 关键逻辑分析

- `syncToCommon` 在 `init()` 和配置更新后调用，将设置同步到 `common` 包
- 磁盘缓存路径为空时使用系统临时目录
- 监控阈值默认较高（90-95%），避免频繁告警

## 6. 关联文件

- `common/disk_cache.go` — 磁盘缓存实现
- `common/performance.go` — 性能监控实现
- `middleware/performance.go` — 性能监控中间件
