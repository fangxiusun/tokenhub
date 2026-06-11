# adaptor.go 代码阅读文档

## 1. 全局总结

本文件实现了 DeepSeek 频道的适配器。DeepSeek 是深度求索公司的 AI 模型服务。该适配器支持 OpenAI 和 Claude 两种请求格式，能够根据中继格式自动选择对应的响应处理器。核心特色是支持 DeepSeek V4 模型的 thinking（思考链）后缀解析和转换。

## 2. 依赖关系

| 包 | 用途 |
|---|------|
| `errors` | 错误创建 |
| `fmt` | 字符串格式化 |
| `io` | I/O 操作 |
| `net/http` | HTTP 请求/响应 |
| `strings` | 字符串处理（后缀检查） |
| `common` | JSON 工具（Marshal） |
| `dto` | 数据传输对象 |
| `relay/channel` | 通用频道工具 |
| `relay/channel/claude` | Claude 适配器复用 |
| `relay/channel/openai` | OpenAI 适配器复用 |
| `relay/common` | RelayInfo 上下文 |
| `relay/constant` | 中继模式常量 |
| `setting/reasoning` | 推理设置（思考后缀解析） |
| `types` | 错误类型、中继格式常量 |
| `gin` | Web 框架 |

## 3. 类型定义

### `Adaptor` 结构体

```go
type Adaptor struct{}
```

无状态适配器。

## 4. 函数详解

### `(a *Adaptor) Init(info *relaycommon.RelayInfo)`
空实现。

### `(a *Adaptor) GetRequestURL(info *relaycommon.RelayInfo) (string, error)`
根据中继格式和模式构建 URL：
- **Claude 格式**: `{base}/anthropic/v1/messages`
- **Completions 模式**: `{base}/beta/completions`（自动追加 `/beta` 后缀）
- **默认**: `{base}/v1/chat/completions`

### `(a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error`
设置 `Authorization: Bearer {apiKey}` 请求头。

### `(a *Adaptor) ConvertOpenAIRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error)`
将 OpenAI 请求转换为 DeepSeek 格式：
1. 调用 `applyDeepSeekV4OpenAIThinkingSuffix` 处理思考后缀
2. 返回修改后的请求

### `(a *Adaptor) ConvertClaudeRequest(c *gin.Context, info *relaycommon.RelayInfo, req *dto.ClaudeRequest) (any, error)`
将 Claude 请求转换为 DeepSeek 格式：
1. 委托 `claude.Adaptor.ConvertClaudeRequest` 进行基础转换
2. 调用 `applyDeepSeekV4ClaudeThinkingSuffix` 处理思考后缀

### `applyDeepSeekV4OpenAIThinkingSuffix(info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) error`
处理 OpenAI 格式的 DeepSeek V4 思考后缀：
1. 从模型名解析思考类型和 effort
2. 设置 `request.THINKING` 字段（JSON 格式）
3. 设置 `request.ReasoningEffort`
4. 更新 `info.UpstreamModelName` 和 `info.ReasoningEffort`

### `applyDeepSeekV4ClaudeThinkingSuffix(info *relaycommon.RelayInfo, request *dto.ClaudeRequest) error`
处理 Claude 格式的 DeepSeek V4 思考后缀：
1. 解析模型名获取思考类型和 effort
2. 设置 `request.Thinking` 字段
3. 设置 `request.OutputConfig`（包含 effort 信息）
4. 更新 relay info 中的相关字段

### `(a *Adaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (any, error)`
委托给 `channel.DoApiRequest`。

### `(a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage any, err *types.NewAPIError)`
根据中继格式选择处理器：
- **Claude 格式**: 委托 `claude.Adaptor.DoResponse`
- **默认**: 委托 `openai.Adaptor.DoResponse`

### `(a *Adaptor) ConvertAudioRequest` / `ConvertImageRequest` / `ConvertEmbeddingRequest` / `ConvertGeminiRequest`
均返回 "not implemented"。

### `(a *Adaptor) ConvertRerankRequest(c *gin.Context, relayMode int, request dto.RerankRequest) (any, error)`
返回 `nil, nil`（空实现，不报错也不做转换）。

### `(a *Adaptor) ConvertOpenAIResponsesRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.OpenAIResponsesRequest) (any, error)`
返回 "not implemented"。

### `(a *Adaptor) GetModelList() []string`
返回 DeepSeek 模型列表。

### `(a *Adaptor) GetChannelName() string`
返回 `"deepseek"`。

## 5. 关键逻辑分析

### 双格式支持
DeepSeek 适配器的独特之处在于同时支持 OpenAI 和 Claude 两种请求/响应格式。这通过 `info.RelayFormat` 进行判断：
- `types.RelayFormatClaude`: 使用 Claude 格式
- 其他: 使用 OpenAI 格式

### DeepSeek V4 思考后缀机制
DeepSeek V4 模型支持在模型名后附加思考类型后缀（如 `deepseek-v4-pro-thinking`），适配器通过 `reasoning.ParseDeepSeekV4ThinkingSuffix` 解析后缀，将思考配置注入请求体。这使得用户可以通过模型名隐式控制推理行为。

### 委托模式
响应处理完全委托给 claude 或 openai 适配器，DeepSeek 本身不实现响应解析逻辑。这减少了代码重复，但要求上游格式与标准 OpenAI/Claude 格式兼容。

### Completions 端点的 /beta 路径
DeepSeek 的 Completions API 使用 `/beta/completions` 路径，适配器会自动在 base URL 后追加 `/beta`（如果尚未存在）。

## 6. 关联文件

- `relay/channel/claude/` - Claude 格式的请求转换和响应处理
- `relay/channel/openai/` - OpenAI 格式的请求转换和响应处理
- `relay/channel/deepseek/constants.go` - 模型列表
- `setting/reasoning/` - DeepSeek V4 思考后缀解析逻辑
- `types/relay_format.go` - 中继格式常量定义
