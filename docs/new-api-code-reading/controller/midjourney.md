# midjourney.go 代码阅读文档

## 1. 全局总结

该文件实现了 Midjourney 任务的后台轮询更新和任务列表查询功能。定期检测未完成的 MJ 任务状态，更新进度、图片 URL 等信息，失败时退还用户额度。

## 2. 依赖关系

- `common` — 通用工具
- `dto` — MidjourneyDto 数据结构
- `logger` — 日志
- `model` — Midjourney 任务模型
- `service` — HTTP 客户端、模型名转换
- `setting` — MJ 通知和转发设置
- `setting/system_setting` — 服务器地址
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

无自定义类型定义。

## 4. 函数详解

### `UpdateMidjourneyTaskBulk()`
后台定时任务（每 15 秒）。流程：获取未完成任务 → 按渠道分组 → 调用上游 API 批量查询 → 更新任务状态 → 失败时退还额度。

### `checkMjTaskNeedUpdate(oldTask, newTask) bool`
检查任务是否需要更新（对比多个字段）。

### `GetAllMidjourney(c *gin.Context)`
管理员获取所有 MJ 任务列表（分页）。支持渠道 ID、MJ ID、时间范围过滤。

### `GetUserMidjourney(c *gin.Context)`
获取当前用户的 MJ 任务列表。

## 5. 关键逻辑分析

- 任务超时判断：超过 1 小时且进度未达 100% 则标记为失败
- 失败任务自动退还额度（`model.IncreaseUserQuota`）
- 支持图片 URL 转发（`MjForwardUrlEnabled`）
- 批量查询使用 `/mj/task/list-by-condition` 接口
- mj_id 为空的任务直接标记为 FAILURE

## 6. 关联文件

- `model/midjourney.go` — MJ 任务模型
- `dto/midjourney.go` — MJ 数据传输对象
