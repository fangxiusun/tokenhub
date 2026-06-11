# utils.go 代码阅读文档

## 1. 全局总结

该文件提供了中间件通用的错误响应工具函数，包括 `abortWithOpenAiMessage`（OpenAI 格式错误）和 `abortWithMidjourneyMessage`（Midjourney 格式错误），用于统一错误响应格式和日志记录。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `fmt` | 格式化日志消息 |
| `github.com/QuantumNous/new-api/common` | 请求 ID 获取、消息格式化 |
| `github.com/QuantumNous/new-api/logger` | 错误日志记录 |
| `github.com/QuantumNous/new-api/types` | 错误码类型 |
| `github.com/gin-gonic/gin` | Gin Web 框架 |

## 3. 类型定义

无自定义类型定义。

## 4. 函数详解

### `abortWithOpenAiMessage(c *gin.Context, statusCode int, message string, code ...types.ErrorCode)`

- **功能**：返回 OpenAI 兼容格式的错误响应并终止请求。
- **参数**：
  - `statusCode`：HTTP 状态码
  - `message`：错误消息
  - `code`：可选的错误码
- **响应格式**：
  ```json
  {
    "error": {
      "message": "... (req_id: xxx)",
      "type": "new_api_error",
      "code": "..."
    }
  }
  ```
- **副作用**：调用 `logger.LogError` 记录错误日志，包含用户 ID。

### `abortWithMidjourneyMessage(c *gin.Context, statusCode int, code int, description string)`

- **功能**：返回 Midjourney 兼容格式的错误响应并终止请求。
- **参数**：
  - `statusCode`：HTTP 状态码
  - `code`：数字错误码
  - `description`：错误描述
- **响应格式**：
  ```json
  {
    "description": "...",
    "type": "new_api_error",
    "code": 123
  }
  ```
- **副作用**：调用 `logger.LogError` 记录错误日志。

## 5. 关键逻辑分析

- **请求 ID 关联**：`abortWithOpenAiMessage` 自动将请求 ID 附加到错误消息中，便于问题追踪。
- **统一错误格式**：两种函数分别适配 OpenAI 和 Midjourney 的错误响应格式，保持 API 兼容性。
- **日志记录**：所有错误响应都通过 `logger.LogError` 记录，便于监控和调试。
- **错误码支持**：`abortWithOpenAiMessage` 支持可选的错误码参数，提供更精细的错误分类。

## 6. 关联文件

- `middleware/jimeng_adapter.go` — 使用 `abortWithOpenAiMessage` 的适配器
- `middleware/request-id.go` — `RequestIdKey` 常量定义
- `logger/logger.go` — `LogError` 日志函数
- `types/error.go` — `ErrorCode` 类型定义
