# system_monitor_windows.go 代码阅读文档

## 1. 全局总结

`system_monitor_windows.go` 是 Windows 平台的磁盘空间信息获取实现文件。它通过 Go 的构建标签（`//go:build windows`）确保仅在 Windows 平台上编译。文件实现了 `GetDiskSpaceInfo()` 函数，通过调用 Windows API `GetDiskFreeSpaceExW`（kernel32.dll）获取指定路径所在磁盘的总空间、可用空间、已用空间及使用百分比。

## 2. 依赖关系

### 标准库依赖
- `os` — `os.TempDir()` 获取系统临时目录作为后备路径
- `syscall` — 调用 Windows API 函数
- `unsafe` — 提供指针转换，用于传递内存地址给 Windows API

### 项目内部依赖
- `GetDiskCachePath()` — 获取缓存目录路径（来自 `common` 包内其他文件）
- `DiskSpaceInfo` — 磁盘空间信息结构体（定义在 `system_monitor.go` 中）

## 3. 类型定义

本文件未定义新的类型，使用了以下类型：

- `syscall.LazyDLL` — 延迟加载的 DLL 引用
- `syscall.LazyProc` — 延迟加载的函数引用
- `DiskSpaceInfo` — 磁盘空间信息结构体（定义在 `system_monitor.go` 中）

## 4. 函数详解

### 4.1 `GetDiskSpaceInfo() DiskSpaceInfo`

**功能**：获取缓存目录所在磁盘的空间信息（Windows 实现）。

**行为**：
1. 调用 `GetDiskCachePath()` 获取缓存目录路径
2. 如果缓存路径为空，回退到系统临时目录 `os.TempDir()`
3. 动态加载 `kernel32.dll`，获取 `GetDiskFreeSpaceExW` 函数地址
4. 将路径字符串转换为 UTF-16 指针
5. 调用 `GetDiskFreeSpaceExW` 获取磁盘空间信息：
   - `freeBytesAvailable` — 可用空间（字节）
   - `totalBytes` — 总空间（字节）
   - `totalFreeBytes` — 总空闲空间（字节）
6. 计算磁盘空间：
   - **Total** = `totalBytes`
   - **Free** = `freeBytesAvailable`
   - **Used** = `totalBytes - totalFreeBytes`
7. 如果总空间大于 0，计算使用百分比
8. 返回 `DiskSpaceInfo` 结构体

**返回值**：`DiskSpaceInfo` 结构体，包含 Total、Free、Used、UsedPercent 字段。

**错误处理**：
- 如果 `syscall.UTF16PtrFromString()` 转换失败，返回空结构体
- 如果 API 调用返回值为 0（失败），返回空结构体

## 5. 关键逻辑分析

### 5.1 构建标签

文件顶部的 `//go:build windows` 确保此文件仅在 Windows 平台上编译。与 `system_monitor_unix.go` 形成互斥编译对。

### 5.2 Windows API 调用

`GetDiskFreeSpaceExW` 是 Windows API 函数，用于获取磁盘空间信息。函数签名：
```c
BOOL GetDiskFreeSpaceExW(
  LPCWSTR lpDirectoryName,
  PULARGE_INTEGER lpFreeBytesAvailable,
  PULARGE_INTEGER lpTotalBytes,
  PULARGE_INTEGER lpTotalFreeBytes
);
```

参数说明：
- `lpDirectoryName` — 目录路径（UTF-16 编码）
- `lpFreeBytesAvailable` — 可用空间（非 root 用户可用）
- `lpTotalBytes` — 总空间
- `lpTotalFreeBytes` — 总空闲空间（包括系统保留）

### 5.3 动态加载 DLL

代码使用 `syscall.NewLazyDLL("kernel32.dll")` 延迟加载 kernel32.dll，避免在程序启动时就加载所有依赖的 DLL。这种方式可以提高启动速度，并在 DLL 不存在时提供更优雅的错误处理。

### 5.4 UTF-16 路径转换

Windows API 使用 UTF-16 编码的字符串。`syscall.UTF16PtrFromString()` 将 Go 字符串转换为 UTF-16 编码的指针，确保路径中的中文字符等能正确传递给 Windows API。

### 5.5 `freeBytesAvailable` vs `totalFreeBytes`

与 Unix 实现类似，这里区分了两种空闲空间：
- `freeBytesAvailable` — 可用空间（类似于 Unix 的 `Bavail`）
- `totalFreeBytes` — 总空闲空间（类似于 Unix 的 `Bfree`）
- `Used` = `totalBytes - totalFreeBytes`，使用总空闲空间来计算已用空间

### 5.6 unsafe.Pointer 使用

`unsafe.Pointer` 用于将 `uint64` 变量的地址转换为 `uintptr`，以便传递给 Windows API。这是 Go 调用 C 风格 API 的标准做法，但需要注意内存安全。

## 6. 关联文件

- **`system_monitor.go`** — 定义 `DiskSpaceInfo` 结构体和 `GetSystemStatus()` 等跨平台接口
- **`system_monitor_unix.go`** — Unix/Linux/macOS 平台的 `GetDiskSpaceInfo()` 实现（互斥编译）
- **`common/env.go`**（推测）— 可能包含 `GetDiskCachePath()` 的定义
