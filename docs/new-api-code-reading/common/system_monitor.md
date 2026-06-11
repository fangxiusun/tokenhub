# system_monitor.go 代码阅读文档

## 1. 全局总结

`system_monitor.go` 是系统监控的核心文件，提供了跨平台的系统状态监控功能。文件定义了 `DiskSpaceInfo` 和 `SystemStatus` 两个数据结构，实现了 CPU、内存和磁盘使用率的实时监控。通过 `atomic.Value` 存储最新的系统状态，支持并发安全读取。监控功能通过 `StartSystemMonitor()` 启动后台协程，根据性能监控配置定期更新系统状态。

## 2. 依赖关系

### 标准库依赖
- `sync/atomic` — 提供原子操作，保证系统状态的并发安全读写
- `time` — 提供定时器和时间计算

### 外部依赖
- `github.com/shirou/gopsutil/cpu` — 获取 CPU 使用率
- `github.com/shirou/gopsutil/mem` — 获取内存使用信息

### 项目内部依赖
- `GetDiskSpaceInfo()` — 获取磁盘空间信息（平台相关实现，在 `system_monitor_unix.go` 或 `system_monitor_windows.go` 中）
- `GetPerformanceMonitorConfig()` — 获取性能监控配置（来自 `common` 包内其他文件）

## 3. 类型定义

### 3.1 `DiskSpaceInfo`

```go
type DiskSpaceInfo struct {
    Total       uint64  `json:"total"`        // 总空间（字节）
    Free        uint64  `json:"free"`         // 可用空间（字节）
    Used        uint64  `json:"used"`         // 已用空间（字节）
    UsedPercent float64 `json:"used_percent"` // 使用百分比
}
```

**用途**：表示磁盘空间信息，带有 JSON 标签，支持序列化/反序列化。

### 3.2 `SystemStatus`

```go
type SystemStatus struct {
    CPUUsage    float64 // CPU 使用率（百分比）
    MemoryUsage float64 // 内存使用率（百分比）
    DiskUsage   float64 // 磁盘使用率（百分比）
}
```

**用途**：表示系统整体状态，包含 CPU、内存、磁盘三项核心指标。

## 4. 函数详解

### 4.1 `init()`

**功能**：包初始化函数，在程序启动时自动执行。

**行为**：将 `latestSystemStatus` 初始化为空的 `SystemStatus{}`，确保在监控启动前有有效的初始值。

### 4.2 `StartSystemMonitor()`

**功能**：启动系统监控后台协程。

**行为**：
1. 启动一个 goroutine
2. 在循环中检查性能监控配置 `GetPerformanceMonitorConfig()`
3. 如果监控未启用，休眠 30 秒后重新检查
4. 如果监控已启用，调用 `updateSystemStatus()` 更新系统状态，然后休眠 5 秒
5. 无限循环持续监控

**使用方式**：通常在服务启动时调用一次，启动后台监控。

### 4.3 `updateSystemStatus()`

**功能**：采集当前系统状态（CPU、内存、磁盘使用率）并存储。

**行为**：
1. 调用 `cpu.Percent(0, false)` 获取 CPU 使用率
   - 参数 `0` 表示采样间隔为 0（立即返回上次采样结果）
   - 参数 `false` 表示获取所有 CPU 的整体使用率
2. 调用 `mem.VirtualMemory()` 获取内存使用信息
3. 调用 `GetDiskSpaceInfo()` 获取磁盘空间信息
4. 将采集到的状态存储到 `latestSystemStatus`

**错误处理**：每个采集步骤都有独立的错误检查，单个指标采集失败不影响其他指标。

### 4.4 `GetSystemStatus() SystemStatus`

**功能**：获取当前系统状态（线程安全）。

**行为**：从 `latestSystemStatus` 原子读取并返回 `SystemStatus` 结构体。

**返回值**：最新的 `SystemStatus` 结构体，包含 CPU、内存、磁盘使用率。

## 5. 关键逻辑分析

### 5.1 原子值存储

`latestSystemStatus` 使用 `atomic.Value` 类型存储，提供以下优势：
- **并发安全**：`Store()` 和 `Load()` 是原子操作，无需额外加锁
- **无锁读取**：`GetSystemStatus()` 可以在高并发下无锁读取最新状态
- **整体替换**：每次更新都是整体替换 `SystemStatus` 结构体，保证读取的一致性

### 5.2 CPU 使用率的采样特性

`cpu.Percent(0, false)` 的行为：
- 第一次调用：可能返回错误或不准确的值（因为没有历史数据）
- 后续调用：返回自上次调用以来的 CPU 使用率
- 代码注释说明了这个特性，提醒开发者在循环中使用会逐渐准确

### 5.3 监控间隔控制

监控间隔由配置动态控制：
- **监控关闭**：每 30 秒检查一次配置（避免频繁检查）
- **监控开启**：每 5 秒更新一次系统状态

这种设计允许运行时动态开关监控功能。

### 5.4 跨平台磁盘信息

`GetDiskSpaceInfo()` 在不同平台上有不同实现：
- Unix/Linux/macOS：使用 `statfs` 系统调用
- Windows：使用 `GetDiskFreeSpaceExW` API

调用方无需关心平台差异，统一使用 `DiskSpaceInfo` 结构体。

### 5.5 错误容忍设计

`updateSystemStatus()` 中每个指标的采集都是独立的：
- CPU 采集失败 → `CPUUsage` 保持为 0
- 内存采集失败 → `MemoryUsage` 保持为 0
- 磁盘采集失败 → `DiskUsage` 保持为 0

单个指标失败不会影响其他指标的更新，提高了监控的鲁棒性。

## 6. 关联文件

- **`system_monitor_unix.go`** — Unix/Linux/macOS 平台的 `GetDiskSpaceInfo()` 实现
- **`system_monitor_windows.go`** — Windows 平台的 `GetDiskSpaceInfo()` 实现
- **`common/performance.go`**（推测）— 可能包含 `GetPerformanceMonitorConfig()` 的配置定义
- **`controller/dashboard.go`**（推测）— 可能调用 `GetSystemStatus()` 在管理面板展示系统状态
- **`common/sys_log.go`** — 系统日志工具，用于输出监控相关日志
