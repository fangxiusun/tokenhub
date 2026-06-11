# dto.go 代码阅读文档

## 1. 全局总结
该文件定义了 AWS Bedrock 渠道特定的数据传输对象，包括 Claude 请求格式和 Nova 请求格式。主要用于将 OpenAI 格式的请求转换为 AWS Bedrock 支持的格式。

## 2. 依赖关系
- 标准库：`context`, `encoding/json`, `io`, `net/http`, `strings`
- 内部包：
  - `github.com/QuantumNous/new-api/common`: 通用工具函数
  - `github.com/QuantumNous/new-api/dto`: 数据传输对象
  - `github.com/QuantumNous/new-api/logger`: 日志记录

## 3. 类型定义
### Claude 请求相关
- `AwsClaudeRequest`: AWS Claude 请求结构体，包含 Anthropic 版本、消息、参数等。

### Nova 请求相关
- `NovaMessage`: Nova 消息结构体，包含角色和内容。
- `NovaContent`: Nova 内容结构体，包含文本。
- `NovaRequest`: Nova 请求结构体，包含模式版本、消息和推理配置。
- `NovaInferenceConfig`: Nova 推理配置，包含最大 token 数、温度、TopP 等。

## 4. 函数详解
1. **`formatRequest`**: 格式化 Claude 请求，设置 Anthropic 版本和 Beta 头。
2. **`convertToNovaRequest`**: 将 OpenAI 请求转换为 Nova 格式。
3. **`parseStopSequences`**: 解析停止序列，支持字符串、字符串数组和接口数组。

## 5. 关键逻辑分析
- **Claude 请求格式**：`AwsClaudeRequest` 符合 AWS Bedrock 的 Claude API 规范。
- **Nova 请求格式**：`NovaRequest` 使用 `messages-v1` 模式版本，支持推理配置。
- **Beta 头处理**：从请求头中提取 `anthropic-beta` 值，转换为 JSON 数组。
- **停止序列解析**：支持多种输入格式，确保停止序列正确处理。

## 6. 关联文件
- `aws/adaptor.go`: 使用这些 DTO 进行请求转换。
- `aws/relay-aws.go`: 使用 `formatRequest` 和 `convertToNovaRequest` 处理请求。
- `dto/dto.go`: 通用数据传输对象定义。