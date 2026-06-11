# relay-baidu.go 代码阅读文档

## 1. 全局总结
该文件实现了百度（Baidu）渠道的请求转换和响应处理，包括将 OpenAI 格式请求转换为百度格式、处理流式/非流式响应、嵌入响应处理以及访问令牌管理。是百度文心一言 API 集成的核心实现。

## 2. 依赖关系
- 标准库：`encoding/json`, `errors`, `fmt`, `io`, `net/http`, `strings`, `sync`, `time`
- 内部包：
  - `github.com/QuantumNous/new-api/common`: 通用工具函数
  - `github.com/QuantumNous/new-api/constant`: 全局常量
  - `github.com/QuantumNous/new-api/dto`: 数据传输对象
  - `relaycommon`: 中继通用配置
  - `relay/helper`: 中继辅助函数
  - `service`: 业务逻辑服务
  - `types`: 类型定义
- 第三方包：
  - `github.com/gin-gonic/gin`: HTTP 框架
  - `github.com/samber/lo`: 泛型工具库

## 3. 类型定义
该文件没有定义新的类型，主要使用 `baidu/dto.go` 中定义的类型。

## 4. 函数详解
### 请求转换函数
1. **`requestOpenAI2Baidu`**: 将 OpenAI 请求转换为百度格式，处理系统消息、参数映射等。
2. **`embeddingRequestOpenAI2Baidu`**: 将 OpenAI 嵌入请求转换为百度格式。

### 响应处理函数
3. **`responseBaidu2OpenAI`**: 将百度聊天响应转换为 OpenAI 格式。
4. **`streamResponseBaidu2OpenAI`**: 将百度流式聊天响应转换为 OpenAI 格式。
5. **`embeddingResponseBaidu2OpenAI`**: 将百度嵌入响应转换为 OpenAI 格式。
6. **`baiduStreamHandler`**: 处理百度流式聊天响应。
7. **`baiduHandler`**: 处理百度非流式聊天响应。
8. **`baiduEmbeddingHandler`**: 处理百度嵌入响应。

### 访问令牌管理函数
9. **`getBaiduAccessToken`**: 获取百度 API 访问令牌，支持缓存和自动刷新。
10. **`getBaiduAccessTokenHelper`**: 从百度 OAuth 端点获取访问令牌。

## 5. 关键逻辑分析
- **系统消息处理**：百度 API 将系统消息放在 `System` 字段中，而不是消息列表中。
- **流式响应**：使用 `helper.StreamScannerHandler` 处理流式响应，逐行解析 JSON。
- **访问令牌缓存**：使用 `sync.Map` 缓存访问令牌，支持自动刷新（过期前 1 小时）。
- **令牌格式**：API Key 格式为 `client_id|client_secret`，用于获取访问令牌。
- **错误处理**：检查响应中的错误码和错误消息，返回适当的错误信息。

## 6. 关联文件
- `baidu/adaptor.go`: 调用这些函数执行请求和处理响应。
- `baidu/constants.go`: 定义模型列表和渠道名称。
- `baidu/dto.go`: 定义请求和响应数据结构。
- `relay/helper/stream.go`: 流式响应处理辅助函数。
- `service/http_client.go`: HTTP 客户端管理。