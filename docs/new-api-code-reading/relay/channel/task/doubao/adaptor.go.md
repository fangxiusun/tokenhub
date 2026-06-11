# adaptor.go 代码阅读文档

## 1. 全局总结

该文件实现了豆包（Doubao/Seedance）视频生成任务的适配器（TaskAdaptor）。它是 relay/channel/task 体系中针对字节跳动豆包视频生成 API 的具体实现，负责将统一请求转换为豆包专有格式、处理上游响应、轮询任务状态、以及将结果转换为 OpenAI 兼容格式。豆包 API 采用 OpenAI 兼容的 content 数组格式，支持文本、图片、视频等多种内容类型。

## 2. 依赖关系

**标准库：**
- `bytes` — 构建请求体的字节缓冲
- `fmt` — 格式化字符串
- `io` — IO 操作
- `net/http` — HTTP 请求
- `strconv` — 字符串转换
- `time` — 时间处理

**项目内部依赖：**
- `github.com/QuantumNous/new-api/common` — JSON 序列化封装
- `github.com/QuantumNous/new-api/constant` — 任务动作常量（TaskActionGenerate）
- `github.com/QuantumNous/new-api/dto` — OpenAI 视频 DTO、BoolValue/IntValue 类型
- `github.com/QuantumNous/new-api/model` — 任务状态常量
- `github.com/QuantumNous/new-api/relay/channel` — 通用任务请求执行
- `github.com/QuantumNous/new-api/relay/channel/task/taskcommon` — 任务基类、元数据反序列化
- `github.com/QuantumNous/new-api/relay/common` — 中继信息、任务请求类型、验证方法
- `github.com/QuantumNous/new-api/service` — HTTP 客户端、错误包装

**第三方依赖：**
- `github.com/gin-gonic/gin` — Web 框架
- `github.com/pkg/errors` — 错误包装
- `github.com/samber/lo` — 泛型工具（ToPtr、Reject）

## 3. 类型定义

### 请求/响应结构体

| 类型名 | 说明 |
|--------|------|
| `ContentItem` | 内容项，支持 text、image_url、video_url、audio_url 类型 |
| `MediaURL` | 媒体 URL 包装器 |
| `requestPayload` | 豆包请求体，包含 model、content 数组、回调 URL、视频参数等 |
| `responsePayload` | 提交响应，仅包含 task ID |
| `responseTask` | 任务查询响应，包含完整的任务状态、视频 URL、用量、错误信息 |

### 适配器结构体

| 类型名 | 说明 |
|--------|------|
| `TaskAdaptor` | 核心适配器，嵌入 `taskcommon.BaseBilling` |

## 4. 函数详解

### 适配器接口方法

| 函数签名 | 说明 |
|----------|------|
| `Init(info)` | 初始化，设置频道类型、基础 URL、API 密钥 |
| `ValidateRequestAndSetAction(c, info)` | 验证请求，仅接受 POST /v1/video/generations |
| `BuildRequestURL(info)` | 构建 URL：`{baseURL}/api/v3/contents/generations/tasks` |
| `BuildRequestHeader(c, req, info)` | 设置 Content-Type、Accept、Authorization 头 |
| `BuildRequestBody(c, info)` | 转换请求为豆包格式并序列化 |
| `DoRequest(c, info, requestBody)` | 委托通用请求执行 |
| `DoResponse(c, resp, info)` | 解析响应，检查 task_id，返回 OpenAI 格式 |
| `FetchTask(baseUrl, key, body, proxy)` | 轮询任务状态：GET `/api/v3/contents/generations/tasks/{task_id}` |
| `GetModelList()` | 返回模型列表 |
| `GetChannelName()` | 返回 "doubao-video" |
| `ParseTaskResult(respBody)` | 解析任务结果，映射状态和 token 用量 |
| `ConvertToOpenAIVideo(originTask)` | 转换存储数据为 OpenAI 格式 |
| `EstimateBilling(c, info)` | 检测视频输入并返回折扣比率 |

### 辅助函数

| 函数签名 | 说明 |
|----------|------|
| `hasVideoInMetadata(metadata)` | 检查 metadata 中是否包含视频输入 |
| `convertToRequestPayload(req)` | 将统一请求转换为豆包格式，处理图片和文本内容 |

## 5. 关键逻辑分析

### 请求构建流程
1. 从 context 获取 `TaskSubmitReq`
2. 将图片 URL 转换为 content 数组中的 image_url 项
3. 通过 `taskcommon.UnmarshalMetadata` 将 metadata 展开到请求体
4. 处理 duration 字段（从 Seconds 字符串转换为 IntValue）
5. 移除 content 中的 text 项，将 prompt 作为最后一个 text 项追加

### 视频输入计费
- `hasVideoInMetadata` 直接检查 metadata 的 content 数组是否包含 video_url 条目
- `EstimateBilling` 在检测到视频输入时返回折扣比率（从 constants.go 获取）

### 状态映射
- pending/queued → Queued (10%)
- processing/running → InProgress (50%)
- succeeded → Success (100%)，附带 video URL 和 token 用量
- failed → Failure (100%)

## 6. 关联文件

- `relay/channel/task/doubao/constants.go` — 模型列表、频道名称、视频输入折扣比率
- `relay/channel/task/taskcommon/` — 元数据反序列化工具
- `relay/common/relay.go` — TaskSubmitReq 等通用类型
- `dto/video.go` — OpenAIVideo 响应 DTO
