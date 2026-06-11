# adaptor.go 代码阅读文档

## 1. 全局总结
智谱 AI（ChatGLM）旧版渠道的适配器实现文件。智谱使用 JWT 认证和专有的 API 格式（/api/paas/v3/model-api/），适配器负责格式转换和请求构建。

## 2. 依赖关系
- **标准库**: errors, fmt, io, net/http
- **内部包**:
  - `github.com/QuantumNous/new-api/dto` — 数据传输对象
  - `github.com/QuantumNous/new-api/relay/channel` — 渠道通用工具
  - `github.com/QuantumNous/new-api/relay/common` — RelayInfo
  - `github.com/QuantumNous/new-api/types` — 错误类型
- **外部依赖**:
  - `github.com/gin-gonic/gin` — HTTP 框架
  - `github.com/samber/lo` — 泛型工具

## 3. 类型定义

### Adaptor 结构体
```go
type Adaptor struct {}
```
空结构体。

## 4. 函数详解

### GetRequestURL(...)
构建智谱 API URL：
- 流式: `/api/paas/v3/model-api/{model}/sse-invoke`
- 非流式: `/api/paas/v3/model-api/{model}/invoke`

### SetupRequestHeader(...)
设置认证头，使用 `getZhipuToken` 生成 JWT Token。

### ConvertOpenAIRequest(...)
将 OpenAI 请求转换为智谱格式：
- TopP >= 1 时修正为 0.99
- 调用 `requestOpenAI2Zhipu` 转换请求

### DoResponse(...)
根据流式状态分发到 `zhipuStreamHandler` 或 `zhipuHandler`。

### 未实现的方法
- `ConvertGeminiRequest` / `ConvertClaudeRequest` / `ConvertAudioRequest` / `ConvertImageRequest` / `ConvertEmbeddingRequest` / `ConvertOpenAIResponsesRequest`

## 5. 关键逻辑分析

1. **JWT 认证**: 智谱使用 API Key（格式: `id.secret`）生成 JWT Token，Token 有有效期并缓存。

2. **TopP 修正**: 智谱 API 不支持 TopP >= 1，自动修正为 0.99。

3. **API 路径格式**: 使用 `/api/paas/v3/model-api/{model}/{method}` 格式，method 为 `invoke` 或 `sse-invoke`。

## 6. 关联文件
- `zhipu/constants.go` — 模型列表
- `zhipu/dto.go` — 请求/响应数据结构
- `zhipu/relay-zhipu.go` — JWT 生成、请求/响应转换
