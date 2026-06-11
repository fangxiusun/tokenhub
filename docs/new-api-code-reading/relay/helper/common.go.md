# common.go 代码阅读文档

## 1. 全局总结

本文件提供了 Relay 模块的通用工具函数，包括 SSE 流式响应写入、WebSocket 消息发送、响应 ID 生成等。是所有 handler 共用的底层工具集。

## 2. 依赖关系

- `common`: JSON 序列化
- `dto`: 响应 DTO
- `types`: 错误类型
- `gorilla/websocket`: WebSocket

## 3. 类型定义

本文件无自定义类型定义。

## 4. 函数详解

### SSE 流式工具

#### `FlushWriter(c) error`
- 刷新 HTTP 响应写入器，带 panic 恢复

#### `SetEventStreamHeaders(c)`
- 设置 SSE 响应头（Content-Type, Cache-Control, Connection 等）
- 幂等：通过 context 标志避免重复设置

#### `ClaudeData(c, resp) error`
- 发送 Claude 格式的 SSE 数据

#### `ClaueChunkData(c, resp, data)`
- 发送 Claude 格式的 SSE chunk 数据

#### `ResponseChunkData(c, resp, data)`
- 发送 Responses API 格式的 SSE chunk 数据

#### `StringData(c, str) error`
- 发送纯文本 SSE 数据

#### `PingData(c) error`
- 发送 SSE ping 注释（`: PING\n\n`）

#### `ObjectData(c, object) error`
- 发送 JSON 对象 SSE 数据

#### `Done(c)`
- 发送 `[DONE]` 终止信号

### WebSocket 工具

#### `WssString(c, ws, str) error`
- 发送 WebSocket 文本消息

#### `WssObject(c, ws, object) error`
- 发送 WebSocket JSON 消息

#### `WssError(c, ws, openaiError)`
- 发送 WebSocket 错误事件

### ID 生成

#### `GetResponseID(c) string`
- 生成 chatcmpl-{requestId} 格式的响应 ID

#### `GetLocalRealtimeID(c) string`
- 生成 evt_{requestId} 格式的实时事件 ID

### 响应生成

#### `GenerateStartEmptyResponse(id, createAt, model, systemFingerprint)`
- 生成流式起始空响应

#### `GenerateStopResponse(id, createAt, model, finishReason)`
- 生成流式停止响应

#### `GenerateFinalUsageResponse(id, createAt, model, usage)`
- 生成流式最终使用量响应

## 5. 关键逻辑分析

1. **安全写入**: 所有写入操作都检查 context 是否已取消
2. **Panic 恢复**: FlushWriter 带 panic 恢复，防止 flush 异常导致崩溃
3. **幂等头设置**: SetEventStreamHeaders 通过 context 标志避免重复设置

## 6. 关联文件

- `common/custom_event.go`: CustomEvent SSE 渲染
- `dto/chat.go`: ChatCompletionsStreamResponse
