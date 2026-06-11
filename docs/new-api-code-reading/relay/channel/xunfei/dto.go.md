# dto.go 代码阅读文档

## 1. 全局总结
讯飞星火渠道的数据传输对象定义文件，定义了请求和响应的数据结构。讯飞的请求格式与 OpenAI 差异较大，使用嵌套的 header/parameter/payload 结构。

## 2. 依赖关系
- **内部包**: `github.com/QuantumNous/new-api/dto`

## 3. 类型定义

### XunfeiMessage
```go
type XunfeiMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}
```

### XunfeiChatRequest
讯飞星火聊天请求，三层嵌套结构：
- **Header**: 包含 `AppId`
- **Parameter.Chat**: 包含 `Domain`、`Temperature`、`TopK`、`MaxTokens`、`Auditing`
- **Payload.Message.Text**: 消息列表

### XunfeiChatResponseTextItem
```go
type XunfeiChatResponseTextItem struct {
    Content string `json:"content"`
    Role    string `json:"role"`
    Index   int    `json:"index"`
}
```

### XunfeiChatResponse
讯飞星火聊天响应：
- **Header**: 包含 `Code`、`Message`、`Sid`、`Status`
- **Payload.Choices**: 包含 `Status`、`Seq`、`Text`（文本列表）
- **Payload.Usage**: 包含 `Text`（使用 `dto.Usage`）

## 4. 函数详解
无函数定义。

## 5. 关键逻辑分析

1. **嵌套结构**: 讯飞的请求格式使用三层嵌套（header/parameter/payload），与 OpenAI 的扁平结构差异很大。

2. **Domain 参数**: `Domain` 字段用于标识使用的模型领域，由 API 版本映射而来。

3. **Status 字段**: 响应中的 `Status` 字段表示消息状态，值为 2 表示结束。

4. **文本数组**: 响应中的 `Text` 是数组格式，即使只有一条消息也包装在数组中。

## 6. 关联文件
- `xunfei/relay-xunfei.go` — 请求/响应转换逻辑
- `xunfei/adaptor.go` — 使用这些 DTO
