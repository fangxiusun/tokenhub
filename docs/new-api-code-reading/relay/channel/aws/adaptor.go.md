# adaptor.go 代码阅读文档

## 1. 全局总结
该文件是 AWS Bedrock 渠道的适配器实现，负责将 OpenAI/Claude 格式的请求转换为 AWS Bedrock API 格式。支持两种客户端模式：API Key 模式和 AKSK（Access Key/Secret Key）模式，并支持 Claude 和 Nova 两种模型系列。

## 2. 依赖关系
- 标准库：`fmt`, `io`, `net/http`, `strings`
- 内部包：
  - `dto`: 数据传输对象
  - `relay/channel`: 渠道通用工具
  - `relay/channel/claude`: Claude 渠道适配器
  - `relaycommon`: 中继通用配置
  - `service`: 业务逻辑服务
  - `types`: 类型定义
- 第三方包：
  - `github.com/aws/aws-sdk-go-v2/service/bedrockruntime`: AWS Bedrock Runtime SDK
  - `github.com/pkg/errors`: 错误处理
  - `github.com/gin-gonic/gin`: HTTP 框架

## 3. 类型定义
### 枚举类型
- `ClientMode`: 客户端模式枚举，包含 `ClientModeApiKey`（API Key 模式）和 `ClientModeAKSK`（AKSK 模式）。

### 结构体
- `Adaptor`: AWS 渠道适配器结构体，包含：
  - `ClientMode`: 客户端模式
  - `AwsClient`: AWS Bedrock Runtime 客户端
  - `AwsModelId`: AWS 模型 ID
  - `AwsReq`: AWS 请求对象
  - `IsNova`: 是否为 Nova 模型

## 4. 函数详解
### 核心函数
1. **`ConvertClaudeRequest`**: 转换 Claude 请求，处理图像 URL 到 Base64 的转换。
2. **`GetRequestURL`**: 根据客户端模式构建请求 URL，API Key 模式使用 AWS REST API。
3. **`SetupRequestHeader`**: 设置请求头，包括 Claude 通用头和授权信息。
4. **`ConvertOpenAIRequest`**: 转换 OpenAI 请求，区分 Nova 和 Claude 模型。
5. **`DoRequest`**: 执行请求，根据客户端模式选择不同的请求方式。
6. **`DoResponse`**: 处理响应，根据客户端模式和模型类型调用不同的处理器。

## 5. 关键逻辑分析
- **双模式支持**：API Key 模式使用 REST API，AKSK 模式使用 AWS SDK。
- **Nova 模型支持**：Nova 模型使用不同的请求格式和响应处理逻辑。
- **图像处理**：Claude 请求中的图像 URL 自动转换为 Base64 格式。
- **跨区域支持**：通过 `getAwsModelID` 函数映射模型 ID，支持跨区域调用。

## 6. 关联文件
- `aws/constants.go`: 定义模型映射和区域配置。
- `aws/dto.go`: AWS 特定的数据传输对象。
- `aws/relay-aws.go`: AWS 请求执行和响应处理。
- `relay/channel/claude/adaptor.go`: Claude 渠道适配器，用于 Claude 请求处理。