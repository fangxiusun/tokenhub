# claude.go 代码阅读文档

## 1. 全局摘要

该文件是 Claude API 的核心 DTO 定义文件，实现了 Claude 请求/响应的完整数据结构。包含消息结构、工具定义、思考配置、使用量统计等，是 Claude 中继适配器的基础。文件定义了约 20 个结构体和 30+ 个方法，处理 Claude API 的各种数据类型转换。

## 2. 依赖

- **标准库**：
  - `encoding/json`：JSON 编解码
  - `fmt`：格式化输出
  - `strings`：字符串操作

- **外部包**：
  - `github.com/QuantumNous/new-api/common`：JSON 工具函数（`Marshal`、`Any2Type`）
  - `github.com/QuantumNous/new-api/types`：类型定义（`TokenCountMeta`、`FileSource`、`ClaudeError`）
  - `github.com/gin-gonic/gin`：HTTP 上下文

## 3. 类型定义

### 消息相关结构体

**ClaudeMetadata**：元数据结构，包含 `UserId` 字段。

**ClaudeMediaMessage**：多媒体消息结构，支持：
- 文本消息（`Text`、`Thinking`）
- 图片消息（`Source` → `ClaudeMessageSource`）
- 工具调用（`Id`、`Name`、`Input`）
- 工具结果（`ToolUseId`、`Content`）
- 缓存控制（`CacheControl`）

**ClaudeMessageSource**：消息源结构：
- `Type`：类型标识
- `MediaType`：媒体类型
- `Data`：内联数据
- `Url`：外部 URL

**ClaudeMessage**：标准消息结构：
- `Role`：角色（user/assistant）
- `Content`：内容（支持字符串或 `[]ClaudeMediaMessage`）

### 工具相关结构体

**Tool**：标准工具定义：
- `Name`：工具名称
- `Description`：描述
- `InputSchema`：输入模式

**InputSchema**：工具输入模式：
- `Type`：类型
- `Properties`：属性
- `Required`：必填字段

**ClaudeWebSearchTool**：网页搜索工具：
- `MaxUses`：最大使用次数
- `UserLocation`：用户位置信息

**ClaudeWebSearchUserLocation**：用户位置结构：
- `Timezone`：时区
- `Country`：国家
- `Region`：地区
- `City`：城市

**ClaudeToolChoice**：工具选择配置：
- `Type`：选择类型
- `Name`：指定工具名称
- `DisableParallelToolUse`：禁用并行工具调用

### 请求/响应结构体

**ClaudeRequest**：Claude 请求结构，包含：
- 基础参数：`Model`、`MaxTokens`、`Temperature`、`TopP`、`TopK`
- 消息：`System`、`Messages`
- 工具：`Tools`、`ToolChoice`
- 高级功能：`Thinking`、`ContextManagement`、`McpServers`
- 透传控制：`InferenceGeo`、`Speed`、`ServiceTier`

**ClaudeResponse**：Claude 响应结构：
- `Content`：响应内容
- `Usage`：使用量统计
- `StopReason`：停止原因
- `Error`：错误信息

### 其他结构体

**Thinking**：思考配置：
- `Type`：类型
- `BudgetTokens`：预算 token 数
- `Display`：显示模式

**ClaudeUsage**：使用量统计：
- `InputTokens`：输入 token
- `OutputTokens`：输出 token
- `CacheCreationInputTokens`：缓存创建输入 token
- `CacheReadInputTokens`：缓存读取输入 token
- `ClaudeCacheCreation5mTokens`：5分钟缓存创建 token
- `ClaudeCacheCreation1hTokens`：1小时缓存创建 token

**ClaudeCacheCreationUsage**：缓存创建使用量：
- `Ephemeral5mInputTokens`：5分钟临时 token
- `Ephemeral1hInputTokens`：1小时临时 token

**ClaudeServerToolUse**：服务器工具使用统计：
- `WebSearchRequests`：网页搜索请求数

## 4. 函数详情

### ClaudeMediaMessage 方法

**SetText(s string)**：设置文本内容。

**GetText() string**：获取文本内容，空指针返回空字符串。

**IsStringContent() bool**：判断 `Content` 是否为字符串类型。

**GetStringContent() string**：获取字符串内容，支持处理 `[]any` 类型的多段文本。

**GetJsonRowString() string**：获取 JSON 序列化字符串。

**SetContent(content any)**：设置任意类型的 `Content`。

**ParseMediaContent() []ClaudeMediaMessage**：解析 `Content` 为多媒体消息数组。

**ToFileSource() types.FileSource**：将源数据转换为 `FileSource`。

### ClaudeMessage 方法

**IsStringContent() bool**：判断 `Content` 是否为字符串。

**GetStringContent() string**：获取字符串内容，支持多段文本拼接。

**SetStringContent(content string)**：设置字符串内容。

**SetContent(content any)**：设置任意类型的 `Content`。

**ParseContent() ([]ClaudeMediaMessage, error)**：解析 `Content` 为多媒体消息数组。

### ClaudeRequest 方法

**GetTokenCountMeta() *types.TokenCountMeta**：获取 token 计数元数据，包含：
1. 提取系统提示词文本
2. 遍历消息提取文本和图片
3. 处理工具定义文本
4. 统计文件元数据

**IsStream(ctx *gin.Context) bool**：判断是否为流式请求。

**SetModelName(modelName string)**：设置模型名称。

**SearchToolNameByToolCallId(toolCallId string) string**：根据工具调用 ID 搜索工具名称。

**AddTool(tool any)**：添加工具到请求。

**GetTools() []any**：获取工具列表。

**GetEfforts() string**：获取输出配置中的 effort 参数。

**IsStringSystem() bool**：判断 `System` 是否为字符串。

**GetStringSystem() string**：获取系统提示词字符串。

**SetStringSystem(system string)**：设置系统提示词。

**ParseSystem() []ClaudeMediaMessage**：解析系统提示词为多媒体消息。

### ClaudeResponse 方法

**SetIndex(i int)**：设置索引。

**GetIndex() int**：获取索引。

**GetClaudeError() *types.ClaudeError**：提取 Claude 错误结构，支持多种错误类型：
- `types.ClaudeError` / `*types.ClaudeError`
- `map[string]interface{}`
- `string`
- 其他未知类型

### Thinking 方法

**GetBudgetTokens() int**：获取预算 token 数。

### ClaudeUsage 方法

**GetCacheCreation5mTokens() int**：获取 5 分钟缓存创建 token 数。

**GetCacheCreation1hTokens() int**：获取 1 小时缓存创建 token 数。

**GetCacheCreationTotalTokens() int**：获取缓存创建总 token 数（兼容新旧格式）。

### 全局函数

**ProcessTools(tools []any) ([]*Tool, []*ClaudeWebSearchTool)**：处理工具列表，分离普通工具和网页搜索工具。

## 5. 关键逻辑分析

1. **内容类型多态**：`Content` 和 `System` 字段支持字符串或结构化数组，通过类型断言动态处理。

2. **Token 计数策略**：`GetTokenCountMeta` 方法遍历所有消息和工具，提取文本内容，同时收集文件元数据用于图片计费。

3. **缓存使用量统计**：支持两种缓存格式（旧版 `CacheCreationInputTokens` 和新版 `ClaudeCacheCreation5mTokens/1hTokens`），确保兼容性。

4. **错误类型提取**：`GetClaudeError` 方法处理多种错误格式，从 JSON map 到结构化错误，确保错误信息可解析。

5. **工具类型分离**：`ProcessTools` 将普通工具和网页搜索工具分开处理，便于不同的计费和处理逻辑。

6. **零值安全**：所有指针类型字段都有空值检查，避免空指针 panic。

## 6. 相关文件

- `relay/claude/`：Claude 中继适配器实现
- `types/claude.go`：Claude 相关类型定义（`ClaudeError`）
- `common/json.go`：JSON 工具函数
- `dto/openai_request.go`：OpenAI 请求 DTO（可能用于格式转换）
- `controller/claude.go`：Claude 控制器