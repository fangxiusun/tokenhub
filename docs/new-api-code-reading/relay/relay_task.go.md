# relay_task.go 代码阅读文档

## 1. 全局总结

本文件实现了异步任务（视频生成、音乐生成等）的提交、查询和回调处理逻辑。核心函数 `RelayTaskSubmit` 完成从请求验证、计费预扣、上游提交到响应解析的完整流程。还包含任务状态查询（fetch）和 Gemini/Vertex 实时状态拉取逻辑。

## 2. 依赖关系

- `model`: 任务数据库模型、渠道模型
- `service`: 计费、错误处理、任务状态转换
- `relay/channel`: 任务适配器接口
- `relay/common`: RelayInfo 上下文
- `relay/helper`: 价格计算、模型映射
- `dto`: 任务 DTO 定义

## 3. 类型定义

### `TaskSubmitResult`
```go
type TaskSubmitResult struct {
    UpstreamTaskID string  // 上游任务 ID
    TaskData       []byte  // 原始响应数据
    Platform       constant.TaskPlatform  // 任务平台
    Quota          int     // 预扣配额
}
```

## 4. 函数详解

### `ResolveOriginTask(c, info) *dto.TaskError`
- **功能**: 处理基于已有任务的提交（remix/continuation）
- **逻辑**: 查找原始任务 → 提取模型名称 → 锁定原始渠道 → 提取 remix 参数（时长、分辨率）

### `RelayTaskSubmit(c, info) (*TaskSubmitResult, *dto.TaskError)`
- **功能**: 完成 task 提交的全部流程
- **流程**: InitChannelMeta → 确定 platform → 创建适配器 → 验证请求 → 模型映射 → 价格计算 → 计费估算 → OtherRatios 应用 → 预扣费 → 构建请求 → 发送请求 → 解析响应 → 提交后计费调整
- **关键点**: 仅首次提交时预扣费，重试时复用已有的 Billing

### `RelayTaskFetch(c, relayMode) *dto.TaskError`
- **功能**: 查询任务状态
- **支持的模式**: SunoFetchByID, SunoFetch, VideoFetchByID
- **路由表**: 使用 `fetchRespBuilders` map 进行模式分发

### `tryRealtimeFetch(task, isOpenAIVideoAPI) []byte`
- **功能**: 从 Gemini/Vertex 实时拉取任务状态
- **条件**: 仅当渠道类型为 Gemini 或 Vertex 时触发
- **逻辑**: 调用 adaptor.FetchTask → 解析响应 → 更新任务状态到数据库

### `TaskModel2Dto(task) *dto.TaskDto`
- **功能**: 将数据库 Task 模型转换为 DTO 格式

## 5. 关键逻辑分析

1. **计费流程**: 价格计算 → OtherRatios 应用 → 预扣费 → 上游提交 → AdjustBillingOnSubmit 调整
2. **Remix 支持**: 从原始任务提取模型名称和计费参数，锁定到原始渠道
3. **实时状态拉取**: Gemini/Vertex 支持 fetch 时直接从上游拉取最新状态，避免仅依赖数据库缓存
4. **状态映射**: `mapTaskStatusToSimple` 将内部状态映射为简化的对外状态（succeeded/failed/queued/processing）

## 6. 关联文件

- `relay/channel/adapter.go`: TaskAdaptor 接口定义
- `model/task.go`: Task 数据库模型
- `service/billing.go`: 计费服务
- `dto/task.go`: 任务 DTO 定义
