# relay-cohere.go 代码阅读文档

## 1. 全局总结

本文件实现了 Cohere 频道的请求转换和响应处理逻辑。包含 OpenAI → Cohere 请求转换、Cohere → OpenAI 响应转换（流式/非流式），以及 Rerank 请求的转换和处理。使用 goroutine 和 channel 实现流式数据的异步处理。

## 2. 依赖关系

| 包 | 用途 |
|---|------|
| `encoding/json` | JSON 序列化/反序列化 |
| `io` | I/O 操作 |
| `net/http` | HTTP 响应处理 |
| `strings` | 字符串处理 |
| `time` | 时间记录 |
| `common` | JSON 工具、安全设置、日志 |
| `dto` | 数据传输对象 |
| `relay/common` | RelayInfo 上下文 |
| `relay/helper` | 流式处理工具 |
| `service` | Token 估算、响应体管理 |
| `types` | 错误类型 |
| `lo` | 工具库（FromPtrOr） |
| `gin` | Web 框架 |

## 3. 类型定义

无类型定义。

## 4. 函数详解

### `requestOpenAI2Cohere(textRequest dto.GeneralOpenAIRequest) *CohereRequest`
将 OpenAI 请求转换为 Cohere 格式：
1. 设置模型、流式标志、最大 token 数（默认 4000）
2. 如果全局安全设置不是 "NONE"，设置 `SafetyMode`
3. 遍历消息列表：
   - `user` 角色 → 设置为 `Message`（当前消息）
   - `assistant` → `CHATBOT`
   - `system` → `SYSTEM`
   - 其他 → `USER`
   - 非 user 消息添加到 `ChatHistory`

### `requestConvertRerank2Cohere(rerankRequest dto.RerankRequest) *CohereRerankRequest`
将标准 Rerank 请求转换为 Cohere 格式，确保 `TopN` 至少为 1。

### `stopReasonCohere2OpenAI(reason string) string`
将 Cohere 的停止原因转换为 OpenAI 格式：
- `COMPLETE` → `stop`
- `MAX_TOKENS` → `max_tokens`
- 其他原样返回

### `cohereStreamHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)`
处理 Cohere 流式响应（使用 goroutine + channel 异步模式）：
1. 启动 goroutine 扫描流数据并推送到 `dataChan`
2. 主协程通过 `c.Stream` 循环读取 channel
3. 对每个事件：
   - 解析为 `CohereResponse`
   - 转换为 `dto.ChatCompletionsStreamResponse`
   - 如果 `IsFinished`，设置 finish_reason 并提取 usage
   - 否则累加响应文本并发送增量内容
4. 流结束时发送 `[DONE]`
5. 如果 usage 未从响应中获取，通过文本估算

### `cohereHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)`
处理 Cohere 非流式响应：
1. 读取并解析响应体
2. 从 `meta.billed_units` 提取 token 用量
3. 构建 OpenAI 格式的 `TextResponse`
4. 序列化并返回

### `cohereRerankHandler(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (*dto.Usage, *types.NewAPIError)`
处理 Cohere Rerank 响应：
1. 读取并解析响应体
2. 如果 `input_tokens` 为 0，使用估算值
3. 构建标准 `dto.RerankResponse`
4. 序列化并返回

## 5. 关键逻辑分析

### 异步流式处理模式
Cohere 的流式处理使用了与其他频道不同的异步模式：
```go
go func() { /* 扫描器推送到 dataChan */ }()
c.Stream(func(w io.Writer) bool {
    select {
    case data := <-dataChan: /* 处理数据 */
    case <-stopChan: /* 结束 */
    }
})
```
这与 Cloudflare 的同步扫描模式形成对比。

### 消息角色映射
OpenAI 的 `assistant` 角色映射为 Cohere 的 `CHATBOT`，这是 Cohere API 的特殊要求。

### 安全模式注入
全局安全设置 `common.CohereSafetySetting` 会影响所有 Cohere 请求的安全模式。

### Token 用量获取策略
优先使用 Cohere API 返回的 `meta.billed_units` 中的精确数据，仅在数据缺失时回退到文本估算。

### 自定义分隔函数
流式扫描器使用自定义分隔函数，以换行符 `\n` 为分隔符（而非默认的 `bufio.ScanLines`），这是因为 Cohere 的 SSE 流格式使用 `\n` 分隔。

## 6. 关联文件

- `relay/channel/cohere/adaptor.go` - 调度本文件中的函数
- `relay/channel/cohere/dto.go` - DTO 定义
- `relay/helper/` - 流式处理工具
- `common/global_setting.go` - `CohereSafetySetting` 全局设置
