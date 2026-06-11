# claude_handler.go 代码阅读文档

## 1. 全局总结

本文件实现了 Claude（Anthropic）原生格式请求的处理入口 `ClaudeHelper`。支持 Claude 请求的模型映射、Thinking 模式适配（扩展思考）、effort level 处理、系统提示注入，以及自动转换为 Responses API 的能力。

## 2. 依赖关系

- `relay/common`: RelayInfo、参数覆盖
- `relay/helper`: 模型映射
- `service`: 计费、响应转换
- `setting/model_setting`: Claude 设置（Thinking 适配器、默认 max_tokens）
- `setting/reasoning`: 思考努力级别后缀处理

## 3. 类型定义

本文件无自定义类型定义。

## 4. 函数详解

### `ClaudeHelper(c *gin.Context, info *relaycommon.RelayInfo) *types.NewAPIError`
- **功能**: 处理 Claude 原生格式请求
- **核心流程**:
  1. 类型断言为 ClaudeRequest
  2. 模型映射
  3. 默认 max_tokens 设置
  4. Thinking 模式适配（effort level / -thinking 后缀）
  5. 系统提示注入
  6. Responses API 路由判断
  7. 构建/发送/解析请求
  8. 计费

## 5. 关键逻辑分析

1. **Effort Level 处理**: 支持 claude-opus-4-6/4-7/4-8 的 effort level 后缀（如 `-high`），自动转换为 Thinking 配置
2. **Thinking 适配器**: 
   - `-thinking` 后缀模型自动启用 extended thinking
   - Opus 4.7/4.8 使用 adaptive thinking + high effort
   - 其他模型使用 enabled thinking + budget_tokens（80% of max_tokens）
3. **Opus 4.7/4.8 特殊处理**: 拒绝非默认 temperature/top_p/top_k，强制设为 nil
4. **Responses API 路由**: 与 TextHelper 类似，支持自动转换为 Responses API
5. **passthrough 模式**: 支持全局和渠道级别的 passthrough

## 6. 关联文件

- `relay/compatible_handler.go`: OpenAI 格式的类似处理
- `relay/chat_completions_via_responses.go`: Responses API 转换
- `setting/model_setting/claude.go`: Claude 特定设置
