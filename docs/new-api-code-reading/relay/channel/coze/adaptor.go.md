# adaptor.go 代码阅读文档

## 1. 全局总结

本文件实现了 Coze 频道的适配器。Coze 是字节跳动推出的 AI Bot 开发平台。该适配器实现了标准的 `channel.Adaptor` 接口，处理请求转换、HTTP 请求执行（包含非流式的轮询等待逻辑）和响应分发。

## 2. 依赖关系

| 包 | 用途 |
|---|------|
| `encoding/json` | JSON 操作 |
| `errors` | 错误创建 |
| `fmt` | 字符串格式化 |
| `io` | I/O 操作 |
| `net/http` | HTTP 请求/响应 |
| `time` | 轮询间隔 |
| `dto` | 数据传输对象 |
| `relay/channel` | 通用频道工具 |
| `relay/common`（别名 `common`） | RelayInfo 上下文 |
| `types` | 错误类型 |
| `gin` | Web 框架 |

## 3. 类型定义

### `Adaptor` 结构体

```go
type Adaptor struct{}
```

无状态适配器，实现 `channel.Adaptor` 接口。

## 4. 函数详解

### `(a *Adaptor) Init(info *relaycommon.RelayInfo)`
空实现。

### `(a *Adaptor) GetRequestURL(info *relaycommon.RelayInfo) (string, error)`
返回 Coze Chat API URL：`{base}/v3/chat`。

### `(a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error`
设置 `Authorization: Bearer {apiKey}` 请求头。

### `(a *Adaptor) ConvertOpenAIRequest(c *gin.Context, info *common.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error)`
将 OpenAI 请求转换为 Coze 格式，调用 `convertCozeChatRequest`。

### `(a *Adaptor) DoRequest(c *gin.Context, info *common.RelayInfo, requestBody io.Reader) (any, error)`
核心请求执行逻辑，根据流式/非流式有不同行为：

**流式模式**：直接委托给 `channel.DoApiRequest`。

**非流式模式**（三步流程）：
1. 发送创建消息请求
2. 解析响应获取 `conversation_id` 和 `chat_id`
3. 轮询检查聊天是否完成（`checkIfChatComplete`，每秒轮询一次）
4. 完成后获取聊天详情（`getChatDetail`）

### `(a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, info *common.RelayInfo) (usage any, err *types.NewAPIError)`
根据流式标志分发：
- 流式: `cozeChatStreamHandler`
- 非流式: `cozeChatHandler`

### 未实现的方法
- `ConvertGeminiRequest`、`ConvertClaudeRequest`、`ConvertAudioRequest`、`ConvertEmbeddingRequest`、`ConvertImageRequest`、`ConvertOpenAIResponsesRequest`、`ConvertRerankRequest` 均返回 "not implemented"

### `(a *Adaptor) GetModelList() []string`
返回 Coze 模型列表。

### `(a *Adaptor) GetChannelName() string`
返回 `"coze"`。

## 5. 关键逻辑分析

### 非流式的轮询模式
Coze 的非流式请求实际上是异步的：先创建消息，然后轮询等待完成，最后获取结果。这导致了较长的响应延迟，但保证了数据完整性。轮询间隔为 1 秒。

### 轮询的错误处理
轮询循环中，如果状态为 `failed`、`canceled` 或 `requires_action`，会立即返回错误，避免无限轮询。

### 上下文数据传递
Coze 适配器通过 `gin.Context` 的 `Set/Get` 方法在不同函数间传递 `conversation_id` 和 `chat_id`，以及 usage 数据。

### 不支持的功能
Coze 适配器不支持 Rerank、Embedding、Image、Audio、Responses 等端点，仅支持 Chat 功能。

## 6. 关联文件

- `relay/channel/coze/relay-coze.go` - 请求转换和响应处理的具体实现
- `relay/channel/coze/dto.go` - Coze 特有的 DTO 定义
- `relay/channel/coze/constants.go` - 模型列表
