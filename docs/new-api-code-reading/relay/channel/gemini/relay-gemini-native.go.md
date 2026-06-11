# relay-gemini-native.go 代码阅读文档

## 1. 全局总结
本文件处理 Gemini 原生模式（RelayModeGemini）下的响应，即不进行格式转换、直接透传 Gemini 原生请求/响应的场景。包含非流式文本生成、嵌入和流式文本生成三种处理器。

## 2. 依赖关系
- **标准库**: `fmt`, `io`, `net/http`
- **项目内部**:
  - `github.com/QuantumNous/new-api/common` — JSON 解析
  - `github.com/QuantumNous/new-api/constant` — 上下文常量
  - `github.com/QuantumNous/new-api/dto` — 数据传输对象
  - `github.com/QuantumNous/new-api/logger` — 日志
  - `github.com/QuantumNous/new-api/relay/common` — RelayInfo
  - `github.com/QuantumNous/new-api/relay/helper` — 流式头设置、流数据写入
  - `github.com/QuantumNous/new-api/service` — 响应体关闭、IO 拷贝
  - `github.com/QuantumNous/new-api/types` — 错误类型
- **外部依赖**: `github.com/gin-gonic/gin`

## 3. 类型定义
无自定义类型定义。

## 4. 函数详解

### `GeminiTextGenerationHandler(c, info, resp) (*dto.Usage, *types.NewAPIError)`
处理 Gemini 原生模式的非流式文本生成响应：
1. 读取完整响应体
2. 解析为 `dto.GeminiChatResponse`
3. 检查空候选和 PromptFeedback.BlockReason，设置管理员拒绝原因
4. 通过 `buildUsageFromGeminiMetadata` 计算 token 使用量
5. 使用 `IOCopyBytesGracefully` 直接透传原始响应体给客户端

### `NativeGeminiEmbeddingHandler(c, resp, info) (*dto.Usage, *types.NewAPIError)`
处理 Gemini 原生模式的嵌入响应：
1. 读取完整响应体
2. 根据 `IsGeminiBatchEmbedding` 标志解析为批量或单次嵌入响应
3. 透传原始响应体
4. 使用估算的 prompt tokens 作为 usage

### `GeminiTextGenerationStreamHandler(c, info, resp) (*dto.Usage, *types.NewAPIError)`
处理 Gemini 原生模式的流式文本生成：
1. 设置 SSE 事件流响应头
2. 复用 `geminiStreamHandler` 框架，回调函数直接将原始 SSE 数据写入客户端
3. 每写入一个数据块后递增 `SendResponseCount`

## 5. 关键逻辑分析

### 原生模式特点
与标准模式不同，原生模式不做格式转换（不转为 OpenAI 格式），直接将 Gemini API 的响应透传给客户端。这适用于客户端本身支持 Gemini 格式的场景。

### 响应体透传
使用 `service.IOCopyBytesGracefully` 直接写入原始响应字节，避免了序列化/反序列化的开销。

### 使用量估算
原生嵌入模式下，由于 Google 尚未明确嵌入模型的计费方式，使用 `service.ResponseText2Usage` 基于估算 prompt tokens 计算 usage。

## 6. 关联文件
- `adaptor.go` — 在 `DoResponse` 中根据 `RelayModeGemini` 调度到本文件的处理器
- `relay-gemini.go` — 提供 `geminiStreamHandler` 和 `buildUsageFromGeminiMetadata` 函数
