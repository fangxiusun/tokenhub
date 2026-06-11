# adaptor.go 代码阅读文档

## 1. 全局总结

该文件实现了 Google Gemini Veo 视频生成任务的适配器（TaskAdaptor）。它负责将统一请求转换为 Gemini Veo 的 `predictLongRunning` 格式、处理上游响应、轮询任务状态、以及将结果转换为 OpenAI 兼容格式。Gemini Veo 使用异步长时运行操作（LongRunning Operation）模式，任务 ID 编码了完整的操作名称路径。

## 2. 依赖关系

**标准库：**
- `bytes` — 字节缓冲
- `fmt` — 格式化
- `io` — IO 操作
- `net/http` — HTTP 请求
- `regexp` — 正则表达式（从操作名中提取模型名）
- `strings` — 字符串处理
- `time` — 时间处理

**项目内部依赖：**
- `github.com/QuantumNous/new-api/common` — JSON 序列化
- `github.com/QuantumNous/new-api/constant` — 任务动作常量（TaskActionTextGenerate、TaskActionGenerate）
- `github.com/QuantumNous/new-api/dto` — OpenAI 视频 DTO
- `github.com/QuantumNous/new-api/model` — 任务状态常量
- `github.com/QuantumNous/new-api/relay/channel` — 通用请求执行
- `github.com/QuantumNous/new-api/relay/channel/task/taskcommon` — 任务基类、元数据工具、LocalTaskID 编解码
- `github.com/QuantumNous/new-api/relay/common` — 中继信息、任务请求
- `github.com/QuantumNous/new-api/service` — HTTP 客户端、错误包装
- `github.com/QuantumNous/new-api/setting/model_setting` — Gemini 版本设置

**第三方依赖：**
- `github.com/gin-gonic/gin` — Web 框架
- `github.com/pkg/errors` — 错误包装

## 3. 类型定义

### 适配器结构体

| 类型名 | 说明 |
|--------|------|
| `TaskAdaptor` | 核心适配器，嵌入 `taskcommon.BaseBilling` |

## 4. 函数详解

### 适配器接口方法

| 函数签名 | 说明 |
|----------|------|
| `Init(info)` | 初始化适配器 |
| `ValidateRequestAndSetAction(c, info)` | 验证请求，设置默认动作为 TaskActionTextGenerate |
| `BuildRequestURL(info)` | 构建 URL：`{baseURL}/{version}/models/{modelName}:predictLongRunning` |
| `BuildRequestHeader(c, req, info)` | 设置 Content-Type、Accept、x-goog-api-key 头 |
| `BuildRequestBody(c, info)` | 构建 Veo 请求体，处理图片输入和参数 |
| `DoRequest(c, info, requestBody)` | 委托通用请求执行 |
| `DoResponse(c, resp, info)` | 解析操作名称，编码为本地任务 ID |
| `FetchTask(baseUrl, key, body, proxy)` | 轮询：GET `{baseUrl}/{version}/{operationName}` |
| `GetModelList()` | 返回 Veo 模型列表 |
| `GetChannelName()` | 返回 "gemini" |
| `EstimateBilling(c, info)` | 返回 seconds 和 resolution 比率 |
| `ParseTaskResult(respBody)` | 解析 operationResponse，映射状态 |
| `ConvertToOpenAIVideo(task)` | 从操作名中提取模型名并转换为 OpenAI 格式 |

### 辅助函数

| 函数签名 | 说明 |
|----------|------|
| `extractModelFromOperationName(name string) string` | 从操作名（如 `models/veo-3.0-generate-001/operations/...`）中提取模型名 |

## 5. 关键逻辑分析

### 任务 ID 编码
Gemini 使用完整操作名作为任务 ID（如 `models/veo-3.0-generate-001/operations/abc123`），通过 `taskcommon.EncodeLocalTaskID` 编码为本地 ID 存储，查询时通过 `DecodeLocalTaskID` 还原。

### 图片输入处理
1. 优先从 multipart 表单中提取 `input_reference` 文件（`ExtractMultipartImage`）
2. 如果没有文件，尝试解析 `req.Images[0]` 中的 base64 或 data URI（`ParseImageInput`）
3. 检测到图片输入时，将 action 从 TextGenerate 改为 Generate

### 版本设置
通过 `model_setting.GetGeminiVersionSetting` 获取 API 版本（如 "v1beta"），用于构建 URL 和轮询。

### 状态映射
- op.Error 非空 → Failure
- op.Done == false → InProgress
- op.Done == true → Success，提取 generatedVideos[0].video.uri

## 6. 关联文件

- `relay/channel/task/gemini/billing.go` — 计费参数解析（时长、分辨率、倍率）
- `relay/channel/task/gemini/dto.go` — Veo 请求/响应 DTO 定义
- `relay/channel/task/gemini/image.go` — 图片输入解析工具
- `relay/channel/task/taskcommon/` — LocalTaskID 编解码、元数据反序列化
- `setting/model_setting/` — Gemini 版本设置
