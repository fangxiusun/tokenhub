# constants.go 代码阅读文档

## 1. 全局总结
Vertex AI 渠道的常量定义文件，包含支持的模型列表和渠道名称。

## 2. 依赖关系
无外部依赖。

## 3. 类型定义

### ModelList
```go
var ModelList = []string{
    "meta/llama3-405b-instruct-maas",
}
```
Vertex AI 自有的模型列表。注意 Claude 和 Gemini 模型列表在各自的渠道包中定义，在 `GetModelList()` 中合并。

### ChannelName
```go
var ChannelName = "vertex-ai"
```

## 4. 函数详解
无函数定义。

## 5. 关键逻辑分析
- 当前 Vertex AI 模型列表仅包含一个开源模型 `meta/llama3-405b-instruct-maas`
- Claude 和 Gemini 模型通过 `adaptor.go` 中的 `GetModelList()` 从各自的渠道包动态合并
- 注释中显示之前有 Claude 和 Gemini 模型被注释掉，说明这些模型现在通过渠道包管理

## 6. 关联文件
- `vertex/adaptor.go` — `GetModelList()` 方法合并多个模型列表
- `relay/channel/claude/constants.go` — Claude 模型列表
- `relay/channel/gemini/constants.go` — Gemini 模型列表
