# logger.go 代码阅读文档

## 1. 全局概述

本文件实现了系统的日志模块，提供分级日志输出（INFO/WARN/ERROR/DEBUG）、日志文件管理、配额格式化显示等功能。日志输出同时写入标准输出和日志文件，并支持自动轮转。

## 2. 依赖关系

- `context` — Go 标准库上下文包
- `fmt` — 格式化输出
- `io` — I/O 操作
- `log` — Go 标准库日志包
- `os` — 操作系统接口
- `path/filepath` — 文件路径处理
- `sync` — 同步原语
- `time` — 时间处理
- `github.com/QuantumNous/new-api/common` — 通用工具
- `github.com/QuantumNous/new-api/setting/operation_setting` — 运营设置
- `github.com/bytedance/gopkg/util/gopool` — 协程池
- `github.com/gin-gonic/gin` — Web 框架

## 3. 类型定义

本文件无自定义类型定义。

## 4. 函数详情

### GetCurrentLogPath

```go
func GetCurrentLogPath() string
```

获取当前日志文件路径（线程安全）。

### SetupLogger

```go
func SetupLogger()
```

设置日志输出：
1. 检查日志目录配置
2. 创建带时间戳的日志文件（格式：`oneapi-YYYYMMDDHHMMSS.log`）
3. 设置 `gin.DefaultWriter` 和 `gin.DefaultErrorWriter` 为多写入器（标准输出 + 日志文件）
4. 关闭旧的日志文件

### LogInfo

```go
func LogInfo(ctx context.Context, msg string)
```

输出 INFO 级别日志。

### LogWarn

```go
func LogWarn(ctx context.Context, msg string)
```

输出 WARN 级别日志。

### LogError

```go
func LogError(ctx context.Context, msg string)
```

输出 ERROR 级别日志。

### LogDebug

```go
func LogDebug(ctx context.Context, msg string, args ...any)
```

输出 DEBUG 级别日志（仅在调试模式下输出）。支持格式化参数。

### LogQuota

```go
func LogQuota(quota int) string
```

将配额值格式化为可读字符串，支持多种显示格式：
- 人民币（CNY）
- 自定义货币
- Token 点数
- 美元（默认）

### FormatQuota

```go
func FormatQuota(quota int) string
```

将配额值格式化为纯数字字符串（不带"额度"后缀）。

### LogJson

```go
func LogJson(ctx context.Context, msg string, obj any)
```

将对象序列化为 JSON 并输出 DEBUG 日志（仅用于测试）。

## 5. 关键逻辑分析

### 日志级别

| 级别 | 常量 | 说明 |
|------|------|------|
| INFO | `loggerINFO` | 一般信息 |
| WARN | `loggerWarn` | 警告信息 |
| ERR | `loggerError` | 错误信息 |
| DEBUG | `loggerDebug` | 调试信息（仅调试模式） |

### 日志格式

```
[LEVEL] YYYY/MM/DD - HH:MM:SS | REQUEST_ID | MESSAGE
```

- `LEVEL` — 日志级别
- `REQUEST_ID` — 请求 ID（从 Context 中获取，无则显示 "SYSTEM"）
- `MESSAGE` — 日志消息

### 日志文件轮转

- 最大日志条目数：`maxLogCount = 1000000`（100 万条）
- 达到上限后自动创建新的日志文件
- 使用 `gopool.Go` 异步执行文件切换，避免阻塞主流程
- 使用 `TryLock` 防止并发切换

### 配额格式化

根据 `operation_setting.GetQuotaDisplayType()` 配置决定显示格式：

| 显示类型 | 格式示例 |
|----------|----------|
| USD | `$0.001234 额度` |
| CNY | `¥0.008765 额度` |
| Custom | `¤0.005432 额度` |
| Tokens | `1234 点额度` |

### 线程安全

- `currentLogPath` 使用 `sync.RWMutex` 保护
- `setupLogLock` 使用 `sync.Mutex` 防止并发设置
- `common.LogWriterMu` 保护日志写入器的并发访问
- `logCount` 使用无锁计数（近似值，不需要精确）

### 日志写入器

- INFO 级别写入 `gin.DefaultWriter`
- WARN/ERROR 级别写入 `gin.DefaultErrorWriter`
- 两个写入器都是 `io.MultiWriter`，同时写入标准输出和日志文件

## 6. 相关文件

- `common/global.go` — `LogDir`、`LogWriterMu` 等全局变量
- `setting/operation_setting.go` — 配额显示类型配置
- `middleware/logger.go` — HTTP 请求日志中间件
- `main.go` — 启动时调用 `SetupLogger()`
