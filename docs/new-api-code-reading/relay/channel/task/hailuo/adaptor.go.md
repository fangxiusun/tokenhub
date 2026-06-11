# adaptor.go 代码阅读文档

## 1. 全局总结

该文件实现了 MiniMax 海螺（Hailuo）视频生成任务的适配器（TaskAdaptor）。它负责将统一请求转换为海螺专有格式、处理上游响应、轮询任务状态、以及将结果转换为 OpenAI 兼容格式。海螺 API 的特点是任务成功后返回的是 file_id，需要额外请求获取实际的视频下载 URL。

## 2. 依赖关系

**标准库：**
- `bytes` — 字节缓冲
- `fmt` — 格式化
- `io` — IO 操作
- `net/http` — HTTP 请求
- `strconv` — 字符串转换
- `strings` — 字符串处理
- `time` — 时间处理

**项目内部依赖：**
- `github.com/QuantumNous/new-api/common` — JSON 序列化
- `github.com/QuantumNous/new-api/model` — 任务状态常量
- `github.com/QuantumNous/new-api/constant` — 任务动作常量
- `github.com/QuantumNous/new-api/dto` — OpenAI 视频 DTO
- `github.com/QuantumNous/new-api/relay/channel` — 通用请求执行
- `github.com/QuantumNous/new-api/relay/channel/task/taskcommon` — 任务基类
- `github.com/QuantumNous/new-api/relay/common` — 中继信息、任务请求
- `github.com/QuantumNous/new-api/service` — HTTP 客户端、错误包装

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
| `ValidateRequestAndSetAction(c, info)` | 验证请求 |
| `BuildRequestURL(info)` | 构建 URL：`{baseURL}/v1/video_generation` |
| `BuildRequestHeader(c, req, info)` | 设置 Content-Type、Accept、Authorization 头 |
| `BuildRequestBody(c, info)` | 转换请求为海螺格式 |
| `DoRequest(c, info, requestBody)` | 委托通用请求执行 |
| `DoResponse(c, resp, info)` | 解析响应，检查状态码，返回 task_id |
| `FetchTask(baseUrl, key, body, proxy)` | 轮询：GET `{baseURL}/v1/query/video_generation?task_id={id}` |
| `GetModelList()` | 返回模型列表 |
| `GetChannelName()` | 返回 "hailuo-video" |
| `ParseTaskResult(respBody)` | 解析任务结果，映射状态 |
| `ConvertToOpenAIVideo(originTask)` | 转换为 OpenAI 格式，处理错误信息 |

### 辅助函数

| 函数签名 | 说明 |
|----------|------|
| `convertToRequestPayload(req, info)` | 转换请求，处理模型配置、时长、分辨率 |
| `parseResolutionFromSize(size, modelConfig)` | 从尺寸字符串解析分辨率 |
| `buildVideoURL(taskID, fileID)` | 通过 file_id 获取实际视频下载 URL |
| `contains(slice, item)` | 字符串切片包含检查 |
| `containsInt(slice, item)` | 整数切片包含检查 |

## 5. 关键逻辑分析

### 视频 URL 获取流程
海螺 API 任务完成后只返回 file_id，不直接返回视频 URL。`buildVideoURL` 方法会额外发送请求到 `/v1/files/retrieve?file_id={id}` 获取实际下载地址。

### 分辨率解析
根据尺寸字符串中的关键词匹配：
- 包含 "1080" → 1080P
- 包含 "768" → 768P
- 包含 "720" → 720P
- 包含 "512" → 512P
- 其他 → 使用模型默认分辨率

### 状态映射
- Preparing/Queueing → InProgress (30%)
- Processing → InProgress (50%)
- Success → Success (100%)
- Fail → Failure (100%)

### 错误处理
`DoResponse` 会检查 `BaseResp.StatusCode` 是否为 `StatusSuccess`（0），非零值视为错误。

## 6. 关联文件

- `relay/channel/task/hailuo/constants.go` — 模型列表、端点常量、状态码、分辨率常量
- `relay/channel/task/hailuo/models.go` — 请求/响应数据结构、模型配置
