# audio_handler.go 代码阅读文档

## 1. 全局总结

本文件实现了音频请求（TTS/语音转写/翻译）的处理入口 `AudioHelper`。处理流程简洁：模型映射 → 请求转换 → 发送 → 响应解析 → 计费。支持音频 token 的特殊计费路径。

## 2. 依赖关系

- `relay/common`: RelayInfo
- `relay/helper`: 模型映射
- `service`: 计费

## 3. 类型定义

本文件无自定义类型定义。

## 4. 函数详解

### `AudioHelper(c *gin.Context, info *relaycommon.RelayInfo) *types.NewAPIError`
- **功能**: 处理音频请求（TTS、Whisper 转写、翻译）
- **流程**: InitChannelMeta → 类型断言 → 深拷贝 → 模型映射 → 获取适配器 → ConvertAudioRequest → DoRequest → DoResponse → 计费
- **计费分支**: 检测 AudioTokens，有则使用 `PostAudioConsumeQuota`，否则使用 `PostTextConsumeQuota`

## 5. 关键逻辑分析

1. **音频计费**: 通过 `CompletionTokenDetails.AudioTokens` 和 `PromptTokensDetails.AudioTokens` 判断是否为音频请求
2. **简洁流程**: 相比 TextHelper，没有 passthrough、StreamOptions 等复杂逻辑
3. **请求转换**: 使用 `ConvertAudioRequest` 返回 `io.Reader`，支持流式数据

## 6. 关联文件

- `relay/channel/adapter.go`: Adaptor.ConvertAudioRequest 接口
- `dto/audio.go`: AudioRequest DTO
