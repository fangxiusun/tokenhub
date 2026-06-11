# embedding.go 代码阅读文档

## 1. 全局摘要

该文件定义了文本嵌入（Embedding）API 的请求和响应数据结构。包含嵌入请求 `EmbeddingRequest`、嵌入选项 `EmbeddingOptions`，以及嵌入响应结构。支持单文本和批量文本输入，提供 token 计数和输入解析功能。

## 2. 依赖

- **标准库**：`strings`（字符串连接）

- **外部包**：
  - `github.com/QuantumNous/new-api/types`：`TokenCountMeta` 类型
  - `github.com/gin-gonic/gin`：HTTP 上下文

## 3. 类型定义

### EmbeddingOptions 结构体
嵌入选项配置，包含：
- `Seed` (int)：随机种子
- `Temperature` (*float64)：温度参数
- `TopK` (int)：Top-K 采样
- `TopP` (*float64)：Top-P 采样
- `FrequencyPenalty` (*float64)：频率惩罚
- `PresencePenalty` (*float64)：存在惩罚
- `NumPredict` (int)：预测 token 数
- `NumCtx` (int)：上下文窗口大小

### EmbeddingRequest 结构体
嵌入请求结构：
- `Model` (string)：模型名称
- `Input` (any)：输入文本（支持字符串或字符串数组）
- `EncodingFormat` (string)：编码格式
- `Dimensions` (*int)：输出维度
- `User` (string)：用户标识
- `Seed` (*float64)：随机种子
- `Temperature` (*float64)：温度参数
- `TopP` (*float64)：Top-P 采样
- `FrequencyPenalty` (*float64)：频率惩罚
- `PresencePenalty` (*float64)：存在惩罚

### EmbeddingResponseItem 结构体
单个嵌入结果：
- `Object` (string)：对象类型
- `Index` (int)：索引
- `Embedding` ([]float64)：嵌入向量

### EmbeddingResponse 结构体
嵌入响应：
- `Object` (string)：对象类型
- `Data` ([]EmbeddingResponseItem)：嵌入结果数组
- `Model` (string)：模型名称
- `Usage` (Usage)：使用量统计

## 4. 函数详情

### GetTokenCountMeta()
```go
func (r *EmbeddingRequest) GetTokenCountMeta() *types.TokenCountMeta
```
**功能**：获取 token 计数元数据。

**逻辑**：
1. 调用 `ParseInput()` 解析输入
2. 将所有文本用换行符连接
3. 返回包含合并文本的元数据

### IsStream()
```go
func (r *EmbeddingRequest) IsStream(c *gin.Context) bool
```
**功能**：判断是否为流式请求。

**返回**：始终返回 `false`（嵌入 API 不支持流式）。

### SetModelName()
```go
func (r *EmbeddingRequest) SetModelName(modelName string)
```
**功能**：设置模型名称（非空时更新）。

### ParseInput()
```go
func (r *EmbeddingRequest) ParseInput() []string
```
**功能**：解析输入字段为字符串数组。

**逻辑**：
1. 检查 `Input` 是否为 `nil`，返回空数组
2. 类型断言处理两种情况：
   - `string`：包装为单元素数组
   - `[]any`：遍历提取字符串元素

## 5. 关键逻辑分析

1. **输入多态性**：`Input` 字段支持字符串或字符串数组，通过 `ParseInput()` 方法统一处理。

2. **流式禁用**：嵌入 API 不支持流式响应，`IsStream()` 始终返回 `false`。

3. **零值处理**：可选参数使用指针类型，确保序列化时能正确处理零值。

4. **Token 计数**：将所有输入文本合并，用于 token 使用量统计。

## 6. 相关文件

- `relay/embedding/`：嵌入中继适配器
- `types/embedding.go`：嵌入相关类型定义
- `controller/embedding.go`：嵌入控制器
- `dto/usage.go`：`Usage` 结构体定义