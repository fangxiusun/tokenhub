# relay-perplexity.go 代码阅读文档

## 1. 全局总结
本文件实现了将 OpenAI 格式的请求转换为 Perplexity 特定格式的逻辑。主要处理消息格式的精简和 MaxTokens 字段的统一。

## 2. 依赖关系
| 包名 | 用途 |
|------|------|
| `github.com/QuantumNous/new-api/dto` | 数据传输对象（请求/消息结构） |

## 3. 类型定义
本文件无额外类型定义。

## 4. 函数详解

### `requestOpenAI2Perplexity`
```go
func requestOpenAI2Perplexity(request dto.GeneralOpenAIRequest) *dto.GeneralOpenAIRequest
```
- **参数**：标准 OpenAI 通用请求
- **返回**：转换后的 Perplexity 格式请求（指针）
- **逻辑**：
  1. 创建新的消息列表，仅保留 `Role` 和 `Content` 字段（丢弃其他字段如 `Name`、`ToolCalls` 等）
  2. 构造新请求，保留以下字段：
     - `Model`、`Stream`、`Temperature`、`TopP`
     - `FrequencyPenalty`、`PresencePenalty`
     - Perplexity 特有字段：`SearchDomainFilter`、`SearchRecencyFilter`、`ReturnImages`、`ReturnRelatedQuestions`、`SearchMode`
  3. 处理 `MaxTokens`：如果原请求设置了 `MaxTokens` 或 `MaxCompletionTokens`，使用 `GetMaxTokens()` 统一获取

## 5. 关键逻辑分析
- **消息精简**：Perplexity 只需要 `Role` 和 `Content`，因此转换时丢弃了 OpenAI 消息中的其他字段（如 `Name`、`ToolCallID`、`ToolCalls` 等）。
- **MaxTokens 统一**：通过 `GetMaxTokens()` 方法统一处理 `MaxTokens` 和 `MaxCompletionTokens`，适配不同客户端的参数命名习惯。
- **搜索参数透传**：Perplexity 特有的搜索相关参数（`SearchDomainFilter`、`SearchRecencyFilter` 等）从原始请求中透传。
- **指针返回**：返回新的指针对象而非修改原请求，避免副作用。

## 6. 关联文件
- `relay/channel/perplexity/adaptor.go` — 调用本函数进行请求转换
- `relay/channel/perplexity/dto.go` — 本文件未定义额外 DTO
- `dto/request.go` — `GeneralOpenAIRequest` 结构体定义
