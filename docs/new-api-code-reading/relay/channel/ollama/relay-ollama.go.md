# relay-ollama.go 代码阅读文档

## 1. 全局总结
该文件是 Ollama 渠道的核心中继处理文件，包含三大功能模块：
1. **请求转换**：将 OpenAI 格式的聊天/生成/嵌入请求转换为 Ollama 原生格式
2. **响应处理**：处理嵌入请求的响应转换
3. **模型管理**：提供 Ollama 模型的拉取、删除、查询版本和获取模型列表等管理功能

## 2. 依赖关系
- `encoding/json`、`fmt`、`io`、`net/http`、`strings`、`time` — 标准库
- `github.com/QuantumNous/new-api/common` — 通用工具函数
- `github.com/QuantumNous/new-api/dto` — 数据传输对象
- `github.com/QuantumNous/new-api/relay/common` — 中继公共信息
- `github.com/QuantumNous/new-api/relay/helper` — 流式扫描器
- `github.com/QuantumNous/new-api/service` — 业务逻辑
- `github.com/QuantumNous/new-api/types` — 错误类型
- `github.com/gin-gonic/gin` — HTTP 框架
- `github.com/samber/lo` — 泛型工具库

## 3. 类型定义
无自定义类型定义，使用 `dto.go` 中定义的 Ollama 专用 DTO。

## 4. 函数详解

### `openAIChatToOllamaChat(c *gin.Context, r *dto.GeneralOpenAIRequest) (*OllamaChatRequest, error)`
- **作用**：将 OpenAI 聊天请求转换为 Ollama 聊天请求
- **关键逻辑**：
  - 处理 ResponseFormat（json/json_schema）
  - 映射模型参数（temperature、top_p、top_k、frequency_penalty、presence_penalty、seed、max_tokens、stop）
  - 转换工具调用格式
  - 处理多模态消息（图片 base64 转换）
  - 处理 tool 消息和 tool_calls

### `openAIToGenerate(c *gin.Context, r *dto.GeneralOpenAIRequest) (*OllamaGenerateRequest, error)`
- **作用**：将 OpenAI completions 请求转换为 Ollama generate 请求
- **关键逻辑**：
  - 处理 Prompt（string/[]any 格式）
  - 处理 Suffix 字段
  - 映射模型参数

### `requestOpenAI2Embeddings(r dto.EmbeddingRequest) *OllamaEmbeddingRequest`
- **作用**：将 OpenAI 嵌入请求转换为 Ollama 嵌入请求
- **关键逻辑**：单输入返回 string，多输入返回数组

### `ollamaEmbeddingHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)`
- **作用**：处理 Ollama 嵌入响应，转换为 OpenAI 格式

### `FetchOllamaModels(baseURL, apiKey string) ([]OllamaModel, error)`
- **作用**：从 Ollama 服务器获取已安装的模型列表
- **API**：GET `/api/tags`

### `PullOllamaModel(baseURL, apiKey, modelName string) error`
- **作用**：非流式拉取模型（同步，30分钟超时）

### `PullOllamaModelStream(baseURL, apiKey, modelName string, progressCallback func(OllamaPullResponse)) error`
- **作用**：流式拉取模型（支持进度回调，1小时超时）
- **逻辑**：使用 StreamScanner 逐行读取拉取进度，支持 error/success 状态检测

### `DeleteOllamaModel(baseURL, apiKey, modelName string) error`
- **作用**：删除本地模型
- **API**：DELETE `/api/delete`

### `FetchOllamaVersion(baseURL, apiKey string) (string, error)`
- **作用**：获取 Ollama 服务器版本
- **API**：GET `/api/version`

## 5. 关键逻辑分析
- **参数映射**：OpenAI 的 `max_tokens` 映射为 Ollama 的 `num_predict`，OpenAI 的 `stop` 支持 string/[]string/[]any 三种格式
- **图片处理**：聊天消息中的图片 URL 通过 `service.GetBase64Data` 下载并转为 base64 编码
- **工具调用**：OpenAI 的 tool_calls 格式转换为 Ollama 的 OllamaToolCall 格式，Arguments 从 JSON 字符串解析为 interface{}
- **模型管理**：提供完整的模型生命周期管理（查询、拉取、删除），支持流式进度回调
- **超时设置**：拉取模型使用较长超时（30分钟/1小时），适应大模型下载场景

## 6. 关联文件
- `relay/channel/ollama/adaptor.go` — 调用本文件中的转换函数
- `relay/channel/ollama/dto.go` — 使用的 DTO 定义
- `relay/channel/ollama/stream.go` — 流式响应处理
- `relay/helper/stream_scanner.go` — 流式扫描器
- `service/image.go` — GetBase64Data 图片下载服务
