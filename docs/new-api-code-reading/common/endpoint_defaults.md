# endpoint_defaults.go 代码阅读文档

## 1. 全局总结

本文件定义了 API 端点的默认配置信息，包括端点类型结构体和内置端点的默认路径/方法映射表。这些默认配置用于为不同 AI 提供商的 API 端点提供标准的请求模板。

## 2. 依赖关系

| 包名 | 用途 |
|------|------|
| `github.com/QuantumNous/new-api/constant` | 端点类型常量定义 |

**被依赖方**：被端点路由和请求处理代码调用，用于获取端点的默认配置。

## 3. 类型定义

### `EndpointInfo` 结构体

```go
type EndpointInfo struct {
    Path   string `json:"path"`
    Method string `json:"method"`
}
```

描述单个端点的默认请求信息。

| 字段 | 类型 | JSON 标签 | 说明 |
|------|------|-----------|------|
| `Path` | `string` | `"path"` | 上游 API 路径 |
| `Method` | `string` | `"method"` | HTTP 请求方法（目前均为 POST） |

**JSON 输出示例**：
```json
{"path":"/v1/chat/completions","method":"POST"}
```

### `defaultEndpointInfoMap` 变量

```go
var defaultEndpointInfoMap = map[constant.EndpointType]EndpointInfo{...}
```

保存内置端点类型的默认 Path 与 Method 映射。

## 4. 函数详解

### `GetDefaultEndpointInfo(et constant.EndpointType) (EndpointInfo, bool)`

**功能**：获取指定端点类型的默认信息。

**参数**：
- `et constant.EndpointType` — 端点类型枚举值

**返回值**：
- `EndpointInfo` — 端点默认信息
- `bool` — 是否存在该端点类型的默认配置

**逻辑**：直接从 `defaultEndpointInfoMap` 中查找并返回。

## 5. 关键逻辑分析

### 内置端点映射表

| 端点类型 | 默认路径 | HTTP 方法 | 说明 |
|----------|----------|-----------|------|
| `EndpointTypeOpenAI` | `/v1/chat/completions` | POST | OpenAI Chat Completions API |
| `EndpointTypeOpenAIResponse` | `/v1/responses` | POST | OpenAI Responses API |
| `EndpointTypeOpenAIResponseCompact` | `/v1/responses/compact` | POST | OpenAI Responses 紧凑版 API |
| `EndpointTypeAnthropic` | `/v1/messages` | POST | Anthropic Messages API |
| `EndpointTypeGemini` | `/v1beta/models/{model}:generateContent` | POST | Google Gemini API |
| `EndpointTypeJinaRerank` | `/v1/rerank` | POST | Jina Rerank API |
| `EndpointTypeImageGeneration` | `/v1/images/generations` | POST | 图像生成 API |
| `EndpointTypeEmbeddings` | `/v1/embeddings` | POST | 文本嵌入 API |

### 路径模板变量

Gemini 端点路径中包含 `{model}` 占位符，表示模型名称需要在运行时动态替换。

### 设计意图

- 所有端点方法目前均为 POST，但结构体设计支持未来扩展其他 HTTP 方法
- JSON 标签使得 `EndpointInfo` 可以直接序列化到 API 输出
- 通过 `GetDefaultEndpointInfo` 函数封装访问逻辑，支持后续添加缓存或验证

## 6. 关联文件

| 文件 | 关联关系 |
|------|----------|
| `constant/endpoint.go` | 定义 `EndpointType` 常量 |
| `endpoint_type.go` | 根据渠道类型确定端点类型 |
| `relay/` | 使用端点信息构建上游请求 |
