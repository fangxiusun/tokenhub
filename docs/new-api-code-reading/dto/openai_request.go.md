# openai_request.go 代码阅读文档

## 1. 全局摘要

该文件是 OpenAI 兼容 API 的核心请求 DTO 定义文件，实现了通用聊天请求 `GeneralOpenAIRequest`、消息结构 `Message`、媒体内容 `MediaContent`，以及 Responses API 的 `OpenAIResponsesRequest`。是整个 API 网关中最重要的请求结构定义文件，支持多种 AI 服务提供商的参数格式。

## 2. 依赖

- **标准库**：
  - `encoding/json`：JSON 编解码
  - `fmt`：格式化输出
  - `strings`：字符串操作

- **外部包**：
  - `github.com/QuantumNous/new-api/common`：JSON 工具函数
  - `github.com/QuantumNous/new-api/types`：类型定义
  - `github.com/samber/lo`：Go 泛型工具库
  - `github.com/gin-gonic/gin`：HTTP 上下文

## 3. 类型定义

### 核心请求结构体

**GeneralOpenAIRequest**：通用 OpenAI 请求结构，包含约 70 个字段：
- **基础参数**：`Model`、`Messages`、`Prompt`、`Stream`
- **生成参数**：`MaxTokens`、`MaxCompletionTokens`、`Temperature`、`TopP`、`TopK`
- **工具参数**：`Tools`、`ToolChoice`、`FunctionCall`
- **流式参数**：`StreamOptions`
- **响应格式**：`ResponseFormat`
- **安全/隐私**：`SafetyIdentifier`、`Store`
- **提供商特有参数**：`ExtraBody`（Gemini）、`SearchParameters`（xAI）、`WebSearchOptions`（Claude）等

**OpenAIResponsesRequest**：Responses API 请求结构：
- `Model`：模型名称
- `Input`：输入数据（原始 JSON）
- `Instructions`：指令
- `Tools`：工具定义
- `Stream`/`StreamOptions`：流式配置
- `Reasoning`：推理配置

### 消息相关结构体

**Message**：消息结构：
- `Role`：角色
- `Content`：内容（支持字符串或数组）
- `Name`：发送者名称
- `ReasoningContent`/`Reasoning`：推理内容
- `ToolCalls`：工具调用（原始 JSON）
- `ToolCallId`：工具调用 ID

**MediaContent**：多媒体内容：
- `Type`：内容类型
- `Text`：文本内容
- `ImageUrl`：图像 URL
- `InputAudio`：音频输入
- `File`：文件数据
- `VideoUrl`：视频 URL
- `CacheControl`：缓存控制

### 工具相关结构体

**ToolCallRequest**：工具调用请求：
- `ID`：调用 ID
- `Type`：类型
- `Function`：函数定义
- `Custom`：自定义数据

**FunctionRequest**：函数定义：
- `Name`：函数名
- `Description`：描述
- `Parameters`：参数定义
- `Arguments`：参数值

**StreamOptions**：流式选项：
- `IncludeUsage`：包含使用量统计
- `IncludeObfuscation`：包含混淆数据

### 媒体类型结构体

**MessageImageUrl**：图像 URL：
- `Url`：URL 地址
- `Detail`：详细程度
- `MimeType`：MIME 类型

**MessageInputAudio**：音频输入：
- `Data`：Base64 编码数据
- `Format`：音频格式

**MessageFile**：文件数据：
- `FileName`：文件名
- `FileData`：文件数据
- `FileId`：文件 ID

**MessageVideoUrl**：视频 URL：
- `Url`：URL 地址

**MediaInput**：媒体输入（Responses API）：
- `Type`：类型
- `Text`：文本
- `ImageUrl`：图像 URL
- `FileUrl`：文件 URL
- `Detail`：详细程度

### 配置结构体

**ResponseFormat**：响应格式：
- `Type`：类型
- `JsonSchema`：JSON 模式

**FormatJsonSchema**：格式化 JSON 模式：
- `Description`：描述
- `Name`：名称
- `Schema`：模式定义
- `Strict`：严格模式

**Reasoning**：推理配置：
- `Effort`：推理努力程度
- `Summary`：摘要

### 常量

**内容类型常量**：
- `ContentTypeText`："text"
- `ContentTypeImageURL`："image_url"
- `ContentTypeInputAudio`："input_audio"
- `ContentTypeFile`："file"
- `ContentTypeVideoUrl`："video_url"

## 4. 函数详情

### GeneralOpenAIRequest 方法

**GetTokenCountMeta() *types.TokenCountMeta**：获取 token 计数元数据：
1. 提取提示词文本
2. 遍历消息提取文本和文件
3. 处理工具定义
4. 返回合并文本和文件元数据

**IsStream(c *gin.Context) bool**：判断是否为流式请求。

**SetModelName(modelName string)**：设置模型名称。

**ToMap() map[string]any**：转换为 map 类型。

**GetSystemRoleName() string**：获取系统角色名称（o1/o3/gpt-5 使用 "developer"）。

**GetMaxTokens() uint**：获取最大 token 数（优先使用 `MaxCompletionTokens`）。

**ParseInput() []string**：解析输入字段为字符串数组。

### Message 方法

**GetReasoningContent() string**：获取推理内容。

**GetPrefix() bool**：获取前缀标识。

**SetPrefix(prefix bool)**：设置前缀标识。

**ParseToolCalls() []ToolCallRequest**：解析工具调用。

**SetToolCalls(toolCalls any)**：设置工具调用。

**StringContent() string**：获取字符串内容。

**SetNullContent()**：设置空内容。

**SetStringContent(content string)**：设置字符串内容。

**SetMediaContent(content []MediaContent)**：设置多媒体内容。

**IsStringContent() bool**：判断是否为字符串内容。

**ParseContent() []MediaContent**：解析内容为多媒体数组：
1. 尝试解析为字符串
2. 尝试解析为数组
3. 根据类型创建对应的 MediaContent

### MediaContent 方法

**GetImageMedia() *MessageImageUrl**：获取图像媒体。

**GetInputAudio() *MessageInputAudio**：获取音频输入。

**GetFile() *MessageFile**：获取文件数据。

**GetVideoUrl() *MessageVideoUrl**：获取视频 URL。

**ToFileSource() types.FileSource**：转换为文件源。

### OpenAIResponsesRequest 方法

**GetTokenCountMeta() *types.TokenCountMeta**：获取 Responses API 的 token 计数。

**IsStream(c *gin.Context) bool**：判断是否为流式请求。

**SetModelName(modelName string)**：设置模型名称。

**GetToolsMap() []map[string]any**：获取工具列表的 map 表示。

**ParseInput() []MediaInput**：解析输入字段为媒体输入数组。

### 全局函数

**IsOpenAIReasoningOModel(modelName string) bool**：判断是否为推理模型（o1/o3/o4）。

**IsOpenAIGPT5Model(modelName string) bool**：判断是否为 GPT-5 模型。

## 5. 关键逻辑分析

1. **内容多态处理**：`Content` 字段支持字符串或结构化数组，通过 `ParseContent()` 方法统一处理。

2. **文件源提取**：`MediaContent.ToFileSource()` 方法根据内容类型提取文件数据，支持图像、音频、文件、视频。

3. **模型角色适配**：`GetSystemRoleName()` 方法根据模型类型返回不同的系统角色名称。

4. **Token 计数策略**：`GetTokenCountMeta()` 方法遍历所有内容，提取文本和文件元数据用于计费。

5. **Responses API 支持**：`OpenAIResponsesRequest` 专门用于 OpenAI Responses API，支持更丰富的输入格式。

6. **缓存控制**：支持 OpenRouter 等提供商的缓存控制参数。

## 6. 相关文件

- `relay/openai/`：OpenAI 中继适配器
- `types/message.go`：消息类型定义
- `common/json.go`：JSON 工具函数
- `dto/openai_response.go`：OpenAI 响应结构
- `controller/chat.go`：聊天控制器