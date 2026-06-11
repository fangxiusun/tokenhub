# rerank.go 代码阅读文档

## 1. 全局总结

本文件实现了 Rerank（重排序）响应的处理逻辑 `RerankHandler`。支持标准 Rerank 响应和 Xinference 特殊格式的转换。

## 2. 依赖关系

- `common`: JSON 反序列化
- `constant`: 渠道类型
- `dto`: RerankResponse DTO
- `relay/channel/xinference`: Xinference 特定响应格式
- `relay/common`: RelayInfo

## 3. 类型定义

本文件无自定义类型定义。

## 4. 函数详解

### `RerankHandler(c, info, resp) (*dto.Usage, *types.NewAPIError)`
- **功能**: 处理 Rerank 响应
- **逻辑**:
  1. 读取响应体
  2. 根据渠道类型选择解析方式
  3. Xinference: 转换 XinRerankResponse → RerankResponse
  4. 其他渠道: 直接解析为 RerankResponse
  5. 设置 PromptTokens = TotalTokens

## 5. 关键逻辑分析

1. **Xinference 适配**: Xinference 的 rerank 响应格式与标准不同，需要转换
2. **Document 处理**: 支持 ReturnDocuments，从原始请求中获取文档内容
3. **Token 统计**: 使用预估的 prompt tokens 作为使用量

## 6. 关联文件

- `relay/channel/xinference/`: Xinference 渠道适配器
- `dto/rerank.go`: RerankResponse DTO
