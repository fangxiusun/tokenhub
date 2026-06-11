# openai_responses_compaction_request.go 代码阅读文档

## 1. 全局摘要

该文件定义了 OpenAI Responses API 的压缩请求结构 `OpenAIResponsesCompactionRequest`，用于处理响应压缩功能。

## 2. 依赖

- **标准库**：
  - `encoding/json`：JSON 编解码
  - `strings`：字符串操作

- **外部包**：
  - `github.com/QuantumNous/new-api/types`：`TokenCountMeta` 类型
  - `github.com/gin-gonic/gin`：HTTP 上下文

## 3. 类型定义

### OpenAIResponsesCompactionRequest 结构体
Responses API 压缩请求结构：
- `Model` (string)：模型名称
- `Input` (json.RawMessage)：输入数据（原始 JSON）
- `Instructions` (json.RawMessage)：指令（原始 JSON）
- `PreviousResponseID` (string)：前一个响应 ID

## 4. 函数详情

### GetTokenCountMeta()
```go
func (r *OpenAIResponsesCompactionRequest) GetTokenCountMeta() *types.TokenCountMeta
```
**功能**：获取 token 计数元数据。

**逻辑**：将 `Instructions` 和 `Input` 字段内容用换行符连接。

### IsStream()
```go
func (r *OpenAIResponsesCompactionRequest) IsStream(c *gin.Context) bool
```
**功能**：判断是否为流式请求。

**返回**：始终返回 `false`（压缩请求不支持流式）。

### SetModelName()
```go
func (r *OpenAIResponsesCompactionRequest) SetModelName(modelName string)
```
**功能**：设置模型名称（非空时更新）。

## 5. 关键逻辑分析

1. **原始数据保留**：使用 `json.RawMessage` 保留原始 JSON 数据，支持延迟解析。

2. **流式禁用**：压缩请求不支持流式响应。

3. **Token 计数**：合并指令和输入文本用于 token 统计。

## 6. 相关文件

- `dto/openai_compaction.go`：压缩响应结构
- `dto/openai_request.go`：通用请求结构
- `relay/openai/`：OpenAI 中继适配器