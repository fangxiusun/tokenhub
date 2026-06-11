# adaptor.go 代码阅读文档

## 1. 全局总结
讯飞星火（Xunfei SparkDesk）渠道的适配器实现文件。讯飞使用 WebSocket 协议进行通信，适配器在 `DoRequest` 中返回空响应，实际通信在 `DoResponse` 中通过 WebSocket 完成。

## 2. 依赖关系
- **标准库**: errors, io, net/http, strings
- **内部包**:
  - `github.com/QuantumNous/new-api/dto` — 数据传输对象
  - `github.com/QuantumNous/new-api/relay/channel` — 渠道通用工具
  - `github.com/QuantumNous/new-api/relay/common` — RelayInfo
  - `github.com/QuantumNous/new-api/types` — 错误类型
- **外部依赖**:
  - `github.com/gin-gonic/gin` — HTTP 框架

## 3. 类型定义

### Adaptor 结构体
```go
type Adaptor struct {
    request *dto.GeneralOpenAIRequest  // 缓存的请求对象
}
```
缓存 OpenAI 请求，供后续 WebSocket 通信使用。

## 4. 函数详解

### ConvertOpenAIRequest(...)
将 OpenAI 请求缓存到 `a.request`，直接返回原请求。

### DoRequest(...)
返回空的 HTTP 响应（状态码 200），因为讯飞使用 WebSocket 而非 HTTP。

### DoResponse(...)
核心处理方法：
1. 解析 API Key（格式: `appId|apiSecret|apiKey`）
2. 根据流式状态调用 `xunfeiStreamHandler` 或 `xunfeiHandler`

### SetupRequestHeader(...)
仅设置标准 API 请求头，无特殊认证头。

### 未实现的方法
- `ConvertGeminiRequest` / `ConvertClaudeRequest` / `ConvertAudioRequest` / `ConvertImageRequest` / `ConvertEmbeddingRequest` / `ConvertOpenAIResponsesRequest`

## 5. 关键逻辑分析

1. **WebSocket 通信模式**: 讯飞星火使用 WebSocket 进行实时通信，而非标准的 HTTP 请求/响应模式。适配器通过缓存请求并在 `DoResponse` 中建立 WebSocket 连接来实现。

2. **API Key 格式**: 三段式格式 `appId|apiSecret|apiKey`，分别用于标识应用、生成签名和认证。

3. **请求缓存**: `ConvertOpenAIRequest` 将请求缓存到结构体字段，供 `DoResponse` 中的 WebSocket 通信使用。

## 6. 关联文件
- `xunfei/constants.go` — 模型列表
- `xunfei/dto.go` — 请求/响应数据结构
- `xunfei/relay-xunfei.go` — WebSocket 通信、签名、请求/响应转换
