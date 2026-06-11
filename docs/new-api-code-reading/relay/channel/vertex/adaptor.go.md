# adaptor.go 代码阅读文档

## 1. 全局总结
Google Vertex AI 渠道的适配器实现文件。Vertex AI 是一个多功能平台，支持三种请求模式：Claude 模型（Anthropic）、Gemini 模型（Google）和开源模型。适配器根据模型名称自动判断请求模式，并复用 claude、gemini、openai 渠道的处理逻辑。

## 2. 依赖关系
- **标准库**: encoding/json, errors, fmt, io, net/http, strings
- **内部包**:
  - `github.com/QuantumNous/new-api/common` — JSON 处理、配置解析
  - `github.com/QuantumNous/new-api/dto` — 数据传输对象
  - `github.com/QuantumNous/new-api/relay/channel` — 渠道通用工具
  - `github.com/QuantumNous/new-api/relay/channel/claude` — Claude 渠道适配器
  - `github.com/QuantumNous/new-api/relay/channel/gemini` — Gemini 渠道适配器
  - `github.com/QuantumNous/new-api/relay/channel/openai` — OpenAI 渠道适配器
  - `github.com/QuantumNous/new-api/relay/common` — RelayInfo
  - `github.com/QuantumNous/new-api/relay/constant` — 中继模式常量
  - `github.com/QuantumNous/new-api/setting/model_setting` — 模型设置
  - `github.com/QuantumNous/new-api/setting/reasoning` — 推理设置
  - `github.com/QuantumNous/new-api/types` — 错误类型
- **外部依赖**:
  - `github.com/gin-gonic/gin` — HTTP 框架
  - `github.com/samber/lo` — 泛型工具

## 3. 类型定义

### 请求模式常量
```go
const (
    RequestModeClaude     = 1  // Claude 模式（Anthropic）
    RequestModeGemini     = 2  // Gemini 模式（Google）
    RequestModeOpenSource = 3  // 开源模型模式
)
```

### claudeModelMap
Claude 模型名称映射表，将标准 Claude 模型名映射为 Vertex AI 格式（如 `claude-3-sonnet-20240229` → `claude-3-sonnet@20240229`）。

### Adaptor 结构体
```go
type Adaptor struct {
    RequestMode        int           // 请求模式（Claude/Gemini/OpenSource）
    AccountCredentials Credentials   // Google 服务账号凭证
}
```

### anthropicVersion 常量
```go
const anthropicVersion = "vertex-2023-10-16"
```
Vertex AI 上 Claude API 的版本号。

## 4. 函数详解

### Init(info *relaycommon.RelayInfo)
根据模型名称自动判断请求模式：
- 以 "claude" 开头 → Claude 模式
- 包含 "llama" 或 "-maas" → 开源模型模式
- 其他 → Gemini 模式

### GetRequestURL(info *relaycommon.RelayInfo) (string, error)
构建请求 URL，根据请求模式和认证方式（服务账号/API Key）使用不同的 URL 构建器。

### ConvertOpenAIRequest(...)
根据请求模式分发转换逻辑：
- Claude 模式: 使用 `claude.RequestOpenAI2ClaudeMessage` 转换
- Gemini 模式: 使用 `gemini.CovertOpenAI2Gemini` 转换
- 开源模型模式: 直接透传请求
- Gemini + imagen 模型: 转换为图片生成请求

### ConvertClaudeRequest(...)
将 Claude 请求包装为 Vertex AI 格式，设置 `anthropicVersion`。

### ConvertGeminiRequest(...)
委托给 `gemini.Adaptor.ConvertGeminiRequest`，并可选移除 `functionResponse.id`。

### DoResponse(...)
根据请求模式和流式状态分发到对应的处理器：
- Claude: `claudeAdaptor.DoResponse`
- Gemini: `gemini.GeminiTextGenerationStreamHandler` 或 `gemini.GeminiChatStreamHandler`
- 开源模型: `openai.OaiStreamHandler` 或 `openai.OpenaiHandler`

### GetModelList()
合并 Vertex AI 自有模型、Claude 模型和 Gemini 模型列表。

## 5. 关键逻辑分析

1. **多模式适配**: Vertex AI 通过 `RequestMode` 字段在同一适配器中支持三种不同的 API 格式，根据模型名称自动选择。

2. **模型名称映射**: Claude 模型名需要从标准格式映射到 Vertex AI 格式（`@` 替代 `-`），通过 `claudeModelMap` 实现。

3. **认证方式**: 支持两种认证：
   - 服务账号（Service Account）：使用 Google Cloud 凭证，通过 JWT 换取 Access Token
   - API Key：直接在 URL 中附加 `key=` 参数

4. **Thinking 模式适配**: 对 Gemini 模型支持 `-thinking`、`-nothinking`、`-thinking-<budget>` 等后缀的智能处理。

5. **Function Response ID 移除**: Vertex AI 不支持 `functionResponse.id` 字段，需要在转换时移除。

## 6. 关联文件
- `vertex/constants.go` — 模型列表常量
- `vertex/dto.go` — Vertex AI Claude 请求结构体
- `vertex/relay-vertex.go` — 模型区域获取逻辑
- `vertex/service_account.go` — Google 服务账号认证
- `vertex/url_builder.go` — URL 构建器
- `relay/channel/claude/` — Claude 渠道实现
- `relay/channel/gemini/` — Gemini 渠道实现
- `relay/channel/openai/` — OpenAI 渠道实现
