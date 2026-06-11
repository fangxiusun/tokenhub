# dto.go 代码阅读文档

## 1. 全局总结

本文件定义了 Dify 频道特有的数据传输对象，包括请求结构、文件结构、元数据结构和响应结构（流式/非流式）。

## 2. 依赖关系

| 包 | 用途 |
|---|------|
| `dto` | 通用 DTO（Usage 类型） |

## 3. 类型定义

### `DifyChatRequest` 结构体

```go
type DifyChatRequest struct {
    Inputs           map[string]interface{} `json:"inputs"`
    Query            string                 `json:"query"`
    ResponseMode     string                 `json:"response_mode"`
    User             string                 `json:"user"`
    AutoGenerateName bool                   `json:"auto_generate_name"`
    Files            []DifyFile             `json:"files"`
}
```

Dify Chat API 请求格式：
- `Inputs`: 输入变量（KV 对）
- `Query`: 用户查询文本
- `ResponseMode`: 响应模式（"blocking" 或 "streaming"）
- `User`: 用户标识
- `AutoGenerateName`: 是否自动生成对话名称
- `Files`: 附件文件列表

### `DifyFile` 结构体

```go
type DifyFile struct {
    Type         string `json:"type"`
    TransferMode string `json:"transfer_mode"`
    URL          string `json:"url,omitempty"`
    UploadFileId string `json:"upload_file_id,omitempty"`
}
```

Dify 文件结构：
- `Type`: 文件类型（如 "image"、MIME 类型）
- `TransferMode`: 传输模式（"remote_url" 或 "local_file"）
- `URL`: 远程文件 URL（remote_url 模式）
- `UploadFileId`: 上传文件 ID（local_file 模式）

### `DifyMetaData` 结构体

```go
type DifyMetaData struct {
    Usage dto.Usage `json:"usage"`
}
```

Dify 响应元数据，包含 token 用量。

### `DifyData` 结构体

```go
type DifyData struct {
    WorkflowId string `json:"workflow_id"`
    NodeId     string `json:"node_id"`
    NodeType   string `json:"node_type"`
    Status     string `json:"status"`
}
```

流式响应中的工作流/节点数据。

### `DifyChatCompletionResponse` 结构体

```go
type DifyChatCompletionResponse struct {
    ConversationId string       `json:"conversation_id"`
    Answer         string       `json:"answer"`
    CreateAt       int64        `json:"create_at"`
    MetaData       DifyMetaData `json:"metadata"`
}
```

Dify 非流式响应格式。

### `DifyChunkChatCompletionResponse` 结构体

```go
type DifyChunkChatCompletionResponse struct {
    Event          string       `json:"event"`
    ConversationId string       `json:"conversation_id"`
    Answer         string       `json:"answer"`
    Data           DifyData     `json:"data"`
    MetaData       DifyMetaData `json:"metadata"`
}
```

Dify 流式响应格式，额外包含 `Event` 和 `Data` 字段。

## 4. 函数详解

无函数定义。

## 5. 关键逻辑分析

### 文件传输的两种模式
Dify 支持两种文件传输方式：
- `remote_url`: 直接引用远程 URL，无需上传
- `local_file`: 先上传文件获取 ID，再通过 ID 引用

### 流式事件类型
流式响应中的 `Event` 字段标识事件类型，常见的有：
- `message` / `agent_message`: 消息增量
- `message_end`: 消息结束
- `workflow_*`: 工作流事件
- `node_*`: 节点事件
- `error`: 错误事件

### Inputs 的灵活性
`Inputs` 使用 `map[string]interface{}` 类型，允许传入任意键值对作为 Dify Bot 的输入变量。

## 6. 关联文件

- `relay/channel/dify/relay-dify.go` - 使用这些 DTO 进行请求构建和响应解析
- `relay/channel/dify/adaptor.go` - 调度使用
