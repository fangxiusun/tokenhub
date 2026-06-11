# chat_completions_via_responses.go 代码阅读文档

## 1. 全局总结

本文件实现了将 Chat Completions 请求自动转换为 Responses API 格式的逻辑。当上游渠道支持 Responses API 但不支持 Chat Completions 时，系统会自动进行格式转换。包含系统提示注入和完整的请求-响应转换链路。

## 2. 依赖关系

- `relay/channel/openai`: Responses 到 Chat 的流式/非流式响应转换器
- `relay/common`: RelayInfo、参数覆盖
- `service`: 请求格式转换
- `dto`: 请求/响应 DTO

## 3. 类型定义

本文件无自定义类型定义。

## 4. 函数详解

### `applySystemPromptIfNeeded(c, info, request)`
- **功能**: 根据渠道设置注入系统提示
- **逻辑**: 检查是否已有系统消息 → 无则追加 / 有则根据配置覆盖

### `chatCompletionsViaResponses(c, info, adaptor, request) (*dto.Usage, *types.NewAPIError)`
- **功能**: 将 Chat Completions 请求转换为 Responses API 格式并发送
- **流程**:
  1. 序列化请求 → 移除禁用字段 → 参数覆盖
  2. 转换为 ResponsesRequest
  3. 临时修改 RelayMode 和 RequestURLPath
  4. 调用 adaptor.ConvertOpenAIResponsesRequest
  5. 发送请求
  6. 根据流式/非流式选择对应的响应转换器

## 5. 关键逻辑分析

1. **格式转换链**: OpenAI Chat → OpenAI Responses → 上游渠道格式
2. **RelayMode 临时修改**: 在转换期间临时将 RelayMode 设为 Responses，完成后恢复
3. **响应转换**: 使用 `OaiResponsesToChatStreamHandler` 和 `OaiResponsesToChatHandler` 将 Responses 格式转回 Chat Completions 格式
4. **双重参数覆盖**: 在转换前后各执行一次参数覆盖

## 6. 关联文件

- `relay/channel/openai/responses_handler.go`: Responses 响应转换器
- `service/convert.go`: 请求格式转换服务
