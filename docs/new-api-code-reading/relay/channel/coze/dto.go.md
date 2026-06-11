# dto.go 代码阅读文档

## 1. 全局总结

本文件定义了 Coze 频道特有的数据传输对象，包括请求结构、响应结构、消息详情结构和错误结构。这些 DTO 用于与 Coze v3 Chat API 进行数据交换。

## 2. 依赖关系

| 包 | 用途 |
|---|------|
| `encoding/json` | JSON RawMessage 类型 |

## 3. 类型定义

### `CozeError` 结构体

```go
type CozeError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}
```

Coze API 错误格式。

### `CozeEnterMessage` 结构体

```go
type CozeEnterMessage struct {
    Role        string          `json:"role"`
    Type        string          `json:"type,omitempty"`
    Content     any             `json:"content,omitempty"`
    MetaData    json.RawMessage `json:"meta_data,omitempty"`
    ContentType string          `json:"content_type,omitempty"`
}
```

发送给 Coze 的消息格式：
- `Role`: 消息角色（"user"）
- `Content`: 消息内容（any 类型，支持多种格式）
- `ContentType`: 内容类型（如 "text"）

### `CozeChatRequest` 结构体

```go
type CozeChatRequest struct {
    BotId              string             `json:"bot_id"`
    UserId             json.RawMessage    `json:"user_id"`
    AdditionalMessages []CozeEnterMessage `json:"additional_messages,omitempty"`
    Stream             bool               `json:"stream,omitempty"`
    CustomVariables    json.RawMessage    `json:"custom_variables,omitempty"`
    AutoSaveHistory    bool               `json:"auto_save_history,omitempty"`
    MetaData           json.RawMessage    `json:"meta_data,omitempty"`
    ExtraParams        json.RawMessage    `json:"extra_params,omitempty"`
    ShortcutCommand    json.RawMessage    `json:"shortcut_command,omitempty"`
    Parameters         json.RawMessage    `json:"parameters,omitempty"`
}
```

Coze v3 Chat API 请求格式，核心字段：
- `BotId`: Bot 标识（从 gin 上下文获取）
- `UserId`: 用户标识
- `AdditionalMessages`: 附加消息列表
- `Stream`: 是否流式输出

### `CozeChatResponse` 结构体

```go
type CozeChatResponse struct {
    Code int                  `json:"code"`
    Msg  string               `json:"msg"`
    Data CozeChatResponseData `json:"data"`
}
```

Coze Chat API 的通用响应格式。

### `CozeChatResponseData` 结构体

```go
type CozeChatResponseData struct {
    Id             string        `json:"id"`
    ConversationId string        `json:"conversation_id"`
    BotId          string        `json:"bot_id"`
    CreatedAt      int64         `json:"created_at"`
    LastError      CozeError     `json:"last_error"`
    Status         string        `json:"status"`
    Usage          CozeChatUsage `json:"usage"`
}
```

Chat 响应数据，包含聊天 ID、会话 ID、状态和用量。

### `CozeChatUsage` 结构体

```go
type CozeChatUsage struct {
    TokenCount  int `json:"token_count"`
    OutputCount int `json:"output_count"`
    InputCount  int `json:"input_count"`
}
```

Coze token 用量统计。

### `CozeChatDetailResponse` 结构体

```go
type CozeChatDetailResponse struct {
    Data   []CozeChatV3MessageDetail `json:"data"`
    Code   int                       `json:"code"`
    Msg    string                    `json:"msg"`
    Detail CozeResponseDetail        `json:"detail"`
}
```

聊天详情响应，包含消息列表。

### `CozeChatV3MessageDetail` 结构体

```go
type CozeChatV3MessageDetail struct {
    Id               string          `json:"id"`
    Role             string          `json:"role"`
    Type             string          `json:"type"`
    BotId            string          `json:"bot_id"`
    ChatId           string          `json:"chat_id"`
    Content          json.RawMessage `json:"content"`
    MetaData         json.RawMessage `json:"meta_data"`
    CreatedAt        int64           `json:"created_at"`
    SectionId        string          `json:"section_id"`
    UpdatedAt        int64           `json:"updated_at"`
    ContentType      string          `json:"content_type"`
    ConversationId   string          `json:"conversation_id"`
    ReasoningContent string          `json:"reasoning_content"`
}
```

V3 版本消息详情，包含：
- `Content`: 消息内容（RawMessage，需要二次解析）
- `Type`: 消息类型（如 "answer"）
- `ReasoningContent`: 推理内容（思考链）

### `CozeResponseDetail` 结构体

```go
type CozeResponseDetail struct {
    Logid string `json:"logid"`
}
```

响应详情，包含日志 ID。

## 4. 函数详解

无函数定义。

## 5. 关键逻辑分析

### RawMessage 的使用
大量字段使用 `json.RawMessage` 类型（如 `Content`、`MetaData`、`CustomVariables` 等），这意味着这些字段的 JSON 结构是灵活的，需要在使用时根据上下文进行二次解析。

### 消息类型区分
`CozeChatV3MessageDetail` 中的 `Type` 字段用于区分消息类型，`"answer"` 类型的消息包含 Bot 的回复内容。

### 推理内容支持
`ReasoningContent` 字段表明 Coze 支持返回模型的推理过程（思考链），这对于推理模型（如 DeepSeek-R1）特别重要。

## 6. 关联文件

- `relay/channel/coze/relay-coze.go` - 使用这些 DTO 进行请求构建和响应解析
- `relay/channel/coze/adaptor.go` - 调度使用
