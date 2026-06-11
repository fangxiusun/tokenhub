# relay_mode.go 代码阅读文档

## 1. 全局总结

本文件定义了 Relay 模块的所有中继模式常量，以及 URL 路径到中继模式的映射函数。中继模式决定了请求的处理流程和使用的 handler。

## 2. 依赖关系

- `net/http`: HTTP 方法常量
- `strings`: 字符串匹配

## 3. 类型定义

### RelayMode 常量 (int)
- `RelayModeUnknown`: 未知模式
- `RelayModeChatCompletions`: Chat Completions
- `RelayModeCompletions`: Completions
- `RelayModeEmbeddings`: Embeddings
- `RelayModeModerations`: Moderations
- `RelayModeImagesGenerations`: 图像生成
- `RelayModeImagesEdits`: 图像编辑
- `RelayModeEdits`: 文本编辑
- `RelayModeMidjourney*`: Midjourney 相关模式（12种）
- `RelayModeAudioSpeech/AudioTranscription/AudioTranslation`: 音频模式
- `RelayModeSuno*`: Suno 相关模式
- `RelayModeVideo*`: 视频相关模式
- `RelayModeRerank`: Rerank
- `RelayModeResponses`: Responses API
- `RelayModeRealtime`: 实时语音
- `RelayModeGemini`: Gemini 原生格式
- `RelayModeResponsesCompact`: Responses 压缩模式

## 4. 函数详解

### `Path2RelayMode(path) int`
- 根据 URL 路径推断 RelayMode
- 支持 `/v1/` 前缀和 `/pg/` playground 前缀

### `Path2RelayModeMidjourney(path) int`
- Midjourney 路径到模式的映射
- 支持 submit/action/modal/shorten/imagine/blend/describe/notify/change/fetch 等

### `Path2RelaySuno(method, path) int`
- Suno 路径到模式的映射
- 根据 HTTP 方法和路径后缀区分 fetch/submit

## 5. 关键逻辑分析

1. **路径匹配**: 使用 HasPrefix/HasSuffix 进行路径匹配
2. **优先级**: ResponsesCompact 在 Responses 之前匹配（更具体的路径优先）
3. **Midjourney 丰富的操作**: 支持 12+ 种 Midjourney 操作模式
4. **Suno 方法区分**: 同一路径根据 GET/POST 区分不同操作

## 6. 关联文件

- `relay/common/relay_info.go`: RelayMode 存储在 RelayInfo 中
- `relay/relay_adaptor.go`: 根据模式选择 handler
