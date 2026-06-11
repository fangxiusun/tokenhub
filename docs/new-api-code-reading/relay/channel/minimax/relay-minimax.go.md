# relay-minimax.go 代码阅读文档

## 1. 全局总结
本文件实现了 MiniMax 渠道的请求 URL 构建逻辑，根据中继模式和响应格式选择不同的 API 端点。

## 2. 依赖关系
- **标准库**: `fmt`
- **项目内部**:
  - `github.com/QuantumNous/new-api/constant` — 渠道基础 URL 常量
  - `github.com/QuantumNous/new-api/relay/common` — RelayInfo
  - `github.com/QuantumNous/new-api/relay/constant` — 中继模式常量
  - `github.com/QuantumNous/new-api/types` — 响应格式常量

## 3. 类型定义
无自定义类型。

## 4. 函数详解

### `GetRequestURL(info) (string, error)`
根据响应格式和中继模式构建请求 URL：

**Claude 格式**:
- `{baseUrl}/anthropic/v1/messages` — 使用 MiniMax 的 Anthropic 兼容端点

**默认格式**（按中继模式）:
- `RelayModeChatCompletions` → `{baseUrl}/v1/text/chatcompletion_v2`
- `RelayModeImagesGenerations` → `{baseUrl}/v1/image_generation`
- `RelayModeAudioSpeech` → `{baseUrl}/v1/t2a_v2`
- 其他模式返回不支持的错误

如果 `ChannelBaseUrl` 为空，使用默认的 MiniMax 渠道基础 URL。

## 5. 关键逻辑分析

### 多端点路由
MiniMax 根据不同能力使用不同的 API 端点：文本生成使用 `chatcompletion_v2`，图像生成使用 `image_generation`，语音合成使用 `t2a_v2`。Claude 格式请求走专用的 Anthropic 兼容端点。

### 默认 URL 降级
当 `ChannelBaseUrl` 为空时，从 `channelconstant.ChannelBaseURLs` 获取默认值，提供了灵活的配置降级机制。

## 6. 关联文件
- `adaptor.go` — 调用 `GetRequestURL` 获取请求 URL
- `constant/constants.go` — 渠道基础 URL 定义
