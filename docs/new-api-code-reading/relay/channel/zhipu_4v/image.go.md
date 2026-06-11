# image.go 代码阅读文档

## 1. 全局总结
智谱 AI v4 版渠道的图片生成响应处理文件。负责将智谱图片 API 响应转换为 OpenAI 格式，支持 URL 和 Base64 两种图片格式。

## 2. 依赖关系
- **标准库**: io, net/http
- **内部包**:
  - `github.com/QuantumNous/new-api/common` — JSON 处理
  - `github.com/QuantumNous/new-api/dto` — Usage 结构体
  - `github.com/QuantumNous/new-api/logger` — 日志
  - `github.com/QuantumNous/new-api/relay/common` — RelayInfo
  - `github.com/QuantumNous/new-api/service` — 图片下载、响应体处理
  - `github.com/QuantumNous/new-api/types` — 错误类型
- **外部依赖**:
  - `github.com/gin-gonic/gin` — HTTP 框架

## 3. 类型定义

### zhipuImageRequest
```go
type zhipuImageRequest struct {
    Model            string `json:"model"`
    Prompt           string `json:"prompt"`
    Quality          string `json:"quality,omitempty"`
    Size             string `json:"size,omitempty"`
    WatermarkEnabled *bool  `json:"watermark_enabled,omitempty"`
    UserID           string `json:"user_id,omitempty"`
}
```

### zhipuImageResponse
```go
type zhipuImageResponse struct {
    Created       *int64            `json:"created,omitempty"`
    Data          []zhipuImageData  `json:"data,omitempty"`
    ContentFilter any               `json:"content_filter,omitempty"`
    Usage         *dto.Usage        `json:"usage,omitempty"`
    Error         *zhipuImageError  `json:"error,omitempty"`
    RequestID     string            `json:"request_id,omitempty"`
    ExtendParam   map[string]string `json:"extendParam,omitempty"`
}
```

### zhipuImageData
```go
type zhipuImageData struct {
    Url      string `json:"url,omitempty"`
    ImageUrl string `json:"image_url,omitempty"`
    B64Json  string `json:"b64_json,omitempty"`
    B64Image string `json:"b64_image,omitempty"`
}
```
智谱图片数据，支持 URL 和 Base64 两种格式。

### openAIImagePayload / openAIImageData
OpenAI 格式的图片响应结构。

## 4. 函数详解

### zhipu4vImageHandler(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (*dto.Usage, *types.NewAPIError)
图片生成响应处理器：
1. 读取并解析智谱响应
2. 检查错误
3. 遍历图片数据：
   - 优先使用 `Url` 或 `ImageUrl`
   - 如果没有 URL，尝试使用 `B64Json` 或 `B64Image`
   - 如果都没有，从 URL 下载并转换为 Base64
4. 构建 OpenAI 格式响应并返回

## 5. 关键逻辑分析

1. **多格式支持**: 智谱图片 API 可能返回 URL 或 Base64 编码的图片，处理器统一转换为 Base64 格式。

2. **URL 字段变体**: 智谱使用 `url` 和 `image_url` 两个字段存储图片 URL，需要依次检查。

3. **Base64 字段变体**: 同样使用 `b64_json` 和 `b64_image` 两个字段。

4. **降级策略**: 如果没有直接的 Base64 数据，尝试从 URL 下载并转换，确保返回可用的图片数据。

5. **Created 时间**: 如果响应中没有 `Created` 字段，使用请求开始时间。

## 6. 关联文件
- `zhipu_4v/adaptor.go` — 在 `DoResponse` 中调用 `zhipu4vImageHandler`
- `zhipu_4v/dto.go` — 相关数据结构
- `service/image.go` — 图片下载工具
