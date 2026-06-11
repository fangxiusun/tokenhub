# adaptor.go 代码阅读文档

## 1. 全局总结
该文件实现了 MokaAI 渠道的适配器（Adaptor），遵循 relay 通用适配器接口。负责 URL 构建、请求头设置、请求格式转换和响应分发。MokaAI 目前仅实现了嵌入（Embedding）功能，其他功能（聊天、图像、音频等）均返回未实现错误。

## 2. 依赖关系
- `errors`、`fmt`、`io`、`net/http`、`strings` — 标准库
- `github.com/QuantumNous/new-api/dto` — 数据传输对象
- `github.com/QuantumNous/new-api/relay/channel` — 渠道公共函数
- `github.com/QuantumNous/new-api/relay/common` — 中继公共信息
- `github.com/QuantumNous/new-api/relay/constant` — 中继常量（RelayMode 等）
- `github.com/QuantumNous/new-api/types` — 错误类型
- `github.com/gin-gonic/gin` — HTTP 框架

## 3. 类型定义

### `Adaptor` (struct)
- **作用**：MokaAI 渠道适配器，实现 relay 适配器接口的所有方法
- **字段**：无字段（空结构体）

## 4. 函数详解

### `(a *Adaptor) Init(info *relaycommon.RelayInfo)`
- **作用**：初始化适配器，当前为空实现

### `(a *Adaptor) GetRequestURL(info *relaycommon.RelayInfo) (string, error)`
- **作用**：构建请求 URL
- **逻辑**：根据模型名前缀 `m3e` 判断是嵌入请求还是聊天请求，分别拼接 `/embeddings` 或 `/chat/` 后缀

### `(a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error`
- **作用**：设置请求头
- **逻辑**：调用通用头部设置函数，然后添加 Bearer Token 认证

### `(a *Adaptor) ConvertOpenAIRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error)`
- **作用**：将 OpenAI 请求转换为 MokaAI 格式
- **逻辑**：仅支持 `RelayModeEmbeddings` 模式，调用 `embeddingRequestOpenAI2Moka` 进行转换

### `(a *Adaptor) ConvertEmbeddingRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.EmbeddingRequest) (any, error)`
- **作用**：直接透传嵌入请求（不做转换）

### `(a *Adaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (any, error)`
- **作用**：执行 HTTP 请求
- **逻辑**：委托给 `channel.DoApiRequest` 通用函数

### `(a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage any, err *types.NewAPIError)`
- **作用**：处理上游响应
- **逻辑**：根据 RelayMode 分发，嵌入模式调用 `mokaEmbeddingHandler`

### `(a *Adaptor) GetModelList() []string` / `GetChannelName() string`
- **作用**：返回模型列表和渠道名称

### 其他未实现方法
- `ConvertGeminiRequest`、`ConvertClaudeRequest`、`ConvertAudioRequest`、`ConvertImageRequest`、`ConvertRerankRequest`、`ConvertOpenAIResponsesRequest` 均返回未实现错误

## 5. 关键逻辑分析
- **适配器模式**：遵循统一的适配器接口，通过 `Init` → `GetRequestURL` → `SetupRequestHeader` → `ConvertRequest` → `DoRequest` → `DoResponse` 的生命周期处理请求
- **URL 路由判断**：通过模型名前缀 `m3e` 区分嵌入和聊天请求的 URL 路径
- **功能限制**：MokaAI 渠道目前仅支持嵌入功能，其他适配器方法均为空实现或 panic

## 6. 关联文件
- `relay/channel/mokaai/relay-mokaai.go` — 提供转换函数和处理器
- `relay/channel/mokaai/constants.go` — 模型列表和渠道名称
- `relay/channel/channel.go` — 通用渠道函数（DoApiRequest、SetupApiRequestHeader）
