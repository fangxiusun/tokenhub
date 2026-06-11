# log.go 代码阅读文档

## 1. 全局总结

该文件实现了系统日志的查询和管理接口，包括获取所有日志、用户日志、按 Token 查询日志、日志统计以及历史日志清理。

## 2. 依赖关系

- `common` — 通用工具、分页查询
- `model` — 日志模型
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

无自定义类型定义。

## 4. 函数详解

### `GetAllLogs(c *gin.Context)`
管理员获取所有日志。支持按类型、时间范围、用户名、Token 名、模型名、渠道、组、requestId 等条件过滤，分页返回。

### `GetUserLogs(c *gin.Context)`
获取当前用户的日志。支持类似的过滤条件（不包括用户名和渠道）。

### `SearchAllLogs` / `SearchUserLogs`
已废弃的搜索接口，返回"该接口已废弃"。

### `GetLogByKey(c *gin.Context)`
按 Token ID 获取日志。

### `GetLogsStat(c *gin.Context)`
管理员日志统计。返回总额度、RPM、TPM。

### `GetLogsSelfStat(c *gin.Context)`
用户日志统计。返回当前用户的额度、RPM、TPM。

### `DeleteHistoryLogs(c *gin.Context)`
删除指定时间戳之前的历史日志（批量删除，每次最多 100 条）。

## 5. 关键逻辑分析

- 支持 `request_id` 和 `upstream_request_id` 双维度请求追踪
- 日志删除使用分批删除（每次 100 条），避免长事务
- 统计接口返回 quota、rpm、tpm 三个维度

## 6. 关联文件

- `model/log.go` — 日志数据模型
