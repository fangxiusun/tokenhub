# task.go 代码阅读文档

## 1. 全局总结

`task.go` 实现了异步任务管理系统，用于处理视频生成、图像生成等需要异步执行的 AI 任务。它定义了任务的数据模型（包含公开数据和私有数据的分离）、状态机（NOT_START → SUBMITTED → QUEUED → IN_PROGRESS → SUCCESS/FAILURE）、查询接口、以及基于 CAS（Compare-And-Swap）的并发安全更新机制。任务系统支持多平台（不同 AI 提供商），并集成了计费上下文用于异步退款和差额结算。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `bytes` | 字节数组比较（快照相等性判断） |
| `database/sql/driver` | 数据库驱动接口（实现 `Scan`/`Value` 方法） |
| `encoding/json` | JSON 序列化/反序列化 |
| `time` | 时间处理 |
| `github.com/QuantumNous/new-api/common` | 通用工具（JSON 操作、随机字符生成） |
| `github.com/QuantumNous/new-api/constant` | 常量定义（TaskPlatform、ChannelType） |
| `github.com/QuantumNous/new-api/dto` | 数据传输对象（视频状态、OpenAI Video 格式） |
| `github.com/QuantumNous/new-api/relay/common` | Relay 信息结构体 |

## 3. 类型定义

### TaskStatus（类型别名）
`string` 类型的任务状态。定义了 7 个状态常量：
- `NOT_START`：未开始
- `SUBMITTED`：已提交
- `QUEUED`：排队中
- `IN_PROGRESS`：进行中
- `FAILURE`：失败
- `SUCCESS`：成功
- `UNKNOWN`：未知

方法 `ToVideoStatus()` 将任务状态映射为 OpenAI Video API 的状态格式。

### Task 结构体
任务核心数据模型：
- `ID`：自增主键
- `CreatedAt`/`UpdatedAt`：时间戳
- `TaskID`：第三方任务 ID（外部系统使用）
- `Platform`：平台标识
- `UserId`/`ChannelId`：关联的用户和渠道
- `Group`：用户分组（用于计费）
- `Quota`：已消耗额度
- `Action`：任务类型（如 song、lyrics 等）
- `Status`：任务状态
- `FailReason`：失败原因
- `SubmitTime`/`StartTime`/`FinishTime`：各阶段时间
- `Progress`：进度百分比
- `Properties`：公开属性（输入、上游模型名、原始模型名）
- `PrivateData`：私有数据（不返回给用户，包含 API Key、上游任务 ID、结果 URL、计费上下文）
- `Data`：通用 JSON 数据

### Properties 结构体
任务公开属性，实现 `Scan`/`Value` 接口支持 GORM 自动 JSON 序列化。

### TaskPrivateData 结构体
任务私有数据，包含：
- `Key`：API Key（Gemini/VertexAI 渠道需要）
- `UpstreamTaskID`：上游真实任务 ID
- `ResultURL`：结果 URL
- `BillingSource`：计费来源（wallet/subscription）
- `SubscriptionId`/`TokenId`：计费关联 ID
- `BillingContext`：计费参数快照

### TaskBillingContext 结构体
计费参数快照，记录提交时的模型价格、分组倍率、模型倍率等，用于轮询阶段重新计算额度。

### SyncTaskQueryParams 结构体
任务查询参数，包含平台、渠道 ID、任务 ID、用户 ID、操作类型、状态、时间范围等。

### taskSnapshot 结构体（私有）
任务快照，用于 CAS 比较。包含状态、进度、时间、失败原因、结果 URL、数据。

### TaskQuotaUsage 结构体
任务额度使用统计，包含模式和计数。

## 4. 函数详解

### 数据设置/获取

#### (t *Task) SetData(data any)
将任意数据序列化为 JSON 并设置到 `Data` 字段。

#### (t *Task) GetData(v any) error
从 `Data` 字段反序列化到指定结构体。

#### (t *Task) GetUpstreamTaskID() string
获取上游任务 ID。优先返回 `PrivateData.UpstreamTaskID`，回退到 `TaskID`（兼容旧数据）。

#### (t *Task) GetResultURL() string
获取结果 URL。优先返回 `PrivateData.ResultURL`，回退到 `FailReason`（兼容旧数据）。

### 任务生成

#### GenerateTaskID() string
生成 `task_` 前缀 + 32 位随机字符的公开任务 ID。

#### InitTask(platform constant.TaskPlatform, relayInfo *commonRelay.RelayInfo) *Task
初始化任务对象。设置渠道信息、上游模型名、原始模型名、API Key（Gemini/VertexAI）、公开任务 ID 等。

### 数据库操作

#### TaskGetAllUserTask(userId int, startIdx, num int, queryParams SyncTaskQueryParams) []*Task
分页查询用户的任务列表，支持多条件过滤。

#### TaskGetAllTasks(startIdx, num int, queryParams SyncTaskQueryParams) []*Task
分页查询所有任务（管理员用），支持多条件过滤。

#### GetTimedOutUnfinishedTasks(cutoffUnix int64, limit int) []*Task
获取超时未完成的任务（提交时间早于截止时间，且未完成/未失败）。

#### GetAllUnFinishSyncTasks(limit int) []*Task
获取所有未完成的任务。

#### GetByOnlyTaskId(taskId string) (*Task, bool, error)
根据第三方任务 ID 查询任务。

#### GetByTaskId(userId int, taskId string) (*Task, bool, error)
根据用户 ID 和任务 ID 查询任务。

#### GetByTaskIds(userId int, taskIds []any) ([]*Task, error)
批量查询用户任务。

#### (Task *Task) Insert() error
插入新任务。

#### (Task *Task) Update() error
更新任务（Save 方法）。

#### (t *Task) UpdateWithStatus(fromStatus TaskStatus) (bool, error)
**CAS 更新**：仅当任务当前状态为 `fromStatus` 时才更新。返回 `(true, nil)` 表示更新成功，`(false, nil)` 表示状态已被其他进程修改。使用 `Select("*").Updates()` 而非 `Save()`，避免 GORM 的 INSERT ON CONFLICT 行为绕过 CAS。

#### TaskBulkUpdate(taskIds []string, params map[string]any) error
按第三方任务 ID 批量更新（无 CAS 保护）。

#### TaskBulkUpdateByID(ids []int64, params map[string]any) error
按主键 ID 批量更新（无 CAS 保护，不用于计费流程）。

### 快照

#### (t *Task) Snapshot() taskSnapshot
创建任务快照，用于 CAS 比较。

#### (s taskSnapshot) Equal(other taskSnapshot) bool
比较两个快照是否相等，使用 `bytes.Equal` 比较 JSON 数据。

### 统计

#### TaskCountAllTasks(queryParams SyncTaskQueryParams) int64
统计所有匹配条件的任务总数。

#### TaskCountAllUserTask(userId int, queryParams SyncTaskQueryParams) int64
统计用户匹配条件的任务总数。

### 格式转换

#### (t *Task) ToOpenAIVideo() *dto.OpenAIVideo
将任务转换为 OpenAI Video API 格式。

## 5. 关键逻辑分析

**CAS 并发控制**：`UpdateWithStatus` 使用数据库行级 WHERE 条件实现 Compare-And-Swap，确保同一任务不会被多个并发进程同时更新到终态。这对于计费相关的状态转换（成功/失败/超时）至关重要，避免重复退款或结算。

**公开/私有数据分离**：`PrivateData` 字段使用 `json:"-"` 标签不返回给用户，但存储在数据库中供内部使用。这保护了 API Key、计费上下文等敏感信息。

**历史兼容**：`GetUpstreamTaskID()` 和 `GetResultURL()` 都有回退逻辑，兼容旧数据格式。

**批量更新的 CAS 警告**：`TaskBulkUpdate` 和 `TaskBulkUpdateByID` 明确标注无 CAS 保护，不应用于计费生命周期流程。

## 6. 关联文件

- `model/task_cas_test.go`：CAS 更新的单元测试和集成测试
- `dto/video.go`：视频状态常量和 OpenAI Video 格式
- `relay/common/relay.go`：RelayInfo 结构体
- `constant/task.go`：TaskPlatform 常量
- `common/utils.go`：JSON 操作和随机字符生成
