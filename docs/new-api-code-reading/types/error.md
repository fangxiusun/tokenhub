# error.go 代码阅读文档

## 1. 全局概述

本文件是系统错误处理的核心文件，定义了统一的错误类型体系。包括 OpenAI 格式错误、Claude 格式错误、NewAPI 统一错误类型、错误码常量以及丰富的错误构造函数。支持错误转换、敏感信息掩码、重试控制等功能。

## 2. 依赖关系

- `encoding/json` — JSON 序列化（用于 `json.RawMessage`）
- `errors` — Go 标准库错误包
- `fmt` — 格式化输出
- `net/http` — HTTP 状态码常量
- `strings` — 字符串操作
- `github.com/QuantumNous/new-api/common` — 通用工具函数

## 3. 类型定义

### OpenAIError 结构体

```go
type OpenAIError struct {
    Message  string          `json:"message"`
    Type     string          `json:"type"`
    Param    string          `json:"param"`
    Code     any             `json:"code"`
    Metadata json.RawMessage `json:"metadata,omitempty"`
}
```

OpenAI 格式的错误响应结构。

### ClaudeError 结构体

```go
type ClaudeError struct {
    Type    string `json:"type,omitempty"`
    Message string `json:"message,omitempty"`
}
```

Claude 格式的错误响应结构。

### ErrorType 类型

```go
type ErrorType string
```

错误类型标识符。

### ErrorCode 类型

```go
type ErrorCode string
```

错误码标识符。

### NewAPIError 结构体

```go
type NewAPIError struct {
    Err            error
    RelayError     any
    skipRetry      bool
    recordErrorLog *bool
    errorType      ErrorType
    errorCode      ErrorCode
    StatusCode     int
    Metadata       json.RawMessage
}
```

系统统一的错误类型，包装了所有错误信息。

### NewAPIErrorOptions 类型

```go
type NewAPIErrorOptions func(*NewAPIError)
```

函数选项模式，用于配置错误构造。

## 4. 函数详情

### 错误构造函数

- `NewError(err error, errorCode ErrorCode, ops ...NewAPIErrorOptions) *NewAPIError` — 创建通用错误
- `NewOpenAIError(err error, errorCode ErrorCode, statusCode int, ops ...NewAPIErrorOptions) *NewAPIError` — 创建 OpenAI 格式错误
- `InitOpenAIError(errorCode ErrorCode, statusCode int, ops ...NewAPIErrorOptions) *NewAPIError` — 初始化 OpenAI 错误
- `NewErrorWithStatusCode(err error, errorCode ErrorCode, statusCode int, ops ...NewAPIErrorOptions) *NewAPIError` — 创建带状态码的错误
- `WithOpenAIError(openAIError OpenAIError, statusCode int, ops ...NewAPIErrorOptions) *NewAPIError` — 从 OpenAI 错误创建
- `WithClaudeError(claudeError ClaudeError, statusCode int, ops ...NewAPIErrorOptions) *NewAPIError` — 从 Claude 错误创建

### 错误方法

- `Unwrap() error` — 支持 `errors.Is` / `errors.As`
- `GetErrorCode() ErrorCode` — 获取错误码
- `GetErrorType() ErrorType` — 获取错误类型
- `Error() string` — 错误消息
- `ErrorWithStatusCode() string` — 带状态码的错误消息
- `MaskSensitiveError() string` — 掩码敏感信息后的错误消息
- `MaskSensitiveErrorWithStatusCode() string` — 带状态码的掩码错误消息
- `SetMessage(message string)` — 设置错误消息
- `ToOpenAIError() OpenAIError` — 转换为 OpenAI 错误格式
- `ToClaudeError() ClaudeError` — 转换为 Claude 错误格式

### 错误判断函数

- `IsChannelError(err *NewAPIError) bool` — 判断是否为渠道错误
- `IsSkipRetryError(err *NewAPIError) bool` — 判断是否跳过重试
- `IsRecordErrorLog(e *NewAPIError) bool` — 判断是否记录错误日志

### 选项函数

- `ErrOptionWithSkipRetry() NewAPIErrorOptions` — 跳过重试选项
- `ErrOptionWithNoRecordErrorLog() NewAPIErrorOptions` — 不记录错误日志选项
- `ErrOptionWithStatusCode(statusCode int) NewAPIErrorOptions` — 设置状态码选项
- `ErrOptionWithHideErrMsg(replaceStr string) NewAPIErrorOptions` — 隐藏错误消息选项

## 5. 关键逻辑分析

### 错误类型分类

| ErrorType | 说明 |
|-----------|------|
| `ErrorTypeNewAPIError` | 系统内部错误 |
| `ErrorTypeOpenAIError` | OpenAI 格式错误 |
| `ErrorTypeClaudeError` | Claude 格式错误 |
| `ErrorTypeMidjourneyError` | Midjourney 错误 |
| `ErrorTypeGeminiError` | Gemini 错误 |
| `ErrorTypeRerankError` | Rerank 错误 |
| `ErrorTypeUpstreamError` | 上游服务错误 |

### 错误码分类

| 类别 | 错误码示例 | 说明 |
|------|-----------|------|
| 通用错误 | `invalid_request`, `sensitive_words_detected` | 通用请求错误 |
| 系统错误 | `count_token_failed`, `model_price_error` | 系统内部错误 |
| 渠道错误 | `channel:no_available_key`, `channel:invalid_key` | 渠道相关错误 |
| 请求错误 | `read_request_body_failed`, `access_denied` | 客户端请求错误 |
| 响应错误 | `bad_response_status_code`, `empty_response` | 上游响应错误 |
| 配额错误 | `insufficient_user_quota` | 配额不足错误 |

### 错误转换逻辑

`ToOpenAIError()` 和 `ToClaudeError()` 方法支持不同错误格式之间的转换：
- 保持原始错误类型信息
- 自动掩码敏感信息（通过 `common.MaskSensitiveInfo`）
- 支持 OpenRouter 的 Metadata 附加

### 深层错误保留

`NewError` 和 `NewOpenAIError` 使用 `errors.As` 检查是否已存在 `NewAPIError`，如果存在则保留深层错误信息，避免错误包装丢失。

### 重试控制

- `skipRetry` 标记用于指示某些错误不应重试（如认证错误）
- `recordErrorLog` 控制是否记录错误日志（默认为 `true`）

## 6. 相关文件

- `relay/` — 中继层构造和处理各种错误
- `controller/` — 控制器将错误转换为 HTTP 响应
- `common/crypto.go` — `MaskSensitiveInfo` 函数
- `types/channel_error.go` — 渠道错误结构体
