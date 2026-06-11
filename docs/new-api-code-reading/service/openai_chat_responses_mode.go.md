# openai_chat_responses_mode.go 代码阅读文档

## 1. 全局总结

该文件提供 Chat Completions 到 Responses API 的路由策略判断。决定是否将 Chat Completions 请求转换为 Responses 格式。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `openaicompat` | 策略实现 |
| `model_setting` | 策略配置 |

## 3. 函数详解

### `ShouldChatCompletionsUseResponsesPolicy(policy, channelID, channelType, model) bool`
根据策略配置判断是否使用 Responses 格式

### `ShouldChatCompletionsUseResponsesGlobal(channelID, channelType, model) bool`
全局策略判断

## 4. 关联文件

- `setting/model_setting` — 策略配置
- `service/openaicompat/` — 实现
