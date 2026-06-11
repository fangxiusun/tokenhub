# rerank_handler.go 代码阅读文档

## 1. 全局总结

本文件实现了 Rerank（重排序）请求的处理入口 `RerankHelper`。Rerank 是一种文档重排序服务，接收查询和文档列表，返回相关性评分。处理流程与其他 handler 类似。

## 2. 依赖关系

- `relay/common`: RelayInfo、参数覆盖
- `relay/helper`: 模型映射
- `service`: 计费

## 3. 类型定义

本文件无自定义类型定义。

## 4. 函数详解

### `RerankHelper(c *gin.Context, info *relaycommon.RelayInfo) *types.NewAPIError`
- **功能**: 处理 Rerank 请求
- **流程**: InitChannelMeta → 类型断言 → 深拷贝 → 模型映射 → 获取适配器 → ConvertRerankRequest → 参数覆盖 → DoRequest → DoResponse → PostTextConsumeQuota

## 5. 关键逻辑分析

1. **passthrough 模式**: 支持全局和渠道级别
2. **参数覆盖**: 支持对转换后的请求进行参数覆盖
3. **简单流程**: 与 EmbeddingHelper 类似

## 6. 关联文件

- `relay/channel/adapter.go`: Adaptor.ConvertRerankRequest 接口
- `relay/common_handler/rerank.go`: Rerank 响应处理
