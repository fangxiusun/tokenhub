# performance.go 代码阅读文档

## 1. 全局总结

该文件实现了系统性能统计和管理接口，包括内存/磁盘/缓存统计、磁盘缓存清理、GC 强制执行、日志文件管理等功能。

## 2. 依赖关系

- `common` — 磁盘缓存、系统状态
- `logger` — 日志路径
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `PerformanceStats` | 性能统计总览 |
| `MemoryStats` | 内存统计 |
| `DiskCacheInfo` | 磁盘缓存目录信息 |
| `PerformanceConfig` | 性能配置 |
| `LogFileInfo` | 日志文件信息 |
| `LogFilesResponse` | 日志文件列表响应 |

## 4. 函数详解

### 性能统计
- `GetPerformanceStats` — 获取完整性能统计（缓存、内存、磁盘、配置）
- `ClearDiskCache` — 清理超过 10 分钟未使用的缓存
- `ResetPerformanceStats` — 重置性能统计计数器
- `ForceGC` — 强制执行垃圾回收

### 日志管理
- `GetLogFiles` — 获取日志文件列表（按文件名降序）
- `CleanupLogFiles` — 清理过期日志（支持按数量或天数）

## 5. 关键逻辑分析

- 缓存统计使用原子计数器，避免每次全量扫描磁盘
- 日志清理会跳过当前活跃的日志文件
- 日志文件名格式：`oneapi-*.log`

## 6. 关联文件

- `common/cache.go` — 磁盘缓存实现
- `logger/` — 日志系统
