# adaptor.go 代码阅读文档

## 1. 全局总结

本文件实现了 Google Vertex AI 视频生成服务的任务适配器（`TaskAdaptor`）。Vertex 适配器是所有任务适配器中最复杂的之一，因为它需要处理 Google Cloud 的认证机制（OAuth2 Access Token）、从操作名称中解析区域/项目/模型信息、支持多种视频编码格式的 base64 回传，以及与 Gemini 任务包的深度集成。适配器支持 Veo 系列视频生成模型。

## 2. 依赖关系

### 标准库
- `bytes` — 缓冲区操作
- `fmt` — 格式化输出
- `io` — IO 操作
- `net/http` — HTTP 请求/响应
- `regexp` — 正则表达式（从操作名提取信息）
- `strings` — 字符串操作
- `time` — 时间处理

### 项目内部包
- `github.com/QuantumNous/new-api/common` — JSON 操作
- `github.com/QuantumNous/new-api/model` — 任务模型类型
- `github.com/QuantumNous/new-api/constant` — 任务动作常量
- `github.com/QuantumNous/new-api/dto` — 数据传输对象
- `github.com/QuantumNous/new-api/relay/channel` — 通用任务 API 请求方法
- `geminitask`（别名）— Gemini 任务工具（Veo 参数解析、图片提取等）
- `taskcommon`（别名）— 任务通用工具
- `vertexcore`（别名）— Vertex 核心功能（认证、URL 构建）
- `relaycommon`（别名）— 中继公共结构体
- `github.com/QuantumNous/new-api/service` — HTTP 客户端、错误包装

### 第三方库
- `github.com/gin-gonic/gin` — Web 框架

## 3. 类型定义

### 结构体

| 名称 | 字段 | 说明 |
|------|------|------|
| `fetchOperationPayload` | `OperationName` | 查询操作状态的请求负载 |
| `submitResponse` | `Name` | 提交任务后的响应，包含操作名称 |
| `operationVideo` | `MimeType`, `BytesBase64Encoded`, `Encoding` | 视频数据结构（base64 编码） |
| `operationResponse` | `Name`, `Done`, `Response`, `Error` | 操作状态查询响应 |
| `TaskAdaptor` | `BaseBilling`（嵌入）, `ChannelType`, `apiKey`, `baseURL` | Vertex 任务适配器主结构体 |

### 正则表达式变量

| 名称 | 模式 | 用途 |
|------|------|------|
| `regionRe` | `locations/([a-z0-9-]+)/` | 从操作名提取区域 |
| `modelRe` | `models/([^/]+)/operations/` | 从操作名提取模型名 |
| `projectRe` | `projects/([^/]+)/locations/` | 从操作名提取项目 ID |

## 4. 函数详解

### `Init(info *relaycommon.RelayInfo)`
- **作用**: 从 `RelayInfo` 初始化通道类型、API Key 和基础 URL。

### `ValidateRequestAndSetAction(c *gin.Context, info *relaycommon.RelayInfo) *dto.TaskError`
- **作用**: 使用 `ValidateBasicTaskRequest` 进行标准校验，默认动作设为 `TaskActionTextGenerate`。

### `BuildRequestURL(info *relaycommon.RelayInfo) (string, error)`
- **作用**: 构建 Vertex AI 的预测 URL。
- **关键逻辑**:
  1. 从 API Key 解析 credentials（JSON 格式的 GCP 凭证）
  2. 获取模型区域（通过 `GetModelRegion` 或默认 `global`）
  3. 调用 `BuildGoogleModelURL` 构建完整 URL，使用 `predictLongRunning` 方法

### `BuildRequestHeader(c *gin.Context, req *http.Request, info *relaycommon.RelayInfo) error`
- **作用**: 设置请求头，包括 OAuth2 Bearer Token 和项目 ID。
- **认证流程**: 解析 credentials → 调用 `AcquireAccessToken` 获取 OAuth2 token → 设置 Authorization 头。

### `EstimateBilling(c *gin.Context, info *relaycommon.RelayInfo) map[string]float64`
- **作用**: 根据视频时长和分辨率计算计费比率。
- **参数来源**: 从 context 获取 `task_request`，使用 Gemini 工具解析 Veo 参数。
- **返回**: `seconds`（时长）和 `resolution`（分辨率比率）。

### `BuildRequestBody(c *gin.Context, info *relaycommon.RelayInfo) (io.Reader, error)`
- **作用**: 将通用请求转换为 Vertex AI 的 Veo 请求格式。
- **处理流程**:
  1. 构建 `VeoInstance`（包含 prompt 和可选图片）
  2. 支持 multipart 图片和 images 数组两种图片来源
  3. 从 metadata 解析 Veo 参数（duration、resolution、aspect_ratio）
  4. 设置默认值（sample_count=1）
  5. 序列化为 `VeoRequestPayload`

### `DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (*http.Response, error)`
- **作用**: 委托给 `channel.DoTaskApiRequest` 发送请求。

### `DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (string, []byte, *dto.TaskError)`
- **作用**: 解析上游响应，返回任务 ID。
- **关键逻辑**: 使用 `EncodeLocalTaskID` 将操作名编码为本地 ID，构建 OpenAI Video 格式响应。

### `buildFetchOperationURL(baseURL, upstreamName string) (string, error)`
- **作用**: 从操作名中提取区域、项目、模型信息，构建查询 URL。
- **URL 方法**: 使用 `fetchPredictOperation`。

### `FetchTask(baseUrl, key string, body map[string]any, proxy string) (*http.Response, error)`
- **作用**: 查询操作状态。
- **流程**: 解码任务 ID → 构建查询 URL → 获取 OAuth2 token → 发送 POST 请求。

### `ParseTaskResult(respBody []byte) (*relaycommon.TaskInfo, error)`
- **作用**: 解析操作状态响应，映射到统一任务状态。
- **状态映射**: 有错误 → 失败；`Done=false` → 进行中（50%）；`Done=true` → 成功。
- **视频数据处理**: 支持三种格式 — `Videos[].BytesBase64Encoded`、`BytesBase64Encoded`、`Video` 字段。
- **MIME 类型推断**: 从 `MimeType` 或 `Encoding` 字段推断，支持 `video/mp4` 等格式。

### `ConvertToOpenAIVideo(task *model.Task) ([]byte, error)`
- **作用**: 将内部任务数据转换为 OpenAI Video 格式。
- **特殊处理**: 使用 `GetUpstreamTaskID()` 获取真实的上游操作名（而非公共任务 ID）。

### 辅助函数

| 函数 | 说明 |
|------|------|
| `extractRegionFromOperationName(name string) string` | 从操作名中用正则提取区域（如 `us-central1`） |
| `extractModelFromOperationName(name string) string` | 从操作名中提取模型名（如 `veo-3.0-generate-001`） |
| `extractProjectFromOperationName(name string) string` | 从操作名中提取项目 ID |

## 5. 关键逻辑分析

### GCP 认证机制
Vertex 适配器使用 Google Cloud 的 OAuth2 认证：
1. API Key 存储的是 GCP 凭证 JSON（包含 service account 信息）
2. 通过 `vertexcore.AcquireAccessToken` 获取短期 Access Token
3. 每次请求都需要重新获取 token（支持代理）

### 操作名解析
Vertex AI 返回的操作名格式为：
```
projects/{project}/locations/{region}/models/{model}/operations/{op_id}
```
适配器使用三个正则表达式从中提取项目、区域和模型信息，用于构建后续查询 URL。

### 视频数据回传
完成的视频以 base64 编码返回，适配器构建 `data:{mime};base64,{data}` 格式的 URL，支持三种字段位置：
1. `response.videos[0].bytesBase64Encoded` — 多视频数组格式
2. `response.bytesBase64Encoded` — 单视频格式
3. `response.video` — 简化格式

### 与 Gemini 任务包的集成
Vertex 适配器大量复用 `geminitask` 包的功能：
- `VeoInstance`、`VeoParameters`、`VeoRequestPayload` — 请求结构体
- `ResolveVeoDuration`、`ResolveVeoResolution` — 参数解析
- `VeoResolutionRatio` — 分辨率计费比率
- `ExtractMultipartImage`、`ParseImageInput` — 图片处理
- `SizeToVeoResolution`、`SizeToVeoAspectRatio` — 尺寸转换

## 6. 关联文件

- `relay/channel/task/vertex/constants.go`（如存在）— 模型列表
- `relay/channel/task/taskcommon/helpers.go` — `EncodeLocalTaskID`、`DecodeLocalTaskID`、`BaseBilling`
- `relay/channel/task/gemini/` — Gemini 任务工具包（Veo 参数处理）
- `relay/channel/vertex/` — Vertex 核心包（认证、URL 构建）
- `relay/common/relay_info.go` — `RelayInfo` 结构体
- `dto/video.go` — `OpenAIVideo` 类型
