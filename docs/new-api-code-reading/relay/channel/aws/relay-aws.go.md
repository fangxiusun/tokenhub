# relay-aws.go 代码阅读文档

## 1. 全局总结
该文件实现了 AWS Bedrock 渠道的请求执行和响应处理，包括 Claude 和 Nova 两种模型系列。支持同步和异步请求，处理流式响应和非流式响应。

## 2. 依赖关系
- 标准库：`context`, `encoding/json`, `fmt`, `io`, `net/http`, `strings`, `time`
- 内部包：
  - `github.com/QuantumNous/new-api/common`: 通用工具函数
  - `github.com/QuantumNous/new-api/dto`: 数据传输对象
  - `relay/channel`: 渠道通用工具
  - `relay/channel/claude`: Claude 渠道适配器
  - `relaycommon`: 中继通用配置
  - `relay/helper`: 中继辅助函数
  - `service`: 业务逻辑服务
  - `types`: 类型定义
  - `model_setting`: 模型设置
- 第三方包：
  - `github.com/gin-gonic/gin`: HTTP 框架
  - `github.com/pkg/errors`: 错误处理
  - `github.com/aws/aws-sdk-go-v2/aws`: AWS SDK
  - `github.com/aws/aws-sdk-go-v2/credentials`: AWS 凭证
  - `github.com/aws/aws-sdk-go-v2/service/bedrockruntime`: AWS Bedrock Runtime SDK
  - `github.com/aws/smithy-go/auth/bearer`: AWS Bearer 认证

## 3. 类型定义
该文件没有定义新的类型，主要使用 `aws/dto.go` 中定义的类型。

## 4. 函数详解
### 工具函数
1. **`getAwsErrorStatusCode`**: 从 AWS SDK 错误中提取 HTTP 状态码。
2. **`newAwsInvokeContext`**: 创建带有超时的 AWS 调用上下文。
3. **`newAwsClient`**: 创建 AWS Bedrock Runtime 客户端，支持 API Key 和 AKSK 模式。
4. **`getAwsRegionPrefix`**: 从区域 ID 中提取前缀（如 `us-east-1` 提取 `us`）。
5. **`awsModelCanCrossRegion`**: 检查模型是否支持跨区域调用。
6. **`awsModelCrossRegion`**: 构建跨区域模型 ID。
7. **`getAwsModelID`**: 获取 AWS 模型 ID，支持模型映射。

### 请求执行函数
8. **`doAwsClientRequest`**: 执行 AWS 客户端请求，支持 Claude 和 Nova 模型。
9. **`buildAwsRequestBody`**: 构建 AWS 请求体，支持 passthrough 模式。

### 响应处理函数
10. **`awsHandler`**: 处理非流式 Claude 响应。
11. **`awsStreamHandler`**: 处理流式 Claude 响应。
12. **`handleNovaRequest`**: 处理 Nova 模型响应，转换为 OpenAI 格式。

## 5. 关键逻辑分析
- **双模式支持**：API Key 模式使用 REST API，AKSK 模式使用 AWS SDK。
- **跨区域调用**：通过模型 ID 映射和区域前缀实现跨区域调用。
- **Passthrough 模式**：支持直接传递原始请求体，跳过格式转换。
- **流式处理**：Claude 模型支持流式响应，Nova 模型仅支持非流式响应。
- **超时控制**：通过上下文超时控制 AWS 调用时间。

## 6. 关联文件
- `aws/adaptor.go`: 调用这些函数执行请求和处理响应。
- `aws/constants.go`: 定义模型映射和区域配置。
- `aws/dto.go`: 定义请求和响应数据结构。
- `relay/channel/claude/relay-claude.go`: Claude 响应处理逻辑。
- `service/http_client.go`: HTTP 客户端管理。