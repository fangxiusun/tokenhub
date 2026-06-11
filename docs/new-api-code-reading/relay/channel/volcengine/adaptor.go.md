# adaptor.go 代码阅读文档

## 1. 全局总结
火山引擎（Volcengine）/ 豆包（Doubao）AI 渠道的适配器实现文件。支持多种中继模式：聊天补全、嵌入、图片生成/编辑、重排序、响应式、语音合成（TTS）。特别地，TTS 功能通过 WebSocket 协议实现。

## 2. 依赖关系
- **标准库**: bytes, encoding/json, errors, fmt, io, net/http, path/filepath, strings
- **内部包**:
  - `github.com/QuantumNous/new-api/constant` — 渠道常量（ChannelBaseURLs、ChannelSpecialBases）
  - `github.com/QuantumNous/new-api/dto` — 数据传输对象
  - `github.com/QuantumNous/new-api/relay/channel` — 渠道通用工具
  - `github.com/QuantumNous/new-api/relay/channel/claude` — Claude 渠道适配器
  - `github.com/QuantumNous/new-api/relay/channel/openai` — OpenAI 渠道适配器
  - `github.com/QuantumNous/new-api/relay/common` — RelayInfo
  - `github.com/QuantumNous/new-api/relay/constant` — 中继模式常量
  - `github.com/QuantumNous/new-api/setting/model_setting` — 模型设置
  - `github.com/QuantumNous/new-api/types` — 错误类型
- **外部依赖**:
  - `github.com/gin-gonic/gin` — HTTP 框架
  - `github.com/samber/lo` — 泛型工具

## 3. 类型定义

### Context Key 常量
```go
const (
    contextKeyTTSRequest     = "volcengine_tts_request"    // TTS 请求上下文键
    contextKeyResponseFormat = "response_format"           // 响应格式上下文键
)
```

### Adaptor 结构体
```go
type Adaptor struct {}
```
空结构体，所有状态通过方法参数传递。

## 4. 函数详解

### Init(info *relaycommon.RelayInfo)
空实现，火山引擎适配器无需初始化。

### GetRequestURL(info *relaycommon.RelayInfo) (string, error)
根据中继模式构建不同的 URL：
- **聊天补全**: `/api/v3/chat/completions`（bot 模型使用 `/api/v3/bots/chat/completions`）
- **嵌入**: `/api/v3/embeddings`
- **图片生成/编辑**: `/api/v3/images/generations`
- **重排序**: `/api/v3/rerank`
- **响应式**: `/api/v3/responses`
- **语音合成**: WebSocket URL `wss://openspeech.bytedance.com/api/v1/tts/ws_binary`
- 支持特殊计划（ChannelSpecialBases）的自定义 URL

### ConvertOpenAIRequest(...)
对 DeepSeek 模型的 `-thinking` 后缀进行智能处理，启用思考模式。

### ConvertClaudeRequest(...)
根据渠道基础 URL 判断使用 Claude 原生格式还是 OpenAI 兼容格式。

### ConvertAudioRequest(...)
构建火山引擎 TTS 请求：
1. 解析 API Key 获取 appID 和 token
2. 映射语音类型和编码格式
3. 支持通过 metadata 覆盖默认配置
4. 设置操作类型（submit 时强制流式）

### ConvertImageRequest(...)
图片生成请求处理，直接透传请求。

### DoRequest(...)
TTS 非流式模式下返回空响应（火山引擎 TTS 的非流式模式不需要 HTTP 请求）。

### DoResponse(...)
根据中继模式分发：
- Claude 格式: 委托给 claude.Adaptor
- TTS 流式: `handleTTSWebSocketResponse`（WebSocket）
- TTS 非流式: `handleTTSResponse`
- 其他: 委托给 openai.Adaptor

### SetupRequestHeader(...)
TTS 模式使用特殊的认证格式 `Bearer;{token}`，其他模式使用标准 `Bearer` 格式。

### detectImageMimeType(filename string) string
根据文件扩展名检测图片 MIME 类型。

## 5. 关键逻辑分析

1. **多协议支持**: 火山引擎适配器同时支持 HTTP 和 WebSocket 两种协议，TTS 使用 WebSocket。

2. **特殊计划路由**: 通过 `ChannelSpecialBases` 支持不同的基础 URL，允许自定义 API 端点。

3. **TTS 请求流程**: 
   - 构建请求并存入 Gin 上下文
   - 在 `DoRequest` 中检查是否需要发送 HTTP 请求
   - 在 `DoResponse` 中处理 WebSocket 响应

4. **DeepSeek 思考模式**: 自动为 DeepSeek 模型处理 `-thinking` 后缀，启用 `THINKING` 字段。

5. **Bot 模型路由**: 以 "bot" 开头的模型使用不同的 API 路径。

## 6. 关联文件
- `volcengine/constants.go` — 模型列表
- `volcengine/protocols.go` — WebSocket 二进制协议实现
- `volcengine/tts.go` — TTS 请求/响应结构和处理
- `relay/channel/openai/` — OpenAI 格式响应处理
- `relay/channel/claude/` — Claude 格式响应处理
