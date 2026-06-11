# openai_response.go 代码阅读文档

## 1. 全局摘要

该文件是 OpenAI 兼容 API 的核心响应 DTO 定义文件，实现了多种响应结构体，包括文本响应、流式响应、嵌入响应、Responses API 响应等。包含使用量统计结构 `Usage`、工具调用响应 `ToolCallResponse`，以及错误提取函数 `GetOpenAIError`。

## 2. 依赖

- **标准库**：
  - `encoding/json`：JSON 编解码
  - `fmt`：格式化输出

- **外部包**：
  - `github.com/QuantumNous/new-api/common`：JSON 工具函数
  - `github.com/QuantumNous/new-api/types`：`OpenAIError` 类型

## 3. 类型定义

### 基础响应结构体

**SimpleResponse**：简单响应：
- `Usage`：使用量统计
- `Error`：错误信息

**TextResponse**：文本响应：
- `Id`：响应 ID
- `Object`：对象类型
- `Created`：创建时间
- `Model`：模型名称
- `Choices`：选择数组
- `Usage`：使用量统计

**OpenAITextResponse**：OpenAI 文本响应：
- 同 `TextResponse`，增加 `Error` 字段

### 嵌入响应结构体

**OpenAIEmbeddingResponse**：嵌入响应：
- `Object`：对象类型
- `Data`：嵌入数据数组
- `Model`：模型名称
- `Usage`：使用量统计

**FlexibleEmbeddingResponse**：灵活嵌入响应（支持不同类型的嵌入数据）

### 流式响应结构体

**ChatCompletionsStreamResponse**：流式响应：
- `Id`：响应 ID
- `Object`：对象类型
- `Created`：创建时间
- `Model`：模型名称
- `SystemFingerprint`：系统指纹
- `Choices`：选择数组
- `Usage`：使用量统计

**ChatCompletionsStreamResponseChoice**：流式选择：
- `Delta`：增量内容
- `Logprobs`：日志概率
- `FinishReason`：完成原因
- `Index`：索引

**ChatCompletionsStreamResponseChoiceDelta**：流式增量：
- `Content`：内容
- `ReasoningContent`：推理内容
- `Role`：角色
- `ToolCalls`：工具调用

### 工具调用响应结构体

**ToolCallResponse**：工具调用响应：
- `Index`：索引（仅在流式响应中非空）
- `ID`：调用 ID
- `Type`：类型
- `Function`：函数响应

**FunctionResponse**：函数响应：
- `Description`：描述
- `Name`：函数名
- `Parameters`：参数定义
- `Arguments`：参数值

### 使用量统计结构体

**Usage**：使用量统计：
- `PromptTokens`：提示 token 数
- `CompletionTokens`：完成 token 数
- `TotalTokens`：总 token 数
- `PromptCacheHitTokens`：缓存命中 token 数
- `InputTokens`：输入 token 数（Claude 格式）
- `OutputTokens`：输出 token 数（Claude 格式）
- `ClaudeCacheCreation5mTokens`：5 分钟缓存创建 token 数
- `ClaudeCacheCreation1hTokens`：1 小时缓存创建 token 数
- `PromptTokensDetails`：提示 token 详情
- `CompletionTokenDetails`：完成 token 详情

**InputTokenDetails**：输入 token 详情：
- `CachedTokens`：缓存 token 数
- `TextTokens`：文本 token 数
- `AudioTokens`：音频 token 数
- `ImageTokens`：图像 token 数

**OutputTokenDetails**：输出 token 详情：
- `TextTokens`：文本 token 数
- `AudioTokens`：音频 token 数
- `ImageTokens`：图像 token 数
- `ReasoningTokens`：推理 token 数

### Responses API 响应结构体

**OpenAIResponsesResponse**：Responses API 响应：
- `ID`：响应 ID
- `Object`：对象类型
- `Status`：状态
- `Output`：输出数组
- `Usage`：使用量统计
- `Error`：错误信息

**ResponsesOutput**：输出内容：
- `Type`：类型
- `ID`：输出 ID
- `Status`：状态
- `Role`：角色
- `Content`：内容数组
- `Quality`：质量
- `Size`：尺寸
- `CallId`：调用 ID
- `Name`：名称
- `Arguments`：参数

**ResponsesOutputContent**：输出内容详情：
- `Type`：类型
- `Text`：文本
- `Annotations`：注释数组

**ResponsesStreamResponse**：Responses 流式响应：
- `Type`：响应类型
- `Response`：完整响应
- `Delta`：增量
- `Item`：输出项
- `OutputIndex`：输出索引
- `ContentIndex`：内容索引

### 其他结构体

**OpenAIVideoResponse**：视频响应：
- `Id`：文件 ID
- `Object`：对象类型
- `Bytes`：文件大小
- `CreatedAt`：创建时间
- `ExpiresAt`：过期时间
- `Filename`：文件名
- `Purpose`：用途

### 常量

**Responses 输出类型**：
- `ResponsesOutputTypeImageGenerationCall`："image_generation_call"

**内置工具类型**：
- `BuildInToolWebSearchPreview`："web_search_preview"
- `BuildInToolFileSearch`："file_search"

**内置调用类型**：
- `BuildInCallWebSearchCall`："web_search_call"

**Responses 输出事件类型**：
- `ResponsesOutputTypeItemAdded`："response.output_item.added"
- `ResponsesOutputTypeItemDone`："response.output_item.done"

## 4. 函数详情

### ChatCompletionsStreamResponse 方法

**IsFinished() bool**：判断流式响应是否完成。

**IsToolCall() bool**：判断是否包含工具调用。

**GetFirstToolCall() *ToolCallResponse**：获取第一个工具调用。

**ClearToolCalls()**：清空工具调用数据。

**Copy() *ChatCompletionsStreamResponse**：复制响应对象。

**GetSystemFingerprint() string**：获取系统指纹。

**SetSystemFingerprint(s string)**：设置系统指纹。

### ChatCompletionsStreamResponseChoiceDelta 方法

**SetContentString(s string)**：设置内容字符串。

**GetContentString() string**：获取内容字符串。

**GetReasoningContent() string**：获取推理内容。

**SetReasoningContent(s string)**：设置推理内容。

### ToolCallResponse 方法

**SetIndex(i int)**：设置索引。

### OpenAIResponsesResponse 方法

**GetOpenAIError() *types.OpenAIError**：提取 OpenAI 错误。

**HasImageGenerationCall() bool**：判断是否包含图像生成调用。

**GetQuality() string**：获取图像生成质量。

**GetSize() string**：获取图像生成尺寸。

### ResponsesOutput 方法

**ArgumentsString() string**：获取函数调用参数字符串。

### 全局函数

**GetOpenAIError(errorField any) *types.OpenAIError**：从动态错误类型中提取 OpenAI 错误结构，支持：
- `types.OpenAIError` / `*types.OpenAIError`
- `map[string]interface{}`
- `string`
- 其他未知类型

**ResponsesArgumentsString(arguments json.RawMessage) string**：将 Responses API 参数转换为字符串格式。

## 5. 关键逻辑分析

1. **使用量统计兼容**：`Usage` 结构体同时支持 OpenAI 和 Claude 的使用量格式，通过字段名映射实现兼容。

2. **错误类型提取**：`GetOpenAIError` 函数处理多种错误格式，确保错误信息可解析。

3. **流式响应处理**：`ChatCompletionsStreamResponse` 提供丰富的流式响应处理方法。

4. **Responses API 支持**：专门的 `OpenAIResponsesResponse` 结构体支持 OpenAI Responses API 的独特格式。

5. **工具调用管理**：提供工具调用的清空、复制、判断等方法，便于流式响应处理。

## 6. 相关文件

- `dto/openai_request.go`：OpenAI 请求结构定义
- `types/error.go`：`OpenAIError` 类型定义
- `relay/openai/`：OpenAI 中继适配器
- `dto/usage.go`：使用量相关结构