# pyro.go 代码阅读文档

## 1. 全局总结

本文件集成了 [Pyroscope](https://pyroscope.io/) 持续性能分析工具。通过环境变量配置 Pyroscope 服务端连接信息，启动后持续采集 CPU、内存分配、goroutine、mutex、block 等多种性能指标并上报到 Pyroscope 服务器。支持基本认证和自定义采样率。

## 2. 依赖关系

| 包名 | 用途 |
|------|------|
| `runtime` | 设置 mutex 和 block profiling 采样率 |
| `github.com/grafana/pyroscope-go` | Pyroscope 客户端 SDK |

**内部依赖**：
- `common/env.go` — `GetEnvOrDefaultString()` 和 `GetEnvOrDefault()` 读取环境变量

## 3. 类型定义

本文件无自定义类型定义。

## 4. 函数详解

### `StartPyroScope() error`

**功能**：初始化并启动 Pyroscope 持续性能分析。

**环境变量**：

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `PYROSCOPE_URL` | `""` | Pyroscope 服务端地址（为空则跳过启动） |
| `PYROSCOPE_APP_NAME` | `"new-api"` | 应用名称 |
| `PYROSCOPE_BASIC_AUTH_USER` | `""` | 基本认证用户名 |
| `PYROSCOPE_BASIC_AUTH_PASSWORD` | `""` | 基本认证密码 |
| `HOSTNAME` | `"new-api"` | 主机名标签 |
| `PYROSCOPE_MUTEX_RATE` | `5` | mutex profiling 采样率 |
| `PYROSCOPE_BLOCK_RATE` | `5` | block profiling 采样率 |

**采集的 Profile 类型**：
- `ProfileCPU` — CPU 使用率
- `ProfileAllocObjects` — 内存分配对象数
- `ProfileAllocSpace` — 内存分配空间
- `ProfileInuseObjects` — 内存使用对象数
- `ProfileInuseSpace` — 内存使用空间
- `ProfileGoroutines` — goroutine 数量
- `ProfileMutexCount` — mutex 竞争次数
- `ProfileMutexDuration` — mutex 等待时长
- `ProfileBlockCount` — block 次数
- `ProfileBlockDuration` — block 等待时长

## 5. 关键逻辑分析

1. **可选集成**：`PYROSCOPE_URL` 为空时直接返回 `nil`，不影响应用启动
2. **mutex/block profiling**：通过 `runtime.SetMutexProfileFraction` 和 `runtime.SetBlockProfileRate` 启用 Go 运行时的锁竞争和阻塞分析，采样率由环境变量控制
3. **标签机制**：通过 `Tags` 传入 `hostname`，便于在 Pyroscope UI 中按主机筛选
4. **Logger 为 nil**：Pyroscope SDK 的日志输出被禁用（`Logger: nil`），调试时可能需要修改
5. **错误传播**：`pyroscope.Start` 失败时返回错误，调用方需处理

## 6. 关联文件

- `common/env.go` — `GetEnvOrDefaultString()`, `GetEnvOrDefault()` 环境变量读取
- `common/pprof.go` — 另一种 CPU profiling 方式（本地文件），与 Pyroscope 互补
- `common/performance_config.go` — 性能监控配置
- `main.go` 或启动入口 — 调用 `StartPyroScope()` 初始化持续分析
