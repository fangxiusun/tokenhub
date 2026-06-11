# task.go 代码阅读文档

## 1. 全局总结

该文件实现了通用任务（Task）的后台轮询和查询接口。任务支持多平台（Kling、Jimeng 等），包含状态管理和用户信息填充。

## 2. 依赖关系

- `common` — 通用工具
- `constant` — 任务平台常量
- `dto` — 任务 DTO
- `model` — 任务模型
- `relay` — 模型到 DTO 转换
- `service` — 任务轮询
- `types` — Set 工具
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

无自定义类型定义。

## 4. 函数详解

### `UpdateTaskBulk()`
后台任务轮询入口，委托给 `service.TaskPollingLoop`。

### `GetAllTask(c *gin.Context)`
管理员获取所有任务。支持按平台、任务 ID、状态、操作类型、时间范围、渠道 ID 过滤。

### `GetUserTask(c *gin.Context)`
获取当前用户的任务列表。

### `tasksToDto(tasks, fillUser)`
将任务模型转换为 DTO，可选填充用户名。

## 5. 关键逻辑分析

- 管理员视图会批量查询用户信息并填充到任务 DTO
- 任务查询支持多维度过滤

## 6. 关联文件

- `model/task.go` — 任务模型
- `service/task.go` — 任务轮询逻辑
