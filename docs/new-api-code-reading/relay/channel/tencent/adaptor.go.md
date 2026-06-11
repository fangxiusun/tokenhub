# adaptor.go 代码阅读文档

## 1. 全局总结
腾讯混元（Tencent Hunyuan）AI 渠道的适配器实现文件。实现了 `channel.Adaptor` 接口，负责将统一的 OpenAI 格式请求转换为腾讯混元 API 格式，处理请求签名、URL 构建、响应转换等工作。腾讯渠道使用 TC3-HMAC-SHA256 签名机制进行身份验证。

## 2. 依赖关系
- **标准库**: errors, fmt, io, net/http, strconv, strings
- **内部包**:
  - `github.com/QuantumNous/new-api/common` — 通用工具函数（时间戳、密钥解析等）
  - `github.com/QuantumNous/new-api/constant` — 上下文键常量
  - `github.com/QuantumNous/new-api/dto` — 数据传输对象
  - `github.com/QuantumNous/new-api/relay/channel` — 渠道通用工具
  - `github.com/QuantumNous/new-api/relay/common` — 中继通用类型（RelayInfo）
  - `github.com/QuantumNous/new-api/types` — 错误类型定义
- **外部依赖**:
  - `github.com/gin-gonic/gin` — HTTP 框架

## 3. 类型定义

### Adaptor 结构体
```go
type Adaptor struct {
    Sign      string  // TC3 签名字符串
    AppID     int64   // 腾讯云应用 ID
    Action    string  // API 操作类型（如 "ChatCompletions"）
    Version   string  // API 版本号（如 "2023-09-01"）
    Timestamp int64   // 请求时间戳
}
```
腾讯混元渠道的主适配器，包含签名和请求元数据。

## 4. 函数详解

### Init(info *relaycommon.RelayInfo)
初始化适配器，设置 Action 为 "ChatCompletions"，Version 为 "2023-09-01"，并获取当前时间戳。

### GetRequestURL(info *relaycommon.RelayInfo) (string, error)
返回渠道基础 URL，末尾追加 `/`。

### SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error
设置请求头，包括标准 API 请求头和腾讯特有的：
- `Authorization`: TC3 签名
- `X-TC-Action`: 操作类型
- `X-TC-Version`: API 版本
- `X-TC-Timestamp`: 请求时间戳

### ConvertOpenAIRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error)
核心转换方法：
1. 从上下文获取 API Key，解析为 appId、secretId、secretKey
2. 调用 `requestOpenAI2Tencent` 转换请求格式
3. 调用 `getTencentSign` 计算请求签名

### DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (any, error)
委托给 `channel.DoApiRequest` 执行实际 HTTP 请求。

### DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage any, err *types.NewAPIError)
根据 `info.IsStream` 分发到流式或非流式处理器：
- 流式: `tencentStreamHandler`
- 非流式: `tencentHandler`

### GetModelList() []string / GetChannelName() string
返回模型列表和渠道名称。

### 未实现的方法
- `ConvertGeminiRequest` / `ConvertClaudeRequest` / `ConvertAudioRequest` / `ConvertImageRequest` / `ConvertEmbeddingRequest` / `ConvertOpenAIResponsesRequest` — 均返回 "not implemented" 错误

## 5. 关键逻辑分析

1. **TC3 签名机制**: 腾讯云 API 使用 TC3-HMAC-SHA256 签名，需要在 `ConvertOpenAIRequest` 中完成签名计算，签名结果存入 `Adaptor.Sign` 字段。
2. **请求格式转换**: OpenAI 格式 → 腾讯混元格式，消息从 `dto.Message` 转换为 `TencentMessage`，角色和内容保持一致。
3. **流式/非流式分发**: `DoResponse` 根据 `IsStream` 标志分发到不同的响应处理器。

## 6. 关联文件
- `tencent/dto.go` — 腾讯混元请求/响应数据结构定义
- `tencent/constants.go` — 模型列表和渠道名称常量
- `tencent/relay-tencent.go` — 请求/响应转换逻辑、签名算法、流式处理器
- `relay/channel/adaptor.go` — 适配器接口定义
- `relay/channel/common.go` — 通用 API 请求工具
