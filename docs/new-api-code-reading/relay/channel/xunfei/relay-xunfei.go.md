# relay-xunfei.go 代码阅读文档

## 1. 全局总结
讯飞星火渠道的核心中继逻辑文件。包含请求格式转换、响应格式转换、WebSocket 通信、HMAC 签名认证、API 版本映射等功能。

## 2. 依赖关系
- **标准库**: crypto/hmac, crypto/sha256, encoding/base64, encoding/json, fmt, io, net/url, strings, time
- **内部包**:
  - `github.com/QuantumNous/new-api/common` — JSON 处理、日志
  - `github.com/QuantumNous/new-api/constant` — 常量定义
  - `github.com/QuantumNous/new-api/dto` — 数据传输对象
  - `github.com/QuantumNous/new-api/relay/helper` — 流式处理工具
  - `github.com/QuantumNous/new-api/types` — 错误类型
- **外部依赖**:
  - `github.com/gin-gonic/gin` — HTTP 框架
  - `github.com/gorilla/websocket` — WebSocket 客户端
  - `github.com/samber/lo` — 泛型工具

## 3. 类型定义
无新增类型定义。

## 4. 函数详解

### requestOpenAI2Xunfei(request dto.GeneralOpenAIRequest, xunfeiAppId string, domain string) *XunfeiChatRequest
将 OpenAI 请求转换为讯飞格式：
- 非 v3.5 模型: 将 system 消息转换为 user+assistant 对话
- 设置 AppId、Domain、Temperature、TopK、MaxTokens

### responseXunfei2OpenAI(response *XunfeiChatResponse) *dto.OpenAITextResponse
将讯飞响应转换为 OpenAI 格式，提取第一个文本项的内容。

### streamResponseXunfei2OpenAI(xunfeiResponse *XunfeiChatResponse) *dto.ChatCompletionsStreamResponse
将讯飞流式响应转换为 OpenAI 流式格式。

### buildXunfeiAuthUrl(hostUrl string, apiKey, apiSecret string) string
构建讯飞认证 URL：
1. 使用 HMAC-SHA256 签名（host、date、request-line）
2. Base64 编码签名
3. 将 host、date、authorization 作为查询参数附加到 URL

### xunfeiStreamHandler(...)
流式 WebSocket 响应处理器：
1. 建立 WebSocket 连接
2. 发送请求
3. 通过 channel 接收响应
4. 使用 Gin 的 Stream 方法逐块发送

### xunfeiHandler(...)
非流式 WebSocket 响应处理器：
1. 建立 WebSocket 连接
2. 发送请求
3. 累积所有响应内容
4. 返回完整响应

### xunfeiMakeRequest(...)
建立 WebSocket 连接并发送请求：
1. 使用 `websocket.Dialer` 建立连接
2. 发送转换后的请求
3. 启动 goroutine 读取响应
4. 返回数据通道和停止通道

### apiVersion2domain(apiVersion string) string
将 API 版本映射为 domain：
- v1.1 → lite
- v2.1 → generalv2
- v3.1 → generalv3
- v3.5 → generalv3.5
- v4.0 → 4.0Ultra

### getXunfeiAuthUrl(c *gin.Context, apiKey string, apiSecret string, modelName string) (string, string)
获取域和认证 URL。

### getAPIVersion(c *gin.Context, modelName string) string
从多个来源获取 API 版本：
1. URL 查询参数 `api-version`
2. 模型名中的版本号
3. Gin 上下文 `api_version`
4. 默认 v1.1

## 5. 关键逻辑分析

1. **WebSocket 通信**: 讯飞星火使用 WebSocket 进行实时通信，通过 channel 在 goroutine 间传递响应数据。

2. **HMAC 签名**: 使用 HMAC-SHA256 对 host、date、request-line 进行签名，签名结果 Base64 编码后作为 authorization 参数。

3. **System 消息转换**: 对于非 v3.5 模型，system 消息被转换为 user+assistant 对话对（user 发送 system 内容，assistant 回复 "Okay"）。

4. **API 版本路由**: 不同版本对应不同的 WebSocket 端点和 domain 参数，通过多种方式自动检测。

5. **流式/非流式统一**: 两种模式都通过 WebSocket 通信，区别在于流式模式逐块发送，非流式模式累积后一次性返回。

## 6. 关联文件
- `xunfei/adaptor.go` — 调用处理器函数
- `xunfei/dto.go` — 数据结构定义
- `xunfei/constants.go` — 模型列表
