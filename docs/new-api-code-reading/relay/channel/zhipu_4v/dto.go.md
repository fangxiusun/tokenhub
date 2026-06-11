# dto.go 代码阅读文档

## 1. 全局总结
智谱 AI v4 版渠道的数据传输对象定义文件。包含聊天响应和流式响应结构体。

## 2. 依赖关系
- **标准库**: time
- **内部包**:
  - `github.com/QuantumNous/new-api/dto` — 数据传输对象
  - `github.com/QuantumNous/new-api/types` — 错误类型

## 3. 类型定义

### ZhipuV4Response
```go
type ZhipuV4Response struct {
    Id                  string                         `json:"id"`
    Created             int64                          `json:"created"`
    Model               string                         `json:"model"`
    TextResponseChoices []dto.OpenAITextResponseChoice `json:"choices"`
    Usage               dto.Usage                      `json:"usage"`
    Error               types.OpenAIError              `json:"error"`
}
```
智谱 v4 非流式响应，与 OpenAI 格式高度兼容。

### ZhipuV4StreamResponse
```go
type ZhipuV4StreamResponse struct {
    Id      string                                    `json:"id"`
    Created int64                                     `json:"created"`
    Choices []dto.ChatCompletionsStreamResponseChoice `json:"choices"`
    Usage   dto.Usage                                 `json:"usage"`
}
```
智谱 v4 流式响应。

### tokenData
```go
type tokenData struct {
    Token      string
    ExpiryTime time.Time
}
```
Token 缓存数据（v4 版本未使用，保留结构）。

## 4. 函数详解
无函数定义。

## 5. 关键逻辑分析
- v4 版响应格式与 OpenAI 高度兼容，直接复用 `dto.OpenAITextResponseChoice` 和 `dto.ChatCompletionsStreamResponseChoice`
- 注释中保留了旧版的自定义结构体定义，说明已迁移到 OpenAI 兼容格式
- 使用 `types.OpenAIError` 处理错误响应

## 6. 关联文件
- `zhipu_4v/adaptor.go` — 使用这些 DTO
- `relay/channel/openai/` — OpenAI 格式处理
