# video.go 代码阅读文档

## 1. 全局摘要

该文件定义了视频生成 API 的请求和响应数据结构。支持视频生成任务的提交和状态查询，包含任务状态、元数据和错误信息。

## 2. 依赖

无外部包依赖。

## 3. 类型定义

### VideoRequest 结构体
视频生成请求：
- `Model` (string)：模型名称
- `Prompt` (string)：文本提示词
- `Image` (string)：图像输入（URL/Base64）
- `Duration` (float64)：视频时长（秒）
- `Width` (int)：视频宽度
- `Height` (int)：视频高度
- `Fps` (int)：帧率
- `Seed` (int)：随机种子
- `N` (int)：生成数量
- `ResponseFormat` (string)：响应格式
- `User` (string)：用户标识
- `Metadata` (map[string]any)：厂商特定参数

### VideoResponse 结构体
视频生成提交响应：
- `TaskId` (string)：任务 ID
- `Status` (string)：任务状态

### VideoTaskResponse 结构体
视频任务状态查询响应：
- `TaskId` (string)：任务 ID
- `Status` (string)：任务状态
- `Url` (string)：视频资源 URL
- `Format` (string)：视频格式
- `Metadata` (*VideoTaskMetadata)：结果元数据
- `Error` (*VideoTaskError)：错误信息

### VideoTaskMetadata 结构体
视频任务元数据：
- `Duration` (float64)：实际视频时长
- `Fps` (int)：实际帧率
- `Width` (int)：实际宽度
- `Height` (int)：实际高度
- `Seed` (int)：使用的随机种子

### VideoTaskError 结构体
视频任务错误信息：
- `Code` (int)：错误代码
- `Message` (string)：错误消息

## 4. 函数详情

无函数定义。

## 5. 关键逻辑分析

1. **异步任务模式**：视频生成采用异步任务模式，提交后返回任务 ID，通过查询获取结果。

2. **元数据支持**：`Metadata` 字段支持厂商特定参数（如 negative_prompt、style 等）。

3. **状态管理**：通过 `Status` 字段管理任务生命周期。

4. **错误处理**：任务失败时通过 `Error` 字段返回详细错误信息。

## 6. 相关文件

- `relay/video/`：视频中继适配器
- `controller/video.go`：视频控制器
- `model/video.go`：视频数据模型