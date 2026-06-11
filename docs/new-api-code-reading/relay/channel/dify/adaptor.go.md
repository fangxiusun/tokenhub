# adaptor.go 代码阅读文档

## 1. 全局总结

本文件实现了 Dify 频道的适配器。Dify 是一个开源的 LLM 应用开发平台。该适配器支持 ChatFlow、Agent、WorkFlow 和 Completion 四种 Bot 类型，通过模型名前缀或配置确定 Bot 类型，构建不同的 API URL。

## 2. 依赖关系

| 包 | 用途 |
|---|------|
| `errors` | 错误创建 |
| `fmt` | 字符串格式化 |
| `io` | I/O 操作 |
| `net/http` | HTTP 请求/响应 |
| `dto` | 数据传输对象 |
| `relay/channel` | 通用频道工具 |
| `relay/common` | RelayInfo 上下文 |
| `types` | 错误类型 |
| `gin` | Web 框架 |

## 3. 类型定义

### Bot 类型常量

```go
const (
    BotTypeChatFlow   = 1 // chatflow default
    BotTypeAgent      = 2
    BotTypeWorkFlow   = 3
    BotTypeCompletion = 4
)
```

Dify Bot 类型枚举：
- `BotTypeChatFlow (1)`: 聊天流（默认）
- `BotTypeAgent (2)`: Agent 类型
- `BotTypeWorkFlow (3)`: 工作流类型
- `BotTypeCompletion (4)`: 补全类型

### `Adaptor` 结构体

```go
type Adaptor struct {
    BotType int
}
```

有状态的适配器，包含 `BotType` 字段。与其他适配器不同，Dify Adaptor 需要维护 Bot 类型状态。

## 4. 函数详解

### `(a *Adaptor) Init(info *relaycommon.RelayInfo)`
初始化适配器，当前硬编码设置 `BotType = BotTypeChatFlow`。注释中的代码显示原本计划根据模型名前缀自动判断 Bot 类型（agent/workflow/chat）。

### `(a *Adaptor) GetRequestURL(info *relaycommon.RelayInfo) (string, error)`
根据 Bot 类型构建不同的 API URL：
- **WorkFlow**: `{base}/v1/workflows/run`
- **Completion**: `{base}/v1/completion-messages`
- **Agent / 默认**: `{base}/v1/chat-messages`

### `(a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error`
设置 `Authorization: Bearer {apiKey}` 请求头。

### `(a *Adaptor) ConvertOpenAIRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error)`
将 OpenAI 请求转换为 Dify 格式，调用 `requestOpenAI2Dify`。

### `(a *Adaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (any, error)`
委托给 `channel.DoApiRequest`。

### `(a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage any, err *types.NewAPIError)`
根据流式标志分发：
- 流式: `difyStreamHandler`
- 非流式: `difyHandler`

### 未实现的方法
- `ConvertGeminiRequest` - 返回 "not implemented"
- `ConvertClaudeRequest` - panic
- `ConvertAudioRequest` - 返回 "not implemented"
- `ConvertImageRequest` - 返回 "not implemented"
- `ConvertEmbeddingRequest` - 返回 "not implemented"
- `ConvertRerankRequest` - 返回 `nil, nil`
- `ConvertOpenAIResponsesRequest` - 返回 "not implemented"

### `(a *Adaptor) GetModelList() []string`
返回 Dify 模型列表（当前为空）。

### `(a *Adaptor) GetChannelName() string`
返回 `"dify"`。

## 5. 关键逻辑分析

### 有状态的适配器
Dify Adaptor 是唯一一个有状态的适配器（包含 `BotType` 字段）。这是因为 Dify 的不同 Bot 类型需要不同的 API 端点。

### Bot 类型判断机制（已注释）
Init 方法中的注释代码显示了原本的设计意图：根据模型名前缀自动判断 Bot 类型。但当前实现硬编码为 ChatFlow，可能是因为 Bot 类型判断逻辑需要更完善的设计。

### 空模型列表
`ModelList` 为空切片，这是因为 Dify 作为平台不直接暴露模型名，模型由 Dify 内部管理。

### Completions 端点
`BotTypeCompletion` 对应的 `/v1/completion-messages` 端点用于非对话式的文本补全任务。

## 6. 关联文件

- `relay/channel/dify/relay-dify.go` - 请求转换和响应处理的具体实现
- `relay/channel/dify/dto.go` - Dify 特有的 DTO 定义
- `relay/channel/dify/constants.go` - 模型列表和 Bot 类型常量
