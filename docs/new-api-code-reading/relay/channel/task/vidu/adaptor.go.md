# adaptor.go 代码阅读文档

## 1. 全局总结

本文件实现了 Vidu 视频生成服务的任务适配器（`TaskAdaptor`）。Vidu 是一个支持多种视频生成模式的服务，包括文生视频、图生视频、首尾帧生视频和参考图生视频。适配器根据输入图片数量自动推断动作类型，并将通用请求格式转换为 Vidu 特有的 API 格式。Vidu 使用 Token 认证方式（而非 Bearer），且任务结果中直接包含视频 URL。

## 2. 依赖关系

### 标准库
- `bytes` — 缓冲区操作
- `fmt` — 格式化输出
- `io` — IO 操作
- `net/http` — HTTP 请求/响应
- `strings` — 字符串操作
- `time` — 时间处理

### 项目内部包
- `github.com/QuantumNous/new-api/common` — JSON 操作
- `github.com/QuantumNous/new-api/constant` — 任务动作和通道类型常量
- `github.com/QuantumNous/new-api/dto` — 数据传输对象
- `github.com/QuantumNous/new-api/model` — 任务模型类型
- `github.com/QuantumNous/new-api/relay/channel` — 通用任务 API 请求方法
- `taskcommon`（别名）— 任务通用工具
- `relaycommon`（别名）— 中继公共结构体
- `github.com/QuantumNous/new-api/service` — HTTP 客户端、错误包装

### 第三方库
- `github.com/gin-gonic/gin` — Web 框架
- `github.com/pkg/errors` — 错误包装

## 3. 类型定义

### 结构体

| 名称 | 字段 | 说明 |
|------|------|------|
| `requestPayload` | `Model`, `Images`, `Prompt`, `Duration`, `Seed`, `Resolution`, `MovementAmplitude`, `Bgm`, `Payload`, `CallbackUrl` | Vidu 上游请求负载 |
| `responsePayload` | `TaskId`, `State`, `Model`, `Images`, `Prompt`, `Duration`, `Seed`, `Resolution`, `Bgm`, `MovementAmplitude`, `Payload`, `CreatedAt` | Vidu 任务提交响应 |
| `taskResultResponse` | `State`, `ErrCode`, `Credits`, `Payload`, `Creations` | Vidu 任务状态查询响应 |
| `creation` | `ID`, `URL`, `CoverURL` | 单个创作结果（视频 URL 和封面） |
| `TaskAdaptor` | `BaseBilling`（嵌入）, `ChannelType`, `baseURL` | Vidu 任务适配器主结构体 |

## 4. 函数详解

### `Init(info *relaycommon.RelayInfo)`
- **作用**: 从 `RelayInfo` 初始化通道类型和基础 URL。
- **注意**: 不存储 `apiKey`，因为 Vidu 使用 Token 认证，key 在请求头中直接使用。

### `ValidateRequestAndSetAction(c *gin.Context, info *relaycommon.RelayInfo) *dto.TaskError`
- **作用**: 校验请求并根据输入自动推断动作类型。
- **动作推断逻辑**:
  1. 首先检查 `metadata["action"]` — 如果指定则使用
  2. 有图片时：
     - 1 张图片 → `TaskActionGenerate`（图生视频）
     - 2 张图片 → `TaskActionFirstTailGenerate`（首尾帧生视频）
     - 3+ 张图片 → `TaskActionReferenceGenerate`（参考图生视频）
  3. 无图片 → `TaskActionTextGenerate`（文生视频）
- **通道限制**: 首尾帧和参考图生视频仅对 `ChannelTypeVidu` 通道生效。

### `BuildRequestBody(c *gin.Context, info *relaycommon.RelayInfo) (io.Reader, error)`
- **作用**: 将通用请求转换为 Vidu 请求格式。
- **特殊处理**: 参考图生视频时，如果模型名包含 `viduq2`，强制去除 `pro`/`turbo` 后缀（API 限制）。

### `BuildRequestURL(info *relaycommon.RelayInfo) (string, error)`
- **作用**: 根据动作类型构建不同的 API 端点。
- **路由映射**:
  - `generate` → `/ent/v2/img2video`
  - `first_tail_generate` → `/ent/v2/start-end2video`
  - `reference_generate` → `/ent/v2/reference2video`
  - 默认（text_generate）→ `/ent/v2/text2video`

### `BuildRequestHeader(c *gin.Context, req *http.Request, info *relaycommon.RelayInfo) error`
- **作用**: 设置请求头。
- **特殊之处**: 使用 `Token` 认证（而非 `Bearer`），格式为 `Authorization: Token {key}`。

### `DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (*http.Response, error)`
- **作用**: 委托给 `channel.DoTaskApiRequest` 发送请求。

### `DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (string, []byte, *dto.TaskError)`
- **作用**: 解析上游响应，返回任务 ID。
- **错误处理**: 如果 `state` 为 `failed`，直接返回错误。
- **响应格式**: 构建 OpenAI Video 格式响应，使用 `PublicTaskID`。

### `FetchTask(baseUrl, key string, body map[string]any, proxy string) (*http.Response, error)`
- **作用**: 查询任务状态，向 `/ent/v2/tasks/{task_id}/creations` 发送 GET 请求。
- **认证**: 使用 `Token` 认证。

### `GetModelList() []string`
- **作用**: 返回支持的模型列表：`viduq2`, `viduq1`, `vidu2.0`, `vidu1.5`。

### `GetChannelName() string`
- **作用**: 返回通道名称 `"vidu"`。

### `convertToRequestPayload(req *relaycommon.TaskSubmitReq, info *relaycommon.RelayInfo) (*requestPayload, error)`
- **作用**: 将通用请求转换为 Vidu 请求负载。
- **默认值**: 模型 `viduq1`、时长 5 秒、分辨率 `1080p`、运动幅度 `auto`、背景音乐 `false`。
- **元数据处理**: 使用 `UnmarshalMetadata` 从 metadata 覆盖默认值。

### `ParseTaskResult(respBody []byte) (*relaycommon.TaskInfo, error)`
- **作用**: 解析任务状态响应，映射到统一状态。
- **状态映射**:
  - `created/queueing` → 已提交
  - `processing` → 进行中
  - `success` → 成功（提取第一个 creation 的 URL）
  - `failed` → 失败（记录 err_code）
  - 其他 → 返回未知状态错误

### `ConvertToOpenAIVideo(originTask *model.Task) ([]byte, error)`
- **作用**: 将内部任务数据转换为 OpenAI Video 格式。
- **数据提取**: 从 `taskResultResponse` 中提取视频 URL 和错误信息。
- **错误处理**: 失败时设置 `Error` 字段。

## 5. 关键逻辑分析

### 智能动作推断
Vidu 适配器的 `ValidateRequestAndSetAction` 实现了基于图片数量的智能动作推断，这是与其他适配器的显著区别。这种设计让客户端无需显式指定动作类型，系统自动根据输入选择合适的生成模式。

### Token 认证
与大多数使用 Bearer Token 的适配器不同，Vidu 使用 `Token` 认证方案（`Authorization: Token {key}`），这在 API 设计上更接近传统的 API Key 认证。

### 视频 URL 直接回传
Vidu 的任务结果直接包含视频 URL（通过 `creations[0].url`），而不需要像 Vertex 那样处理 base64 编码数据。这简化了结果处理逻辑。

### 模型名规范化
在参考图生视频模式下，Vidu API 要求使用纯净的 `viduq2` 模型名（不能带 `pro`/`turbo` 后缀），适配器在 `BuildRequestBody` 中自动处理这一限制。

### 多种视频生成模式
Vidu 支持四种视频生成模式，通过不同的 API 端点实现：
1. 文生视频（`text2video`）— 纯文本提示
2. 图生视频（`img2video`）— 单张参考图
3. 首尾帧生视频（`start-end2video`）— 两张图片作为首尾帧
4. 参考图生视频（`reference2video`）— 多张参考图

## 6. 关联文件

- `relay/channel/task/taskcommon/helpers.go` — `DefaultString`、`DefaultInt`、`UnmarshalMetadata`、`BaseBilling`
- `relay/common/relay_info.go` — `TaskSubmitReq` 结构体（含 `Images`、`Metadata` 等字段）
- `constant/` — 任务动作常量（`TaskActionGenerate`、`TaskActionTextGenerate` 等）
- `dto/video.go` — `OpenAIVideo` 类型
- `model/task.go` — `TaskStatus*` 常量、`Task` 模型
- `relay/channel/task/channel/task_api.go` — `DoTaskApiRequest` 通用请求方法
