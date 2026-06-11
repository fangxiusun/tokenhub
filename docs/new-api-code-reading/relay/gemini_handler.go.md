# gemini_handler.go 代码阅读文档

## 1. 全局总结

本文件实现了 Google Gemini 原生格式请求的处理入口，包括聊天请求 `GeminiHelper` 和嵌入请求 `GeminiEmbeddingHandler`。支持 Thinking 模式适配、系统提示注入、批量嵌入请求处理。

## 2. 依赖关系

- `relay/channel/gemini`: Gemini 特定的 Thinking 适配器
- `relay/common`: RelayInfo、参数覆盖
- `relay/helper`: 模型映射
- `service`: 计费

## 3. 类型定义

本文件无自定义类型定义。

## 4. 函数详解

### `isNoThinkingRequest(req *dto.GeminiChatRequest) bool`
- 判断请求是否明确禁用了思考（ThinkingBudget == 0）

### `trimModelThinking(modelName string) string`
- 去除模型名称中的 `-nothinking`、`-thinking`、`-thinking-N` 后缀

### `GeminiHelper(c *gin.Context, info *relaycommon.RelayInfo) *types.NewAPIError`
- **功能**: 处理 Gemini 原生格式的聊天请求
- **特殊逻辑**:
  - Thinking 适配器：自动为没有 ThinkingConfig 的请求添加思考配置
  - 无思考请求：检测 `-nothinking` 后缀并尝试获取对应价格
  - 系统提示注入到 SystemInstructions

### `GeminiEmbeddingHandler(c *gin.Context, info *relaycommon.RelayInfo) *types.NewAPIError`
- **功能**: 处理 Gemini 嵌入请求（单条和批量）
- **支持**: `:embedContent` 和 `:batchEmbedContents` 两种端点

## 5. 关键逻辑分析

1. **Thinking 适配**: 通过 `gemini.ThinkingAdaptor` 自动配置思考预算
2. **无思考价格**: 当检测到 ThinkingBudget=0 时，尝试使用 `-nothinking` 后缀的模型价格
3. **系统提示**: 注入到 `SystemInstructions.Parts`，支持追加和覆盖模式
4. **空系统提示清理**: 如果 SystemInstructions 所有 Part 的 Text 为空，则设为 nil
5. **批量嵌入**: 支持 Gemini 的批量嵌入 API，自动提取所有输入文本

## 6. 关联文件

- `relay/channel/gemini/`: Gemini 渠道适配器
- `dto/gemini.go`: Gemini 请求/响应 DTO
