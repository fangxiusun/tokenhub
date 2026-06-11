# audio.go 代码阅读文档

## 1. 全局摘要

该文件定义了音频处理相关的请求和响应数据传输对象（DTO）。主要包含音频请求结构体 `AudioRequest`、音频响应结构体 `AudioResponse`，以及 Whisper 语音识别的详细 JSON 响应结构体。文件实现了请求参数的序列化/反序列化、流式判断、模型名称设置等辅助方法。

## 2. 依赖

- **标准库**：
  - `encoding/json`：JSON 编解码
  - `strings`：字符串操作（用于模型名称判断）

- **外部包**：
  - `github.com/QuantumNous/new-api/types`：项目内部类型定义，提供 `TokenCountMeta` 等类型
  - `github.com/gin-gonic/gin`：Gin Web 框架，用于 HTTP 上下文处理

## 3. 类型定义

### AudioRequest 结构体
音频请求的核心数据结构，包含以下字段：

**通用字段**：
- `Model` (string)：模型名称
- `Input` (string)：输入文本内容
- `Voice` (string)：语音类型
- `Instructions` (string)：可选指令
- `ResponseFormat` (string)：响应格式
- `Speed` (*float64)：语速（指针类型，支持零值）
- `StreamFormat` (string)：流格式
- `Metadata` (json.RawMessage)：元数据

**vllm-omini 特有字段**（使用 `json.RawMessage` 存储原始 JSON）：
- `TaskType`：任务类型
- `Language`：语言
- `RefAudio`：参考音频
- `RefText`：参考文本
- `XVectorOnlyMode`：仅 X 向量模式
- `MaxNewTokens`：最大新 token 数
- `InitialCodecChunkFrames`：初始编码块帧数

### AudioResponse 结构体
简单的音频响应结构，仅包含 `Text` 字段。

### WhisperVerboseJSONResponse 结构体
Whisper 语音识别详细响应，包含：
- `Task`：任务类型
- `Language`：语言
- `Duration`：持续时间
- `Text`：识别文本
- `Segments`：分段信息数组

### Segment 结构体
Whisper 分段详情，包含：
- `Id`、`Seek`：分段标识
- `Start`、`End`：时间戳
- `Text`：文本内容
- `Tokens`：token 数组
- `Temperature`、`AvgLogprob`、`CompressionRatio`、`NoSpeechProb`：统计指标

## 4. 函数详情

### GetTokenCountMeta()
```go
func (r *AudioRequest) GetTokenCountMeta() *types.TokenCountMeta
```
**功能**：获取请求的 token 计数元数据。

**逻辑**：
1. 创建 `TokenCountMeta` 结构，设置 `CombineText` 为 `Input`
2. 默认 `TokenType` 为 `TokenTypeTextNumber`
3. 如果模型名称包含 "gpt"，则设置为 `TokenTypeTokenizer`
4. 返回元数据指针

### IsStream()
```go
func (r *AudioRequest) IsStream(c *gin.Context) bool
```
**功能**：判断请求是否为流式请求。

**逻辑**：当 `StreamFormat` 字段等于 "sse" 时返回 `true`。

### SetModelName()
```go
func (r *AudioRequest) SetModelName(modelName string)
```
**功能**：设置模型名称。

**逻辑**：仅当 `modelName` 不为空字符串时更新 `Model` 字段。

## 5. 关键逻辑分析

1. **Token 计数策略**：根据模型名称自动选择不同的 token 计数方式。GPT 系列模型使用分词器计数，其他模型使用文本数字计数。

2. **流式判断机制**：通过 `StreamFormat` 字段判断是否使用 Server-Sent Events (SSE) 流式响应。

3. **vllm-omini 兼容性**：使用 `json.RawMessage` 类型存储原始 JSON 数据，保持与上游 API 的兼容性，避免在 DTO 层进行数据转换。

4. **零值处理**：`Speed` 字段使用指针类型，确保在序列化时能正确处理零值（如语速为 0.0）。

## 6. 相关文件

- `types/token.go`：包含 `TokenCountMeta` 和 `TokenType` 的定义
- `relay/audio/`：音频中继处理逻辑，使用这些 DTO 进行请求/响应转换
- `controller/audio.go`：音频相关的 HTTP 处理器
- `middleware/request.go`：请求处理中间件，可能调用 `GetTokenCountMeta()`