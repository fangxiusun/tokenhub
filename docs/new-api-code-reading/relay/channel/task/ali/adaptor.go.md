# adaptor.go 代码阅读文档

## 1. 全局总结

该文件实现了阿里通义万相（Wanx）视频生成任务的适配器（TaskAdaptor）。它是 relay/channel/task 体系中针对阿里云视频生成 API 的具体实现，负责将统一的视频生成请求转换为阿里专有格式、处理上游响应、轮询任务状态、以及将结果转换为 OpenAI 兼容格式。支持文生视频（t2v）、图生视频（i2v）、首尾帧生视频（kf2v）等多种生成模式。

## 2. 依赖关系

**标准库：**
- `bytes` — 构建请求体的字节缓冲
- `fmt` — 格式化字符串
- `io` — IO 操作（Reader、ReadAll）
- `net/http` — HTTP 请求构建与执行
- `strconv` — 字符串到整数转换
- `strings` — 字符串处理（前缀匹配、大小写转换等）

**项目内部依赖：**
- `github.com/QuantumNous/new-api/common` — JSON 序列化/反序列化封装、时间戳工具
- `github.com/QuantumNous/new-api/dto` — OpenAI 视频响应 DTO、IntValue/BoolValue 类型
- `github.com/QuantumNous/new-api/logger` — 日志记录
- `github.com/QuantumNous/new-api/model` — 任务状态常量定义
- `github.com/QuantumNous/new-api/relay/channel` — 通用任务 API 请求执行
- `github.com/QuantumNous/new-api/relay/channel/task/taskcommon` — 任务基类（BaseBilling）
- `github.com/QuantumNous/new-api/relay/common` — 中继信息、任务请求/响应通用类型、验证方法
- `github.com/QuantumNous/new-api/service` — HTTP 客户端代理、错误包装

**第三方依赖：**
- `github.com/gin-gonic/gin` — Web 框架上下文
- `github.com/pkg/errors` — 错误包装
- `github.com/samber/lo` — 泛型切片工具（Contains）

## 3. 类型定义

### 请求/响应结构体

| 类型名 | 说明 |
|--------|------|
| `AliVideoRequest` | 阿里通义万相视频生成请求，包含 model、input、parameters 三个顶层字段 |
| `AliVideoInput` | 视频输入参数，支持 prompt、img_url、首尾帧 URL、音频 URL、反向提示词、特效模板 |
| `AliVideoParameters` | 视频参数，包含分辨率（resolution）、尺寸（size）、时长、prompt 智能改写、水印、音频、随机种子 |
| `AliVideoResponse` | 阿里 API 响应，包含 output、request_id、code、message、usage |
| `AliVideoOutput` | 输出信息，包含 task_id、task_status、时间戳、视频 URL、错误码 |
| `AliUsage` | 使用统计，包含 duration、video_count、SR |
| `AliMetadata` | 额外元数据结构，用于从 metadata 中提取阿里特有参数 |

### 适配器结构体

| 类型名 | 说明 |
|--------|------|
| `TaskAdaptor` | 核心适配器，嵌入 `taskcommon.BaseBilling`，包含 ChannelType、apiKey、baseURL 字段 |

### 变量

| 变量名 | 说明 |
|--------|------|
| `size480p` | 480P 分辨率对应的尺寸列表 |
| `size720p` | 720P 分辨率对应的尺寸列表 |
| `size1080p` | 1080P 分辨率对应的尺寸列表 |

## 4. 函数详解

### 适配器接口方法

| 函数签名 | 说明 |
|----------|------|
| `Init(info *relaycommon.RelayInfo)` | 初始化适配器，设置频道类型、基础 URL、API 密钥 |
| `ValidateRequestAndSetAction(c, info) *dto.TaskError` | 验证请求并设置动作，委托给 `ValidateMultipartDirect` |
| `BuildRequestURL(info) (string, error)` | 构建上游 API URL：`{baseURL}/api/v1/services/aigc/video-generation/video-synthesis` |
| `BuildRequestHeader(c, req, info) error` | 设置请求头：Authorization（Bearer）、Content-Type、X-DashScope-Async（enable） |
| `BuildRequestBody(c, info) (io.Reader, error)` | 从 context 获取任务请求，转换为阿里格式并序列化 |
| `DoRequest(c, info, requestBody) (*http.Response, error)` | 委托给 `channel.DoTaskApiRequest` 执行 HTTP 请求 |
| `DoResponse(c, resp, info) (taskID, taskData, taskErr)` | 解析上游响应，检查错误，转换为 OpenAI 格式并返回 |
| `FetchTask(baseUrl, key, body, proxy) (*http.Response, error)` | 轮询任务状态，发送 GET 请求到 `/api/v1/tasks/{task_id}` |
| `GetModelList() []string` | 返回支持的模型列表（来自 constants.go） |
| `GetChannelName() string` | 返回频道名称 "ali" |
| `ParseTaskResult(respBody) (*relaycommon.TaskInfo, error)` | 解析任务结果，将阿里状态映射为内部状态 |
| `ConvertToOpenAIVideo(task) ([]byte, error)` | 将存储的任务数据转换为 OpenAI 格式的视频响应 |
| `EstimateBilling(c, info) map[string]float64` | 根据请求参数计算计费比率（时长、分辨率） |

### 辅助函数

| 函数签名 | 说明 |
|----------|------|
| `sizeToResolution(size string) (string, error)` | 将尺寸字符串（如 "1920*1080"）转换为分辨率标签（如 "1080P"） |
| `ProcessAliOtherRatios(aliReq) (map[string]float64, error)` | 根据模型和分辨率计算计费倍率（其他比率） |
| `convertToAliRequest(info, req) (*AliVideoRequest, error)` | 将统一任务请求转换为阿里专有请求格式 |
| `convertAliStatus(aliStatus string) string` | 将阿里任务状态映射为 OpenAI 视频状态 |

## 5. 关键逻辑分析

### 请求构建流程
1. 从 `gin.Context` 中获取已解析的 `TaskSubmitReq`
2. 通过 `convertToAliRequest` 转换为阿里格式，处理模型映射、分辨率映射、时长默认值
3. 从 metadata 中提取额外参数（如 negative_prompt、seed 等），但禁止通过 metadata 修改 model
4. 设置 `X-DashScope-Async: enable` 头部以启用异步任务模式

### 分辨率映射策略
- 文生视频（t2v）：支持直接指定尺寸（如 "1920*1080"）或分辨率标签（如 "1080p"）
- 图生视频（i2v）：根据模型版本设置默认分辨率（wan2.5/2.6 默认 1080P，wan2.2-i2v-flash 默认 720P）
- 尺寸到分辨率转换通过预定义的尺寸列表实现

### 计费倍率计算
- `ProcessAliOtherRatios` 维护了每个模型在不同分辨率下的倍率映射
- `EstimateBilling` 在验证请求后、价格计算前调用，返回 `seconds` 和 `resolution-*` 类型的比率

### 状态映射
- PENDING → Queued, RUNNING → InProgress, SUCCEEDED → Success
- FAILED/CANCELED/UNKNOWN → Failure

## 6. 关联文件

- `relay/channel/task/ali/constants.go` — 模型列表和频道名称常量
- `relay/channel/task/taskcommon/` — 任务基类 BaseBilling
- `relay/common/relay.go` — RelayInfo、TaskSubmitReq 等通用类型
- `relay/channel/task_api.go` — DoTaskApiRequest 通用请求执行函数
- `dto/video.go` — OpenAIVideo 响应 DTO
- `service/task.go` — TaskErrorWrapper 错误包装函数
