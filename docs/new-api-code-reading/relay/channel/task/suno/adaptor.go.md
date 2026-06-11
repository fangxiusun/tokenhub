# adaptor.go 代码阅读文档

## 1. 全局总结

本文件实现了 Suno 音乐生成服务的任务适配器（`TaskAdaptor`）。Suno 适配器与视频类适配器有显著差异：它使用专用的批量轮询机制（`UpdateSunoTasks`）而非单任务轮询，因此 `ParseTaskResult` 方法不适用。适配器支持两种动作：音乐生成（`generate`）和歌词生成（`lyrics`），通过 URL 路径参数区分。

## 2. 依赖关系

### 标准库
- `bytes` — 缓冲区操作
- `fmt` — 格式化输出
- `io` — IO 操作
- `net/http` — HTTP 请求/响应
- `strings` — 字符串操作

### 项目内部包
- `github.com/QuantumNous/new-api/common` — JSON 操作、body 缓存
- `github.com/QuantumNous/new-api/constant` — Suno 动作常量
- `github.com/QuantumNous/new-api/dto` — 数据传输对象
- `github.com/QuantumNous/new-api/relay/channel` — 通用任务 API 请求方法
- `taskcommon`（别名）— 基础计费结构体
- `relaycommon`（别名）— 中继公共结构体
- `github.com/QuantumNous/new-api/service` — 错误包装

### 第三方库
- `github.com/gin-gonic/gin` — Web 框架

## 3. 类型定义

### 结构体

| 名称 | 字段 | 说明 |
|------|------|------|
| `TaskAdaptor` | `BaseBilling`（嵌入）, `ChannelType` | Suno 任务适配器，嵌入 `BaseBilling` 提供默认计费方法 |

注意：Suno 适配器不存储 `apiKey` 和 `baseURL`，因为这些信息在 `BuildRequestURL` 和 `BuildRequestHeader` 中通过 `info` 参数获取。

## 4. 函数详解

### `ParseTaskResult([]byte) (*relaycommon.TaskInfo, error)`
- **作用**: 直接返回错误，表明 Suno 不使用单任务轮询。
- **原因**: Suno 使用批量轮询路径（`service.UpdateSunoTasks`），通过专用的 `/fetch` API 获取任务状态，与视频类适配器的逐任务轮询机制不同。

### `Init(info *relaycommon.RelayInfo)`
- **作用**: 仅初始化通道类型。

### `ValidateRequestAndSetAction(c *gin.Context, info *relaycommon.RelayInfo) *dto.TaskError`
- **作用**: 从 URL 路径参数获取 action，解析请求体，执行动作校验，将请求存入 context。
- **关键逻辑**: action 从 URL 路径（如 `/suno/submit/generate`）提取，而非请求体。

### `BuildRequestURL(info *relaycommon.RelayInfo) (string, error)`
- **作用**: 构建上游 URL，格式为 `{baseURL}/suno/submit/{action}`。

### `BuildRequestHeader(c *gin.Context, req *http.Request, info *relaycommon.RelayInfo) error`
- **作用**: 设置 Content-Type、Accept 和 Authorization 头。

### `BuildRequestBody(c *gin.Context, info *relaycommon.RelayInfo) (io.Reader, error)`
- **作用**: 从 context 获取已解析的请求对象，序列化为 JSON。
- **与 Sora 的区别**: 不需要处理 multipart，因为 Suno 只接受 JSON 请求。

### `DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (*http.Response, error)`
- **作用**: 委托给 `channel.DoTaskApiRequest` 发送请求。

### `DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (string, []byte, *dto.TaskError)`
- **作用**: 解析上游响应，返回上游任务 ID。
- **特殊处理**: 构建公共响应对象，用 `PublicTaskID` 替换上游 ID 返回给客户端。
- **返回值**: 返回上游任务 ID（用于后续轮询），不返回 taskData。

### `FetchTask(baseUrl, key string, body map[string]any, proxy string) (*http.Response, error)`
- **作用**: 向 `/suno/fetch` 发送 POST 请求，用于批量获取任务状态。
- **与视频适配器的区别**: 使用 POST 方法而非 GET，因为 Suno 的 fetch API 接受任务 ID 列表。

### `GetModelList() []string` / `GetChannelName() string`
- **作用**: 返回模型列表和通道名称常量。

### `actionValidate(c *gin.Context, sunoRequest *dto.SunoSubmitReq, action string) error`
- **作用**: 根据动作类型校验请求参数。
- **校验规则**:
  - `generate` 动作：如果 `Mv` 为空，默认设为 `chirp-v3-0`
  - `lyrics` 动作：`Prompt` 不能为空
  - 其他动作：返回 `invalid_action` 错误

## 5. 关键逻辑分析

### 批量轮询机制
Suno 适配器最显著的特点是不使用 `ParseTaskResult`，而是通过 `FetchTask` 发送到 `/suno/fetch` 端点进行批量轮询。这是因为 Suno 的 API 设计支持一次查询多个任务状态。

### Action 路由
Suno 的动作通过 URL 路径参数传递（`/suno/submit/{action}`），而非请求体。`ValidateRequestAndSetAction` 从 `c.Param("action")` 获取并转大写处理。

### 默认模型版本
`actionValidate` 中为音乐生成动作设置默认模型版本 `chirp-v3-0`，确保即使客户端未指定也能正常工作。

### 安全设计
与 Sora 类似，Suno 适配器也使用 `PublicTaskID` 机制，将上游任务 ID 替换为公共 ID 返回给客户端。

## 6. 关联文件

- `relay/channel/task/suno/models.go` — 模型列表和通道名称
- `service/suno_task.go` — `UpdateSunoTasks` 批量轮询实现
- `dto/suno.go` — `SunoSubmitReq`、`SunoDataResponse` 类型
- `constant/suno.go` — `SunoActionMusic`、`SunoActionLyrics` 常量
- `relay/channel/task/taskcommon/helpers.go` — `BaseBilling` 结构体
- `relay/common/relay_info.go` — `RelayInfo` 结构体
