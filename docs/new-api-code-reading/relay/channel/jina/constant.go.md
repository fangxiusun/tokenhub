# constant.go 代码阅读文档

## 1. 全局总结
本文件定义了 Jina 渠道的常量，包括渠道名称和支持的模型列表。

## 2. 依赖关系
无外部依赖。

## 3. 类型定义

### `ModelList` 变量
```go
var ModelList = []string{
    "jina-clip-v1",
    "jina-reranker-v2-base-multilingual",
    "jina-reranker-m0",
}
```
Jina 渠道支持的模型：
- `jina-clip-v1` — CLIP 多模态嵌入模型
- `jina-reranker-v2-base-multilingual` — 多语言重排序模型 v2
- `jina-reranker-m0` — 重排序模型 m0

### `ChannelName` 变量
```go
var ChannelName = "jina"
```
渠道标识名称。

## 4. 函数详解
无函数定义。

## 5. 关键逻辑分析
Jina 的模型列表覆盖了嵌入（CLIP）和重排序两类能力，与其渠道定位一致。`jina-clip-v1` 是一个视觉-语言嵌入模型，支持图像和文本的跨模态嵌入。

## 6. 关联文件
- `adaptor.go` — 使用 `ModelList` 和 `ChannelName`
