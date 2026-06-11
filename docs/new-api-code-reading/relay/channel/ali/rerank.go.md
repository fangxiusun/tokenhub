# rerank.go 代码阅读文档

## 1. 全局总结
该文件实现了阿里云（Ali）渠道的重排序（Rerank）功能，包括请求转换和响应处理。重排序用于对文档进行相关性排序，是搜索和信息检索的重要功能。

## 2. 依赖关系
- 标准库：`encoding/json`, `io`, `net/http`
- 内部包：
  - `github.com/QuantumNous/new-api/dto`: 数据传输对象
  - `relaycommon`: 中继通用配置
  - `service`: 业务逻辑服务
  - `types`: 类型定义
- 第三方包：
  - `github.com/gin-gonic/gin`: HTTP 框架

## 3. 类型定义
该文件没有定义新的类型，主要使用 `ali/dto.go` 中定义的 `AliRerankRequest` 和 `AliRerankResponse`。

## 4. 函数详解
1. **`ConvertRerankRequest`**: 将 OpenAI 格式的重排序请求转换为阿里云格式，设置默认的 `ReturnDocuments` 参数。
2. **`RerankHandler`**: 处理阿里云重排序响应，转换为 OpenAI 格式并返回给客户端。

## 5. 关键逻辑分析
- **默认参数处理**：如果 `ReturnDocuments` 未指定，默认设置为 `true`。
- **响应格式转换**：将阿里云的 `AliRerankResponse` 转换为 OpenAI 的 `dto.RerankResponse` 格式。
- **错误处理**：检查响应中的错误码，返回适当的错误信息。
- **使用量计算**：从阿里云响应中提取 token 使用量，转换为 OpenAI 格式。

## 6. 关联文件
- `ali/adaptor.go`: 调用 `ConvertRerankRequest` 和 `RerankHandler` 处理重排序请求。
- `ali/dto.go`: 定义 `AliRerankRequest` 和 `AliRerankResponse` 结构体。
- `dto/dto.go`: 通用重排序请求和响应结构体。