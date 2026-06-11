# adaptor.go 代码阅读文档

## 1. 全局总结

本文件实现了 Sora 视频生成服务的任务适配器（`TaskAdaptor`），是整个任务中继系统中 Sora 通道的核心实现。它遵循统一的 `TaskAdaptor` 接口，负责请求校验、请求构建（支持 JSON 和 multipart/form-data 两种格式）、上游请求发送、响应解析、任务状态轮询以及结果转换。Sora 适配器的一个特殊之处在于支持"Remix"操作（基于已有视频生成变体）和 multipart 文件上传。

## 2. 依赖关系

### 标准库
- `bytes` — 缓冲区操作
- `fmt` — 格式化输出
- `io` — IO 操作
- `mime/multipart` — multipart 表单构建
- `net/http` — HTTP 请求/响应
- `net/textproto` — MIME 头部构造
- `strconv` — 字符串转数字
- `strings` — 字符串操作

### 项目内部包
- `github.com/QuantumNous/new-api/common` — JSON 操作、body 缓存、multipart 解析
- `github.com/QuantumNous/new-api/constant` — 任务动作常量
- `github.com/QuantumNous/new-api/dto` — 数据传输对象（TaskError）
- `github.com/QuantumNous/new-api/model` — 任务状态常量
- `github.com/QuantumNous/new-api/relay/channel` — 通用任务 API 请求方法
- `taskcommon`（别名）— 任务通用工具和基础计费结构体
- `relaycommon`（别名）— 中继公共结构体和工具
- `github.com/QuantumNous/new-api/service` — HTTP 客户端、错误包装

### 第三方库
- `github.com/gin-gonic/gin` — Web 框架
- `github.com/pkg/errors` — 错误包装
- `github.com/tidwall/sjson` — JSON 字节级操作

## 3. 类型定义

### 结构体

| 名称 | 字段 | 说明 |
|------|------|------|
| `ContentItem` | `Type`, `Text`, `ImageURL` | 多模态内容项，支持文本和图片 URL |
| `ImageURL` | `URL` | 图片 URL 封装 |
| `responseTask` | `ID`, `TaskID`, `Object`, `Model`, `Status`, `Progress`, `CreatedAt`, `CompletedAt`, `ExpiresAt`, `Seconds`, `Size`, `RemixedFromVideoID`, `Error` | Sora 上游响应的任务结构体，`TaskID` 用于兼容旧接口 |
| `TaskAdaptor` | `BaseBilling`（嵌入）, `ChannelType`, `apiKey`, `baseURL` | Sora 任务适配器主结构体，嵌入 `BaseBilling` 提供默认计费方法 |

### 嵌入结构体

- `taskcommon.BaseBilling` — 提供 `EstimateBilling`、`AdjustBillingOnSubmit`、`AdjustBillingOnComplete` 的默认（无操作）实现，Sora 适配器重写了 `EstimateBilling`。

## 4. 函数详解

### `Init(info *relaycommon.RelayInfo)`
- **作用**: 从 `RelayInfo` 初始化适配器的通道类型、API Key 和基础 URL。
- **参数**: `info` — 中继信息结构体，包含通道配置。

### `validateRemixRequest(c *gin.Context) *dto.TaskError`
- **作用**: 校验 Remix 请求，确保 `prompt` 字段非空，并将解析后的请求存入 context。
- **返回**: 错误时返回 `TaskError`。

### `ValidateRequestAndSetAction(c *gin.Context, info *relaycommon.RelayInfo) *dto.TaskError`
- **作用**: 根据 action 类型分派校验逻辑。Remix 动作走专用校验，其他走通用 multipart 校验。
- **关键逻辑**: 区分 `TaskActionRemix` 和默认路径。

### `EstimateBilling(c *gin.Context, info *relaycommon.RelayInfo) map[string]float64`
- **作用**: 根据用户请求的 `seconds` 和 `size` 计算计费比率。
- **返回**: `map[string]float64`，包含 `seconds`（时长）和 `size`（分辨率系数）。
- **计费规则**:
  - 默认时长 4 秒，默认分辨率 `720x1280`（系数 1）
  - `1792x1024` 或 `1024x1792` 分辨率系数为 1.666667
  - Remix 路径跳过此方法（由 `ResolveOriginTask` 设置）

### `BuildRequestURL(info *relaycommon.RelayInfo) (string, error)`
- **作用**: 构建上游请求 URL。
- **逻辑**: Remix 返回 `/v1/videos/{id}/remix`，其他返回 `/v1/videos`。

### `BuildRequestHeader(c *gin.Context, req *http.Request, info *relaycommon.RelayInfo) error`
- **作用**: 设置请求头，包括 Bearer Token 和原始 Content-Type。

### `BuildRequestBody(c *gin.Context, info *relaycommon.RelayInfo) (io.Reader, error)`
- **作用**: 构建请求体，是本文件最复杂的方法。
- **JSON 路径**: 解析 body map，替换 `model` 字段为上游模型名。
- **Multipart 路径**: 重建 multipart 表单，替换 `model` 字段，处理文件上传（含 Content-Type 自动检测）。
- **其他**: 直接返回原始 body。

### `DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (*http.Response, error)`
- **作用**: 委托给 `channel.DoTaskApiRequest` 发送请求。

### `DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (string, []byte, *dto.TaskError)`
- **作用**: 解析上游响应，返回任务 ID。
- **关键逻辑**: 支持 `id` 和 `task_id` 双字段兼容；用 `PublicTaskID` 替换上游 ID 返回给客户端。

### `FetchTask(baseUrl, key string, body map[string]any, proxy string) (*http.Response, error)`
- **作用**: 轮询任务状态，向 `/v1/videos/{task_id}` 发送 GET 请求。

### `GetModelList() []string` / `GetChannelName() string`
- **作用**: 返回模型列表和通道名称常量。

### `ParseTaskResult(respBody []byte) (*relaycommon.TaskInfo, error)`
- **作用**: 解析上游轮询响应，映射到统一任务状态。
- **状态映射**: `queued/pending` → 队列中, `processing/in_progress` → 进行中, `completed` → 成功, `failed/cancelled` → 失败。

### `ConvertToOpenAIVideo(task *model.Task) ([]byte, error)`
- **作用**: 将内部任务数据转换为 OpenAI Video 格式，使用 `sjson` 设置 `id` 字段。

## 5. 关键逻辑分析

### Multipart 请求重建
`BuildRequestBody` 方法中的 multipart 处理是核心复杂度所在：
1. 解析原始 multipart 表单
2. 重建新的 multipart writer
3. 替换 `model` 字段
4. 遍历文件字段，自动检测 Content-Type（通过读取前 512 字节嗅探）
5. 重新设置请求的 Content-Type 头

### 公共任务 ID 机制
Sora 适配器在 `DoResponse` 中使用 `info.PublicTaskID` 替换上游返回的任务 ID，实现了上游 ID 与客户端可见 ID 的解耦，增强了安全性。

### Remix 操作支持
Remix 是 Sora 的特色功能，允许基于已有视频生成变体。适配器通过 `info.Action` 判断是否走 Remix 路径，影响 URL 构建、请求校验和计费逻辑。

## 6. 关联文件

- `relay/channel/task/sora/constants.go` — 模型列表和通道名称常量
- `relay/channel/task/taskcommon/helpers.go` — `BaseBilling` 嵌入结构体
- `relay/channel/task/channel/task_api.go` — `DoTaskApiRequest` 通用请求方法
- `relay/common/relay_info.go` — `RelayInfo` 结构体定义
- `relay/common/task_submit.go` — `TaskSubmitReq` 和 `ValidateMultipartDirect`
- `dto/task.go` — `TaskError` 类型
- `model/task.go` — `TaskStatus*` 常量
