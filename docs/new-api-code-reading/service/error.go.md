# error.go 代码阅读文档

## 1. 全局总结

该文件提供统一的错误处理和包装功能，包括 Midjourney、Claude、OpenAI、Task 等格式的错误包装器，以及上游响应错误的解析和状态码重映射。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | JSON 解析、敏感信息掩码 |
| `dto` | 错误响应结构体 |
| `types` | 错误类型定义 |
| `logger` | 日志记录 |

## 3. 类型定义

无自定义类型，使用 `dto` 和 `types` 包中的类型。

## 4. 函数详解

### `MidjourneyErrorWrapper(code, desc) *dto.MidjourneyResponse`
Midjourney 错误包装

### `ClaudeErrorWrapper(err, code, statusCode) *dto.ClaudeErrorWithStatusCode`
Claude 错误包装，隐藏上游详细错误信息

### `RelayErrorHandler(ctx, resp, showBodyWhenFail) *types.NewAPIError`
上游响应错误解析：
1. 读取响应体
2. 尝试解析为通用错误格式
3. 尝试转换为 OpenAI 错误格式
4. 构建 NewAPIError

### `ResetStatusCode(newApiErr, statusCodeMappingStr)`
状态码重映射，支持 JSON 配置的映射规则

### `TaskErrorWrapper(err, code, statusCode) *dto.TaskError`
任务错误包装

### `TaskErrorFromAPIError(apiErr) *dto.TaskError`
将 NewAPIError 转换为 TaskError

## 5. 关键逻辑分析

1. **错误隐藏**：上游详细错误信息被掩码处理
2. **状态码重映射**：支持灵活的状态码转换配置
3. **多格式兼容**：支持 OpenAI、Claude、Gemini 等多种错误格式

## 6. 关联文件

- `types/error.go` — 错误类型定义
- `dto/error.go` — 错误响应结构体
