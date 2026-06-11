# system_monitor_unix.go 代码阅读文档

## 1. 全局总结

`system_monitor_unix.go` 是 Unix/Linux/macOS 平台的磁盘空间信息获取实现文件。它通过 Go 的构建标签（`//go:build !windows`）确保仅在非 Windows 平台上编译。文件实现了 `GetDiskSpaceInfo()` 函数，利用 Unix 系统调用 `statfs` 获取指定路径所在磁盘的总空间、可用空间、已用空间及使用百分比。

## 2. 依赖关系

### 标准库依赖
- `os` — `os.TempDir()` 获取系统临时目录作为后备路径

### 外部依赖
- `golang.org/x/sys/unix` — 提供 `unix.Statfs()` 系统调用和 `unix.Statfs_t` 结构体

### 项目内部依赖
- `GetDiskCachePath()` — 获取缓存目录路径（来自 `common` 包内其他文件）
- `DiskSpaceInfo` — 磁盘空间信息结构体（定义在 `system_monitor.go` 中）

## 3. 类型定义

本文件未定义新的类型，使用了以下类型：

- `unix.Statfs_t` — Unix 文件系统状态结构体，包含块大小、块数量、可用块数等字段
- `DiskSpaceInfo` — 磁盘空间信息结构体（定义在 `system_monitor.go` 中）

## 4. 函数详解

### 4.1 `GetDiskSpaceInfo() DiskSpaceInfo`

**功能**：获取缓存目录所在磁盘的空间信息（Unix/Linux/macOS 实现）。

**行为**：
1. 调用 `GetDiskCachePath()` 获取缓存目录路径
2. 如果缓存路径为空，回退到系统临时目录 `os.TempDir()`
3. 调用 `unix.Statfs()` 获取文件系统状态信息
4. 计算磁盘空间：
   - **Total** = `stat.Blocks * stat.Bsize`（总空间）
   - **Free** = `stat.Bavail * stat.Bsize`（对非 root 用户可用的空闲空间）
   - **Used** = `Total - stat.Bfree * bsize`（已用空间）
5. 如果总空间大于 0，计算使用百分比 `UsedPercent = Used / Total * 100`
6. 返回 `DiskSpaceInfo` 结构体

**返回值**：`DiskSpaceInfo` 结构体，包含 Total、Free、Used、UsedPercent 字段。

**错误处理**：如果 `unix.Statfs()` 调用失败，返回空的 `DiskSpaceInfo{}`（所有字段为零值）。

## 5. 关键逻辑分析

### 5.1 构建标签

文件顶部的 `//go:build !windows` 确保此文件仅在非 Windows 平台上编译。这是 Go 的跨平台编译机制，与 `system_monitor_windows.go` 形成互斥编译对。

### 5.2 磁盘空间计算

`unix.Statfs_t` 结构体的关键字段：
- `Bsize` — 文件系统块大小（字节）
- `Blocks` — 文件系统总块数
- `Bavail` — 对非超级用户可用的空闲块数
- `Bfree` — 空闲块总数（包括超级用户保留的块）

计算逻辑：
```
Total = Blocks * Bsize        // 总磁盘空间
Free  = Bavail * Bsize        // 可用空间（非 root 用户可用）
Used  = Total - Bfree * Bsize // 已用空间
```

**注意**：`Bavail` 和 `Bfree` 的区别：
- `Bfree` 是总空闲块数（包括 root 保留块）
- `Bavail` 是非 root 用户可用的空闲块数
- 这里使用 `Bavail` 作为 `Free`，`Bfree` 用于计算 `Used`，这种设计更贴近实际可用空间

### 5.3 FreeBSD 兼容性

代码注释中提到"显式转换以兼容 FreeBSD，其字段类型为 int64"。`unix.Statfs_t` 的字段在不同 Unix 系统中可能有不同的类型（如 `int32`、`int64`、`uint32`），通过 `uint64()` 显式转换确保在所有平台上正确计算。

### 5.4 路径回退策略

如果 `GetDiskCachePath()` 返回空字符串，函数回退到 `os.TempDir()`，确保始终能获取到磁盘空间信息。

## 6. 关联文件

- **`system_monitor.go`** — 定义 `DiskSpaceInfo` 结构体和 `GetSystemStatus()` 等跨平台接口
- **`system_monitor_windows.go`** — Windows 平台的 `GetDiskSpaceInfo()` 实现（互斥编译）
- **`common/env.go`**（推测）— 可能包含 `GetDiskCachePath()` 的定义
