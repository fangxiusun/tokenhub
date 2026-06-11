# dto.go 代码阅读文档

## 1. 全局总结

本文件定义了 Cloudflare Workers AI 频道特有的数据传输对象（DTO），包括请求结构和响应结构，用于与 Cloudflare API 进行数据交换。

## 2. 依赖关系

| 包 | 用途 |
|---|------|
| `dto` | 通用数据传输对象（Message 类型） |

## 3. 类型定义

### `CfRequest` 结构体

```go
type CfRequest struct {
    Messages    []dto.Message `json:"messages,omitempty"`
    Lora        string        `json:"lora,omitempty"`
    MaxTokens   uint          `json:"max_tokens,omitempty"`
    Prompt      string        `json:"prompt,omitempty"`
    Raw         bool          `json:"raw,omitempty"`
    Stream      bool          `json:"stream,omitempty"`
    Temperature *float64      `json:"temperature,omitempty"`
}
```

Cloudflare Workers AI 的请求格式。用于旧版 Completions API 的请求转换：
- `Messages`: 消息列表（标准 OpenAI 格式）
- `Lora`: 可选的 LoRA 适配器标识
- `MaxTokens`: 最大生成 token 数
- `Prompt`: 文本提示词（旧版 Completions API 使用）
- `Raw`: 是否返回原始模型输出
- `Stream`: 是否启用流式输出
- `Temperature`: 温度参数（指针类型，支持显式零值）

### `CfAudioResponse` 结构体

```go
type CfAudioResponse struct {
    Result CfSTTResult `json:"result"`
}
```

Cloudflare 语音转文本（STT）API 的响应格式。

### `CfSTTResult` 结构体

```go
type CfSTTResult struct {
    Text string `json:"text"`
}
```

STT 结果数据，包含识别出的文本内容。

## 4. 函数详解

无函数定义。

## 5. 关键逻辑分析

### 温度参数使用指针类型
`Temperature *float64` 使用指针类型是为了遵循 Rule 6（保留显式零值），当用户显式设置 `temperature: 0` 时不会被 `omitempty` 吞掉。

### 三种 DTO 的用途区分
- `CfRequest`: 用于将通用 Completions 请求转换为 Cloudflare 格式
- `CfAudioResponse` / `CfSTTResult`: 用于解析 Cloudflare STT API 的响应

## 6. 关联文件

- `relay/channel/cloudflare/relay_cloudflare.go` - 使用 `CfRequest` 构建请求，使用 `CfAudioResponse` 解析 STT 响应
- `relay/channel/cloudflare/adaptor.go` - `ConvertOpenAIRequest` 中调用转换函数生成 `CfRequest`
