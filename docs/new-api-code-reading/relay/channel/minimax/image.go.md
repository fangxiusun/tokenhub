# image.go 代码阅读文档

## 1. 全局总结
本文件实现了 MiniMax 渠道的图像生成功能，包括请求转换、宽高比计算、响应格式标准化和响应处理。支持将 OpenAI ImageRequest 转换为 MiniMax 格式，并将 MiniMax 响应转换为 OpenAI 格式。

## 2. 依赖关系
- **标准库**: `fmt`, `io`, `net/http`, `strconv`, `strings`
- **项目内部**:
  - `github.com/QuantumNous/new-api/common` — JSON 处理
  - `github.com/QuantumNous/new-api/dto` — 数据传输对象
  - `github.com/QuantumNous/new-api/relay/common` — RelayInfo
  - `github.com/QuantumNous/new-api/service` — 响应体关闭
  - `github.com/QuantumNous/new-api/types` — 错误类型
- **外部依赖**: `github.com/gin-gonic/gin`

## 3. 类型定义

### `MiniMaxImageRequest` 结构体
```go
type MiniMaxImageRequest struct {
    Model           string `json:"model"`
    Prompt          string `json:"prompt"`
    AspectRatio     string `json:"aspect_ratio,omitempty"`
    ResponseFormat  string `json:"response_format,omitempty"`
    N               int    `json:"n,omitempty"`
    PromptOptimizer *bool  `json:"prompt_optimizer,omitempty"`
    AigcWatermark   *bool  `json:"aigc_watermark,omitempty"`
}
```
MiniMax 图像生成请求，包含模型、提示词、宽高比、响应格式、生成数量、提示词优化器和水印开关。

### `MiniMaxImageResponse` 结构体
```go
type MiniMaxImageResponse struct {
    ID   string `json:"id"`
    Data struct {
        ImageURLs   []string `json:"image_urls"`
        ImageBase64 []string `json:"image_base64"`
    } `json:"data"`
    Metadata map[string]any `json:"metadata"`
    BaseResp struct {
        StatusCode int    `json:"status_code"`
        StatusMsg  string `json:"status_msg"`
    } `json:"base_resp"`
}
```
MiniMax 图像生成响应，包含 ID、图像数据（URL 和 Base64 两种格式）、元数据和状态信息。

## 4. 函数详解

### `oaiImage2MiniMaxImageRequest(request) MiniMaxImageRequest`
将 OpenAI ImageRequest 转换为 MiniMax 格式：
- 标准化响应格式
- 设置默认模型为 `image-01`
- 从 `Extra["aspect_ratio"]` 或 `Size` 计算宽高比
- 从 `Extra["prompt_optimizer"]` 提取提示词优化器设置

### `aspectRatioFromImageRequest(request) string`
从 OpenAI ImageRequest 计算宽高比：
1. 优先从 `Extra["aspect_ratio"]` 获取
2. 根据 `Size` 字符串映射（如 `1024x1024` → `1:1`）
3. 解析 `Size` 字符串计算宽高比，使用 GCD 约分
4. 只返回 MiniMax 支持的比例（`1:1`, `16:9`, `4:3`, `3:2`, `2:3`, `3:4`, `9:16`, `21:9`）

### `parseImageSize(size) (width, height int, ok bool)`
解析 `WxH` 格式的尺寸字符串。

### `reduceAspectRatio(width, height) string`
使用 GCD 约分计算最简宽高比。

### `gcd(a, b) int`
计算最大公约数。

### `normalizeMiniMaxResponseFormat(responseFormat) string`
标准化响应格式：`""` / `"url"` → `"url"`，`"b64_json"` / `"base64"` → `"base64"`。

### `responseMiniMax2OpenAIImage(response, info) (*dto.ImageResponse, error)`
将 MiniMax 响应转换为 OpenAI 格式，同时处理 URL 和 Base64 两种图像数据，复制 metadata。

### `miniMaxImageHandler(c, resp, info) (*dto.Usage, *types.NewAPIError)`
图像生成响应处理主函数：
1. 读取并解析响应
2. 检查 `BaseResp.StatusCode` 是否为 0
3. 转换为 OpenAI 格式
4. 写入 JSON 响应

## 5. 关键逻辑分析

### 宽高比计算策略
MiniMax 使用宽高比（`aspect_ratio`）而非具体尺寸。转换器首先尝试从 `Extra` 字段获取显式比例，然后根据 `Size` 映射，最后尝试解析和约分。只有 MiniMax 支持的特定比例才会被传递。

### GCD 约分
使用欧几里得算法计算最大公约数，将任意尺寸约分为最简比例。例如 `1536x1024` → `3:2`。

### 响应格式标准化
`normalizeMiniMaxResponseFormat` 将 OpenAI 的多种格式名称统一为 MiniMax 支持的 `"url"` 或 `"base64"`。

### 元数据透传
MiniMax 图像响应中的 `metadata` 字段被序列化后传递到 OpenAI 响应的 `metadata` 字段，保留了厂商特有的信息。

## 6. 关联文件
- `adaptor.go` — 在 `ConvertImageRequest` 中调用 `oaiImage2MiniMaxImageRequest`
- `constants.go` — 包含图像模型（`image-01`, `image-01-live`）
