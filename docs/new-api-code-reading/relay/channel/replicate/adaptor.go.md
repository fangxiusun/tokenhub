# adaptor.go 代码阅读文档

## 1. 全局总结
本文件实现了 Replicate AI 平台的渠道适配器。Replicate 是一个 AI 模型托管平台，本适配器主要支持图像生成任务（使用 FLUX 模型）。适配器将 OpenAI 格式的图像请求转换为 Replicate 的预测（Prediction）API 格式，处理文件上传、尺寸映射、响应解析等完整流程。

## 2. 依赖关系
| 包名 | 用途 |
|------|------|
| `bytes` | 缓冲区操作（文件上传） |
| `encoding/json` | JSON 处理 |
| `errors` | 错误创建 |
| `fmt` | 格式化字符串 |
| `io` | IO 操作 |
| `mime/multipart` | 多部分表单处理（文件上传） |
| `net/http` | HTTP 请求/响应 |
| `net/textproto` | MIME 头部处理 |
| `strconv` | 字符串转换（尺寸解析） |
| `strings` | 字符串操作 |
| `github.com/QuantumNous/new-api/common` | JSON 封装、时间戳 |
| `github.com/QuantumNous/new-api/constant` | 常量定义（频道基础 URL） |
| `github.com/QuantumNous/new-api/dto` | 数据传输对象 |
| `github.com/QuantumNous/new-api/relay/channel` | 渠道公共功能 |
| `github.com/QuantumNous/new-api/relay/common` | 中继公共类型 |
| `github.com/QuantumNous/new-api/relay/constant` | 中继常量（RelayModeImagesEdits） |
| `github.com/QuantumNous/new-api/service` | 服务层（HTTP 客户端、图片下载） |
| `github.com/QuantumNous/new-api/types` | 类型定义 |
| `github.com/gin-gonic/gin` | Gin HTTP 框架 |
| `github.com/samber/lo` | 辅助函数库 |

## 3. 类型定义

### `Adaptor`
Replicate 渠道适配器结构体，无字段。

## 4. 函数详解

### `Init`
```go
func (a *Adaptor) Init(info *relaycommon.RelayInfo)
```
- 空实现

### `GetRequestURL`
```go
func (a *Adaptor) GetRequestURL(info *relaycommon.RelayInfo) (string, error)
```
- **逻辑**：
  1. 检查 `info` 非 nil
  2. 如果 `ChannelBaseUrl` 为空，使用默认的 Replicate 基础 URL
  3. 如果 `RequestURLPath` 为空，返回基础 URL
  4. 否则拼接完整 URL

### `SetupRequestHeader`
```go
func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error
```
- **逻辑**：
  1. 验证 `info` 和 `ApiKey` 非空
  2. 设置标准 `Authorization: Bearer` 头
  3. 设置 `Prefer: wait` 头（Replicate 同步等待特性）
  4. 设置默认的 `Content-Type` 和 `Accept` 为 `application/json`

### `ConvertImageRequest`
```go
func (a *Adaptor) ConvertImageRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.ImageRequest) (any, error)
```
- **说明**：核心转换方法，将 OpenAI 图像请求转换为 Replicate 预测格式
- **逻辑**：
  1. 验证 prompt 非空（支持从表单字段获取）
  2. 确定模型名称（优先级：`UpstreamModelName` > `request.Model` > 默认 `ModelFlux11Pro`）
  3. 设置请求路径为 `/v1/models/{model}/predictions`
  4. 构造输入负载：
     - 基础字段：`prompt`
     - 尺寸映射：调用 `mapOpenAISizeToFlux` 转换
     - 输出格式、数量、质量参数
  5. 处理图像编辑模式（`RelayModeImagesEdits`）：上传参考图片
  6. 合并 `ExtraFields` 和 `Extra` 中的额外参数
  7. 返回 `{"input": {...}}` 格式

### `DoRequest`
```go
func (a *Adaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (any, error)
```
- 委托给 `channel.DoApiRequest`

### `DoResponse`
```go
func (a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (any, *types.NewAPIError)
```
- **说明**：处理 Replicate 预测响应，转换为 OpenAI 图像响应格式
- **逻辑**：
  1. 读取并解析响应为 `PredictionResponse`
  2. 检查预测错误和状态
  3. 提取输出 URL（支持单个字符串和数组）
  4. 根据 `ResponseFormat` 决定返回 URL 或 Base64
  5. 如果需要 Base64，调用 `downloadImagesToBase64` 下载并转换
  6. 构造 `ImageResponse` 并写入响应

### `ConvertOpenAIRequest` / `ConvertRerankRequest` / `ConvertEmbeddingRequest` / `ConvertAudioRequest` / `ConvertOpenAIResponsesRequest` / `ConvertClaudeRequest` / `ConvertGeminiRequest`
- 均返回 "not implemented" 错误，说明该渠道仅支持图像生成

### `downloadImagesToBase64`
```go
func downloadImagesToBase64(urls []string) ([]string, error)
```
- 遍历 URL 列表，调用 `service.GetImageFromUrl` 下载并转换为 Base64

### `mapOpenAISizeToFlux`
```go
func mapOpenAISizeToFlux(size string) (aspect string, width int, height int, ok bool)
```
- **说明**：将 OpenAI 的尺寸字符串（如 "1024x1024"）映射为 FLUX 支持的宽高比
- **支持的宽高比**：1:1、16:9、9:16、3:2、2:3、4:5、5:4、3:4、4:3
- **回退策略**：如果不匹配已知比例，使用 `custom` 模式并归一化尺寸

### `reduceRatio` / `gcd`
- 辅助函数，用于计算和约简宽高比

### `normalizeFluxDimension`
```go
func normalizeFluxDimension(value int) int
```
- 将尺寸值归一化到 256-1440 范围内，步长为 32

### `uploadFileFromForm`
```go
func uploadFileFromForm(c *gin.Context, info *relaycommon.RelayInfo, fieldCandidates ...string) (string, error)
```
- **说明**：从表单中提取图片文件并上传到 Replicate
- **逻辑**：
  1. 解析多部分表单
  2. 按候选字段名查找文件
  3. 构造上传请求（multipart/form-data）
  4. 发送到 `/v1/files` 端点
  5. 返回上传后的文件 URL

## 5. 关键逻辑分析
- **同步等待模式**：通过 `Prefer: wait` 头让 Replicate 同步返回结果，简化了异步轮询逻辑。
- **尺寸映射**：实现了 OpenAI 尺寸到 FLUX 宽高比的完整映射，支持标准比例和自定义尺寸。
- **文件上传**：支持图像编辑模式，需要先将参考图片上传到 Replicate。
- **响应格式**：支持 URL 和 Base64 两种返回格式，通过 `ResponseFormat` 控制。
- **Extra 参数合并**：支持通过 `ExtraFields` 和 `Extra` 传递自定义参数，增强了灵活性。
- **nil 检查**：所有公共方法都进行了 nil 检查，提高了健壮性。

## 6. 关联文件
- `relay/channel/replicate/constants.go` — 模型列表和渠道常量
- `relay/channel/replicate/dto.go` — Replicate API 数据结构
- `relay/channel/adapter.go` — `Adaptor` 接口定义
- `relay/channel/channel.go` — 公共请求发送逻辑
- `service/` — HTTP 客户端和图片下载工具
- `constant/` — 频道基础 URL 定义
