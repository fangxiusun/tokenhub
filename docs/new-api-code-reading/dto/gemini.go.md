# gemini.go 代码阅读文档

## 1. 全局摘要

该文件是 Google Gemini API 的核心 DTO 定义文件，实现了 Gemini 聊天、图像生成、嵌入等 API 的完整数据结构。包含请求/响应结构体、工具配置、安全设置、使用量统计等。文件支持 snake_case 和 camelCase 两种 JSON 命名风格的自动转换。

## 2. 依赖

- **标准库**：
  - `encoding/json`：JSON 编解码
  - `strings`：字符串操作

- **外部包**：
  - `github.com/QuantumNous/new-api/common`：JSON 工具函数
  - `github.com/QuantumNous/new-api/logger`：日志记录
  - `github.com/QuantumNous/new-api/types`：类型定义
  - `github.com/gin-gonic/gin`：HTTP 上下文

## 3. 类型定义

### 核心请求/响应结构体

**GeminiChatRequest**：Gemini 聊天请求：
- `Requests`：批量请求支持
- `Contents`：对话内容数组
- `SafetySettings`：安全设置
- `GenerationConfig`：生成配置
- `Tools`：工具定义（原始 JSON）
- `ToolConfig`：工具配置
- `SystemInstructions`：系统指令
- `CachedContent`：缓存内容标识

**GeminiChatResponse**：Gemini 聊天响应：
- `Candidates`：候选响应数组
- `PromptFeedback`：提示反馈
- `UsageMetadata`：使用量元数据

### 内容相关结构体

**GeminiChatContent**：对话内容：
- `Role`：角色
- `Parts`：内容部分数组

**GeminiPart**：内容部分：
- `Text`：文本内容
- `Thought`：思考标识
- `InlineData`：内联数据
- `FunctionCall`：函数调用
- `FunctionResponse`：函数响应
- `FileData`：文件数据
- `ExecutableCode`：可执行代码
- `CodeExecutionResult`：代码执行结果

**GeminiInlineData**：内联数据：
- `MimeType`：MIME 类型
- `Data`：Base64 编码数据

### 配置相关结构体

**GeminiChatGenerationConfig**：生成配置：
- `Temperature`、`TopP`、`TopK`：采样参数
- `MaxOutputTokens`：最大输出 token
- `StopSequences`：停止序列
- `ResponseMimeType`：响应 MIME 类型
- `ResponseSchema`/`ResponseJsonSchema`：响应模式
- `ThinkingConfig`：思考配置
- `SpeechConfig`/`ImageConfig`：语音/图像配置

**GeminiThinkingConfig**：思考配置：
- `IncludeThoughts`：包含思考内容
- `ThinkingBudget`：思考预算
- `ThinkingLevel`：思考级别

**ToolConfig**：工具配置：
- `FunctionCallingConfig`：函数调用配置
- `RetrievalConfig`：检索配置
- `IncludeServerSideToolInvocations`：包含服务器端工具调用

### 工具相关结构体

**GeminiChatTool**：工具定义：
- `GoogleSearch`：Google 搜索
- `GoogleSearchRetrieval`：Google 搜索检索
- `CodeExecution`：代码执行
- `FunctionDeclarations`：函数声明
- `URLContext`：URL 上下文

**FunctionCall**：函数调用：
- `FunctionName`：函数名
- `Arguments`：参数

**GeminiFunctionResponse**：函数响应：
- `Name`：函数名
- `Response`：响应数据
- `WillContinue`：是否继续
- `Scheduling`：调度信息

### 图像生成相关结构体

**GeminiImageRequest**：图像生成请求：
- `Instances`：实例数组
- `Parameters`：参数配置

**GeminiImageResponse**：图像生成响应：
- `Predictions`：预测结果数组

**GeminiImagePrediction**：图像预测：
- `MimeType`：MIME 类型
- `BytesBase64Encoded`：Base64 编码图像数据
- `SafetyAttributes`：安全属性

### 嵌入相关结构体

**GeminiEmbeddingRequest**：嵌入请求：
- `Model`：模型名称
- `Content`：内容
- `TaskType`：任务类型
- `OutputDimensionality`：输出维度

**GeminiBatchEmbeddingRequest**：批量嵌入请求：
- `Requests`：请求数组

**GeminiEmbeddingResponse**：嵌入响应：
- `Embedding`：嵌入向量

**GeminiBatchEmbeddingResponse**：批量嵌入响应：
- `Embeddings`：嵌入向量数组

### 使用量统计结构体

**GeminiUsageMetadata**：使用量元数据：
- `PromptTokenCount`：提示 token 数
- `CandidatesTokenCount`：候选 token 数
- `TotalTokenCount`：总 token 数
- `ThoughtsTokenCount`：思考 token 数
- `CachedContentTokenCount`：缓存内容 token 数

## 4. 函数详情

### GeminiChatRequest 方法

**UnmarshalJSON(data []byte) error**：自定义 JSON 反序列化，支持 snake_case 和 camelCase 字段。

**GetTokenCountMeta() *types.TokenCountMeta**：获取 token 计数元数据：
1. 提取所有文本内容
2. 识别内联数据（图片/音频/视频）
3. 返回合并文本和文件元数据

**IsStream(c *gin.Context) bool**：判断是否为流式请求：
- 检查查询参数 `alt=sse`
- 检查 URL 路径是否包含 `streamGenerateContent`

**SetModelName(modelName string)**：空操作（Gemini 请求无模型字段）。

**GetTools() []GeminiChatTool**：解析工具列表（支持数组和单对象）。

**SetTools(tools []GeminiChatTool)**：设置工具列表。

### GeminiThinkingConfig 方法

**UnmarshalJSON(data []byte) error**：自定义反序列化，支持 snake_case 字段。

**SetThinkingBudget(budget int)**：设置思考预算。

### GeminiInlineData 方法

**ToFileSource() types.FileSource**：转换为文件源。

**UnmarshalJSON(data []byte) error**：自定义反序列化，支持 `mime_type` 和 `mimeType`。

### GeminiPart 方法

**UnmarshalJSON(data []byte) error**：自定义反序列化，支持 `inline_data` 和 `inlineData`。

### GeminiChatGenerationConfig 方法

**UnmarshalJSON(data []byte) error**：自定义反序列化，支持所有 snake_case 字段。

### GeminiEmbeddingRequest 方法

**IsStream(c *gin.Context) bool**：始终返回 `false`。

**GetTokenCountMeta() *types.TokenCountMeta**：获取文本内容的 token 计数。

**SetModelName(modelName string)**：设置模型名称。

### GeminiBatchEmbeddingRequest 方法

**IsStream(c *gin.Context) bool**：始终返回 `false`。

**GetTokenCountMeta() *types.TokenCountMeta**：合并所有请求的文本。

**SetModelName(modelName string)**：批量设置模型名称。

## 5. 关键逻辑分析

1. **JSON 命名风格兼容**：所有主要结构体实现自定义 `UnmarshalJSON` 方法，优先解析 snake_case 字段，回退到 camelCase。

2. **流式请求识别**：支持两种流式判断方式：查询参数 `alt=sse` 和 URL 路径包含 `streamGenerateContent`。

3. **文件类型识别**：根据 MIME 类型前缀自动识别文件类型（image/、audio/、video/）。

4. **工具解析灵活性**：`GetTools()` 方法支持工具列表为数组或单对象两种格式。

5. **批量请求支持**：`GeminiBatchEmbeddingRequest` 支持批量嵌入请求的合并处理。

6. **零值安全**：可选参数使用指针类型，确保序列化时能正确处理零值。

## 6. 相关文件

- `relay/gemini/`：Gemini 中继适配器实现
- `types/gemini.go`：Gemini 相关类型定义
- `common/json.go`：JSON 工具函数
- `dto/usage.go`：`Usage` 结构体定义
- `controller/gemini.go`：Gemini 控制器