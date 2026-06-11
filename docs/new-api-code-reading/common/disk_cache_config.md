# disk_cache_config.go 代码阅读文档

## 1. 全局总结

该文件实现了磁盘缓存的配置管理和统计功能。它定义了磁盘缓存的配置结构体、全局配置实例、配置读写函数，以及详细的缓存统计信息收集和更新机制。文件使用原子操作和读写锁确保并发安全，是磁盘缓存子系统的核心配置模块。

## 2. 依赖关系

| 包名 | 用途 |
|------|------|
| `sync` | 读写互斥锁（`sync.RWMutex`） |
| `sync/atomic` | 原子操作，用于无锁的统计计数更新 |

内部依赖：`GetDiskCacheInfo()` 函数（在 `disk_cache.go` 中定义）

## 3. 类型定义

### 3.1 DiskCacheConfig 结构体

```go
type DiskCacheConfig struct {
    Enabled     bool
    ThresholdMB int
    MaxSizeMB   int
    Path        string
}
```

**用途**：定义磁盘缓存的完整配置参数。

**字段说明**：
- `Enabled`: 是否启用磁盘缓存功能
- `ThresholdMB`: 触发磁盘缓存的请求体大小阈值（MB），小于此阈值的请求不使用磁盘缓存
- `MaxSizeMB`: 磁盘缓存的最大总大小（MB），防止缓存无限增长
- `Path`: 磁盘缓存目录路径，为空时使用系统临时目录

### 3.2 DiskCacheStats 结构体

```go
type DiskCacheStats struct {
    ActiveDiskFiles         int64 `json:"active_disk_files"`
    CurrentDiskUsageBytes   int64 `json:"current_disk_usage_bytes"`
    ActiveMemoryBuffers     int64 `json:"active_memory_buffers"`
    CurrentMemoryUsageBytes int64 `json:"current_memory_usage_bytes"`
    DiskCacheHits           int64 `json:"disk_cache_hits"`
    MemoryCacheHits         int64 `json:"memory_cache_hits"`
    DiskCacheMaxBytes       int64 `json:"disk_cache_max_bytes"`
    DiskCacheThresholdBytes int64 `json:"disk_cache_threshold_bytes"`
}
```

**用途**：收集和报告磁盘缓存的运行时统计信息，支持 JSON 序列化用于 API 暴露。

**字段说明**：
- `ActiveDiskFiles`: 当前活跃的磁盘缓存文件数量
- `CurrentDiskUsageBytes`: 当前磁盘缓存总大小（字节）
- `ActiveMemoryBuffers`: 当前内存缓存的数量
- `CurrentMemoryUsageBytes`: 当前内存缓存总大小（字节）
- `DiskCacheHits`: 磁盘缓存命中次数
- `MemoryCacheHits`: 内存缓存命中次数
- `DiskCacheMaxBytes`: 磁盘缓存最大限制（字节）
- `DiskCacheThresholdBytes`: 磁盘缓存阈值（字节）

## 4. 函数详解

### 4.1 配置管理函数

#### GetDiskCacheConfig() DiskCacheConfig

**功能**：获取完整的磁盘缓存配置

**实现**：使用读锁保护配置读取，返回配置的副本

#### SetDiskCacheConfig(config DiskCacheConfig)

**功能**：设置磁盘缓存配置

**实现**：使用写锁保护配置写入

#### IsDiskCacheEnabled() bool

**功能**：检查磁盘缓存是否启用

**实现**：使用读锁读取 `Enabled` 字段

#### GetDiskCacheThresholdBytes() int64

**功能**：获取磁盘缓存触发阈值（字节）

**实现**：
- 使用读锁读取 `ThresholdMB`
- 使用位左移 `<< 20` 将 MB 转换为字节

#### GetDiskCacheMaxSizeBytes() int64

**功能**：获取磁盘缓存最大大小（字节）

**实现**：
- 使用读锁读取 `MaxSizeMB`
- 使用位左移 `<< 20` 将 MB 转换为字节

#### GetDiskCachePath() string

**功能**：获取磁盘缓存目录路径

**实现**：使用读锁读取 `Path` 字段

### 4.2 统计信息函数

#### GetDiskCacheStats() DiskCacheStats

**功能**：获取完整的缓存统计信息快照

**实现**：
1. 使用原子操作读取所有统计字段
2. 调用 `GetDiskCacheMaxSizeBytes()` 和 `GetDiskCacheThresholdBytes()` 获取配置相关的统计
3. 返回统计信息的副本

#### IncrementDiskFiles(size int64)

**功能**：增加磁盘文件计数和使用量

**参数**：
- `size`: 新增文件的大小（字节）

**实现**：
- 使用 `atomic.AddInt64` 原子增加文件计数
- 使用 `atomic.AddInt64` 原子增加使用量

#### DecrementDiskFiles(size int64)

**功能**：减少磁盘文件计数和使用量

**参数**：
- `size`: 删除文件的大小（字节）

**实现**：
- 使用 `atomic.AddInt64` 原子减少计数
- 如果结果小于 0，使用 `atomic.StoreInt64` 重置为 0（防止负数）
- 对使用量执行相同的保护逻辑

#### IncrementMemoryBuffers(size int64)

**功能**：增加内存缓存计数和使用量

**参数**：
- `size`: 新增内存缓冲区的大小（字节）

**实现**：与 `IncrementDiskFiles` 类似，使用原子操作更新内存缓存统计

#### DecrementMemoryBuffers(size int64)

**功能**：减少内存缓存计数和使用量

**参数**：
- `size`: 释放内存缓冲区的大小（字节）

**实现**：与 `DecrementDiskFiles` 类似，使用原子操作更新内存缓存统计

#### IncrementDiskCacheHits()

**功能**：增加磁盘缓存命中次数计数

**实现**：使用 `atomic.AddInt64` 原子增加命中计数

#### IncrementMemoryCacheHits()

**功能**：增加内存缓存命中次数计数

**实现**：使用 `atomic.AddInt64` 原子增加命中计数

#### ResetDiskCacheStats()

**功能**：重置命中统计信息（不重置当前使用量）

**实现**：使用 `atomic.StoreInt64` 将命中计数重置为 0

#### ResetDiskCacheUsage()

**功能**：重置磁盘缓存使用量统计（用于清理缓存后）

**实现**：使用 `atomic.StoreInt64` 将文件计数和使用量重置为 0

#### SyncDiskCacheStats()

**功能**：从实际磁盘状态同步统计信息

**实现**：
1. 调用 `GetDiskCacheInfo()` 获取实际文件数量和总大小
2. 使用 `atomic.StoreInt64` 同步到统计变量

#### IsDiskCacheAvailable(requestSize int64) bool

**功能**：检查是否可以创建新的磁盘缓存

**参数**：
- `requestSize`: 请求的数据大小（字节）

**返回值**：`true` 表示可以创建新缓存

**实现**：
1. 检查磁盘缓存是否启用
2. 获取最大缓存大小
3. 使用原子操作读取当前使用量
4. 判断 `当前使用量 + 请求大小 <= 最大大小`

## 5. 关键逻辑分析

### 5.1 并发安全设计

文件采用两种并发控制机制：
- **读写锁（RWMutex）**：保护配置读写，允许多个读操作并发，写操作独占
- **原子操作（atomic）**：用于统计计数的更新，避免锁的开销

这种混合设计兼顾了性能和安全性：
- 配置读取频繁但写入较少，适合读写锁
- 统计计数更新频繁且需要高性能，适合原子操作

### 5.2 MB 到字节的转换

使用位左移 `<< 20` 而非乘法 `* 1024 * 1024`：
- 位运算性能更高
- 代码更简洁
- 语义更清晰（1 MB = 2^20 字节）

### 5.3 防止负数计数

`DecrementDiskFiles` 和 `DecrementMemoryBuffers` 函数包含负数检查：
- 使用 `atomic.AddInt64` 返回新值
- 如果结果小于 0，使用 `atomic.StoreInt64` 重置为 0
- 防止统计计数因并发操作或逻辑错误变为负数

### 5.4 统计信息的双重来源

磁盘缓存统计有两个来源：
- **实时更新**：通过 `Increment`/`Decrement` 函数实时维护
- **同步校正**：通过 `SyncDiskCacheStats` 从磁盘实际状态同步

这种设计允许：
- 高性能的实时统计
- 定期校正统计与实际状态的偏差

### 5.5 配置变更响应

`GetDiskCacheConfig` 返回配置的副本而非引用，确保：
- 调用方获得一致的配置快照
- 配置变更不会影响已获取的快照

## 6. 关联文件

- `disk_cache.go`：磁盘缓存的文件操作实现
- `setting/performance_setting.go`：磁盘缓存配置的管理和持久化
- `controller/` 目录：通过 API 暴露缓存统计信息
- `relay/` 目录：使用磁盘缓存处理大请求体
