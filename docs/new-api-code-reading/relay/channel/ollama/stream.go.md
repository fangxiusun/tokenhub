# stream.go 代码阅读文档

## 1. 全局总结
该文件负责处理 Ollama 渠道的响应，包括流式（SSE）和非流式两种模式。流式处理器将 Ollama 的 NDJSON 流转换为 OpenAI 兼容的 SSE 格式，支持文本内容、思维链（Thinking）和工具调用。非流式处理器将 Ollama 的多行 NDJSON 响应聚合为单个 OpenAI 格式响应。

## 2. 依赖关系
- `encoding/json`、`fmt`、`io`、`net/http`、`strings`、`time` — 标准库
- `github.com/QuantumNous/new-api/common` — 通用工具函数
- `github.com/QuantumNous/new-api/dto` — 数据传输对象
- `github.com/QuantumNous/new-api/logger` — 日志
- `github.com/QuantumNous/new-api/relay/common` — 中继公共信息
- `github.com/QuantumNous/new-api/relay/helper` — 流式处理辅助函数
- `github.com/QuantumNous/new-api/service` — 业务逻辑
- `github.com/QuantumNous/new-api/types` — 错误类型
- `github.com/gin-gonic/gin` — HTTP 框架

## 3. 类型定义

### `ollamaChatStreamChunk` (struct)
- **作用**：Ollama 流式响应的单个 chunk 结构
- **字段**：
  - `Model` — 模型名称
  - `CreatedAt` — 创建时间（RFC3339 格式）
  - `Message` — 聊天消息（包含 Role、Content、Thinking、ToolCalls）
  - `Response` — 生成模式的响应文本
  - `Done` — 是否完成
  - `DoneReason` — 完成原因
  - 统计字段：TotalDuration、LoadDuration、PromptEvalCount、EvalCount 等

## 4. 函数详解

### `toUnix(ts string) int64`
- **作用**：将 RFC3339 时间戳字符串转换为 Unix 时间戳
- **逻辑**：支持 RFC3339Nano 和 RFC3339 两种格式，解析失败返回当前时间

### `ollamaStreamHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)`
- **作用**：处理 Ollama 流式响应
- **关键逻辑**：
  1. 发送初始空响应（start frame）
  2. 使用 StreamScanner 逐行读取 NDJSON
  3. 非完成 chunk：提取 content/thinking/tool_calls，构建 OpenAI SSE delta
  4. 完成 chunk：提取 usage，发送 stop frame 和 usage frame，发送 `[DONE]`
  5. Thinking 内容处理：支持 JSON 字符串和原始字符串两种格式

### `ollamaChatHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)`
- **作用**：处理 Ollama 非流式聊天响应
- **关键逻辑**：
  1. 读取完整响应体
  2. 按行分割，逐行解析 NDJSON（兼容多行和单行格式）
  3. 聚合所有 content 和 thinking 内容
  4. 构建完整的 OpenAI 格式响应
  5. 提取 usage 统计信息

### `contentPtr(s string) *string`
- **作用**：将字符串转为指针，空字符串返回 nil

## 5. 关键逻辑分析
- **NDJSON 解析**：Ollama 使用 NDJSON（换行分隔的 JSON）格式，与 OpenAI 的 SSE 格式不同，需要逐行解析并转换
- **Thinking 支持**：流式和非流式都支持思维链内容，Thinking 字段可能是 JSON 编码的字符串或原始字符串，需要双重解析
- **Tool Calls**：流式模式下为每个工具调用生成唯一 ID（`call_0`、`call_1`...），将 arguments 从 interface{} 序列化为 JSON 字符串
- **Usage 提取**：从最后一个 chunk 的 `prompt_eval_count` 和 `eval_count` 提取 token 使用量
- **兼容性**：非流式处理器兼容单行 JSON 和多行 NDJSON 两种响应格式
- **时间解析**：Ollama 返回的时间格式不统一，`toUnix` 函数提供多格式兼容

## 6. 关联文件
- `relay/channel/ollama/adaptor.go` — 调用流式/非流式处理器
- `relay/channel/ollama/dto.go` — `ollamaChatStreamChunk` 相关 DTO 定义
- `relay/channel/ollama/relay-ollama.go` — 嵌入响应处理
- `relay/helper/stream_scanner.go` — StreamScanner 流式扫描器
- `relay/helper/response.go` — SSE 响应辅助函数
