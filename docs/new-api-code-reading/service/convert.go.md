# convert.go 代码阅读文档

## 1. 全局总结

该文件实现了不同 AI API 格式之间的请求和响应转换，包括：
- Claude → OpenAI 请求转换
- OpenAI → Claude 响应转换（流式/非流式）
- Gemini → OpenAI 请求转换
- OpenAI → Gemini 响应转换（流式/非流式）

是多格式 API 网关的核心转换层。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | JSON 序列化、类型转换 |
| `dto` | 请求/响应数据结构 |
| `openrouter` | OpenRouter 特定格式 |
| `relaycommon` | RelayInfo、ClaudeConvertInfo |
| `reasonmap` | 停止原因映射 |
| `lo` | 指针工具函数 |

## 3. 类型定义

无自定义类型，使用 `dto` 包中的请求/响应结构体。

## 4. 函数详解

### `ClaudeToOpenAIRequest(claudeRequest, info) (*dto.GeneralOpenAIRequest, error)`
- 转换 Claude 消息格式到 OpenAI 格式
- 处理：system 消息、工具调用、图片、thinking 模式
- 支持 OpenRouter 特殊处理

### `StreamResponseOpenAI2Claude(openAIResponse, info) []*dto.ClaudeResponse`
- 流式 OpenAI 响应转 Claude 格式
- 处理：message_start、content_block_start/delta/stop、message_delta/stop
- 状态机管理：跟踪当前内容块类型和索引

### `ResponseOpenAI2Claude(openAIResponse, info) *dto.ClaudeResponse`
- 非流式 OpenAI 响应转 Claude 格式

### `GeminiToOpenAIRequest(geminiRequest, info) (*dto.GeneralOpenAIRequest, error)`
- 转换 Gemini 请求到 OpenAI 格式
- 处理：Contents、SystemInstructions、Tools、GenerationConfig

### `ResponseOpenAI2Gemini` / `StreamResponseOpenAI2Gemini`
- OpenAI 响应转 Gemini 格式

## 5. 关键逻辑分析

1. **状态机设计**：流式转换使用 ClaudeConvertInfo 跟踪状态
2. **内容块索引**：thinking/text/tools 使用不同的索引管理
3. **工具调用处理**：支持并行工具调用的索引偏移
4. **缓存 token 处理**：OpenRouter Claude 特殊的缓存 token 计算
5. **停止原因映射**：OpenAI finish_reason → Claude stop_reason

## 6. 关联文件

- `dto/` — 请求/响应数据结构
- `relay/common/relay_info.go` — ClaudeConvertInfo
- `relay/reasonmap` — 停止原因映射
