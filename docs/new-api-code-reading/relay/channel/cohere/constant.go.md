# constant.go 代码阅读文档

## 1. 全局总结

本文件定义了 Cohere 频道的模型列表和频道名称常量。模型涵盖 Cohere 的 Command 系列聊天模型和 Rerank 重排序模型。

## 2. 依赖关系

无外部依赖。

## 3. 类型定义

### `ModelList` 变量

```go
var ModelList = []string{
    "command-a-03-2025",
    "command-r", "command-r-plus",
    "command-r-08-2024", "command-r-plus-08-2024",
    "c4ai-aya-23-35b", "c4ai-aya-23-8b",
    "command-light", "command-light-nightly", "command", "command-nightly",
    "rerank-english-v3.0", "rerank-multilingual-v3.0",
    "rerank-english-v2.0", "rerank-multilingual-v2.0",
}
```

Cohere 支持的模型列表，分为两类：
- **聊天模型**: Command 系列（command-r、command-r-plus 等）和 C4AI Aya 系列
- **Rerank 模型**: rerank-english 和 rerank-multilingual 的多个版本

### `ChannelName` 变量

```go
var ChannelName = "cohere"
```

频道标识名称。

## 4. 函数详解

无函数定义。

## 5. 关键逻辑分析

### 模型分类
模型列表包含两种不同类型的模型，分别对应不同的 API 端点：
- 聊天模型使用 `/v1/chat` 端点
- Rerank 模型使用 `/v1/rerank` 端点

### 多语言支持
Rerank 模型区分英文和多语言版本，版本号从 v2.0 到 v3.0。

## 6. 关联文件

- `relay/channel/cohere/adaptor.go` - 返回此模型列表
