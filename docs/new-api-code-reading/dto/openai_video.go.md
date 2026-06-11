# openai_video.go 代码阅读文档

## 1. 全局摘要

该文件定义了 OpenAI 视频生成 API 的数据结构，包括视频对象 `OpenAIVideo`、视频错误 `OpenAIVideoError`，以及视频状态常量。

## 2. 依赖

- **标准库**：
  - `strconv`：字符串转换
  - `strings`：字符串操作

## 3. 类型定义

### 视频状态常量
- `VideoStatusUnknown`："unknown" - 未知状态
- `VideoStatusQueued`："queued" - 排队中
- `VideoStatusInProgress`："in_progress" - 进行中
- `VideoStatusCompleted`："completed" - 已完成
- `VideoStatusFailed`："failed" - 失败

### OpenAIVideo 结构体
视频对象结构：
- `ID` (string)：视频 ID
- `TaskID` (string)：任务 ID（兼容旧接口，待废弃）
- `Object` (string)：对象类型
- `Model` (string)：模型名称
- `Status` (string)：状态（使用视频状态常量）
- `Progress` (int)：进度百分比
- `CreatedAt` (int64)：创建时间戳
- `CompletedAt` (int64)：完成时间戳
- `ExpiresAt` (int64)：过期时间戳
- `Seconds` (string)：视频时长
- `Size` (string)：视频尺寸
- `RemixedFromVideoID` (string)：混音来源视频 ID
- `Error` (*OpenAIVideoError)：错误信息
- `Metadata` (map[string]any)：元数据

### OpenAIVideoError 结构体
视频错误结构：
- `Message` (string)：错误消息
- `Code` (string)：错误代码

## 4. 函数详情

### SetProgressStr()
```go
func (m *OpenAIVideo) SetProgressStr(progress string)
```
**功能**：设置进度百分比（从字符串解析）。

**逻辑**：移除百分号后解析为整数。

### SetMetadata()
```go
func (m *OpenAIVideo) SetMetadata(k string, v any)
```
**功能**：设置元数据键值对。

**逻辑**：初始化 `Metadata` map（如果为 nil）后设置值。

### NewOpenAIVideo()
```go
func NewOpenAIVideo() *OpenAIVideo
```
**功能**：创建新的视频对象实例。

**返回**：初始化为默认状态的视频对象。

## 5. 关键逻辑分析

1. **状态管理**：使用字符串常量定义视频状态，便于状态比较和转换。

2. **进度解析**：支持从字符串解析进度百分比（如 "50%"）。

3. **元数据管理**：延迟初始化 `Metadata` map，避免不必要的内存分配。

4. **兼容性设计**：`TaskID` 字段用于兼容旧接口，标记为待废弃。

## 6. 相关文件

- `dto/openai_response.go`：`OpenAIVideoResponse` 结构定义
- `relay/video/`：视频中继适配器
- `controller/video.go`：视频控制器