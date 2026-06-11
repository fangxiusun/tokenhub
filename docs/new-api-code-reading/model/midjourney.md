# midjourney.go 代码阅读文档

## 1. 全局总结
该文件定义了 Midjourney 任务的数据模型和相关的数据库操作函数。Midjourney 是一个用于管理 AI 图像生成任务的模型，记录了用户提交的图像生成请求及其状态、进度、结果等信息。该文件提供了任务的查询、插入、更新、批量更新和计数等功能，是 Midjourney 任务管理的核心数据层组件。

## 2. 依赖关系
- **GORM**: 用于数据库 ORM 操作，包括模型定义、查询构建、数据持久化等。
- **package model**: 当前包，提供数据库连接（DB）和跨数据库兼容性工具。

## 3. 类型定义
### Midjourney 结构体
表示 Midjourney 任务的数据模型，包含以下字段：
- `Id`: 任务 ID（主键）
- `Code`: 任务代码
- `UserId`: 用户 ID（建立索引）
- `Action`: 任务动作类型（如生成、变体等）
- `MjId`: Midjourney 任务唯一标识符（建立索引）
- `Prompt`: 用户输入的提示词
- `PromptEn`: 英文提示词
- `Description`: 任务描述
- `State`: 任务状态
- `SubmitTime`: 提交时间戳（建立索引）
- `StartTime`: 开始时间戳（建立索引）
- `FinishTime`: 完成时间戳（建立索引）
- `ImageUrl`: 生成图像的 URL
- `VideoUrl`: 生成视频的 URL
- `VideoUrls`: 多个视频 URL 的列表
- `Status`: 任务状态（建立索引）
- `Progress`: 任务进度（建立索引）
- `FailReason`: 失败原因
- `ChannelId`: 渠道 ID
- `Quota`: 配额消耗
- `Buttons`: 按钮配置
- `Properties`: 扩展属性（JSON 格式）

### TaskQueryParams 结构体
用于封装查询参数，支持按渠道 ID、MjId、时间范围进行过滤：
- `ChannelID`: 渠道 ID
- `MjID`: Midjourney 任务 ID
- `StartTimestamp`: 开始时间戳
- `EndTimestamp`: 结束时间戳

## 4. 函数详解
### 查询函数
- **GetAllUserTask(userId int, startIdx int, num int, queryParams TaskQueryParams) []*Midjourney**
  - 查询指定用户的所有任务，支持分页和条件过滤（MjId、时间范围）。
  - 返回任务列表，按 ID 降序排列。

- **GetAllTasks(startIdx int, num int, queryParams TaskQueryParams) []*Midjourney**
  - 查询所有任务（管理员视图），支持分页和条件过滤（渠道 ID、MjId、时间范围）。
  - 返回任务列表，按 ID 降序排列。

- **GetAllUnFinishTasks() []*Midjourney**
  - 查询所有未完成的任务（进度不为 "100%"）。
  - 返回未完成任务列表。

- **GetByOnlyMJId(mjId string) *Midjourney**
  - 根据 MjId 查询单个任务。
  - 返回第一个匹配的任务，未找到返回 nil。

- **GetByMJId(userId int, mjId string) *Midjourney**
  - 根据用户 ID 和 MjId 查询单个任务。
  - 返回第一个匹配的任务，未找到返回 nil。

- **GetByMJIds(userId int, mjIds []string) []*Midjourney**
  - 根据用户 ID 和多个 MjId 批量查询任务。
  - 返回匹配的任务列表。

- **GetMjByuId(id int) *Midjourney**
  - 根据任务 ID 查询单个任务。
  - 返回第一个匹配的任务，未找到返回 nil。

### 更新函数
- **UpdateProgress(id int, progress string) error**
  - 更新指定任务的进度字段。

- **(midjourney *Midjourney) Insert() error**
  - 插入新的 Midjourney 任务记录。

- **(midjourney *Midjourney) Update() error**
  - 更新整个 Midjourney 任务记录（使用 Save 方法）。

- **(midjourney *Midjourney) UpdateWithStatus(fromStatus string) (bool, error)**
  - 条件更新：只有当任务当前状态为 fromStatus 时才执行更新（CAS 操作）。
  - 返回是否成功更新（是否赢得竞争）。

### 批量更新函数
- **MjBulkUpdate(mjIds []string, params map[string]any) error**
  - 根据 MjId 列表批量更新任务字段。

- **MjBulkUpdateByTaskIds(taskIDs []int, params map[string]any) error**
  - 根据任务 ID 列表批量更新任务字段。

### 计数函数
- **CountAllTasks(queryParams TaskQueryParams) int64**
  - 统计所有任务数量（管理员视图），支持条件过滤。
  - 返回总数。

- **CountAllUserTask(userId int, queryParams TaskQueryParams) int64**
  - 统计指定用户的任务数量，支持条件过滤。
  - 返回总数。

## 5. 关键逻辑分析
- **查询构建模式**：使用 GORM 的链式查询构建器，通过 `Where` 方法动态添加过滤条件，实现灵活的查询。
- **条件更新（CAS）**：`UpdateWithStatus` 方法通过 `Where("status = ?", fromStatus)` 确保只有在状态未被其他进程修改的情况下才更新，避免并发冲突。
- **批量更新**：通过 `Updates` 方法一次性更新多个记录的指定字段，提高性能。
- **分页查询**：使用 `Limit` 和 `Offset` 实现分页，通过 `Order("id desc")` 按 ID 降序排列。
- **错误处理**：查询失败时返回 nil 或空值，调用方需要检查返回值是否为 nil。

## 6. 关联文件
- **model/main.go**: 提供全局数据库连接 DB 和跨数据库兼容性工具（如 commonGroupCol、commonTrueVal 等）。
- **controller/midjourney.go**: 可能包含处理 Midjourney 相关 HTTP 请求的控制器。
- **router/midjourney.go**: 可能包含 Midjourney 相关的路由定义。
- **service/midjourney.go**: 可能包含 Midjourney 业务逻辑服务层。
- **common/json.go**: 提供 JSON 序列化/反序列化工具，用于处理 Midjourney 结构体中的 JSON 字段（如 Properties）。