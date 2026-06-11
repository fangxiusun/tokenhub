# dto.go 代码阅读文档

## 1. 全局总结
腾讯混元 API 的数据传输对象（DTO）定义文件。定义了请求和响应的所有结构体，包括消息格式、聊天请求、错误信息、使用量统计和响应体。

## 2. 依赖关系
无外部依赖，纯类型定义文件。

## 3. 类型定义

### TencentMessage
```go
type TencentMessage struct {
    Role    string `json:"Role"`
    Content string `json:"Content"`
}
```
腾讯混元消息格式。注意 JSON 标签使用大写 `Role` 和 `Content`，与 OpenAI 的小写格式不同。

### TencentChatRequest
腾讯混元聊天请求体，字段包括：
- `Model *string` — 模型名称（hunyuan-lite/standard/standard-256K/pro）
- `Messages []*TencentMessage` — 聊天上下文，最多 40 条消息
- `Stream *bool` — 是否流式调用
- `TopP *float64` — 多样性控制 [0.0, 1.0]
- `Temperature *float64` — 温度参数 [0.0, 2.0]

### TencentError
```go
type TencentError struct {
    Code    int    `json:"Code"`
    Message string `json:"Message"`
}
```
腾讯 API 错误信息结构。

### TencentUsage
```go
type TencentUsage struct {
    PromptTokens     int `json:"PromptTokens"`
    CompletionTokens int `json:"CompletionTokens"`
    TotalTokens      int `json:"TotalTokens"`
}
```
Token 使用量统计。

### TencentResponseChoices
响应选项结构：
- `FinishReason` — 流式结束标志（"stop" 表示尾包）
- `Messages` — 同步模式返回的内容
- `Delta` — 流模式返回的增量内容

### TencentChatResponse
完整的聊天响应体，包含：
- `Choices` — 响应选项数组
- `Created` — Unix 时间戳
- `Id` — 会话 ID
- `Usage` — Token 用量
- `Error` — 错误信息（可能为 null）
- `ReqID` — 唯一请求 ID

### TencentChatResponseSB
```go
type TencentChatResponseSB struct {
    Response TencentChatResponse `json:"Response,omitempty"`
}
```
非流式响应的包装结构，外层包裹 `Response` 字段。

## 4. 函数详解
无函数定义。

## 5. 关键逻辑分析
1. **字段命名差异**: 腾讯混元 API 使用大写 JSON 字段名（如 `Role`、`Content`），与 OpenAI 的小写格式不同，需要在转换时注意。
2. **流式/非流式双模式**: `TencentResponseChoices` 同时包含 `Messages`（同步）和 `Delta`（流式）字段，根据调用模式取用不同字段。
3. **错误处理**: `Error` 字段可能返回 null，需要在解析时做空值判断。

## 6. 关联文件
- `tencent/relay-tencent.go` — 使用这些 DTO 进行请求/响应转换
- `tencent/adaptor.go` — 在 `ConvertOpenAIRequest` 中构建 `TencentChatRequest`
