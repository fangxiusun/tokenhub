# dto.go 代码阅读文档

## 1. 全局总结
智谱 AI 旧版渠道的数据传输对象定义文件。

## 2. 依赖关系
- **标准库**: time
- **内部包**: `github.com/QuantumNous/new-api/dto`

## 3. 类型定义

### ZhipuMessage
```go
type ZhipuMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}
```

### ZhipuRequest
```go
type ZhipuRequest struct {
    Prompt      []ZhipuMessage `json:"prompt"`
    Temperature *float64       `json:"temperature,omitempty"`
    TopP        float64        `json:"top_p,omitempty"`
    RequestId   string         `json:"request_id,omitempty"`
    Incremental bool           `json:"incremental,omitempty"`
}
```
智谱请求使用 `prompt` 字段（而非 `messages`）。

### ZhipuResponseData
```go
type ZhipuResponseData struct {
    TaskId     string         `json:"task_id"`
    RequestId  string         `json:"request_id"`
    TaskStatus string         `json:"task_status"`
    Choices    []ZhipuMessage `json:"choices"`
    dto.Usage  `json:"usage"`
}
```

### ZhipuResponse
```go
type ZhipuResponse struct {
    Code    int               `json:"code"`
    Msg     string            `json:"msg"`
    Success bool              `json:"success"`
    Data    ZhipuResponseData `json:"data"`
}
```

### ZhipuStreamMetaResponse
流式响应的元数据，包含 `RequestId`、`TaskId`、`TaskStatus` 和 `Usage`。

### zhipuTokenData
```go
type zhipuTokenData struct {
    Token      string
    ExpiryTime time.Time
}
```
缓存的 JWT Token 数据。

## 4. 函数详解
无函数定义。

## 5. 关键逻辑分析
- 智谱旧版 API 使用 `prompt` 而非 `messages`
- 响应包含 `TaskId` 和 `TaskStatus` 用于异步任务管理
- 流式响应使用 `meta:` 前缀传递元数据

## 6. 关联文件
- `zhipu/relay-zhipu.go` — 使用这些 DTO
