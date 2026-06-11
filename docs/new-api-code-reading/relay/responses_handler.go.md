# responses_handler.go 代码阅读文档

## 1. 全局总结

本文件实现了 OpenAI Responses API 格式请求的处理入口 `ResponsesHelper`。支持标准 Responses 请求和 Compaction（压缩）请求两种模式，处理请求转换、发送、响应解析和计费。

## 2. 依赖关系

- `relay/common`: RelayInfo、参数覆盖
- `relay/helper`: 模型映射、价格计算
- `service`: 计费

## 3. 类型定义

本文件无自定义类型定义。

## 4. 函数详解

### `ResponsesHelper(c *gin.Context, info *relaycommon.RelayInfo) *types.NewAPIError`
- **功能**: 处理 OpenAI Responses API 请求
- **支持的请求类型**:
  - `dto.OpenAIResponsesRequest`: 标准 Responses 请求
  - `dto.OpenAIResponsesCompactionRequest`: 压缩请求（自动转换为标准格式）
- **特殊逻辑**:
  - Compaction 模式仅支持 OpenAI 和 Codex 渠道
  - Compaction 模式需要重新计算价格（基于预估 prompt tokens）
  - gpt-4o-audio 模型使用音频计费路径

## 5. 关键逻辑分析

1. **请求类型适配**: 自动将 CompactionRequest 转换为 ResponsesRequest
2. **Compaction 计费**: 使用 `ModelPriceHelper` 重新计算价格，基于预估的 prompt tokens
3. **passthrough 模式**: 支持全局和渠道级别的 passthrough
4. **音频检测**: gpt-4o-audio 模型使用 `PostAudioConsumeQuota`

## 6. 关联文件

- `relay/chat_completions_via_responses.go`: Chat Completions 到 Responses 的转换
- `dto/responses.go`: Responses API DTO
