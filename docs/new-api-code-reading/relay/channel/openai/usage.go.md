# usage.go 代码阅读文档

## 1. 全局总结
该文件负责 OpenAI 渠道的 usage（使用量）后处理，特别是从不同渠道的非标准响应位置提取 cached_tokens（缓存 token）信息。支持 DeepSeek、智谱（Zhipu）、Moonshot 和 OpenAI（llama.cpp）四种渠道的特殊 cached_tokens 提取逻辑。

## 2. 依赖关系
- `github.com/QuantumNous/new-api/common` — 通用工具函数
- `github.com/QuantumNous/new-api/constant` — 渠道类型常量
- `github.com/QuantumNous/new-api/dto` — 数据传输对象
- `github.com/QuantumNous/new-api/relay/common` — 中继公共信息

## 3. 类型定义
无自定义类型定义。

## 4. 函数详解

### `applyUsagePostProcessing(info *relaycommon.RelayInfo, usage *dto.Usage, responseBody []byte)`
- **作用**：根据渠道类型进行 usage 后处理
- **参数**：中继信息、使用量、原始响应体
- **逻辑**：根据 `ChannelType` 分发到不同的处理分支：
  - **DeepSeek**：从 `prompt_cache_hit_tokens` 提取
  - **智谱 v4**：依次尝试 `input_tokens_details.cached_tokens` → 响应体标准位置 → `prompt_cache_hit_tokens`
  - **Moonshot**：依次尝试 `input_tokens_details.cached_tokens` → choices 中的 `cached_tokens` → 标准位置 → `prompt_cache_hit_tokens`
  - **OpenAI（llama.cpp）**：从 `timings.cache_n` 提取

### `extractCachedTokensFromBody(body []byte) (int, bool)`
- **作用**：从响应体标准位置提取 cached_tokens
- **提取位置**：
  1. `usage.prompt_tokens_details.cached_tokens`
  2. `usage.cached_tokens`
  3. `usage.prompt_cache_hit_tokens`

### `extractMoonshotCachedTokensFromBody(body []byte) (int, bool)`
- **作用**：从 Moonshot 的非标准位置提取 cached_tokens
- **提取位置**：`choices[].usage.cached_tokens`

### `extractLlamaCachedTokensFromBody(body []byte) (int, bool)`
- **作用**：从 llama.cpp 的非标准位置提取 cached_tokens
- **提取位置**：`timings.cache_n`

## 5. 关键逻辑分析
- **多渠道兼容**：不同渠道的 cached_tokens 位置各不相同，此文件通过渠道类型分发和多级 fallback 确保所有渠道的缓存 token 都能被正确提取
- **Fallback 链**：每个渠道都有多个提取位置的 fallback 逻辑，确保即使响应格式变化也能提取到 cached_tokens
- **指针类型处理**：使用 `*int` 指针类型区分"字段不存在"和"值为 0"的情况
- **严格 JSON 解析**：使用局部 struct 只解析需要的字段，避免完整反序列化的开销
- **DeepSeek 简化**：DeepSeek 直接从 `prompt_cache_hit_tokens` 提取，无需多级 fallback

## 6. 关联文件
- `relay/channel/openai/relay-openai.go` — OaiStreamHandler 和 OpenaiHandler 中调用
- `relay/channel/openai/relay_image.go` — 图像处理器中调用
- `relay/channel/openai/relay_responses.go` — Responses 处理器中使用
- `dto/usage.go` — Usage 结构体定义
- `constant/channel.go` — ChannelType 常量定义
