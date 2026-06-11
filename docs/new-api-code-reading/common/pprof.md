# pprof.go 代码阅读文档

## 1. 全局总结

本文件实现了 CPU 使用率的定时监控功能。当 CPU 使用率超过阈值（80%）时，自动采集 pprof 性能分析文件并保存到本地 `./pprof/` 目录。监控以无限循环方式运行，每 30 秒检测一次 CPU 状态，每次采样持续 10 秒。

## 2. 依赖关系

| 包名 | 用途 |
|------|------|
| `fmt` | 格式化 pprof 文件名 |
| `os` | 文件/目录操作（创建、状态检查） |
| `runtime/pprof` | Go 内置 CPU profiling |
| `time` | 定时和采样持续时间控制 |
| `github.com/shirou/gopsutil/cpu` | 跨平台 CPU 使用率采集 |

**内部依赖**：
- `common/sys_log.go` — `SysLog()` 函数记录错误日志

## 3. 类型定义

本文件无自定义类型定义。

## 4. 函数详解

### `Monitor()`

**功能**：定时监控 CPU 使用率，超过阈值时自动采集 pprof 文件。

**执行流程**：
1. 进入无限循环
2. 使用 `cpu.Percent(time.Second, false)` 采集 1 秒的 CPU 使用率（单核聚合）
3. 若 CPU 使用率 > 80%：
   - 检查 `./pprof` 目录是否存在，不存在则创建
   - 创建带时间戳的 pprof 文件：`./pprof/cpu-YYYYMMDDHHmmss.pprof`
   - 启动 CPU profiling（`pprof.StartCPUProfile`）
   - 等待 10 秒进行采样
   - 停止 profiling 并关闭文件
4. 每次检测间隔 30 秒

## 5. 关键逻辑分析

1. **硬编码阈值**：CPU 阈值固定为 80%，未使用 `performance_config.go` 中的配置，存在不一致
2. **采样时长与注释不符**：注释写 "profile for 30 seconds"，实际 `time.Sleep` 为 10 秒
3. **错误处理**：目录创建失败、文件创建失败、profiling 启动失败均通过 `SysLog` 记录并 `continue` 跳过本轮
4. **panic 风险**：`cpu.Percent` 返回错误时直接 `panic`，在生产环境中可能导致进程崩溃
5. **文件命名**：使用时间戳格式 `20060102150405` 确保文件名唯一
6. **goroutine 设计**：该函数设计为在独立 goroutine 中运行（无限循环），不返回

## 6. 关联文件

- `common/performance_config.go` — 性能监控配置（CPU 阈值配置），但 `Monitor()` 未使用
- `common/sys_log.go` — 错误日志记录
- `main.go` 或启动入口 — 可能通过 `go common.Monitor()` 启动监控协程
