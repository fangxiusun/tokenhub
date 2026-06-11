# task_polling.go 代码阅读文档

## 1. 全局总结

该文件实现异步任务的轮询系统，每 15 秒检查未完成的任务，支持 Suno、视频、Midjourney 等平台的任务状态更新。包含超时任务清理和任务失败退款逻辑。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | 日志 |
| `constant` | 任务常量 |
| `dto` | 任务响应结构体 |
| `model` | 任务操作 |
| `relaycommon` | RelayInfo、TaskInfo |
| `taskcommon` | 进度常量 |
| `lo` | 条件工具 |

## 3. 类型定义

### `TaskPollingAdaptor` 接口
- `Init(info)` — 初始化
- `FetchTask(baseURL, key, body, proxy)` — 获取任务
- `ParseTaskResult(body)` — 解析结果
- `AdjustBillingOnComplete(task, taskResult) int` — 完成时计费调整

## 4. 核心函数

### `TaskPollingLoop()`
主轮询循环（15秒间隔）：
1. 清理超时任务
2. 查询未完成任务
3. 按平台分发更新

### `sweepTimedOutTasks(ctx)`
超时任务清理：
- 检查 TaskTimeoutMinutes 配置
- 使用 CAS 防止覆盖正常状态
- 旧系统遗留任务不退款

### `DispatchPlatformUpdate(platform, taskChannelM, taskM)`
按平台分发更新

### `updateSunoTasks(ctx, channelId, taskIds, taskM)`
Suno 任务批量更新

### `updateVideoSingleTask(ctx, adaptor, ch, taskId, taskM)`
单个视频任务更新：
1. 获取任务状态
2. 更新任务进度
3. 成功时：结算计费
4. 失败时：退款

### `settleTaskBillingOnComplete(ctx, adaptor, task, taskResult)`
任务完成时的计费调整：
1. 按次计费跳过差额结算
2. 优先使用 adaptor 计算的额度
3. 回退到 token 重算
4. 无调整时保持预扣额度

## 5. 关键逻辑分析

1. **CAS 更新**：使用 UpdateWithStatus 防止并发覆盖
2. **超时处理**：旧系统遗留任务特殊处理（不退款）
3. **计费优先级**：adaptor 调整 > token 重算 > 保持预扣
4. **429 处理**：速率限制错误不认为是任务失败

## 6. 关联文件

- `task_billing.go` — 任务计费
- `relay/channel/task/` — 任务适配器
