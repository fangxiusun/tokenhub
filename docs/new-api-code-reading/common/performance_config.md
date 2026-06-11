# performance_config.go 代码阅读文档

## 1. 全局总结

本文件定义了性能监控的配置结构体和全局配置管理。使用 `sync/atomic.Value` 实现线程安全的配置读写，支持在运行时动态更新性能监控阈值（CPU/内存/磁盘）。文件非常简洁（33 行），仅包含配置定义、默认值初始化和 getter/setter。

## 2. 依赖关系

| 包名 | 用途 |
|------|------|
| `sync/atomic` | `atomic.Value` 提供线程安全的配置存储 |

**内部依赖**：无。

## 3. 类型定义

| 类型名 | 字段 | 类型 | 说明 |
|--------|------|------|------|
| `PerformanceMonitorConfig` | `Enabled` | `bool` | 是否启用性能监控 |
| | `CPUThreshold` | `int` | CPU 使用率阈值（百分比） |
| | `MemoryThreshold` | `int` | 内存使用率阈值（百分比） |
| | `DiskThreshold` | `int` | 磁盘使用率阈值（百分比） |

## 4. 函数详解

### `init()`
包初始化函数，在程序启动时自动执行。将默认配置存储到 `performanceMonitorConfig` 原子变量中：
- `Enabled`: `true`
- `CPUThreshold`: `90`
- `MemoryThreshold`: `90`
- `DiskThreshold`: `90`

### `GetPerformanceMonitorConfig() PerformanceMonitorConfig`
获取当前性能监控配置。通过 `atomic.Value.Load()` 读取并类型断言为 `PerformanceMonitorConfig`。

### `SetPerformanceMonitorConfig(config PerformanceMonitorConfig)`
设置性能监控配置。通过 `atomic.Value.Store()` 原子写入新配置。

## 5. 关键逻辑分析

1. **线程安全设计**：使用 `sync/atomic.Value` 而非互斥锁，实现无锁的读写操作。适合"读多写少"的配置场景
2. **默认阈值 90%**：CPU/内存/磁盘阈值均默认 90%，属于保守阈值，仅在极端情况下触发告警
3. **运行时可调**：通过 `SetPerformanceMonitorConfig` 可在不重启服务的情况下动态调整阈值（如通过管理 API）
4. **无依赖**：该文件不依赖项目内其他模块，是独立的配置管理单元

## 6. 关联文件

- `controller/performance.go` — 性能监控控制器，可能调用 getter/setter 管理配置
- `common/pprof.go` — CPU profiling 监控，硬编码阈值为 80%（与此处配置独立）
- `setting/` — 其他系统配置管理模块
