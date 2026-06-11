# constant.go 代码阅读文档

## 1. 全局总结
Xinference 渠道的常量定义文件，包含支持的模型列表和渠道名称。

## 2. 依赖关系
无外部依赖。

## 3. 类型定义

### ModelList
```go
var ModelList = []string{
    "bge-reranker-v2-m3",
    "jina-reranker-v2",
}
```
Xinference 支持的重排序（Rerank）模型。

### ChannelName
```go
var ChannelName = "xinference"
```

## 4. 函数详解
无函数定义。

## 5. 关键逻辑分析
- Xinference 渠道专注于重排序（Rerank）功能
- 支持 BGE 和 Jina 两种重排序模型
- 文件名使用 `constant.go`（单数），与其他渠道的 `constants.go`（复数）不同

## 6. 关联文件
- `xinference/dto.go` — 重排序响应数据结构
