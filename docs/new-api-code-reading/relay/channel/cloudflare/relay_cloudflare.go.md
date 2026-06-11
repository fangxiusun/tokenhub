# relay_cloudflare.go 代码阅读文档

## 1. 全局总结

本文件实现了 Cloudflare Workers AI 频道的请求转换和响应处理逻辑。包含将 OpenAI 格式请求转换为 Cloudflare 格式的函数，以及流式/非流式响应处理器和语音转文本（STT）处理器。所有响应都被转换为 OpenAI 兼容格式返回给下游。

## 2. 依赖关系

| 包 | 用途 |
|---|------|
| `bufio` | 流式扫描器的分隔策略 |
| `encoding/json` | JSON 序列化/反序列化 |
| `io` | I/O 操作（ReadAll） |
| `net/http` | HTTP 响应处理 |
| `strings` | 字符串处理（前缀/后缀裁剪） |
| `time` | 时间记录（首响应时间） |
| `dto` | 通用数据传输对象 |
| `logger` | 日志记录 |
| `relay/common` | RelayInfo 上下文 |
| `relay/helper` | 流式处理工具（SetEventStreamHeaders、ObjectData、Done 等） |
| `service` | 业务工具（ResponseText2Usage、CloseResponseBodyGracefully） |
| `types` | 错误类型 |
| `lo` | 工具库（FromPtrOr） |
| `gin` | Web 框架 |

## 3. 类型定义

无类型定义，本文件仅包含函数。

## 4. 函数详解

### `convertCf2CompletionsRequest(textRequest dto.GeneralOpenAIRequest) *CfRequest`
将 OpenAI 格式的请求转换为 Cloudflare 旧版 Completions 格式：
- 提取 `Prompt` 字段（从 `textRequest.Prompt` 中取字符串值）
- 设置 `MaxTokens`、`Stream`、`Temperature`

### `cfStreamHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*types.NewAPIError, *dto.Usage)`
处理 Cloudflare 流式响应：
1. 创建流式扫描器，逐行读取 `data: ` 前缀的 SSE 事件
2. 对每个事件解析为 `dto.ChatCompletionsStreamResponse`
3. 设置 `role: assistant`，累加响应文本
4. 设置响应 ID 和模型名称
5. 记录首响应时间
6. 流结束后根据响应文本估算 usage（token 用量）
7. 如果需要包含 usage，发送最终 usage 事件
8. 发送 `[DONE]` 标记

### `cfHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*types.NewAPIError, *dto.Usage)`
处理 Cloudflare 非流式响应：
1. 读取完整响应体
2. 解析为 `dto.TextResponse`
3. 设置模型名称
4. 累加所有 choice 的响应文本
5. 通过文本估算 usage
6. 设置 usage 和响应 ID
7. 序列化并返回 JSON 响应

### `cfSTTHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*types.NewAPIError, *dto.Usage)`
处理 Cloudflare 语音转文本响应：
1. 读取并解析为 `CfAudioResponse`
2. 提取识别文本，构建 `dto.AudioResponse`
3. 序列化并返回 JSON 响应
4. 基于识别文本估算 usage

## 5. 关键逻辑分析

### Token 用量估算机制
Cloudflare API 不直接返回 token 用量数据，因此使用 `service.ResponseText2Usage` 通过响应文本来估算 token 数量。这是一种近似方法，实际 token 数可能与估算值有偏差。

### 流式处理流程
```
读取行 -> 去除 "data: " 前缀 -> 检查 "[DONE]" -> 解析 JSON -> 设置元数据 -> 渲染到客户端
```

### 错误处理策略
流式处理中单条消息解析失败不会中断整个流（`continue`），而是记录错误后继续处理。非流式处理中，任何解析错误都会返回 `ErrorCodeBadResponseBody` 错误。

### 响应体关闭
所有处理函数在读取完响应体后都调用 `service.CloseResponseBodyGracefully(resp)` 确保资源释放。

## 6. 关联文件

- `relay/channel/cloudflare/adaptor.go` - 调用本文件中的处理函数
- `relay/channel/cloudflare/dto.go` - `CfAudioResponse` 等 DTO 定义
- `relay/helper/` - 流式处理工具函数
- `relay/service/` - Token 用量估算和响应体管理
