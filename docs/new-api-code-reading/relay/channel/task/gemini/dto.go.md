# dto.go 代码阅读文档

## 1. 全局总结

该文件定义了 Gemini Veo 视频生成 API 的请求和响应数据传输对象（DTO）。这些结构体被 Gemini 和 Vertex 两个适配器共用，描述了 Veo `predictLongRunning` 端点的请求体、实例、参数、以及操作响应格式。

## 2. 依赖关系

无外部依赖，纯结构体定义文件。

## 3. 类型定义

### 请求结构体

| 类型名 | 说明 |
|--------|------|
| `VeoImageInput` | 图片输入，包含 base64 编码数据和 MIME 类型 |
| `VeoInstance` | 单个 Veo 请求实例，包含 prompt 和可选的 image |
| `VeoParameters` | Veo 参数块，包含采样数、时长、宽高比、分辨率、反向提示词、人物生成、存储 URI、压缩质量、调整模式、随机种子、是否生成音频 |
| `VeoRequestPayload` | 顶层请求体，包含 instances 数组和 parameters |

### 响应结构体

| 类型名 | 说明 |
|--------|------|
| `submitResponse` | 提交响应，仅包含操作名称（name） |
| `operationVideo` | 操作视频，包含 MIME 类型、base64 数据、编码格式 |
| `operationResponse` | 操作响应，包含名称、完成状态、响应数据（视频列表/生成视频响应）、错误信息 |

### VeoParameters 字段说明

| 字段名 | 类型 | 说明 |
|--------|------|------|
| `SampleCount` | `int` | 采样数量 |
| `DurationSeconds` | `int` | 视频时长（秒） |
| `AspectRatio` | `string` | 宽高比（如 "16:9"） |
| `Resolution` | `string` | 分辨率（如 "720p"、"1080p"、"4k"） |
| `NegativePrompt` | `string` | 反向提示词 |
| `PersonGeneration` | `string` | 人物生成控制 |
| `StorageUri` | `string` | 输出存储 URI |
| `CompressionQuality` | `string` | 压缩质量 |
| `ResizeMode` | `string` | 调整模式 |
| `Seed` | `*int` | 随机种子（指针类型以区分零值和未设置） |
| `GenerateAudio` | `*bool` | 是否生成音频（指针类型） |

## 4. 函数详解

无函数定义。

## 5. 关键逻辑分析

### 指针类型使用
`Seed` 和 `GenerateAudio` 使用指针类型（`*int`、`*bool`），符合项目 Rule 6 的要求：可选标量字段必须使用指针类型以区分"未设置"和"设置为零值"。

### operationResponse 嵌套结构
响应结构支持多种视频返回格式：
- `Videos[]` — 内联视频数据（base64 编码）
- `GenerateVideoResponse.GeneratedVideos[].Video.URI` — 远程视频 URI
- `Video` — 单个视频字符串

实际使用中主要通过 `GenerateVideoResponse` 获取视频 URI。

### Vertex 共用
这些 DTO 同时被 Gemini 和 Vertex 两个适配器使用（通过不同的 baseURL 区分），体现了代码复用设计。

## 6. 关联文件

- `relay/channel/task/gemini/adaptor.go` — 使用这些 DTO 构建请求和解析响应
- `relay/channel/task/vertex/` — Vertex 适配器也使用这些 DTO
