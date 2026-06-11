# relay-coze.go 代码阅读文档

## 1. 全局总结

本文件实现了 Coze 频道的请求转换和响应处理逻辑。包含 OpenAI → Coze 请求转换、流式/非流式响应处理、事件驱动的流式事件处理，以及辅助的轮询和消息获取功能。

## 2. 依赖关系

| 包 | 用途 |
|---|------|
| `bufio` | 流式扫描 |
| `encoding/json` | JSON 操作 |
| `errors` | 错误创建 |
| `fmt` | 字符串格式化 |
| `io` | I/O 操作 |
| `net/http` | HTTP 请求/响应 |
| `strings` | 字符串处理 |
| `common` | JSON 工具、日志、时间戳 |
| `dto` | 数据传输对象 |
| `relay/common` | RelayInfo 上下文 |
| `relay/helper` | 流式处理工具 |
| `service` | HTTP 客户端、Token 估算 |
| `types` | 错误类型 |
| `lo` | 工具库（FromPtrOr） |
| `gin` | Web 框架 |

## 3. 类型定义

无类型定义。

## 4. 函数详解

### `convertCozeChatRequest(c *gin.Context, request dto.GeneralOpenAIRequest) *CozeChatRequest`
将 OpenAI 请求转换为 Coze 格式：
1. 只提取 `user` 角色的消息，转换为 `CozeEnterMessage`
2. 设置 `ContentType: "text"`（TODO: 支持更多内容类型）
3. 如果没有设置 `User`，使用 response ID 作为用户标识
4. 从 gin 上下文获取 `bot_id`

### `cozeChatHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)`
处理 Coze 非流式响应：
1. 读取并解析响应体
2. 检查 `Code != 0` 表示错误
3. 从 gin 上下文获取 usage 数据（由轮询阶段设置）
4. 遍历 `Data` 数组，找到 `type == "answer"` 的消息作为响应内容
5. 构建 OpenAI 格式响应并返回

### `cozeChatStreamHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)`
处理 Coze 流式响应（事件驱动模式）：
1. 逐行扫描 SSE 流
2. 识别 `event:` 和 `data:` 前缀的行
3. 遇到空行时处理当前事件
4. 委托 `handleCozeEvent` 处理具体事件
5. 流结束后如果 usage 未获取，通过文本估算

### `handleCozeEvent(c *gin.Context, event string, data string, ...)`
处理具体的 Coze 流式事件：
- **`conversation.chat.completed`**: 解析 usage 数据，发送停止响应
- **`conversation.message.delta`**: 解析增量消息，转换为 OpenAI 流式 chunk 并发送
- **`error`**: 记录错误日志

### `checkIfChatComplete(a *Adaptor, c *gin.Context, info *relaycommon.RelayInfo) (error, bool)`
轮询检查聊天是否完成：
1. 发送 GET 请求到 `/v3/chat/retrieve`
2. 解析响应状态
3. `completed` → 设置 usage 到 gin 上下文，返回 true
4. `failed/canceled/requires_action` → 返回错误
5. 其他 → 返回 false（继续轮询）

### `getChatDetail(a *Adaptor, c *gin.Context, info *relaycommon.RelayInfo) (*http.Response, error)`
获取聊天消息详情：
1. 发送 GET 请求到 `/v3/chat/message/list`
2. 返回原始 HTTP 响应

### `doRequest(req *http.Request, info *relaycommon.RelayInfo) (*http.Response, error)`
执行 HTTP 请求，支持代理配置：
- 如果频道配置了代理，创建代理 HTTP 客户端
- 否则使用默认客户端

## 5. 关键逻辑分析

### 事件驱动的流式处理
Coze 的流式响应使用标准的 SSE 格式，但事件类型丰富。处理流程是：
```
扫描行 → 识别 event/data → 空行触发处理 → handleCozeEvent 分发
```

### 三种事件类型
1. `conversation.chat.completed`: 聊天完成，包含最终 usage
2. `conversation.message.delta`: 消息增量，包含文本片段
3. `error`: 错误事件

### 消息内容的二次解析
`CozeChatV3MessageDetail.Content` 是 `json.RawMessage` 类型，需要先解析为字符串再使用。

### 非流式的三步流程
非流式请求实际上是一个异步过程：
1. 创建聊天（返回 chat_id）
2. 轮询状态（每秒一次）
3. 获取消息详情

这种设计是因为 Coze 的 Bot 执行可能需要较长时间。

### Dify 兼容的 thinking 标签转换
流式处理中没有直接处理 thinking 标签，但在 `streamResponseDify2OpenAI` 中有类似的处理逻辑（这实际上是 Dify 的代码，可能被误包含在此文件中，但实际是在 `relay-dify.go` 中）。

### 代理支持
`doRequest` 函数支持通过频道设置配置 HTTP 代理，这对于需要代理访问 Coze API 的场景很有用。

## 6. 关联文件

- `relay/channel/coze/adaptor.go` - 调度本文件中的函数
- `relay/channel/coze/dto.go` - DTO 定义
- `relay/helper/` - 流式处理工具
- `service/http_client.go` - HTTP 客户端和代理支持
