# url_builder.go 代码阅读文档

## 1. 全局总结
Vertex AI 渠道的 URL 构建器，负责根据不同的模型类型（Google/Anthropic/开源）和认证方式构建正确的 API 端点 URL。

## 2. 依赖关系
- **标准库**: fmt, strings

## 3. 类型定义

### 常量
```go
const (
    DefaultAPIVersion    = "v1"
    OpenSourceAPIVersion = "v1beta1"
    PublisherGoogle      = "google"
    PublisherAnthropic   = "anthropic"
)
```

## 4. 函数详解

### normalizeVertexBaseURL(baseURL string) string
清理基础 URL：去除首尾空格和末尾斜杠。

### normalizeVertexRegion(region string) string
标准化区域名：空值默认为 `"global"`。

### appendVertexAPIVersion(baseURL, version string) string
向基础 URL 追加 API 版本，避免重复追加。

### BuildAPIBaseURL(baseURL, version, projectID, region string) string
构建 Vertex AI API 基础 URL：
- 如果提供了自定义 baseURL：拼接版本号和项目/区域路径
- 如果未提供 baseURL：
  - 有 projectID: `https://{region}-aiplatform.googleapis.com/{version}/projects/{projectID}/locations/{region}`
  - 无 projectID: `https://{region}-aiplatform.googleapis.com/{version}`
  - "global" 区域使用 `aiplatform.googleapis.com`（无区域前缀）

### BuildPublisherModelURL(baseURL, version, projectID, region, publisher, modelName, action string) string
构建发布者模型 URL：`{base}/publishers/{publisher}/models/{modelName}:{action}`

### BuildGoogleModelURL(...)
Google 模型 URL 构建器（publisher 为 "google"）。

### BuildAnthropicModelURL(...)
Anthropic 模型 URL 构建器（publisher 为 "anthropic"）。

### BuildOpenSourceChatCompletionsURL(baseURL, projectID, region string) string
开源模型的 Chat Completions URL：`{base}/endpoints/openapi/chat/completions`

## 5. 关键逻辑分析

1. **URL 结构**: Vertex AI 的 URL 格式为 `https://{region}-aiplatform.googleapis.com/{version}/projects/{projectID}/locations/{region}/publishers/{publisher}/models/{model}:{action}`

2. **全局区域**: 当区域为 "global" 时，URL 中不包含区域前缀，直接使用 `aiplatform.googleapis.com`

3. **版本管理**: 使用 `v1` 作为默认版本，开源模型使用 `v1beta1`

4. **动作后缀**: 模型 URL 以 `:{action}` 结尾，如 `:generateContent`、`:streamGenerateContent?alt=sse` 等

## 6. 关联文件
- `vertex/adaptor.go` — 在 `getRequestUrl` 中调用各种构建器
- `vertex/relay-vertex.go` — 提供区域信息
