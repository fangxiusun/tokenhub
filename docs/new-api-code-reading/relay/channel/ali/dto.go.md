# dto.go 代码阅读文档

## 1. 全局总结
该文件定义了阿里云（Ali）渠道特定的数据传输对象（DTO），包括请求和响应结构体。这些结构体用于在阿里云 API 之间传递数据，支持聊天、嵌入、图像生成、重排序等功能。

## 2. 依赖关系
- 标准库：`strings`
- 内部包：
  - `github.com/QuantumNous/new-api/dto`: 通用数据传输对象
  - `github.com/QuantumNous/new-api/logger`: 日志记录
  - `service`: 业务逻辑服务
- 第三方包：
  - `github.com/gin-gonic/gin`: HTTP 框架

## 3. 类型定义
### 聊天相关
- `AliMessage`: 消息结构体，包含 `Content`（任意类型）和 `Role`（字符串）。
- `AliMediaContent`: 媒体内容结构体，包含 `Image`（Base64 图像）和 `Text`（文本）。
- `AliInput`: 输入结构体，包含 `Prompt` 和 `Messages`。
- `AliParameters`: 参数结构体，包含 `TopP`、`TopK`、`Seed`、`EnableSearch` 等。
- `AliChatRequest`: 聊天请求结构体，包含模型、输入和参数。

### 嵌入相关
- `AliEmbeddingRequest`: 嵌入请求结构体。
- `AliEmbedding`: 嵌入结果结构体。
- `AliEmbeddingResponse`: 嵌入响应结构体。

### 错误和使用量
- `AliError`: 错误结构体，包含错误码、消息和请求 ID。
- `AliUsage`: 使用量结构体，包含输入、输出和总 token 数。

### 图像相关
- `TaskResult`: 任务结果结构体，包含 Base64 图像、URL、错误码和消息。
- `AliOutput`: 输出结构体，包含任务 ID、状态、文本、结果和选择。
- `AliImageRequest`: 图像请求结构体。
- `AliImageParameters`: 图像参数结构体，包含尺寸、数量、步数、水印等。
- `AliImageInput`: 图像输入结构体，包含提示词和消息。
- `WanImageInput`: Wan 模型图像输入结构体。
- `WanImageParameters`: Wan 模型图像参数结构体。

### 重排序相关
- `AliRerankParameters`: 重排序参数结构体。
- `AliRerankInput`: 重排序输入结构体。
- `AliRerankRequest`: 重排序请求结构体。
- `AliRerankResponse`: 重排序响应结构体。

## 4. 函数详解
1. **`(o *AliOutput) ChoicesToOpenAIImageDate`**: 将选择转换为 OpenAI 图像数据格式。
2. **`(o *AliOutput) ResultToOpenAIImageDate`**: 将结果转换为 OpenAI 图像数据格式。
3. **`(p *AliImageParameters) PromptExtendValue()`**: 获取提示词扩展值，默认返回 false。

## 5. 关键逻辑分析
- **灵活的输入类型**：`AliInput` 支持 `Prompt` 和 `Messages` 两种输入方式，适应不同的 API 要求。
- **图像格式转换**：支持 Base64 和 URL 两种图像格式，自动根据响应格式转换。
- **Wan 模型支持**：专门为 Wan 系列图像模型定义了输入和参数结构体。
- **错误处理**：统一错误结构体，便于错误传播和处理。

## 6. 关联文件
- `ali/adaptor.go`: 使用这些 DTO 进行请求构建和响应处理。
- `ali/image.go`: 图像请求和响应处理逻辑。
- `ali/rerank.go`: 重排序请求处理逻辑。
- `dto/dto.go`: 通用数据传输对象定义。