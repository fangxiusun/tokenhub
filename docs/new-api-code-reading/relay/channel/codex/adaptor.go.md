# adaptor.go 代码阅读文档

## 1. 全局总结

本文件实现了 Codex（ChatGPT/Codex CLI）频道的适配器。该适配器专门支持 OpenAI Responses API 格式（`/v1/responses` 和 `/v1/responses/compact`），不支持传统的 Chat Completions、Embeddings 等端点。它处理 OAuth 密钥解析、系统提示词注入、请求头定制等 Codex 特有的逻辑。

## 2. 依赖关系

| 包 | 用途 |
|---|------|
| `encoding/json` | JSON RawMessage 操作 |
| `errors` | 错误创建 |
| `io` | I/O 操作 |
| `net/http` | HTTP 请求/响应 |
| `strings` | 字符串处理（前缀检查、空白修剪） |
| `common` | JSON 工具函数（Marshal/Unmarshal） |
| `dto` | 数据传输对象 |
| `relay/channel` | 通用频道工具 |
| `relay/channel/openai` | OpenAI Responses 处理器复用 |
| `relay/common` | RelayInfo 上下文、URL 构建 |
| `relay/constant` | 中继模式常量 |
| `types` | 错误类型 |
| `gin` | Web 框架 |

## 3. 类型定义

### `Adaptor` 结构体

```go
type Adaptor struct{}
```

无状态适配器，实现 `channel.Adaptor` 接口。

## 4. 函数详解

### `(a *Adaptor) Init(info *relaycommon.RelayInfo)`
空实现，无需初始化。

### `(a *Adaptor) GetRequestURL(info *relaycommon.RelayInfo) (string, error)`
构建 Codex API 的请求 URL：
- 仅支持 `RelayModeResponses` 和 `RelayModeResponsesCompact` 模式
- Responses 模式: `{base}/backend-api/codex/responses`
- ResponsesCompact 模式: `{base}/backend-api/codex/responses/compact`
- 其他模式返回错误

### `(a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error`
设置 Codex 特有的请求头：
1. 调用通用工具设置基础头
2. 解析 API Key 为 OAuth JSON 对象（必须以 `{` 开头）
3. 设置 `Authorization: Bearer {accessToken}`
4. 设置 `chatgpt-account-id: {accountId}`
5. 设置 `OpenAI-Beta: responses=experimental`（如未设置）
6. 设置 `originator: codex_cli_rs`（如未设置）
7. 强制设置 `Content-Type: application/json`（Codex 后端对 Content-Type 严格校验）
8. 流式请求设置 `Accept: text/event-stream`

### `(a *Adaptor) ConvertOpenAIRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error)`
直接返回错误 — Codex 不支持 `/v1/chat/completions` 端点。

### `(a *Adaptor) ConvertOpenAIResponsesRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.OpenAIResponsesRequest) (any, error)`
核心请求转换逻辑：
1. **系统提示词注入**：
   - 如果频道配置了 `SystemPrompt` 且请求没有 `instructions`，直接设置
   - 如果配置了 `SystemPromptOverride`，将系统提示词与现有 instructions 合并
   - 如果两者都为空，确保 `instructions` 至少为空字符串（Codex 后端要求此字段存在）
2. **Compact 模式**：直接返回请求，不做额外修改
3. **普通模式**：强制 `store: false`，移除 `max_output_tokens` 和 `temperature`

### `(a *Adaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (any, error)`
委托给 `channel.DoApiRequest` 执行 HTTP 请求。

### `(a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage any, err *types.NewAPIError)`
响应处理分发：
- 非 Responses 模式返回错误
- Compact 模式: `openai.OaiResponsesCompactionHandler`
- 流式: `openai.OaiResponsesStreamHandler`
- 非流式: `openai.OaiResponsesHandler`

### `(a *Adaptor) ConvertAudioRequest` / `ConvertImageRequest` / `ConvertRerankRequest` / `ConvertEmbeddingRequest` / `ConvertClaudeRequest` / `ConvertGeminiRequest`
均返回 "endpoint not supported" 错误。

### `(a *Adaptor) GetModelList() []string`
返回 Codex 支持的模型列表。

### `(a *Adaptor) GetChannelName() string`
返回 `"codex"`。

## 5. 关键逻辑分析

### OAuth 密钥格式要求
Codex 频道要求 API Key 必须是 JSON 对象格式（`ParseOAuthKey`），包含 `access_token` 和 `account_id` 等字段。这与标准的 Bearer Token 认证不同。

### 系统提示词注入策略
实现了灵活的系统提示词管理：
- 空 instructions → 直接设置
- 有 instructions + Override=true → 前置拼接
- 有 instructions + Override=false → 不修改
- 有 instructions + Override=true + 空值 → 设置

### 请求字段清理
对于非 Compact 模式，强制移除 `max_output_tokens` 和 `temperature`，因为 Codex 后端不接受这些参数。

### Content-Type 强制设置
Codex 后端对 Content-Type 格式严格校验，客户端可能发送带 charset 的值（如 `application/json; charset=utf-8`），会被上游拒绝。因此强制设置为精确的 `application/json`。

## 6. 关联文件

- `relay/channel/codex/oauth_key.go` - OAuth 密钥解析
- `relay/channel/codex/constants.go` - 模型列表和频道名称
- `relay/channel/openai/` - OpenAI Responses 处理器复用
- `relay/common/` - URL 构建工具（GetFullRequestURL）
