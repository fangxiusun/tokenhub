# policy.go 代码阅读文档

## 1. 全局总结

本文件定义了 Chat Completions 到 Responses API 转换的策略判断逻辑。通过读取全局配置中的 `ChatCompletionsToResponsesPolicy`，判断指定渠道+模型是否需要执行格式转换。文件仅包含两个简短的决策函数，属于策略模式的轻量实现。

## 2. 依赖关系

| 包 | 用途 |
|---|---|
| `github.com/QuantumNous/new-api/setting/model_setting` | 获取全局设置和策略类型定义 |

内部依赖 `regex.go` 中的 `matchAnyRegex` 函数进行模型名正则匹配。

## 3. 类型定义

本文件未定义新类型。

## 4. 函数详解

### `ShouldChatCompletionsUseResponsesPolicy(policy, channelID, channelType, model) bool`

**签名**:
```go
func ShouldChatCompletionsUseResponsesPolicy(
    policy model_setting.ChatCompletionsToResponsesPolicy,
    channelID int, channelType int, model string,
) bool
```

**功能**: 判断给定策略下，指定渠道和模型是否应使用 Responses API。

**逻辑**:
1. 调用 `policy.IsChannelEnabled(channelID, channelType)` 检查该渠道是否在策略的启用列表中
2. 若渠道启用，调用 `matchAnyRegex(policy.ModelPatterns, model)` 检查模型名是否匹配任一模式
3. 两个条件都满足才返回 `true`

### `ShouldChatCompletionsUseResponsesGlobal(channelID, channelType, model) bool`

**签名**:
```go
func ShouldChatCompletionsUseResponsesGlobal(
    channelID int, channelType int, model string,
) bool
```

**功能**: 便捷函数，直接使用全局设置中的策略进行判断。

**逻辑**: 调用 `model_setting.GetGlobalSettings()` 获取全局配置，提取 `ChatCompletionsToResponsesPolicy` 后委托给 `ShouldChatCompletionsUseResponsesPolicy`。

## 5. 关键逻辑分析

1. **两层过滤**: 先按渠道过滤（`IsChannelEnabled`），再按模型名过滤（正则匹配），实现精确控制哪些渠道的哪些模型需要转换。
2. **全局策略**: 通过 `model_setting.GetGlobalSettings()` 获取统一策略，确保所有请求使用同一套规则，避免配置分散。

## 6. 关联文件

- `new-api/service/openaicompat/regex.go` — 提供 `matchAnyRegex` 正则匹配函数
- `new-api/setting/model_setting/` — `ChatCompletionsToResponsesPolicy` 类型定义和 `GetGlobalSettings()` 函数
- `new-api/service/openaicompat/chat_to_responses.go` — 实际的格式转换逻辑
