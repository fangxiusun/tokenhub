# global.go 代码阅读文档

## 1. 全局总结

该文件定义全局模型设置，包括透传请求开关、Thinking 模型黑名单、Chat Completions 到 Responses API 的转换策略。

## 2. 依赖关系

- `slices` — 切片操作
- `strings` — 字符串处理
- `github.com/QuantumNous/new-api/setting/config` — 全局配置管理器

## 3. 类型定义

| 结构体 | 字段 | 类型 | 说明 |
|--------|------|------|------|
| `ChatCompletionsToResponsesPolicy` | `Enabled` | `bool` | 是否启用转换 |
| | `AllChannels` | `bool` | 是否对所有渠道生效 |
| | `ChannelIDs` | `[]int` | 指定渠道 ID 列表 |
| | `ChannelTypes` | `[]int` | 指定渠道类型列表 |
| | `ModelPatterns` | `[]string` | 模型匹配模式 |
| `GlobalSettings` | `PassThroughRequestEnabled` | `bool` | 是否启用透传请求 |
| | `ThinkingModelBlacklist` | `[]string` | Thinking 模型黑名单 |
| | `ChatCompletionsToResponsesPolicy` | `ChatCompletionsToResponsesPolicy` | 转换策略 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `IsChannelEnabled` | `func (p ChatCompletionsToResponsesPolicy) IsChannelEnabled(channelID int, channelType int) bool` | 检查渠道是否启用转换 |
| `GetGlobalSettings` | `func GetGlobalSettings() *GlobalSettings` | 获取全局设置 |
| `ShouldPreserveThinkingSuffix` | `func ShouldPreserveThinkingSuffix(modelName string) bool` | 判断模型是否应保留 Thinking 后缀 |

## 5. 关键逻辑分析

- `IsChannelEnabled` 优先级：`AllChannels` > `ChannelIDs` > `ChannelTypes`
- 黑名单默认包含 `moonshotai/kimi-k2-thinking` 和 `kimi-k2-thinking`
- `ShouldPreserveThinkingSuffix` 使用精确匹配

## 6. 关联文件

- `relay/handler.go` — 使用转换策略
- `middleware/distributor.go` — 使用透传请求开关
