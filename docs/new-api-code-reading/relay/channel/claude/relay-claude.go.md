# relay-claude.go 代码阅读文档

## 1. 全局总结
该文件实现了 Claude 渠道的请求转换和响应处理，包括将 OpenAI 格式请求转换为 Claude 格式、处理流式/非流式响应、工具调用转换以及使用量计算。是 Claude API 集成的核心实现，支持复杂的工具调用和思考模式。

## 2. 依赖关系
- 标准库：`encoding/json`, `fmt`, `io`, `net/http`, `strings`
- 内部包：
  - `github.com/QuantumNous/new-api/common`: 通用工具函数
  - `github.com/QuantumNous/new-api/constant`: 全局常量
  - `github.com/QuantumNous/new-api/dto`: 数据传输对象
  - `github.com/QuantumNous/new-api/logger`: 日志记录
  - `relay/channel/openrouter`: OpenRouter 渠道适配器
  - `relaycommon`: 中继通用配置
  - `relay/helper`: 中继辅助函数
  - `relay/reasonmap`: 原因映射
  - `service`: 业务逻辑服务
  - `model_setting`: 模型设置
  - `setting/reasoning`: 推理设置
  - `types`: 类型定义
- 第三方包：
  - `github.com/gin-gonic/gin`: HTTP 框架
  - `github.com/tidwall/gjson`: JSON 查询
  - `github.com/tidwall/sjson`: JSON 设置

## 3. 类型定义
### 结构体
- `ClaudeResponseInfo`: Claude 响应信息结构体，包含响应 ID、创建时间、模型、响应文本、使用量和完成状态。

## 4. 函数详解
### 请求转换函数
1. **`RequestOpenAI2ClaudeMessage`**: 将 OpenAI 请求转换为 Claude 格式，支持工具调用、网络搜索、思考模式等。
2. **`mapToolChoice`**: 映射 OpenAI 工具选择格式为 Claude 格式。

### 响应处理函数
3. **`StreamResponseClaude2OpenAI`**: 将 Claude 流式响应转换为 OpenAI 格式。
4. **`ResponseClaude2OpenAI`**: 将 Claude 非流式响应转换为 OpenAI 格式。
5. **`HandleStreamResponseData`**: 处理流式响应数据。
6. **`HandleStreamFinalResponse`**: 处理流式响应最终数据。
7. **`HandleClaudeResponseData`**: 处理 Claude 响应数据。
8. **`ClaudeStreamHandler`**: Claude 流式响应处理器。
9. **`ClaudeHandler`**: Claude 非流式响应处理器。

### 使用量处理函数
10. **`cacheCreationTokensForOpenAIUsage`**: 计算缓存创建 token 数。
11. **`buildOpenAIStyleUsageFromClaudeUsage`**: 将 Claude 使用量转换为 OpenAI 格式。
12. **`buildMessageDeltaPatchUsage`**: 构建消息增量补丁使用量。
13. **`patchClaudeMessageDeltaUsageData`**: 修补消息增量使用量数据。

### 辅助函数
14. **`stopReasonClaude2OpenAI`**: 将 Claude 停止原因转换为 OpenAI 完成原因。
15. **`maybeMarkClaudeRefusal`**: 标记 Claude 拒绝响应。
16. **`FormatClaudeResponseInfo`**: 格式化 Claude 响应信息。

## 5. 关键逻辑分析
- **工具调用转换**：将 OpenAI 工具调用格式转换为 Claude 格式，包括输入模式映射。
- **网络搜索支持**：通过 `web_search_20250305` 工具类型支持网络搜索功能。
- **思考模式**：支持 Claude 的思考模式，包括 `enabled` 和 `adaptive` 类型。
- **努力级别**：支持通过模型后缀（如 `-thinking`）和 `reasoning_effort` 参数控制思考努力程度。
- **消息格式化**：处理系统消息、工具结果、图像/文档内容等复杂消息格式。
- **使用量计算**：支持 Claude 特有的缓存 token 计算和使用量格式转换。
- **流式处理**：支持流式响应的增量处理和最终使用量汇总。

## 6. 关联文件
- `claude/adaptor.go`: 调用这些函数执行请求和处理响应。
- `claude/constants.go`: 定义模型列表和渠道名称。
- `dto/dto.go`: 通用数据传输对象定义。
- `relay/channel/openai/adaptor.go`: OpenAI 适配器，用于响应格式转换。
- `relay/helper/stream.go`: 流式响应处理辅助函数。
- `model_setting/claude.go`: Claude 模型设置。
- `setting/reasoning/reasoning.go`: 推理设置。