# adaptor.go 代码阅读文档

## 1. 全局总结
智谱 AI v4 版渠道的适配器实现文件。支持 OpenAI 兼容格式、Claude 格式、嵌入和图片生成等多种中继模式。使用标准的 Bearer Token 认证。

## 2. 依赖关系
- **标准库**: errors, fmt, io, net/http
- **内部包**:
  - `github.com/QuantumNous/new-api/constant` — 渠道常量（ChannelBaseURLs、ChannelSpecialBases）
  - `github.com/QuantumNous/new-api/dto` — 数据传输对象
  - `github.com/QuantumNous/new-api/relay/channel` — 渠道通用工具
  - `github.com/QuantumNous/new-api/relay/channel/claude` — Claude 渠道适配器
  - `github.com/QuantumNous/new-api/relay/channel/openai` — OpenAI 渠道适配器
  - `github.com/QuantumNous/new-api/relay/common` — RelayInfo
  - `github.com/QuantumNous/new-api/relay/constant` — 中继模式常量
  - `github.com/QuantumNous/new-api/types` — 错误类型
- **外部依赖**:
  - `github.com/gin-gonic/gin` — HTTP 框架
  - `github.com/samber/lo` — 泛型工具

## 3. 类型定义

### Adaptor 结构体
```go
type Adaptor struct {}
```
空结构体。

## 4. 函数详解

### GetRequestURL(...)
根据中继模式构建 URL：
- **Claude 格式**: `/api/anthropic/v1/messages`（支持特殊计划）
- **嵌入**: `/api/paas/v4/embeddings`
- **图片生成**: `/api/paas/v4/images/generations`
- **默认（聊天补全）**: `/api/paas/v4/chat/completions`

### SetupRequestHeader(...)
设置标准 Bearer Token 认证头。

### ConvertOpenAIRequest(...)
将 OpenAI 请求转换为智谱 v4 格式：
- TopP >= 1 时修正为 0.99
- 调用 `requestOpenAI2Zhipu` 转换请求

### ConvertClaudeRequest(...)
直接透传 Claude 请求。

### ConvertEmbeddingRequest(...)
直接透传嵌入请求。

### DoResponse(...)
根据中继格式和模式分发：
- Claude 格式: `claude.Adaptor.DoResponse`
- 图片生成: `zhipu4vImageHandler`
- 其他: `openai.Adaptor.DoResponse`

### 未实现的方法
- `ConvertGeminiRequest` / `ConvertAudioRequest` / `ConvertOpenAIResponsesRequest`

## 5. 关键逻辑分析

1. **多格式支持**: 同时支持 OpenAI 和 Claude 两种 API 格式，通过 `RelayFormat` 判断。

2. **v4 API**: 使用 `/api/paas/v4/` 前缀，区别于旧版的 `/api/paas/v3/`。

3. **特殊计划路由**: 通过 `ChannelSpecialBases` 支持自定义 API 端点。

4. **TopP 修正**: 与旧版相同，TopP >= 1 时修正为 0.99。

## 6. 关联文件
- `zhipu_4v/constants.go` — 模型列表
- `zhipu_4v/dto.go` — 响应数据结构
- `zhipu_4v/image.go` — 图片生成处理
- `zhipu_4v/relay-zhipu_v4.go` — 请求转换
- `relay/channel/claude/` — Claude 格式处理
- `relay/channel/openai/` — OpenAI 格式处理
