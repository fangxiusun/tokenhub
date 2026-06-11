# suno.go 代码阅读文档

## 1. 全局摘要

该文件定义了 Suno 音乐生成 API 的请求和响应数据结构。Suno 是一个 AI 音乐生成平台，该文件支持与 Suno API 和 SunoGo API 的交互。

## 2. 依赖

- **标准库**：`encoding/json`

## 3. 类型定义

### Suno API 请求结构体

**SunoSubmitReq**：Suno 提交请求：
- `GptDescriptionPrompt` (string)：GPT 描述提示词
- `Prompt` (string)：提示词
- `Mv` (string)：音乐视频
- `Title` (string)：标题
- `Tags` (string)：标签
- `ContinueAt` (float64)：继续位置
- `TaskID` (string)：任务 ID
- `ContinueClipId` (string)：继续片段 ID
- `MakeInstrumental` (bool)：生成纯音乐

### Suno API 响应结构体

**SunoDataResponse**：Suno 数据响应：
- `TaskID` (string)：任务 ID
- `Action` (string)：任务类型（song、lyrics、description-mode）
- `Status` (string)：任务状态（submitted、queueing、processing、success、failed）
- `FailReason` (string)：失败原因
- `SubmitTime` (int64)：提交时间
- `StartTime` (int64)：开始时间
- `FinishTime` (int64)：完成时间
- `Data` (json.RawMessage)：原始数据

**SunoSong**：歌曲数据：
- `ID` (string)：歌曲 ID
- `VideoURL` (string)：视频 URL
- `AudioURL` (string)：音频 URL
- `ImageURL` (string)：图片 URL
- `ImageLargeURL` (string)：大图 URL
- `MajorModelVersion` (string)：主要模型版本
- `ModelName` (string)：模型名称
- `Status` (string)：状态
- `Title` (string)：标题
- `Text` (string)：歌词
- `Metadata` (SunoMetadata)：元数据

**SunoMetadata**：元数据：
- `Tags` (string)：标签
- `Prompt` (string)：提示词
- `GPTDescriptionPrompt` (interface{})：GPT 描述提示词
- `AudioPromptID` (interface{})：音频提示 ID
- `Duration` (interface{})：时长
- `ErrorType` (interface{})：错误类型
- `ErrorMessage` (interface{})：错误消息

**SunoLyrics**：歌词数据：
- `ID` (string)：歌词 ID
- `Status` (string)：状态
- `Title` (string)：标题
- `Text` (string)：歌词内容

### SunoGo API 请求结构体

**SunoGoAPISubmitReq**：SunoGo API 提交请求：
- `CustomMode` (bool)：自定义模式
- `Input` (SunoGoAPISubmitReqInput)：输入数据
- `NotifyHook` (string)：通知钩子

**SunoGoAPISubmitReqInput**：SunoGo API 输入数据：
- `GptDescriptionPrompt` (string)：GPT 描述提示词
- `Prompt` (string)：提示词
- `Mv` (string)：音乐视频
- `Title` (string)：标题
- `Tags` (string)：标签
- `ContinueAt` (float64)：继续位置
- `TaskID` (string)：任务 ID
- `ContinueClipId` (string)：继续片段 ID
- `MakeInstrumental` (bool)：生成纯音乐

### SunoGo API 响应结构体

**GoAPITaskResponse[T any]**：通用任务响应（泛型）：
- `Code` (int)：状态码
- `Message` (string)：消息
- `Data` (T)：响应数据
- `ErrorMessage` (string)：错误消息

**GoAPITaskResponseData**：任务响应数据：
- `TaskID` (string)：任务 ID

**GoAPIFetchResponseData**：获取响应数据：
- `TaskID` (string)：任务 ID
- `Status` (string)：状态
- `Input` (string)：输入
- `Clips` (map[string]SunoSong)：片段映射

## 4. 函数详情

无函数定义。

## 5. 关键逻辑分析

1. **双 API 支持**：同时支持 Suno 原生 API 和 SunoGo API。

2. **任务状态管理**：通过 `Status` 字段管理音乐生成任务的生命周期。

3. **多媒体输出**：支持音频、视频、图片等多种输出格式。

4. **泛型响应**：`GoAPITaskResponse[T any]` 使用泛型支持不同类型的响应数据。

5. **数据库兼容**：`SunoDataResponse` 包含 GORM 标签，支持数据库存储。

## 6. 相关文件

- `relay/suno/`：Suno 中继适配器
- `model/suno.go`：Suno 数据模型
- `controller/suno.go`：Suno 控制器