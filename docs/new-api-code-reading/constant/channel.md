# channel.go 代码阅读文档

## 1. 全局概述

本文件是渠道（Channel）系统的核心常量定义文件，包含：
- 所有支持的 AI 服务商渠道类型常量
- 各渠道的默认 Base URL
- 渠道类型名称映射表
- 特殊渠道的 Base URL 配置（如编码计划）
- 渠道名称查询函数

## 2. 依赖关系

本文件无外部依赖。

## 3. 类型定义

### ChannelSpecialBase 结构体

```go
type ChannelSpecialBase struct {
    ClaudeBaseURL string
    OpenAIBaseURL string
}
```

用于存储特殊渠道的双协议 Base URL（Claude 格式和 OpenAI 格式）。

## 4. 函数详情

### GetChannelTypeName

```go
func GetChannelTypeName(channelType int) string
```

根据渠道类型 ID 返回对应的渠道名称。如果找不到对应的名称，返回 "Unknown"。

## 5. 关键逻辑分析

### 渠道类型常量

定义了 58 种渠道类型（从 `ChannelTypeUnknown = 0` 到 `ChannelTypeCodex = 57`），涵盖：

| 类别 | 渠道 |
|------|------|
| 国际主流 | OpenAI, Anthropic, Google Gemini, AWS Bedrock, Mistral, Cohere |
| 国内厂商 | 百度, 智谱, 阿里, 讯飞, 360, 腾讯, Moonshot, DeepSeek, 火山引擎 |
| 中转/聚合 | OpenRouter, OhMyGPT, API2GPT, AIGC2D, AIProxy |
| 特殊能力 | Midjourney, SunoAPI, Dify, Kling, Jimeng, Vidu, Sora |
| 本地部署 | Ollama, Xinference |

### ChannelBaseURLs 切片

与渠道类型常量按索引一一对应的默认 Base URL 切片。部分渠道（如 OpenAI 的特殊变体）使用空字符串表示需要自定义配置。

### ChannelTypeNames 映射

`map[int]string` 类型，将渠道类型 ID 映射为可读名称，用于前端展示和日志输出。

### ChannelSpecialBases 映射

存储需要双协议支持的特殊渠道配置：
- `glm-coding-plan` / `glm-coding-plan-international` — 智谱编码计划
- `kimi-coding-plan` — Kimi 编码计划
- `doubao-coding-plan` — 豆包编码计划

每个条目包含 Claude 格式和 OpenAI 格式的不同 Base URL。

### ChannelTypeDummy 哨兵值

与 `APITypeDummy` 类似，`ChannelTypeDummy` 仅用于统计渠道总数，注释明确指示不要在其后添加新渠道。

## 6. 相关文件

- `constant/api_type.go` — API 类型常量，与渠道类型一一对应
- `model/channel.go` — 渠道数据模型
- `relay/channel/` — 各渠道适配器实现
- `middleware/distributor.go` — 渠道分发逻辑
