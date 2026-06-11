# dto.go 代码阅读文档

## 1. 全局总结
该文件定义了百度（Baidu）渠道特定的数据传输对象，包括请求和响应结构体。这些结构体用于在百度文心一言 API 之间传递数据，支持聊天、嵌入和访问令牌管理。

## 2. 依赖关系
- 标准库：`encoding/json`, `time`
- 内部包：
  - `github.com/QuantumNous/new-api/dto`: 通用数据传输对象

## 3. 类型定义
### 聊天相关
- `BaiduMessage`: 消息结构体，包含角色和内容。
- `BaiduChatRequest`: 聊天请求结构体，包含消息、温度、TopP、惩罚分数等参数。
- `BaiduChatResponse`: 聊天响应结构体，包含 ID、对象、结果、使用量等。
- `BaiduChatStreamResponse`: 流式聊天响应结构体，继承自 `BaiduChatResponse`，添加句子 ID 和结束标志。

### 嵌入相关
- `BaiduEmbeddingRequest`: 嵌入请求结构体，包含输入文本数组。
- `BaiduEmbeddingData`: 嵌入数据结构体，包含对象、嵌入向量和索引。
- `BaiduEmbeddingResponse`: 嵌入响应结构体，包含 ID、对象、数据数组和使用量。

### 错误和访问令牌
- `Error`: 错误结构体，包含错误码和错误消息。
- `BaiduAccessToken`: 访问令牌结构体，包含令牌、错误信息和过期时间。
- `BaiduTokenResponse`: 令牌响应结构体，包含过期时间和访问令牌。

## 4. 函数详解
该文件没有定义任何函数。

## 5. 关键逻辑分析
- **参数映射**：`BaiduChatRequest` 包含百度特有的参数，如 `PenaltyScore`（惩罚分数）、`DisableSearch`（禁用搜索）等。
- **流式支持**：`BaiduChatStreamResponse` 添加了流式响应所需的 `SentenceId` 和 `IsEnd` 字段。
- **访问令牌管理**：`BaiduAccessToken` 包含过期时间管理，支持令牌缓存和自动刷新。
- **错误处理**：统一错误结构体，便于错误传播和处理。

## 6. 关联文件
- `baidu/adaptor.go`: 使用这些 DTO 进行请求构建和响应处理。
- `baidu/relay-baidu.go`: 请求转换和响应处理逻辑。
- `dto/dto.go`: 通用数据传输对象定义。