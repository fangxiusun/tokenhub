# performance.go 代码阅读文档

## 1. 全局总结

该文件实现了系统性能检查中间件 `SystemPerformanceCheck`，在处理 Relay 接口请求前检查 CPU、内存、磁盘使用率是否超过阈值。超载时返回 503 Service Unavailable 错误，根据请求路径自动选择 OpenAI 或 Claude 格式的错误响应。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `fmt` | 格式化错误消息 |
| `net/http` | HTTP 状态码 |
| `strings` | 路径前缀判断 |
| `github.com/QuantumNous/new-api/common` | 性能监控配置和系统状态获取 |
| `github.com/QuantumNous/new-api/types` | 错误类型定义 |
| `github.com/gin-gonic/gin` | Gin Web 框架 |

## 3. 类型定义

无自定义类型定义。

## 4. 函数详解

### `SystemPerformanceCheck() gin.HandlerFunc`

- **功能**：创建系统性能检查中间件。
- **逻辑**：
  - 请求路径以 `/v1/messages` 开头 → 使用 Claude 格式错误响应
  - 其他路径 → 使用 OpenAI 格式错误响应
  - 调用 `checkSystemPerformance()` 检查系统状态

### `checkSystemPerformance() *types.NewAPIError`

- **功能**：检查系统性能指标是否超过阈值。
- **检查项**：
  1. CPU 使用率 > `CPUThreshold` → 返回 `system_cpu_overloaded` 错误
  2. 内存使用率 > `MemoryThreshold` → 返回 `system_memory_overloaded` 错误
  3. 磁盘使用率 > `DiskThreshold` → 返回 `system_disk_overloaded` 错误
- **返回**：超载时返回 `*types.NewAPIError`，正常时返回 `nil`

## 5. 关键逻辑分析

- **条件检查**：仅在对应阈值配置 > 0 时才检查该指标，允许管理员选择性启用监控。
- **双格式错误响应**：根据 API 类型（OpenAI/Claude）返回不同格式的错误，保持与上游 API 的兼容性。
- **配置驱动**：阈值通过 `common.GetPerformanceMonitorConfig()` 获取，可在运行时动态调整。
- **优雅降级**：性能检查失败时返回 503 而非 500，明确表示服务暂时不可用。

## 6. 关联文件

- `common/performance.go` — `GetPerformanceMonitorConfig` 和 `GetSystemStatus` 实现
- `types/error.go` — `NewAPIError` 类型定义
