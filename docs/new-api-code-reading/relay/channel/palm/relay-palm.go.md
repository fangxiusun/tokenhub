# relay-palm.go 代码阅读文档

## 1. 全局总结
本文件实现了 PaLM API 响应到 OpenAI 格式的转换逻辑，以及流式和非流式响应的处理函数。是 PaLM 渠道的核心响应处理模块，负责将 PaLM 特有的响应格式统一转换为系统内部使用的 OpenAI 标准格式。

## 2. 依赖关系
| 包名 | 用途 |
|------|------|
| `encoding/json` | JSON 序列化/反序列化 |
| `io` | IO 操作 |
| `net/http` | HTTP 响应处理 |
| `github.com/QuantumNous/new-api/common` | JSON 封装（Marshal）、时间戳、自定义事件 |
| `github.com/QuantumNous/new-api/constant` | 常量定义（FinishReasonStop） |
| `github.com/QuantumNous/new-api/dto` | 数据传输对象（OpenAI 响应格式） |
| `github.com/QuantumNous/new-api/relay/common` | 中继公共类型 |
| `github.com/QuantumNous/new-api/relay/helper` | 辅助函数（响应 ID 生成、SSE 头设置） |
| `github.com/QuantumNous/new-api/service` | 服务层（用量计算、响应体关闭、IO 拷贝） |
| `github.com/QuantumNous/new-api/types` | 错误类型定义 |
| `github.com/gin-gonic/gin` | Gin HTTP 框架 |

## 3. 类型定义
本文件无额外类型定义，使用同包 `dto.go` 中定义的结构体。

## 4. 函数详解

### `responsePaLM2OpenAI`
```go
func responsePaLM2OpenAI(response *PaLMChatResponse) *dto.OpenAITextResponse
```
- **参数**：PaLM 聊天响应
- **返回**：转换后的 OpenAI 文本响应
- **逻辑**：遍历 PaLM 的 `Candidates`，将每个候选转换为 OpenAI 的 `Choice` 结构，设置角色为 "assistant"，结束原因为 "stop"

### `streamResponsePaLM2OpenAI`
```go
func streamResponsePaLM2OpenAI(palmResponse *PaLMChatResponse) *dto.ChatCompletionsStreamResponse
```
- **参数**：PaLM 聊天响应
- **返回**：转换后的 OpenAI 流式响应块
- **逻辑**：取第一个候选的内容，构造流式响应块，模型固定为 "palm2"

### `palmStreamHandler`
```go
func palmStreamHandler(c *gin.Context, resp *http.Response) (*types.NewAPIError, string)
```
- **参数**：Gin 上下文、HTTP 响应
- **返回**：错误信息（如有）、响应文本
- **逻辑**：
  1. 在 goroutine 中读取完整响应体（PaLM 不支持真正的流式）
  2. 解析为 `PaLMChatResponse`
  3. 转换为 OpenAI 流式格式
  4. 通过 `c.Stream` 发送 SSE 事件
  5. 最后发送 `data: [DONE]` 结束信号

### `palmHandler`
```go
func palmHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)
```
- **参数**：Gin 上下文、中继信息、HTTP 响应
- **返回**：用量信息（如有）、错误信息（如有）
- **逻辑**：
  1. 读取完整响应体
  2. 解析为 `PaLMChatResponse`
  3. 检查错误（PaLM 错误在响应体中）
  4. 转换为 OpenAI 格式
  5. 计算 token 用量
  6. 写入 JSON 响应

## 5. 关键逻辑分析
- **伪流式处理**：PaLM API 不支持真正的 SSE 流式响应，`palmStreamHandler` 通过 goroutine 一次性读取完整响应后模拟流式发送。
- **错误处理策略**：PaLM 的错误信息嵌入在响应体的 `Error` 字段中（而非 HTTP 状态码），需要在解析后检查 `palmResponse.Error.Code`。
- **用量计算**：使用 `service.ResponseText2Usage` 通过文本内容估算 token 用量，而非从响应中直接获取。
- **JSON 封装**：所有 JSON 操作使用 `common.Marshal` 而非标准库，遵循项目规则。
- **响应体关闭**：多处调用 `service.CloseResponseBodyGracefully` 确保资源正确释放。

## 6. 关联文件
- `relay/channel/palm/adaptor.go` — 调用本文件中的处理函数
- `relay/channel/palm/dto.go` — PaLM 数据结构定义
- `relay/channel/palm/constants.go` — 渠道常量
- `relay/helper/response.go` — `GetResponseID`、`SetEventStreamHeaders` 辅助函数
- `relay/channel/channel.go` — 公共请求发送逻辑
- `service/` — 用量计算和响应处理工具函数
