# task_video.go 代码阅读文档

## 1. 全局总结

该文件实现了视频生成任务的后台轮询更新功能。支持多平台（Kling、Jimeng 等）的视频任务状态同步。

## 2. 依赖关系

- `common` — 通用工具
- `constant` — 任务平台常量
- `dto` — 任务数据结构
- `logger` — 日志
- `model` — 任务模型
- `relay` — 适配器
- `relay/channel` — 渠道适配器
- `relay/common` — RelayInfo
- `setting/ratio_setting` — 模型后缀配置

## 3. 类型定义

无。

## 4. 函数详解

### `UpdateVideoTaskAll(ctx, platform, taskChannelM, taskM) error`
视频任务批量更新入口。按渠道分组处理未完成的视频任务。

### `updateVideoTaskAll(ctx, platform, channelId, taskIds, taskM) error`
单渠道视频任务更新。调用上游 API 查询任务状态并更新本地记录。

## 5. 关键逻辑分析

- 任务超时判断：超过 1 小时且未完成则标记失败
- 失败任务自动退还额度
- 支持多平台的适配器模式

## 6. 关联文件

- `model/task.go` — 任务模型
- `controller/task.go` — 通用任务管理
