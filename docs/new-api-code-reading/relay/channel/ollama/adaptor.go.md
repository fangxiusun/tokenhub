# adaptor.go 代码阅读文档

## 1. 全局总结
该文件实现了 Ollama 渠道的适配器。Ollama 是本地部署的 LLM 服务，其 API 格式与 OpenAI 不同，适配器负责将 OpenAI/Claude 格式的请求转换为 Ollama 原生格式，并将 Ollama 的响应转换回 OpenAI 格式。支持聊天补全、文本生成、嵌入和 Claude 格式转换。

## 2. 依赖关系
- `errors`、`io`、`net/http`、`strings` — 标准库
- `github.com/QuantumNous/new-api/dto` — 数据传输对象
- `github.com/QuantumNous/new-api/relay/channel` — 渠道公共函数
- `github.com/QuantumNous/new-api/relay/channel/openai` — OpenAI 适配器（Claude→OpenAI→Ollama 转换链）
- `github.com/QuantumNous/new-api/relay/common` — 中继公共信息
- `github.com/QuantumNous/new-api/relay/constant` — 中继常量
- `github.com/QuantumNous/new-api/types` — 错误类型
- `github.com/gin-gonic/gin` — HTTP 框架

## 3. 类型定义

### `Adaptor` (struct)
- **作用**：Ollama 渠道适配器
- **字段**：无字段（空结构体）

## 4. 函数详解

### `(a *Adaptor) Init(info *relaycommon.RelayInfo)`
- **作用**：初始化适配器，当前为空实现

### `(a *Adaptor) GetRequestURL(info *relaycommon.RelayInfo) (string, error)`
- **作用**：构建请求 URL
- **逻辑**：根据请求类型选择 Ollama API 路径：
  - 嵌入 → `/api/embed`
  - 文本生成（completions） → `/api/generate`
  - 聊天补全 → `/api/chat`

### `(a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error`
- **作用**：设置请求头，添加 Bearer Token 认证

### `(a *Adaptor) ConvertOpenAIRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error)`
- **作用**：将 OpenAI 请求转换为 Ollama 格式
- **逻辑**：
  - 文本生成模式 → 调用 `openAIToGenerate` 转换为 `OllamaGenerateRequest`
  - 聊天模式 → 调用 `openAIChatToOllamaChat` 转换为 `OllamaChatRequest`

### `(a *Adaptor) ConvertClaudeRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.ClaudeRequest) (any, error)`
- **作用**：将 Claude 请求转换为 Ollama 格式
- **逻辑**：Claude → OpenAI → Ollama 三步转换链

### `(a *Adaptor) ConvertEmbeddingRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.EmbeddingRequest) (any, error)`
- **作用**：转换嵌入请求
- **逻辑**：调用 `requestOpenAI2Embeddings` 转换

### `(a *Adaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (any, error)`
- **作用**：执行 HTTP 请求，委托给 `channel.DoApiRequest`

### `(a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage any, err *types.NewAPIError)`
- **作用**：处理上游响应
- **逻辑**：根据模式分发：
  - 嵌入 → `ollamaEmbeddingHandler`
  - 流式 → `ollamaStreamHandler`
  - 非流式 → `ollamaChatHandler`

### `(a *Adaptor) GetModelList() []string` / `GetChannelName() string`
- **作用**：返回模型列表和渠道名称

## 5. 关键逻辑分析
- **三步转换链**：Claude 请求通过 Claude→OpenAI→Ollama 的转换链处理，复用了 openai 适配器的转换逻辑
- **URL 路径路由**：通过 `RequestURLPath` 和 `RelayMode` 双重判断选择正确的 Ollama API 端点
- **流式/非流式分发**：DoResponse 根据 `IsStream` 标志选择不同的处理器
- **功能限制**：不支持图像、音频和 Responses API

## 6. 关联文件
- `relay/channel/ollama/relay-ollama.go` — 请求转换函数和模型管理 API
- `relay/channel/ollama/stream.go` — 流式和非流式响应处理器
- `relay/channel/ollama/dto.go` — Ollama 专用 DTO 定义
- `relay/channel/ollama/constants.go` — 模型列表和渠道名称
- `relay/channel/openai/adaptor.go` — Claude→OpenAI 转换委托
