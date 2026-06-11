# image.go 代码阅读文档

## 1. 全局总结
本文件处理即梦（Jimeng）图像生成 API 的响应，将即梦特有的响应格式转换为标准 OpenAI 图像响应格式。支持 Base64 和 URL 两种图像返回方式。

## 2. 依赖关系
- **标准库**: `encoding/json`, `fmt`, `io`, `net/http`
- **项目内部**:
  - `github.com/QuantumNous/new-api/dto` — 数据传输对象
  - `github.com/QuantumNous/new-api/relay/common` — RelayInfo
  - `github.com/QuantumNous/new-api/service` — 响应体关闭
  - `github.com/QuantumNous/new-api/types` — 错误类型
- **外部依赖**: `github.com/gin-gonic/gin`

## 3. 类型定义

### `ImageResponse` 结构体
```go
type ImageResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    struct {
        BinaryDataBase64 []string `json:"binary_data_base64"`
        ImageUrls        []string `json:"image_urls"`
        RephraseResult   string   `json:"rephraser_result"`
        RequestID        string   `json:"request_id"`
    } `json:"data"`
    RequestID   string `json:"request_id"`
    Status      int    `json:"status"`
    TimeElapsed string `json:"time_elapsed"`
}
```
即梦 API 的图像生成响应结构，包含状态码、消息、图像数据（Base64 和 URL 两种格式）、请求 ID 和耗时信息。

## 4. 函数详解

### `responseJimeng2OpenAIImage(_, response, info) *dto.ImageResponse`
将即梦响应转换为 OpenAI 图像响应格式：
- 遍历 `BinaryDataBase64` 数组，转换为 `dto.ImageData{B64Json: base64Data}`
- 遍历 `ImageUrls` 数组，转换为 `dto.ImageData{Url: imageUrl}`
- 使用 `info.StartTime` 作为 `Created` 时间戳

### `jimengImageHandler(c, resp, info) (*dto.Usage, *types.NewAPIError)`
图像生成响应处理主函数：
1. 读取完整响应体
2. 解析为 `ImageResponse`
3. 检查响应码（`Code == 10000` 表示成功），失败时返回即梦错误
4. 调用 `responseJimeng2OpenAIImage` 转换格式
5. 序列化为 JSON 并写入响应
6. 返回空 usage（即梦图像生成不返回 token 使用量）

## 5. 关键逻辑分析

### 错误码判断
即梦 API 使用 `Code` 字段表示业务状态，`10000` 为成功状态码。非成功状态返回包含即梦错误码和消息的 OpenAI 兼容错误。

### 图像数据双格式
即梦支持同时返回 Base64 和 URL 两种格式的图像数据，两者都可能出现在同一个响应中。转换函数将两种格式统一追加到 OpenAI 的 `Data` 数组中。

### Usage 返回
图像生成不返回 token 使用量，返回空 `dto.Usage{}`。这是即梦 API 的特性，与 OpenAI 的图像生成 API 行为一致。

## 6. 关联文件
- `adaptor.go` — 在 `DoResponse` 中调用 `jimengImageHandler`
- `constants.go` — 模型列表（`jimeng_high_aes_general_v21_L`）
