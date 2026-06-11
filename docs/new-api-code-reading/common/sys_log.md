# sys_log.go 代码阅读文档

## 1. 全局总结

`sys_log.go` 是系统日志工具文件，提供了项目级别的日志输出函数。它封装了 Gin 框架的日志写入器（`gin.DefaultWriter` 和 `gin.DefaultErrorWriter`），通过读写锁（`sync.RWMutex`）保证并发安全，支持日志文件轮转时的安全切换。文件共定义了 4 个日志函数：`SysLog`（系统信息日志）、`SysError`（系统错误日志）、`FatalLog`（致命错误日志并退出进程）和 `LogStartupSuccess`（启动成功提示日志，带彩色输出）。

## 2. 依赖关系

### 标准库依赖
- `fmt` — 格式化输出日志内容
- `os` — `os.Exit(1)` 用于致命错误后终止进程
- `sync` — 提供 `sync.RWMutex` 读写锁
- `time` — 提供时间戳格式化

### 外部依赖
- `github.com/gin-gonic/gin` — 使用 Gin 框架的 `gin.DefaultWriter` 和 `gin.DefaultErrorWriter` 作为日志输出目标

### 项目内部依赖
- `GetNetworkIps()` — 获取网络 IP 地址列表（来自 `common` 包内其他文件）
- `IsRunningInContainer()` — 判断是否在容器中运行（来自 `common` 包内其他文件）
- `SystemName`、`Version` — 系统名称和版本号常量（来自 `common` 包内其他文件）

## 3. 类型定义

本文件未定义新的类型，但操作了以下外部类型的值：

- `gin.DefaultWriter`（`io.Writer`）— Gin 默认标准输出写入器
- `gin.DefaultErrorWriter`（`io.Writer`）— Gin 默认错误输出写入器
- `LogWriterMu`（`sync.RWMutex`）— 全局读写锁，保护日志写入器的并发访问

## 4. 函数详解

### 4.1 `SysLog(s string)`

**功能**：向标准输出写入系统级信息日志。

**行为**：
1. 获取当前时间
2. 获取读锁 `LogWriterMu.RLock()`
3. 向 `gin.DefaultWriter` 格式化输出 `[SYS] 时间 | 消息` 格式的日志
4. 释放读锁

**日志格式**：`[SYS] 2006/01/02 - 15:04:05 | 消息内容`

### 4.2 `SysError(s string)`

**功能**：向错误输出写入系统级错误日志。

**行为**：与 `SysLog` 类似，但输出目标为 `gin.DefaultErrorWriter`。

### 4.3 `FatalLog(v ...any)`

**功能**：向错误输出写入致命错误日志，然后终止进程。

**行为**：
1. 获取当前时间
2. 获取读锁
3. 向 `gin.DefaultErrorWriter` 格式化输出 `[FATAL] 时间 | 消息` 格式的日志
4. 释放读锁
5. 调用 `os.Exit(1)` 终止进程

**注意事项**：
- 使用可变参数 `v ...any`，允许传入任意类型的值
- `os.Exit(1)` 不会执行 `defer`，调用方需注意资源清理

### 4.4 `LogStartupSuccess(startTime time.Time, port string)`

**功能**：在服务启动成功后，输出格式化的启动成功提示信息（带 ANSI 彩色码）。

**行为**：
1. 计算启动耗时（毫秒）
2. 获取网络 IP 列表
3. 获取读锁（使用 `defer` 延迟释放）
4. 输出服务名称、版本号和启动耗时（绿色字体）
5. 如果不在容器中运行，输出本地访问地址 `http://localhost:{port}/`
6. 遍历所有网络 IP，输出网络访问地址
7. 所有输出均通过 `gin.DefaultWriter`（标准输出）

**输出格式示例**：
```
  new-api v1.0.0  ready in 123 ms

  ➜  Local:   http://localhost:3000/
  ➜  Network: http://192.168.1.100:3000/

```

**ANSI 颜色码**：`\033[32m` 绿色，`\033[1m` 加粗，`\033[0m` 重置。

## 5. 关键逻辑分析

### 5.1 读写锁并发控制

`LogWriterMu` 是一个 `sync.RWMutex`，用于保护 `gin.DefaultWriter` 和 `gin.DefaultErrorWriter` 的并发安全。设计模式如下：

- **读锁（RLock）**：所有日志写入操作（`SysLog`、`SysError`、`FatalLog`、`LogStartupSuccess`）都使用读锁，允许多个日志操作并发执行
- **写锁（Lock）**：日志文件轮转时（在其他文件中实现），需要获取写锁来替换写入器并关闭旧文件

这种设计保证了：
1. 正常日志写入不会互相阻塞
2. 日志轮转时能安全地替换写入器，不会出现数据竞争

### 5.2 日志输出格式统一

所有日志函数遵循统一的时间格式 `2006/01/02 - 15:04:05`（Go 的参考时间格式），前缀标识日志级别（`[SYS]` 或 `[FATAL]`），便于日志过滤和分析。

### 5.3 FatalLog 的进程终止

`FatalLog` 在写入日志后调用 `os.Exit(1)`，这是不可恢复的操作。调用方需要确保在调用前已完成所有必要的清理工作，因为 `defer` 函数不会被执行。

### 5.4 启动信息的条件输出

`LogStartupSuccess` 根据 `IsRunningInContainer()` 的返回值决定是否显示本地访问地址。在容器环境中，通常不需要本地访问地址，因为服务通过端口映射暴露。

## 6. 关联文件

- **`common/json.go`** — JSON 序列化工具函数
- **`common/env.go`**（推测）— 可能包含 `GetNetworkIps()`、`IsRunningInContainer()`、`SystemName`、`Version` 等定义
- **`common/system_monitor.go`** — 系统监控功能，使用 `SysLog` 输出监控信息
- **`common/performance.go`**（推测）— 可能包含 `GetPerformanceMonitorConfig()` 的配置
- **`router/router.go`**（推测）— 在服务启动时调用 `LogStartupSuccess` 输出启动信息
