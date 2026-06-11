# 中转层详细设计 (`relay/`)

## 1. 概述
中转层是系统的核心 AI API 中转/代理系统。它处理请求转换、提供商适配和响应处理，支持 40+ AI 提供商。这是系统最复杂的部分，采用了适配器模式和工厂模式。

## 2. 文件详细说明

### 2.1 顶层文件
- **`relay_adaptor.go`** -- 工厂函数：
  - `GetAdaptor(apiType int)`: 根据 API 类型常量实例化具体的适配器。使用 switch 语句映射 33 种渠道类型。
  - `GetTaskAdaptor(platform constant.TaskPlatform)`: 根据任务平台实例化任务适配器。
  - **设计模式**: 工厂模式，每次调用创建新实例（线程安全），无单例状态。
- **`relay_task.go`** -- 异步任务生命周期：
  - `RelayTaskSubmit`: 任务提交（验证、预扣费、提交、调整计费）
  - `RelayTaskFetch`: 任务状态查询
  - `RelayTaskPolling`: 任务轮询（后台任务）
  - 计费集成：预扣费、提交调整、完成调整
- **`compatible_handler.go`** -- `TextHelper()`: 主要的 OpenAI 格式聊天补全处理器。
- **`claude_handler.go`** -- `ClaudeHelper()`: 原生 Claude `/v1/messages` 处理器。
- **`gemini_handler.go`** -- `GeminiHelper()` + `GeminiEmbeddingHandler()`: 原生 Gemini 处理器。
- **`responses_handler.go`** -- `ResponsesHelper()`: OpenAI Responses API 处理器。
- **`chat_completions_via_responses.go`** -- 桥接器：将 chat/completions 路由到 Responses API 后端。
- **`embedding_handler.go`** -- `EmbeddingHelper()`: 嵌入请求处理器。
- **`image_handler.go`** -- `ImageHelper()`: 图像生成/编辑处理器。
- **`audio_handler.go`** -- `AudioHelper()`: TTS/STT 处理器。
- **`rerank_handler.go`** -- `RerankHelper()`: 重排序处理器。
- **`mjproxy_handler.go`** -- Midjourney 代理处理器。
- **`websocket.go`** -- WebSocket/实时处理器。
- **`param_override_error.go`** -- 参数覆盖错误包装。

### 2.2 `relay/channel/` -- 适配器层

#### 共享基础设施
- **`adapter.go`** -- 核心接口定义：
  ```go
  type Adaptor interface {
      Init(info *relaycommon.RelayInfo)
      GetRequestURL(info *relaycommon.RelayInfo) (string, error)
      SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error
      ConvertOpenAIRequest(...) (any, error)
      ConvertRerankRequest(...) (any, error)
      ConvertEmbeddingRequest(...) (any, error)
      ConvertAudioRequest(...) (any, error)
      ConvertImageRequest(...) (any, error)
      ConvertOpenAIResponsesRequest(...) (any, error)
      ConvertClaudeRequest(...) (any, error)
      ConvertGeminiRequest(...) (any, error)
      DoRequest(...) (any, error)
      DoResponse(...) (usage any, err *types.NewAPIError)
      GetModelList() []string
      GetChannelName() string
  }
  ```
  - **设计模式**: 适配器模式，统一不同提供商的接口差异。
  - **关键设计**: 每个适配器必须处理所有入站格式转换（OpenAI, Claude, Gemini），无论上游提供商是什么。
- **`api_request.go`** -- 共享 HTTP 管道：
  - `DoApiRequest()`: 主要的 API 请求执行函数
  - `DoFormRequest()`: 表单请求执行
  - `DoWssRequest()`: WebSocket 请求执行
  - `DoTaskApiRequest()`: 任务 API 请求执行
  - 头部覆盖处理：支持 `{api_key}` 和 `{client_header:...}` 占位符

#### 提供商适配器（32 个包）
| 包名 | 提供商 | 说明 |
|------|--------|------|
| `openai/` | OpenAI (+ Azure, OpenRouter, Xinference, 360, LingYiWanWu) | 最常用的适配器，支持多种 OpenAI 兼容 API |
| `claude/` | Anthropic Claude | 原生 Claude Messages API |
| `gemini/` | Google Gemini | 原生 Gemini API |
| `vertex/` | Google Vertex AI | Google Cloud Vertex AI |
| `palm/` | Google PaLM | Google PaLM API |
| `ali/` | 阿里云/通义 | 阿里云 DashScope API |
| `aws/` | AWS Bedrock | AWS Bedrock API |
| `baidu/` | 百度 (v1) | 百度文心一言 API |
| `baidu_v2/` | 百度 (v2) | 百度文心一言 API v2 |
| `cloudflare/` | Cloudflare Workers AI | Cloudflare AI API |
| `codex/` | Codex | Codex API |
| `cohere/` | Cohere | Cohere API |
| `coze/` | Coze | Coze API |
| `deepseek/` | DeepSeek | DeepSeek API |
| `dify/` | Dify | Dify ChatFlow API |
| `jimeng/` | 即梦 | 即梦 API |
| `jina/` | Jina | Jina API |
| `minimax/` | MiniMax | MiniMax API |
| `mistral/` | Mistral | Mistral API |
| `mokaai/` | MokaAI | MokaAI API |
| `moonshot/` | Moonshot (使用 Claude 协议) | Moonshot API |
| `ollama/` | Ollama | Ollama 本地模型 |
| `openrouter/` | OpenRouter (仅常量，使用 OpenAI 适配器) | OpenRouter API |
| `perplexity/` | Perplexity | Perplexity API |
| `replicate/` | Replicate | Replicate API |
| `siliconflow/` | SiliconFlow | SiliconFlow API |
| `submodel/` | Submodel | Submodel API |
| `tencent/` | 腾讯云 | 腾讯混元 API |
| `volcengine/` | 火山引擎 | 火山引擎 API |
| `xai/` | xAI (Grok) | xAI API |
| `xinference/` | Xinference (仅常量，使用 OpenAI 适配器) | Xinference API |
| `xunfei/` | 讯飞 | 讯飞星火 API |
| `zhipu/` | 智谱 AI (v1) | 智谱 GLM API |
| `zhipu_4v/` | 智谱 AI (v4) | 智谱 GLM API v4 |
| `ai360/` | 360 AI (仅常量，使用 OpenAI 适配器) | 360 API |
| `lingyiwanwu/` | 零一万物 (仅常量，使用 OpenAI 适配器) | 零一万物 API |

#### 任务适配器（10 个包）
| 包名 | 提供商 | 说明 |
|------|--------|------|
| `task/ali/` | 阿里云异步任务 | 阿里云视频/音频生成 |
| `task/doubao/` | 豆包/火山引擎视频 | 豆包视频生成 |
| `task/gemini/` | Gemini 异步生成 | Gemini 视频生成 |
| `task/hailuo/` | 海螺 (MiniMax 视频) | 海螺视频生成 |
| `task/jimeng/` | 即梦异步任务 | 即梦视频生成 |
| `task/kling/` | Kling 视频 | Kling 视频生成 |
| `task/sora/` | Sora/OpenAI 视频 | OpenAI Sora 视频生成 |
| `task/suno/` | Suno 音乐 | Suno 音乐生成 |
| `task/vertex/` | Vertex AI 异步任务 | Vertex AI 视频生成 |
| `task/vidu/` | Vidu 视频 | Vidu 视频生成 |

### 2.3 子目录
- **`relay/common/`** -- `RelayInfo` 结构体（中心上下文对象），`StreamStatus`, `Billing`, `Override`, `RequestConversion`, `OutboundBody`, `RelayUtils`。
- **`relay/constant/`** -- `relay_mode.go`: 所有中转模式常量和路径到模式的映射。
- **`relay/helper/`** -- 模型映射、定价、流扫描、计费表达式、请求验证的辅助函数。
- **`relay/common_handler/`** -- 共享响应处理器（如 `rerank.go`）。
- **`relay/reasonmap/`** -- 推理映射工具。

## 3. 适配器模式详解

### 3.1 接口设计
所有适配器必须实现 `Adaptor` 接口，定义了标准化的请求处理流程：
1. **Init**: 初始化适配器，设置渠道信息
2. **GetRequestURL**: 根据渠道信息构造请求 URL
3. **SetupRequestHeader**: 设置请求头（认证、自定义头等）
4. **Convert*Request**: 将标准格式（OpenAI/Claude/Gemini）转换为厂商特定格式
5. **DoRequest**: 发送 HTTP 请求到上游
6. **DoResponse**: 处理上游响应，提取使用量信息

### 3.2 工厂模式
在 `relay/relay_adaptor.go` 中通过 `GetAdaptor(apiType int)` 工厂函数，根据渠道类型常量实例化具体适配器。

### 3.3 请求流程
1. **控制器** 解析请求，确定中转格式，创建 `RelayInfo`
2. **处理器** 调用 `GetAdaptor()`，`adaptor.Init()`，`adaptor.Convert*Request()`
3. **共享 `DoApiRequest()`** 处理 URL 构造、头部设置、HTTP 执行
4. **`adaptor.DoResponse()`** 解析上游响应并提取使用量
5. **服务层** 处理计费和配额扣减

### 3.4 RelayInfo 结构体
`RelayInfo` 是传递整个中转流水线的中心上下文对象，包含：
- 用户/Token 信息：userId, quota, group, token 详情
- 渠道元数据：`*ChannelMeta`（渠道类型、API 密钥、基础 URL、参数/头部覆盖、流支持标志）
- 请求上下文：中转模式、中转格式、原始/上游模型名称、流标志、URL 路径
- 计费状态：`PriceData`, `Billing` 会话, `ForcePreConsume`, 订阅信息
- 转换链：`RequestConversionChain []types.RelayFormat` 跟踪格式转换序列

---

## 关联文件列表

### 中转层核心文件
- `relay/relay_adaptor.go` - 工厂函数
- `relay/relay_task.go` - 异步任务生命周期
- `relay/compatible_handler.go` - OpenAI 格式处理器
- `relay/claude_handler.go` - Claude 格式处理器
- `relay/gemini_handler.go` - Gemini 格式处理器
- `relay/responses_handler.go` - Responses API 处理器
- `relay/chat_completions_via_responses.go` - Chat-Responses 桥接
- `relay/embedding_handler.go` - 嵌入处理器
- `relay/image_handler.go` - 图像处理器
- `relay/audio_handler.go` - 音频处理器
- `relay/rerank_handler.go` - 重排序处理器
- `relay/mjproxy_handler.go` - Midjourney 代理
- `relay/websocket.go` - WebSocket 处理器
- `relay/param_override_error.go` - 参数覆盖错误

### 适配器接口定义
- `relay/channel/adapter.go` - Adaptor 和 TaskAdaptor 接口
- `relay/channel/api_request.go` - 共享 HTTP 请求工具

### 提供商适配器（32 个）
- `relay/channel/openai/adaptor.go` - OpenAI 适配器
- `relay/channel/claude/adaptor.go` - Claude 适配器
- `relay/channel/gemini/adaptor.go` - Gemini 适配器
- `relay/channel/vertex/adaptor.go` - Vertex AI 适配器
- `relay/channel/palm/adaptor.go` - PaLM 适配器
- `relay/channel/ali/adaptor.go` - 阿里云适配器
- `relay/channel/aws/adaptor.go` - AWS Bedrock 适配器
- `relay/channel/baidu/adaptor.go` - 百度适配器
- `relay/channel/baidu_v2/adaptor.go` - 百度 v2 适配器
- `relay/channel/cloudflare/adaptor.go` - Cloudflare 适配器
- `relay/channel/codex/adaptor.go` - Codex 适配器
- `relay/channel/cohere/adaptor.go` - Cohere 适配器
- `relay/channel/coze/adaptor.go` - Coze 适配器
- `relay/channel/deepseek/adaptor.go` - DeepSeek 适配器
- `relay/channel/dify/adaptor.go` - Dify 适配器
- `relay/channel/jimeng/adaptor.go` - 即梦适配器
- `relay/channel/jina/adaptor.go` - Jina 适配器
- `relay/channel/minimax/adaptor.go` - MiniMax 适配器
- `relay/channel/mistral/adaptor.go` - Mistral 适配器
- `relay/channel/mokaai/adaptor.go` - MokaAI 适配器
- `relay/channel/moonshot/adaptor.go` - Moonshot 适配器
- `relay/channel/ollama/adaptor.go` - Ollama 适配器
- `relay/channel/openrouter/adaptor.go` - OpenRouter 适配器
- `relay/channel/perplexity/adaptor.go` - Perplexity 适配器
- `relay/channel/replicate/adaptor.go` - Replicate 适配器
- `relay/channel/siliconflow/adaptor.go` - SiliconFlow 适配器
- `relay/channel/submodel/adaptor.go` - Submodel 适配器
- `relay/channel/tencent/adaptor.go` - 腾讯适配器
- `relay/channel/volcengine/adaptor.go` - 火山引擎适配器
- `relay/channel/xai/adaptor.go` - xAI 适配器
- `relay/channel/xinference/adaptor.go` - Xinference 适配器
- `relay/channel/xunfei/adaptor.go` - 讯飞适配器
- `relay/channel/zhipu/adaptor.go` - 智谱适配器
- `relay/channel/zhipu_4v/adaptor.go` - 智谱 v4 适配器
- `relay/channel/ai360/adaptor.go` - 360 适配器
- `relay/channel/lingyiwanwu/adaptor.go` - 零一万物适配器

### 任务适配器（10 个）
- `relay/channel/task/ali/adaptor.go` - 阿里云任务适配器
- `relay/channel/task/doubao/adaptor.go` - 豆包任务适配器
- `relay/channel/task/gemini/adaptor.go` - Gemini 任务适配器
- `relay/channel/task/hailuo/adaptor.go` - 海螺任务适配器
- `relay/channel/task/jimeng/adaptor.go` - 即梦任务适配器
- `relay/channel/task/kling/adaptor.go` - Kling 任务适配器
- `relay/channel/task/sora/adaptor.go` - Sora 任务适配器
- `relay/channel/task/suno/adaptor.go` - Suno 任务适配器
- `relay/channel/task/vertex/adaptor.go` - Vertex 任务适配器
- `relay/channel/task/vidu/adaptor.go` - Vidu 任务适配器

### 中转公共组件
- `relay/common/relay_info.go` - RelayInfo 结构体
- `relay/constant/relay_mode.go` - 中转模式常量
- `relay/helper/` - 辅助函数目录
- `relay/common_handler/` - 共享处理器目录
- `relay/reasonmap/` - 推理映射目录
