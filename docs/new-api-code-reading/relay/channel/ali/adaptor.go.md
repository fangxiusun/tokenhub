# adaptor.go 代码阅读文档

## 1. 全局总结
该文件是阿里云（Ali）渠道的适配器实现，负责将 OpenAI 格式的请求转换为阿里云 API 格式，并处理响应。适配器支持多种中继模式（聊天、嵌入、图像生成、重排序等），并实现了 `channel.Adaptor` 接口。

## 2. 依赖关系
- 标准库：`errors`, `fmt`, `io`, `net/http`, `strings`
- 内部包：
  - `github.com/QuantumNous/new-api/common`: 通用工具函数
  - `github.com/QuantumNous/new-api/dto`: 数据传输对象
  - `github.com/QuantumNous/new-api/relay/channel`: 渠道通用工具
  - `github.com/QuantumNous/new-api/relay/channel/claude`: Claude 渠道适配器
  - `github.com/QuantumNous/new-api/relay/channel/openai`: OpenAI 渠道适配器
  - `relaycommon`: 中继通用配置
  - `relay/constant`: 中继常量
  - `service`: 业务逻辑服务
  - `model_setting`: 模型设置
  - `types`: 类型定义
- 第三方包：
  - `github.com/gin-gonic/gin`: HTTP 框架
  - `github.com/samber/lo`: 泛型工具库

## 3. 类型定义
### 结构体
- `Adaptor`: 阿里云适配器结构体，包含一个布尔字段 `IsSyncImageModel`，用于标识是否为同步图像模型。

### 常量
- `aliAnthropicMessagesModelsEnv`: 环境变量名，用于配置支持 Anthropic Messages API 的模型。
- `defaultAliAnthropicMessagesModels`: 默认模型列表，包含 `"qwen,deepseek-v4,kimi,glm,minimax-m"`。

### 变量
- `syncModels`: 同步图像模型列表，包含 `"z-image"`, `"qwen-image"`, `"wan2.6"`。

## 4. 函数详解
### 核心函数
1. **`supportsAliAnthropicMessages(modelName string) bool`**
   - 检查模型是否支持 Anthropic Messages API。
   - 通过环境变量配置模型模式列表，匹配模型名称。

2. **`aliAnthropicMessagesModelPatterns() []string`**
   - 从环境变量获取支持的模型模式列表。

3. **`isSyncImageModel(modelName string) bool`**
   - 判断模型是否为同步图像模型，委托给 `model_setting.IsSyncImageModel`。

4. **`ConvertClaudeRequest`**
   - 转换 Claude 请求：如果模型支持 Anthropic Messages API，直接返回原请求；否则转换为 OpenAI 格式。

5. **`GetRequestURL`**
   - 根据中继格式和模式构建请求 URL，支持多种端点（聊天、嵌入、图像生成、重排序等）。

6. **`SetupRequestHeader`**
   - 设置请求头，包括授权信息、SSE 启用标志和插件信息。

7. **`ConvertOpenAIRequest`**
   - 将 OpenAI 请求转换为阿里云格式，调用 `requestOpenAI2Ali`。

8. **`ConvertImageRequest`**
   - 转换图像请求，支持同步和异步图像生成，处理表单和 JSON 输入。

9. **`DoResponse`**
   - 处理响应：根据中继格式和模式调用不同的处理器（Claude、OpenAI、图像、重排序等）。

## 5. 关键逻辑分析
- **双模式支持**：支持 OpenAI 格式和 Anthropic Messages API 格式，根据模型动态选择。
- **图像处理**：区分同步和异步图像模型，使用不同的 API 端点和处理逻辑。
- **错误处理**：在关键步骤返回明确的错误信息，便于调试。
- **性能优化**：使用 `lo` 库进行集合操作，提高代码可读性。

## 6. 关联文件
- `ali/constants.go`: 定义模型列表和渠道名称。
- `ali/dto.go`: 阿里云特定的数据传输对象。
- `ali/image.go`: 图像请求和响应处理。
- `ali/text.go`: 文本请求转换。
- `ali/rerank.go`: 重排序请求处理。
- `relay/channel/openai/adaptor.go`: OpenAI 渠道适配器，用于响应格式转换。
- `relay/channel/claude/adaptor.go`: Claude 渠道适配器，用于 Anthropic Messages API 支持。