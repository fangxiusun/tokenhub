# relay-zhipu_v4.go 代码阅读文档

## 1. 全局总结
智谱 AI v4 版渠道的请求转换文件。将 OpenAI 格式请求转换为智谱 v4 格式，处理消息内容的多模态格式和参数兼容性。

## 2. 依赖关系
- **标准库**: strings
- **内部包**: `github.com/QuantumNous/new-api/dto`

## 3. 类型定义
无新增类型定义。

## 4. 函数详解

### requestOpenAI2Zhipu(request dto.GeneralOpenAIRequest) *dto.GeneralOpenAIRequest
将 OpenAI 请求转换为智谱 v4 格式：
1. 遍历消息列表：
   - 纯文本消息: 直接复制 Role 和 Content
   - 多模态消息: 解析内容数组，处理图片 URL 的 Base64 前缀
2. 处理 Stop 参数（支持 string 或 []string 类型）
3. 构建新的请求对象，包含所有必要参数
4. 处理 MaxTokens 和 MaxCompletionTokens 的兼容性

## 5. 关键逻辑分析

1. **Base64 图片处理**: 如果图片 URL 以 `data:image/` 开头，去除 URL 前缀（如 `data:image/png;base64,`），只保留纯 Base64 数据。

2. **消息格式统一**: 将 `dto.Message` 转换为标准格式，确保 Role、Content、ToolCalls、ToolCallId 字段正确传递。

3. **Stop 参数兼容**: 智谱 v4 的 Stop 字段使用数组格式，需要处理 string → []string 的转换。

4. **MaxTokens 处理**: 检查 `MaxTokens` 和 `MaxCompletionTokens`，使用 `GetMaxTokens()` 方法获取最终值。

5. **THINKING 支持**: 传递 `THINKING` 字段以支持思考模式。

## 6. 关联文件
- `zhipu_4v/adaptor.go` — 在 `ConvertOpenAIRequest` 中调用
- `zhipu_4v/dto.go` — 相关数据结构
- `dto/message.go` — 消息内容解析方法
