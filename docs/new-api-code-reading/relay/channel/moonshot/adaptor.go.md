# adaptor.go 代码阅读文档

## 1. 全局总结
该文件实现了 Moonshot（月之暗面/Kimi）渠道的适配器。Moonshot 同时支持 OpenAI 和 Claude 两种请求格式，适配器通过委托 claude 和 openai 子适配器来处理不同格式的转换和响应。支持聊天补全、嵌入、重排序和图像生成功能。

## 2. 依赖关系
- `errors`、`fmt`、`io`、`net/http`、`strings` — 标准库
- `github.com/QuantumNous/new-api/common` — 通用工具函数
- `github.com/QuantumNous/new-api/constant` — 渠道类型常量
- `github.com/QuantumNous/new-api/dto` — 数据传输对象
- `github.com/QuantumNous/new-api/relay/channel` — 渠道公共函数
- `github.com/QuantumNous/new-api/relay/channel/claude` — Claude 适配器（委托）
- `github.com/QuantumNous/new-api/relay/channel/openai` — OpenAI 适配器（委托）
- `github.com/QuantumNous/new-api/relay/common` — 中继公共信息
- `github.com/QuantumNous/new-api/relay/constant` — 中继常量
- `github.com/QuantumNous/new-api/types` — 错误类型
- `github.com/gin-gonic/gin` — HTTP 框架

## 3. 类型定义

### `Adaptor` (struct)
- **作用**：Moonshot 渠道适配器
- **字段**：无字段（空结构体）

## 4. 函数详解

### `(a *Adaptor) Init(info *relaycommon.RelayInfo)`
- **作用**：初始化适配器，当前为空实现

### `(a *Adaptor) GetRequestURL(info *relaycommon.RelayInfo) (string, error)`
- **作用**：构建请求 URL
- **逻辑**：
  1. 检查 `ChannelSpecialBases` 特殊基础 URL 配置
  2. 根据 `RelayFormat`（Claude/OpenAI）选择不同的路径
  3. 默认模式下根据 `RelayMode` 区分 rerank/embeddings/chat/completions 路径

### `(a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error`
- **作用**：设置请求头，添加 Bearer Token 认证

### `(a *Adaptor) ConvertOpenAIRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error)`
- **作用**：转换 OpenAI 请求
- **逻辑**：对 `kimi-k2.6` 模型强制将 temperature 设为 1.0（该模型仅支持 temperature=1）

### `(a *Adaptor) ConvertClaudeRequest(c *gin.Context, info *relaycommon.RelayInfo, req *dto.ClaudeRequest) (any, error)`
- **作用**：转换 Claude 格式请求
- **逻辑**：委托给 `claude.Adaptor` 处理

### `(a *Adaptor) ConvertImageRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.ImageRequest) (any, error)`
- **作用**：转换图像请求
- **逻辑**：委托给 `openai.Adaptor` 处理

### `(a *Adaptor) ConvertRerankRequest(c *gin.Context, relayMode int, request dto.RerankRequest) (any, error)`
- **作用**：直接透传重排序请求

### `(a *Adaptor) ConvertEmbeddingRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.EmbeddingRequest) (any, error)`
- **作用**：直接透传嵌入请求

### `(a *Adaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (any, error)`
- **作用**：执行 HTTP 请求，委托给 `channel.DoApiRequest`

### `(a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage any, err *types.NewAPIError)`
- **作用**：处理上游响应
- **逻辑**：根据 `RelayFormat` 分发：
  - Claude 格式 → 委托 `claude.Adaptor.DoResponse`
  - 其他格式 → 委托 `openai.Adaptor.DoResponse`

### `getUpstreamModelName(info *relaycommon.RelayInfo, fallback string) string`
- **作用**：获取上游模型名，优先使用 info 中的值

### `isTemperatureOneOnlyModel(model string) bool`
- **作用**：判断是否为仅支持 temperature=1 的模型（当前只有 `kimi-k2.6`）

## 5. 关键逻辑分析
- **双格式支持**：Moonshot 同时支持 OpenAI 和 Claude 格式，通过 `RelayFormat` 动态分发到对应的子适配器
- **委托模式**：Claude 请求转换和响应处理直接委托给 claude 适配器，图像请求委托给 openai 适配器，避免重复实现
- **特殊模型处理**：`kimi-k2.6` 模型有 temperature 限制，适配器自动修正
- **特殊基础 URL**：支持通过 `ChannelSpecialBases` 配置特殊的 API 基础 URL

## 6. 关联文件
- `relay/channel/moonshot/constants.go` — 模型列表和渠道名称
- `relay/channel/claude/adaptor.go` — Claude 格式处理委托
- `relay/channel/openai/adaptor.go` — OpenAI 格式处理委托
- `relay/channel/openai/usage.go` — Moonshot 特殊的 cached_tokens 提取逻辑
