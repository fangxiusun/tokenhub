# compatible_handler.go 代码阅读文档

## 1. 全局总结

本文件实现了 OpenAI 兼容格式的文本请求处理入口 `TextHelper`。这是最常用的 handler，处理 `/v1/chat/completions` 等文本类请求。它支持 passthrough 模式、模型映射、系统提示注入、参数覆盖、Responses API 转换等核心功能。

## 2. 依赖关系

- `relay/common`: RelayInfo、参数覆盖、出站请求体
- `relay/helper`: 模型映射、价格计算
- `service`: 计费、响应转换
- `setting/model_setting`: 全局设置
- `setting/ratio_setting`: 音频比率

## 3. 类型定义

本文件无自定义类型定义。

## 4. 函数详解

### `TextHelper(c *gin.Context, info *relaycommon.RelayInfo) *types.NewAPIError`
- **功能**: 处理 OpenAI 兼容格式的文本类请求
- **流程**:
  1. InitChannelMeta
  2. 类型断言为 GeneralOpenAIRequest
  3. 深拷贝请求
  4. 模型映射
  5. 处理 StreamOptions（支持/强制配置）
  6. 获取适配器
  7. 判断是否走 Responses API 路径（chatCompletionsViaResponses）
  8. 构建请求体（passthrough 或 convert）
  9. 系统提示注入
  10. 移除禁用字段
  11. 参数覆盖
  12. 发送请求
  13. 处理响应
  14. 计费

## 5. 关键逻辑分析

1. **Responses API 路由**: 当满足条件时（非 passthrough、支持 Responses API），自动将 chat completions 请求转换为 Responses API 格式
2. **StreamOptions 处理**:
   - 不支持时设为 nil
   - 支持时根据配置决定是否强制 include_usage
3. **系统提示注入**: 支持追加和覆盖两种模式，覆盖时会设置 `ContextKeySystemPromptOverride`
4. **passthrough 模式**: 直接转发原始请求体，不做格式转换
5. **音频计费**: 检测音频 token 并使用对应的计费路径

## 6. 关联文件

- `relay/chat_completions_via_responses.go`: Responses API 转换逻辑
- `relay/channel/adapter.go`: Adaptor 接口
- `relay/common/override.go`: 参数覆盖实现
