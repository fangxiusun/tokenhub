# embedding_handler.go 代码阅读文档

## 1. 全局总结

本文件实现了 Embedding（向量嵌入）请求的处理入口 `EmbeddingHelper`。处理流程与其他 handler 类似：模型映射 → 请求转换 → 发送 → 响应解析 → 计费。

## 2. 依赖关系

- `relay/common`: RelayInfo、参数覆盖
- `relay/helper`: 模型映射
- `service`: 计费

## 3. 类型定义

本文件无自定义类型定义。

## 4. 函数详解

### `EmbeddingHelper(c *gin.Context, info *relaycommon.RelayInfo) *types.NewAPIError`
- **功能**: 处理 Embedding 请求
- **流程**: InitChannelMeta → 类型断言 → 深拷贝 → 模型映射 → 获取适配器 → ConvertEmbeddingRequest → 参数覆盖 → DoRequest → DoResponse → PostTextConsumeQuota

## 5. 关键逻辑分析

1. **简单流程**: 相比 TextHelper，没有 passthrough、StreamOptions、系统提示等复杂逻辑
2. **参数覆盖**: 支持对转换后的请求进行参数覆盖
3. **计费**: 使用 PostTextConsumeQuota 进行文本计费

## 6. 关联文件

- `relay/channel/adapter.go`: Adaptor.ConvertEmbeddingRequest 接口
- `dto/embedding.go`: EmbeddingRequest DTO
