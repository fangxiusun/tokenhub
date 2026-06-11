# channel_upstream_update.go 代码阅读文档

## 1. 全局总结

该文件实现了渠道上游模型自动同步功能。定期检测各渠道上游提供商的可用模型列表，与本地配置的模型列表进行对比，自动发现新增和删除的模型，并支持自动同步或手动确认后应用变更。

## 2. 依赖关系

- `common` — 通用工具、环境变量
- `constant` — 渠道类型、基础 URL
- `dto` — 渠道其他设置
- `model` — 渠道模型、数据库操作
- `relay/channel/gemini` — Gemini 模型获取
- `relay/channel/ollama` — Ollama 模型获取
- `service` — 通知、代理缓存
- `gin-gonic/gin` — HTTP 框架
- `samber/lo` — 集合操作

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `applyChannelUpstreamModelUpdatesRequest` | 应用模型更新的请求体 |
| `applyAllChannelUpstreamModelUpdatesResult` | 批量应用结果 |
| `detectChannelUpstreamModelUpdatesResult` | 检测结果 |
| `upstreamModelUpdateChannelSummary` | 通知中的渠道摘要 |

## 4. 函数详解

### 模型集合操作
- `normalizeModelNames` — 模型名去重和修剪
- `mergeModelNames` — 合并两个模型列表（去重）
- `subtractModelNames` — 从列表中移除指定模型
- `intersectModelNames` — 取交集
- `applySelectedModelChanges` — 应用选中的添加/删除变更

### 核心检测逻辑
- `collectPendingUpstreamModelChangesFromModels` — 对比本地和上游模型，考虑忽略列表和模型映射
- `collectPendingUpstreamModelChanges` — 获取上游模型并计算变更
- `checkAndPersistChannelUpstreamModelUpdates` — 检测并持久化变更，支持自动同步

### 上游模型获取
- `fetchChannelUpstreamModelIDs(channel)` — 根据渠道类型调用对应 API 获取模型列表

### HTTP 处理器
- `DetectChannelUpstreamModelUpdates` — 检测单个渠道的模型变更
- `ApplyChannelUpstreamModelUpdates` — 应用单个渠道的模型变更
- `DetectAllChannelUpstreamModelUpdates` — 批量检测
- `ApplyAllChannelUpstreamModelUpdates` — 批量应用

### 后台任务
- `StartChannelUpstreamModelUpdateTask` — 启动定时检测任务
- `runChannelUpstreamModelUpdateTaskOnce` — 执行一次全量检测

## 5. 关键逻辑分析

- 检测间隔有最小限制（默认 300 秒），避免频繁请求上游
- 忽略列表支持精确匹配和正则匹配（`regex:` 前缀）
- 模型映射的源模型不会被误删
- 通知有 24 小时抑制窗口，避免重复通知
- 任务使用 `atomic.Bool` 防止并发执行
- 批量操作使用游标分页（`id > lastID`）避免内存问题

## 6. 关联文件

- `controller/channel.go` — 渠道 CRUD
- `model/channel.go` — 渠道模型
- `dto/option.go` — `ChannelOtherSettings` 定义
