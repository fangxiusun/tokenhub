# constant.go 代码阅读文档

## 1. 全局总结
本文件定义了 Gemini 渠道的常量，包括支持的模型列表、安全设置类别列表和渠道名称。这些常量被适配器和其他 Gemini 处理函数使用。

## 2. 依赖关系
无外部依赖，仅定义包级变量。

## 3. 类型定义

### `ModelList` 变量
```go
var ModelList = []string{...}
```
Gemini 渠道支持的所有模型名称列表，包含：
- **稳定版本**: `gemini-2.5-flash`, `gemini-2.5-pro`, `gemini-2.0-flash` 等
- **最新版本**: `gemini-flash-latest`, `gemini-pro-latest` 等
- **预览版本**: `gemini-2.5-flash-preview-tts`, `gemini-3-pro-preview`, `nano-banana-pro-preview` 等
- **Gemma 模型**: `gemma-3-1b-it` 到 `gemma-3n-e2b-it`
- **嵌入模型**: `gemini-embedding-001`, `gemini-embedding-2-preview`
- **Imagen 模型**: `imagen-4.0-generate-001` 等
- **Veo 模型**: `veo-2.0-generate-001` 到 `veo-3.1-fast-generate-preview`
- **其他**: `aqa`

### `SafetySettingList` 变量
```go
var SafetySettingList = []string{...}
```
Gemini API 的安全设置类别列表，用于构建 `safetySettings` 参数：
- `HARM_CATEGORY_HARASSMENT` — 骚扰内容
- `HARM_CATEGORY_HATE_SPEECH` — 仇恨言论
- `HARM_CATEGORY_SEXUALLY_EXPLICIT` — 色情内容
- `HARM_CATEGORY_DANGEROUS_CONTENT` — 危险内容

注：`HARM_CATEGORY_CIVIC_INTEGRITY` 已被弃用，未包含在列表中。

### `ChannelName` 变量
```go
var ChannelName = "google gemini"
```
渠道标识名称。

## 4. 函数详解
无函数定义。

## 5. 关键逻辑分析

### 模型列表组织
模型列表按类别分组，覆盖了 Gemini 生态系统的多种模型类型：文本生成、多模态、嵌入、图片生成（Imagen）、视频生成（Veo）、开源（Gemma）等。该列表直接用于渠道注册和模型匹配。

### 安全设置
安全设置列表在 `relay-gemini.go` 中被遍历，每个类别从配置中获取阈值，构建 `safetySettings` 数组传递给 Gemini API。

## 6. 关联文件
- `adaptor.go` — 使用 `ModelList` 和 `ChannelName`
- `relay-gemini.go` — 使用 `SafetySettingList` 构建安全设置
