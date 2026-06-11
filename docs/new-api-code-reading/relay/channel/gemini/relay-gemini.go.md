# relay-gemini.go 代码阅读文档

## 1. 全局总结
这是 Gemini 渠道最核心的文件（1808 行），实现了 OpenAI 格式到 Gemini 格式的完整请求转换、Gemini 响应到 OpenAI 格式的转换，以及流式/非流式、聊天/嵌入/图片生成的响应处理。包含 thinking 适配器、函数调用转换、schema 清理、MIME 类型校验、token 使用量计算等关键逻辑。

## 2. 依赖关系
- **标准库**: `context`, `encoding/json`, `errors`, `fmt`, `io`, `net/http`, `strconv`, `strings`, `time`, `unicode/utf8`
- **项目内部**:
  - `github.com/QuantumNous/new-api/common` — JSON 包装、UUID、指针工具
  - `github.com/QuantumNous/new-api/constant` — 常量定义
  - `github.com/QuantumNous/new-api/dto` — 数据传输对象
  - `github.com/QuantumNous/new-api/logger` — 日志
  - `github.com/QuantumNous/new-api/relay/channel/openai` — OpenAI 流式处理
  - `github.com/QuantumNous/new-api/relay/common` — RelayInfo
  - `github.com/QuantumNous/new-api/relay/helper` — 流扫描、响应 ID、空响应生成
  - `github.com/QuantumNous/new-api/service` — 文件处理、响应工具
  - `github.com/QuantumNous/new-api/setting/model_setting` — 模型设置
  - `github.com/QuantumNous/new-api/setting/reasoning` — 推理努力级别
  - `github.com/QuantumNous/new-api/types` — 错误类型
- **外部依赖**: `github.com/gin-gonic/gin`, `github.com/samber/lo`

## 3. 类型定义

### `geminiSupportedMimeTypes` 变量
```go
var geminiSupportedMimeTypes = map[string]bool{...}
```
Gemini 支持的 MIME 类型白名单，包含 PDF、音频（mpeg/mp3/wav）、图片（png/jpeg/webp/heic/heif）、纯文本、视频（mov/mpeg/mp4/mpg/avi/wmv/flv 等）。

### `geminiOpenAPISchemaAllowedFields` 变量
```go
var geminiOpenAPISchemaAllowedFields = map[string]struct{}{...}
```
Gemini 函数参数 schema 中允许的字段白名单（基于 OpenAPI Schema 子集），包括 `type`, `properties`, `items`, `anyOf`, `enum`, `required` 等。

### 常量
- `thoughtSignatureBypassValue` — thought signature 绕过值
- `pro25MinBudget` / `pro25MaxBudget` — gemini-2.5-pro 思考预算范围（128~32768）
- `flash25MaxBudget` — gemini-2.5-flash 最大预算（24576）
- `flash25LiteMinBudget` / `flash25LiteMaxBudget` — gemini-2.5-flash-lite 预算范围（512~24576）
- `geminiFunctionSchemaMaxDepth` — schema 递归最大深度（64）

### `GeminiModelsResponse` 结构体
```go
type GeminiModelsResponse struct {
    Models        []dto.GeminiModel
    NextPageToken string
}
```
Gemini 模型列表 API 的响应结构。

## 4. 函数详解

### 模型判断函数
- **`isNew25ProModel(modelName)`** — 判断是否为新的 gemini-2.5-pro 模型（排除旧版预览）
- **`is25FlashLiteModel(modelName)`** — 判断是否为 gemini-2.5-flash-lite 模型

### 思考预算计算
- **`clampThinkingBudget(modelName, budget)`** — 根据模型名称将思考预算限制在允许范围内
- **`clampThinkingBudgetByEffort(modelName, effort)`** — 根据 effort 级别（high/medium/low/minimal）计算思考预算，分别对应 80%/50%/20%/5% 的最大预算

### `ThinkingAdaptor(geminiRequest, info, oaiRequest...)`
思考适配器，根据模型名称后缀配置 `ThinkingConfig`：
- `-thinking-<budget>`: 指定预算值
- `-thinking`: 自动根据 max_tokens 百分比或 reasoningEffort 计算预算
- `-nothinking`: 禁用思考（预算为 0）
- effort 后缀（如 `-high`）: 使用 ThinkingLevel

### `CovertOpenAI2Gemini(c, textRequest, info) (*GeminiChatRequest, error)`
核心转换函数，将 OpenAI 请求转为 Gemini 请求：
1. 基础参数映射：temperature, topP, maxOutputTokens, seed, stopSequences
2. 处理 `extra_body.google` 中的 thinking_config 和 image_config（snake_case 到 camelCase 转换）
3. 调用 ThinkingAdaptor 配置思考参数
4. 构建 safetySettings
5. 转换 tools：处理 googleSearch/codeExecution/urlContext 特殊工具，清理函数参数 schema
6. 转换 tool_choice 为 Gemini toolConfig
7. 处理 response_format（json_schema/json_object → ResponseMimeType + ResponseSchema）
8. 遍历 messages：
   - system/developer → SystemInstructions
   - tool/function → FunctionResponse（尝试解析 JSON 对象/数组/纯文本）
   - assistant/tool_calls → FunctionCall（含 ThoughtSignature）
   - 文本内容：解析 markdown 嵌入图片（`![image](data:...)`）
   - 多媒体内容：转为 InlineData（含 MIME 类型白名单校验）
9. 角色映射：assistant → model

### 辅助函数
- **`parseStopSequences(stop)`** — 解析停止序列，支持 string/string slice/interface slice
- **`hasFunctionCallContent(call)`** — 检查 FunctionCall 是否有实际内容
- **`getSupportedMimeTypesList()`** — 获取支持的 MIME 类型列表

### Schema 清理函数
- **`cleanFunctionParameters(params)`** — 递归清理函数参数中不支持的字段
- **`cleanFunctionParametersWithDepth(params, depth)`** — 带深度限制的递归清理，保留 Gemini 支持的字段子集
- **`cleanFunctionParametersShallow(params)`** — 浅层清理，超过最大深度时截断递归
- **`normalizeGeminiSchemaTypeAndNullable(schema)`** — 标准化 type 为大写（OBJECT/ARRAY/STRING 等），处理 null 类型转 nullable
- **`removeAdditionalPropertiesWithDepth(schema, depth)`** — 递归删除 additionalProperties 和 title/$schema（最大深度 5）

### 转义处理
- **`unescapeString(s)`** — 手动实现 JSON 字符串反转义
- **`unescapeMapOrSlice(data)`** — 递归反转义 map/slice 中的字符串

### 响应转换函数
- **`getResponseToolCall(item)`** — 将 Gemini FunctionCall 转为 OpenAI ToolCallResponse
- **`buildUsageFromGeminiMetadata(metadata, fallbackPromptTokens)`** — 从 Gemini UsageMetadata 构建 OpenAI Usage，包含 prompt/completion/reasoning/cached/audio/text/image tokens 的详细统计
- **`responseGeminiChat2OpenAI(c, response)`** — 非流式响应转换，处理 InlineData（图片/媒体）、FunctionCall、Thought、ExecutableCode、CodeExecutionResult
- **`streamResponseGeminiChat2OpenAI(geminiResponse)`** — 流式响应转换，逻辑类似但返回 ChatCompletionsStreamResponse

### 流式处理
- **`handleStream(c, info, resp)`** — 序列化并发送流式数据块
- **`handleFinalStream(c, info, resp)`** — 发送最终流式响应
- **`geminiStreamHandler(c, info, resp, callback)`** — 通用流式处理框架，扫描 SSE 数据、更新 usage、调用回调
- **`GeminiChatStreamHandler(c, info, resp)`** — 聊天流式处理，管理 tool call index 映射、首空响应发送、停止响应发送
- **`GeminiChatHandler(c, info, resp)`** — 非流式聊天处理，支持 OpenAI/Claude/Gemini 三种输出格式

### 其他处理器
- **`GeminiEmbeddingHandler(c, info, resp)`** — 嵌入响应处理，转换为 OpenAI 嵌入格式
- **`GeminiImageHandler(c, info, resp)`** — 图片生成响应处理，过滤 RAI 过滤的图片，计算 token 使用量（每张图 258 tokens）

### 工具函数
- **`FetchGeminiModels(baseURL, apiKey, proxyURL)`** — 分页获取 Gemini 可用模型列表
- **`convertToolChoiceToGeminiConfig(toolChoice)`** — 将 OpenAI tool_choice 转为 Gemini FunctionCallingConfig

## 5. 关键逻辑分析

### Thinking 适配器机制
通过模型名后缀（`-thinking`, `-nothinking`, `-thinking-<budget>`, effort 后缀）自动配置 Gemini 的思考功能。支持两种配置路径：自动根据 max_tokens 百分比计算，或用户通过 `extra_body.google.thinking_config` 显式指定。

### OpenAI→Gemini 消息转换
消息转换是最复杂的部分：需要处理 5 种角色（system/assistant/user/tool/function），将 OpenAI 的多模态内容（文本、图片 URL、base64 图片、markdown 嵌入图片）统一转换为 Gemini 的 Parts 格式。tool/function 消息需要尝试解析 JSON 内容。

### Thought Signature 绕过
当开启 FunctionCallThoughtSignature 且渠道类型为 Gemini/VertexAi 时，在 assistant 消息的第一个 FunctionCall 或文本 Part 上附加固定签名值，绕过 Gemini 的签名验证。

### Schema 清理
递归清理函数参数 schema，只保留 Gemini 支持的 OpenAPI Schema 字段子集。处理深度限制（64 层），超过后进行浅层截断。标准化 type 值为大写，处理 null 类型到 nullable 转换。

### 内存优化
响应转换使用 `strings.Builder` 并预估增长量（`Grow`），避免图片 base64 数据导致的多次内存分配。流式处理中使用 `inlineGrow` 预计算 buffer 大小。

## 6. 关联文件
- `adaptor.go` — 调用本文件的核心转换函数
- `relay-gemini-native.go` — Gemini 原生模式响应处理
- `constant.go` — 模型列表和安全设置常量
- `relay/channel/openai/` — OpenAI 流式处理工具
- `dto/gemini.go` — Gemini 相关 DTO 定义
- `setting/model_setting/` — Gemini 模型配置
