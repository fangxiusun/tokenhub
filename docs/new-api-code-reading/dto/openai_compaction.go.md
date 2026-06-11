# openai_compaction.go 代码阅读文档

## 1. 全局摘要

该文件定义了 OpenAI Responses API 的压缩响应结构。包含 `OpenAIResponsesCompactionResponse` 结构体，用于处理 OpenAI 的响应压缩功能。

## 2. 依赖

- **标准库**：`encoding/json`

- **外部包**：
  - `github.com/QuantumNous/new-api/types`：`OpenAIError` 类型

## 3. 类型定义

### OpenAIResponsesCompactionResponse 结构体
OpenAI Responses 压缩响应结构：
- `ID` (string)：响应 ID
- `Object` (string)：对象类型
- `CreatedAt` (int)：创建时间戳
- `Output` (json.RawMessage)：原始输出数据
- `Usage` (*Usage)：使用量统计
- `Error` (any)：错误信息（支持多种类型）

## 4. 函数详情

### GetOpenAIError()
```go
func (o *OpenAIResponsesCompactionResponse) GetOpenAIError() *types.OpenAIError
```
**功能**：提取 OpenAI 错误结构。

**逻辑**：调用全局 `GetOpenAIError` 函数处理动态错误类型。

## 5. 关键逻辑分析

1. **原始数据保留**：`Output` 使用 `json.RawMessage` 保留原始 JSON 数据，支持延迟解析。

2. **错误类型兼容**：`Error` 使用 `any` 类型，通过 `GetOpenAIError` 函数统一处理多种错误格式。

## 6. 相关文件

- `dto/openai_response.go`：OpenAI 响应结构定义
- `relay/openai/`：OpenAI 中继适配器
- `dto/error.go`：错误处理结构