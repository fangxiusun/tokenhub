# relay-openai.go 代码阅读文档

## 1. 全局总结
该文件是 OpenAI 渠道的核心中继处理文件，实现了流式和非流式两种响应处理模式。流式处理器（OaiStreamHandler）将上游 SSE 流转换为客户端格式，支持 OpenAI/Claude/Gemini 三种输出格式，处理 thinking 内容转换。非流式处理器（OpenaiHandler）处理标准 HTTP 响应，支持 OpenRouter Enterprise 格式和多种输出格式转换。

## 2. 依赖关系
- `fmt`、`io`、`net/http`、`strings` — 标准库
- `github.com/QuantumNous/new-api/common` — 通用工具函数
- `github.com/QuantumNous/new-api/constant` — 常量定义
- `github.com/QuantumNous/new-api/dto` — 数据传输对象
- `github.com/QuantumNous/new-api/logger` — 日志
- `github.com/QuantumNous/new-api/relay/channel/openrouter` — OpenRouter 适配
- `github.com/QuantumNous/new-api/relay/common` — 中继公共信息
- `github.com/QuantumNous/new-api/relay/helper` — 流式处理辅助
- `github.com/QuantumNous/new-api/service` — 业务逻辑
- `github.com/QuantumNous/new-api/types` — 错误类型
- `github.com/gin-gonic/gin` — HTTP 框架

## 3. 类型定义
无自定义类型定义。

## 4. 函数详解

### `sendStreamData(c *gin.Context, info *relaycommon.RelayInfo, data string, forceFormat bool, thinkToContent bool) error`
- **作用**：发送流式数据，支持 think_to_content 转换
- **关键逻辑**：
  - `thinkToContent=false`：直接发送原始数据
  - `thinkToContent=true`：将 reasoning_content 转换为 content，用 `<think>`/`</think>` 标签包裹
  - 首次思考内容：发送 `<think>\n` + 内容
  - 思考到内容转换：发送 `\n</think>\n`
  - 后续思考内容：直接转换为 content

### `OaiStreamHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)`
- **作用**：处理 OpenAI 流式响应
- **关键逻辑**：
  1. 使用 StreamScannerHandler 处理 SSE 流
  2. 对每个数据帧：发送格式化数据 + 提取 token 信息
  3. 音频模型特殊处理：从倒数第二个帧提取 usage
  4. 处理最后一个响应帧
  5. 如果没有 stream usage，通过文本内容估算 token
  6. 发送最终 usage 和 [DONE]
  7. 调用 applyUsagePostProcessing 进行渠道特定的 usage 后处理

### `OpenaiHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)`
- **作用**：处理 OpenAI 非流式响应
- **关键逻辑**：
  1. 读取完整响应体
  2. OpenRouter Enterprise 格式特殊处理
  3. 解析 OpenAI 响应
  4. 检测 content_filter 完成原因
  5. 如果 PromptTokens 为 0，估算 usage
  6. 根据 RelayFormat 转换响应格式（OpenAI/Claude/Gemini）
  7. 如果 usage 被修改，重新序列化响应体

## 5. 关键逻辑分析
- **Think to Content 转换**：通过 `ThinkingContentInfo` 状态机跟踪思考内容的发送状态，实现 `<think>` 标签的正确包裹和闭合
- **音频模型特殊处理**：音频模型的 usage 在倒数第二个 SSE 帧中，而非最后一个帧
- **Usage 后处理**：`applyUsagePostProcessing` 根据渠道类型（DeepSeek、智谱、Moonshot、OpenAI）提取不同位置的 cached_tokens
- **OpenRouter Enterprise**：响应体被包装在 `data` 字段中，需要先提取再解析
- **Content Filter**：检测到 `content_filter` 完成原因时，设置上下文标记用于后续审计
- **格式转换链**：非流式响应通过 `ResponseOpenAI2Claude` 和 `ResponseOpenAI2Gemini` 转换为对应的格式

## 6. 关联文件
- `relay/channel/openai/adaptor.go` — DoResponse 中分发到这些处理器
- `relay/channel/openai/helper.go` — HandleStreamFormat、handleLastResponse 等辅助函数
- `relay/channel/openai/usage.go` — applyUsagePostProcessing usage 后处理
- `relay/channel/openrouter/adaptor.go` — OpenRouter Enterprise 响应格式
- `relay/helper/stream_scanner.go` — StreamScannerHandler 流式扫描器
- `service/response_convert.go` — 格式转换服务
