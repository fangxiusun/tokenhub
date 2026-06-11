# rerank.go 代码阅读文档

## 1. 全局摘要

该文件定义了文本重排序（Rerank）API 的请求和响应数据结构。重排序功能用于根据查询对文档进行相关性排序。

## 2. 依赖

- **标准库**：
  - `fmt`：格式化输出
  - `strings`：字符串操作

- **外部包**：
  - `github.com/QuantumNous/new-api/types`：`TokenCountMeta` 类型
  - `github.com/gin-gonic/gin`：HTTP 上下文

## 3. 类型定义

### RerankRequest 结构体
重排序请求结构：
- `Documents` ([]any)：文档数组
- `Query` (string)：查询文本
- `Model` (string)：模型名称
- `TopN` (*int)：返回前 N 个结果
- `ReturnDocuments` (*bool)：是否返回文档内容
- `MaxChunkPerDoc` (*int)：每文档最大分块数
- `OverLapTokens` (*int)：重叠 token 数

### RerankResponseResult 结构体
重排序结果：
- `Document` (any)：文档内容
- `Index` (int)：原始索引
- `RelevanceScore` (float64)：相关性分数

### RerankDocument 结构体
重排序文档：
- `Text` (any)：文档文本

### RerankResponse 结构体
重排序响应：
- `Results` ([]RerankResponseResult)：结果数组
- `Usage` (Usage)：使用量统计

## 4. 函数详情

### IsStream()
```go
func (r *RerankRequest) IsStream(c *gin.Context) bool
```
**功能**：判断是否为流式请求。

**返回**：始终返回 `false`（重排序不支持流式）。

### GetTokenCountMeta()
```go
func (r *RerankRequest) GetTokenCountMeta() *types.TokenCountMeta
```
**功能**：获取 token 计数元数据。

**逻辑**：将所有文档和查询文本用换行符连接。

### SetModelName()
```go
func (r *RerankRequest) SetModelName(modelName string)
```
**功能**：设置模型名称（非空时更新）。

### GetReturnDocuments()
```go
func (r *RerankRequest) GetReturnDocuments() bool
```
**功能**：获取是否返回文档内容。

**逻辑**：空指针返回 `false`，否则返回指针值。

## 5. 关键逻辑分析

1. **灵活文档格式**：`Documents` 使用 `[]any` 类型，支持多种文档格式。

2. **流式禁用**：重排序 API 不支持流式响应。

3. **零值安全**：使用指针类型处理可选参数，避免空指针 panic。

## 6. 相关文件

- `relay/rerank/`：重排序中继适配器
- `controller/rerank.go`：重排序控制器
- `types/rerank.go`：重排序类型定义