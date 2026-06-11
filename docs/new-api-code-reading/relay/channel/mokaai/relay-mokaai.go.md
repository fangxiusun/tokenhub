# relay-mokaai.go 代码阅读文档

## 1. 全局总结
该文件是 MokaAI 渠道的中继处理核心，负责将 OpenAI 格式的嵌入（Embedding）请求转换为 MokaAI（百度千帆）的格式，并将 MokaAI 的响应转换回 OpenAI 格式。同时提供了一个嵌入请求的 HTTP 处理器，用于读取上游响应、转换格式并返回给客户端。

## 2. 依赖关系
- `encoding/json` — JSON 序列化/反序列化
- `io` — IO 操作
- `net/http` — HTTP 响应处理
- `github.com/QuantumNous/new-api/common` — 通用工具函数（JSON 序列化）
- `github.com/QuantumNous/new-api/dto` — 数据传输对象（请求/响应结构体）
- `github.com/QuantumNous/new-api/relay/common` — 中继公共信息（RelayInfo）
- `github.com/QuantumNous/new-api/service` — 业务逻辑（关闭响应体、IO 拷贝）
- `github.com/QuantumNous/new-api/types` — 错误类型定义
- `github.com/gin-gonic/gin` — HTTP 框架

## 3. 类型定义
无自定义类型定义。使用 `dto.EmbeddingRequest`、`dto.EmbeddingResponse`、`dto.OpenAIEmbeddingResponse`、`dto.OpenAIEmbeddingResponseItem` 等外部 DTO。

## 4. 函数详解

### `embeddingRequestOpenAI2Moka(request dto.GeneralOpenAIRequest) *dto.EmbeddingRequest`
- **作用**：将 OpenAI 格式的通用请求转换为 MokaAI 嵌入请求格式
- **参数**：OpenAI 通用请求（输入可能是 string、[]string 或 []interface{}）
- **返回**：MokaAI 格式的嵌入请求
- **逻辑**：通过类型 switch 处理三种不同的输入格式，统一转为 `[]string`

### `embeddingResponseMoka2OpenAI(response *dto.EmbeddingResponse) *dto.OpenAIEmbeddingResponse`
- **作用**：将 MokaAI 的嵌入响应转换为 OpenAI 标准格式
- **参数**：MokaAI 嵌入响应
- **返回**：OpenAI 格式的嵌入响应
- **逻辑**：遍历 MokaAI 的 data 数组，逐一映射到 OpenAI 的响应结构中

### `mokaEmbeddingHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)`
- **作用**：处理 MokaAI 嵌入请求的完整 HTTP 响应流程
- **参数**：Gin 上下文、中继信息、上游 HTTP 响应
- **返回**：使用量信息、错误（可为 nil）
- **逻辑**：读取响应体 → 反序列化为 MokaAI 格式 → 转换为 OpenAI 格式 → 序列化并写入客户端响应

## 5. 关键逻辑分析
- **输入格式兼容**：`embeddingRequestOpenAI2Moka` 通过类型 switch 兼容三种输入格式（string、[]string、[]interface{}），确保不同客户端调用方式都能正确转换
- **错误处理**：所有错误统一包装为 `types.NewAPIError`，通过 `ErrorCodeBadResponseBody` 标识响应体解析失败
- **优雅关闭**：使用 `service.CloseResponseBodyGracefully` 确保响应体被正确关闭，避免资源泄漏
- **内容透传**：转换后的响应通过 `service.IOCopyBytesGracefully` 写回客户端，保持 HTTP 状态码不变

## 6. 关联文件
- `relay/channel/mokaai/adaptor.go` — 调用本文件中的转换函数和处理器
- `relay/channel/mokaai/constants.go` — 定义模型列表和渠道名称
- `dto/embedding.go` — 嵌入请求/响应的 DTO 定义
- `relay/common/relay_info.go` — RelayInfo 结构体定义
