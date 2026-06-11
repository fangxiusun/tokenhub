# dto.go 代码阅读文档

## 1. 全局总结
该文件定义了 Ollama 渠道专用的数据传输对象（DTO），包括聊天请求/响应、生成请求、嵌入请求/响应、模型管理请求/响应等结构体，用于与 Ollama API 进行序列化和反序列化。

## 2. 依赖关系
- `encoding/json` — JSON RawMessage 类型

## 3. 类型定义

### `OllamaChatMessage` (struct)
- **作用**：Ollama 聊天消息结构
- **字段**：Role（角色）、Content（内容）、Images（图片 base64 数组）、ToolCalls（工具调用）、ToolName（工具名称）、Thinking（思维链内容）

### `OllamaToolFunction` (struct)
- **作用**：工具函数定义
- **字段**：Name（函数名）、Description（描述）、Parameters（参数，interface{} 类型）

### `OllamaTool` (struct)
- **作用**：工具定义
- **字段**：Type（工具类型）、Function（函数定义）

### `OllamaToolCall` (struct)
- **作用**：工具调用结果
- **字段**：Function（包含 Name 和 Arguments）

### `OllamaChatRequest` (struct)
- **作用**：Ollama 聊天 API 请求
- **字段**：Model、Messages、Tools、Format（格式化输出）、Stream、Options（模型参数）、KeepAlive（保持连接）、Think（思维链开关）

### `OllamaGenerateRequest` (struct)
- **作用**：Ollama 文本生成 API 请求
- **字段**：Model、Prompt、Suffix、Images、Format、Stream、Options、KeepAlive、Think

### `OllamaEmbeddingRequest` (struct)
- **作用**：Ollama 嵌入 API 请求
- **字段**：Model、Input（支持单个字符串或数组）、Options、Dimensions

### `OllamaEmbeddingResponse` (struct)
- **作用**：Ollama 嵌入 API 响应
- **字段**：Error、Model、Embeddings（嵌入向量数组）、PromptEvalCount

### `OllamaTagsResponse` (struct)
- **作用**：Ollama 模型列表 API 响应
- **字段**：Models（模型数组）

### `OllamaModel` (struct)
- **作用**：Ollama 模型信息
- **字段**：Name、Size、Digest、ModifiedAt、Details

### `OllamaModelDetail` (struct)
- **作用**：Ollama 模型详细信息
- **字段**：ParentModel、Format、Family、Families、ParameterSize、QuantizationLevel

### `OllamaPullRequest` (struct)
- **作用**：拉取模型请求
- **字段**：Name（模型名）、Stream（是否流式）

### `OllamaPullResponse` (struct)
- **作用**：拉取模型响应
- **字段**：Status（状态）、Digest、Total（总大小）、Completed（已完成大小）

### `OllamaDeleteRequest` (struct)
- **作用**：删除模型请求
- **字段**：Name（模型名）

## 4. 函数详解
无函数定义。

## 5. 关键逻辑分析
- **多模态支持**：`OllamaChatMessage` 支持图片（Images 字段）和工具调用（ToolCalls 字段）
- **思维链支持**：通过 `Thinking` 字段（json.RawMessage）支持 Ollama 的思维链功能
- **灵活输入**：`OllamaEmbeddingRequest.Input` 使用 interface{} 类型，支持单个字符串或字符串数组
- **模型管理**：提供了完整的模型管理 DTO（拉取、删除、查询版本）

## 6. 关联文件
- `relay/channel/ollama/relay-ollama.go` — 使用这些 DTO 进行请求转换和 API 调用
- `relay/channel/ollama/stream.go` — 使用 `ollamaChatStreamChunk`（内联定义）进行流式处理
