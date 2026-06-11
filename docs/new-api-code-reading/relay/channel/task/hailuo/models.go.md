# models.go 代码阅读文档

## 1. 全局总结

该文件定义了 MiniMax 海螺视频生成 API 的所有数据模型（请求/响应结构体）以及模型配置信息。包含视频生成请求、提交响应、任务查询响应、错误信息、文件检索响应等结构体，以及每个支持模型的默认配置（分辨率、时长、功能支持等）。

## 2. 依赖关系

无外部依赖，纯结构体和配置定义文件。

## 3. 类型定义

### 请求结构体

| 类型名 | 说明 |
|--------|------|
| `SubjectReference` | 主体参考，包含类型（"character"）和参考图片数组 |
| `VideoRequest` | 视频生成请求，包含模型、提示词、时长、分辨率、回调 URL、水印、首尾帧图片、主体参考等 |
| `QueryTaskRequest` | 任务查询请求，仅包含 task_id |

### 响应结构体

| 类型名 | 说明 |
|--------|------|
| `VideoResponse` | 视频生成提交响应，包含 task_id 和基础响应 |
| `BaseResp` | 基础响应，包含状态码和状态消息 |
| `QueryTaskResponse` | 任务查询响应，包含 task_id、状态、file_id、视频尺寸、基础响应 |
| `ErrorInfo` | 错误信息 |
| `TaskStatusInfo` | 任务状态信息 |
| `RetrieveFileResponse` | 文件检索响应，包含文件对象和基础响应 |
| `FileObject` | 文件对象，包含 file_id、大小、创建时间、文件名、用途、下载 URL |

### 配置结构体

| 类型名 | 说明 |
|--------|------|
| `ModelConfig` | 模型配置，包含名称、默认分辨率、支持的时长列表、支持的分辨率列表、功能开关 |

### VideoRequest 字段说明

| 字段名 | 类型 | 说明 |
|--------|------|------|
| `Model` | `string` | 模型名称 |
| `Prompt` | `string` | 文本提示词 |
| `PromptOptimizer` | `*bool` | 是否启用提示词优化 |
| `FastPretreatment` | `*bool` | 是否启用快速预处理 |
| `Duration` | `*int` | 视频时长 |
| `Resolution` | `string` | 分辨率 |
| `CallbackURL` | `string` | 回调 URL |
| `AigcWatermark` | `*bool` | 是否添加 AIGC 水印 |
| `FirstFrameImage` | `string` | 首帧图片（图生视频/首尾帧生视频） |
| `LastFrameImage` | `string` | 尾帧图片（首尾帧生视频） |
| `SubjectReference` | `[]SubjectReference` | 主体参考（主体参考生视频） |

## 4. 函数详解

| 函数签名 | 说明 |
|----------|------|
| `GetModelConfig(model string) ModelConfig` | 根据模型名获取配置，不存在则返回默认配置 |

## 5. 关键逻辑分析

### 模型配置详情

| 模型 | 默认分辨率 | 支持时长 | 支持分辨率 | 提示词优化 | 快速预处理 |
|------|-----------|---------|-----------|-----------|-----------|
| MiniMax-Hailuo-2.3 | 768P | 6, 10 | 768P, 1080P | ✓ | ✓ |
| MiniMax-Hailuo-2.3-Fast | 768P | 6, 10 | 768P, 1080P | ✓ | ✓ |
| MiniMax-Hailuo-02 | 768P | 6, 10 | 512P, 768P, 1080P | ✓ | ✓ |
| T2V-01-Director | 768P | 6 | 768P, 1080P | ✓ | ✗ |
| T2V-01 | 720P | 6 | 720P | ✓ | ✗ |
| I2V-01-Director | 720P | 6 | 720P, 1080P | ✓ | ✗ |
| I2V-01-live | 720P | 6 | 720P, 1080P | ✓ | ✗ |
| I2V-01 | 720P | 6 | 720P, 1080P | ✓ | ✗ |
| S2V-01 | 720P | 6 | 720P | ✓ | ✗ |

### 默认配置回退
当模型名不在预定义配置中时，`GetModelConfig` 返回一个使用全局默认值（720P、6秒）的配置，并默认启用提示词优化、禁用快速预处理。

## 6. 关联文件

- `relay/channel/task/hailuo/adaptor.go` — 使用这些数据结构和 `GetModelConfig`
- `relay/channel/task/hailuo/constants.go` — 分辨率常量定义
