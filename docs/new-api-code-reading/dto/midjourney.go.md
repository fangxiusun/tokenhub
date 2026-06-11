# midjourney.go 代码阅读文档

## 1. 全局摘要

该文件定义了 Midjourney 图像生成 API 的请求和响应数据结构。包含 Midjourney 请求/响应结构体、任务 DTO、状态结构体、按钮定义等。主要用于与 Midjourney 代理服务进行交互，处理图像生成、变体、放大等操作。

## 2. 依赖

无外部包依赖，仅使用 Go 标准类型。

## 3. 类型定义

### 请求结构体

**SwapFaceRequest**：换脸请求：
- `SourceBase64`：源图像 Base64 数据
- `TargetBase64`：目标图像 Base64 数据

**MidjourneyRequest**：Midjourney 请求：
- `Prompt`：提示词
- `CustomId`：自定义 ID
- `BotType`：机器人类型
- `NotifyHook`：通知钩子
- `Action`：操作类型
- `Index`：索引
- `State`：状态
- `TaskId`：任务 ID
- `Base64Array`：Base64 图像数组
- `Content`：内容
- `MaskBase64`：遮罩 Base64 数据

### 响应结构体

**MidjourneyResponse**：Midjourney 响应：
- `Code`：状态码
- `Description`：描述信息
- `Properties`：属性数据
- `Result`：结果数据

**MidjourneyUploadResponse**：上传响应：
- `Code`：状态码
- `Description`：描述信息
- `Result`：结果 URL 数组

**MidjourneyResponseWithStatusCode**：带状态码的响应：
- `StatusCode`：HTTP 状态码
- `Response`：响应内容

### 任务 DTO 结构体

**MidjourneyDto**：Midjourney 任务数据传输对象：
- `MjId`：Midjourney 任务 ID
- `Action`：操作类型
- `CustomId`：自定义 ID
- `BotType`：机器人类型
- `Prompt`/`PromptEn`：提示词（中文/英文）
- `Description`：描述
- `State`：状态
- `SubmitTime`/`StartTime`/`FinishTime`：时间戳
- `ImageUrl`：图像 URL
- `VideoUrl`：视频 URL
- `VideoUrls`：视频 URL 数组
- `Status`：状态
- `Progress`：进度
- `FailReason`：失败原因
- `Buttons`：按钮数组
- `MaskBase64`：遮罩数据
- `Properties`：属性数据

**ImgUrls**：URL 结构：
- `Url`：URL 地址

### 状态相关结构体

**MidjourneyStatus**：状态结构：
- `Status`：状态码

**MidjourneyWithoutStatus**：无状态的任务结构：
- `Id`：数据库 ID
- `Code`：状态码
- `UserId`：用户 ID
- `Action`：操作类型
- `MjId`：Midjourney 任务 ID
- `Prompt`/`PromptEn`：提示词
- `Description`：描述
- `State`：状态
- `SubmitTime`/`StartTime`/`FinishTime`：时间戳
- `ImageUrl`：图像 URL
- `Progress`：进度
- `FailReason`：失败原因
- `ChannelId`：渠道 ID

### 按钮和属性结构体

**ActionButton**：操作按钮：
- `CustomId`：自定义 ID
- `Emoji`：表情符号
- `Label`：标签
- `Type`：类型
- `Style`：样式

**Properties**：属性数据：
- `FinalPrompt`：最终提示词
- `FinalZhPrompt`：最终中文提示词

## 4. 函数详情

该文件主要为数据结构定义，无复杂函数实现。

## 5. 关键逻辑分析

1. **任务状态管理**：通过 `Status`、`Progress`、`FailReason` 字段管理任务生命周期。

2. **多语言支持**：`Prompt` 和 `PromptEn` 字段支持中英文提示词。

3. **多媒体支持**：支持图像（`ImageUrl`）和视频（`VideoUrl`、`VideoUrls`）输出。

4. **数据库兼容**：`MidjourneyWithoutStatus` 结构体包含 GORM 标签，支持数据库存储。

5. **按钮交互**：`Buttons` 字段支持动态按钮定义，用于用户交互操作。

## 6. 相关文件

- `relay/midjourney/`：Midjourney 中继适配器
- `model/midjourney.go`：Midjourney 数据模型
- `controller/midjourney.go`：Midjourney 控制器
- `types/midjourney.go`：Midjourney 类型定义