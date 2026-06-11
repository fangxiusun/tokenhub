# task.go 代码阅读文档

## 1. 全局摘要

该文件定义了通用任务管理的数据结构，包括任务错误 `TaskError`、任务响应 `TaskResponse`、任务 DTO `TaskDto`，以及获取请求 `FetchReq`。用于管理异步任务的生命周期。

## 2. 依赖

- **标准库**：`encoding/json`

## 3. 类型定义

### TaskError 结构体
任务错误结构：
- `Code` (string)：错误代码
- `Message` (string)：错误消息
- `Data` (any)：错误数据
- `StatusCode` (int)：HTTP 状态码（不序列化）
- `LocalError` (bool)：本地错误标识（不序列化）
- `Error` (error)：Go 错误对象（不序列化）

### TaskData 接口
任务数据类型约束：
```go
type TaskData interface {
    SunoDataResponse | []SunoDataResponse | string | any
}
```

### TaskResponse[T TaskData] 结构体
通用任务响应（泛型）：
- `Code` (string)：状态码
- `Message` (string)：消息
- `Data` (T)：响应数据

### TaskDto 结构体
任务数据传输对象：
- `ID` (int64)：数据库 ID
- `CreatedAt` (int64)：创建时间
- `UpdatedAt` (int64)：更新时间
- `TaskID` (string)：任务 ID
- `Platform` (string)：平台
- `UserId` (int)：用户 ID
- `Group` (string)：分组
- `ChannelId` (int)：渠道 ID
- `Quota` (int)：配额消耗
- `Action` (string)：操作类型
- `Status` (string)：状态
- `FailReason` (string)：失败原因
- `ResultURL` (string)：结果 URL
- `SubmitTime` (int64)：提交时间
- `StartTime` (int64)：开始时间
- `FinishTime` (int64)：完成时间
- `Progress` (string)：进度
- `Properties` (any)：属性数据
- `Username` (string)：用户名
- `Data` (json.RawMessage)：原始数据

### FetchReq 结构体
获取请求：
- `IDs` ([]string)：任务 ID 数组

### 常量

**TaskSuccessCode**："success" - 成功状态码

## 4. 函数详情

### IsSuccess()
```go
func (t *TaskResponse[T]) IsSuccess() bool
```
**功能**：判断任务是否成功。

**逻辑**：检查 `Code` 是否等于 `TaskSuccessCode`。

## 5. 关键逻辑分析

1. **泛型支持**：`TaskResponse[T TaskData]` 使用泛型支持不同类型的响应数据。

2. **错误结构设计**：`TaskError` 包含序列化字段和非序列化字段，便于内部错误处理。

3. **任务状态管理**：通过 `Status` 字段管理任务生命周期。

4. **多平台支持**：`Platform` 字段支持不同任务平台的标识。

5. **配额管理**：`Quota` 字段记录任务消耗的配额。

## 6. 相关文件

- `dto/suno.go`：Suno 任务相关结构
- `model/task.go`：任务数据模型
- `controller/task.go`：任务控制器