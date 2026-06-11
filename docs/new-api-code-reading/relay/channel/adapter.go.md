# adapter.go 代码阅读文档

## 1. 全局总结

本文件定义了 Relay 渠道适配器的核心接口：`Adaptor`（同步请求）和 `TaskAdaptor`（异步任务），以及 `OpenAIVideoConverter`（视频格式转换）接口。所有渠道适配器都必须实现这些接口。

## 2. 依赖关系

- `dto`: 请求/响应 DTO
- `model`: 任务模型
- `relay/common`: RelayInfo
- `types`: 错误类型

## 3. 类型定义

### `Adaptor` 接口
同步请求适配器接口，所有渠道必须实现：
- `Init(info)`: 初始化适配器
- `GetRequestURL(info) (string, error)`: 构建请求 URL
- `SetupRequestHeader(c, req, info) error`: 设置请求头
- `ConvertOpenAIRequest(c, info, request) (any, error)`: 转换 OpenAI 格式请求
- `ConvertClaudeRequest(c, info, request) (any, error)`: 转换 Claude 格式请求
- `ConvertGeminiRequest(c, info, request) (any, error)`: 转换 Gemini 格式请求
- `ConvertRerankRequest(c, relayMode, request) (any, error)`: 转换 Rerank 请求
- `ConvertEmbeddingRequest(c, info, request) (any, error)`: 转换 Embedding 请求
- `ConvertAudioRequest(c, info, request) (io.Reader, error)`: 转换 Audio 请求
- `ConvertImageRequest(c, info, request) (any, error)`: 转换 Image 请求
- `ConvertOpenAIResponsesRequest(c, info, request) (any, error)`: 转换 Responses 请求
- `DoRequest(c, info, requestBody) (any, error)`: 发送请求
- `DoResponse(c, resp, info) (usage, err)`: 处理响应
- `GetModelList() []string`: 获取支持的模型列表
- `GetChannelName() string`: 获取渠道名称

### `TaskAdaptor` 接口
异步任务适配器接口：
- `Init(info)`: 初始化
- `ValidateRequestAndSetAction(c, info) *dto.TaskError`: 验证请求并设置操作
- `EstimateBilling(c, info) map[string]float64`: 估算计费比率
- `AdjustBillingOnSubmit(info, taskData) map[string]float64`: 提交后调整计费
- `AdjustBillingOnComplete(task, taskResult) int`: 完成后调整计费
- `BuildRequestURL(info) (string, error)`: 构建请求 URL
- `BuildRequestHeader(c, req, info) error`: 设置请求头
- `BuildRequestBody(c, info) (io.Reader, error)`: 构建请求体
- `DoRequest(c, info, requestBody) (*http.Response, error)`: 发送请求
- `DoResponse(c, resp, info) (taskID, taskData, err)`: 处理响应
- `FetchTask(baseUrl, key, body, proxy) (*http.Response, error)`: 轮询任务状态
- `ParseTaskResult(respBody) (*TaskInfo, error)`: 解析任务结果

### `OpenAIVideoConverter` 接口
视频格式转换接口：
- `ConvertToOpenAIVideo(originTask) ([]byte, error)`: 将任务结果转换为 OpenAI Video 格式

## 4. 关键逻辑分析

1. **多格式支持**: Adaptor 接口支持 OpenAI、Claude、Gemini 三种原生请求格式的转换
2. **任务生命周期**: TaskAdaptor 定义了完整的任务生命周期：验证 → 估算计费 → 提交 → 调整计费 → 轮询 → 完成计费
3. **计费三阶段**: EstimateBilling（预估）→ AdjustBillingOnSubmit（提交后调整）→ AdjustBillingOnComplete（完成后调整）

## 5. 关联文件

- `relay/channel/openai/`: OpenAI 适配器实现
- `relay/channel/claude/`: Claude 适配器实现
- `relay/channel/gemini/`: Gemini 适配器实现
