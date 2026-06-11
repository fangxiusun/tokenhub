# constants.go 代码阅读文档

## 1. 全局总结
火山引擎渠道的常量定义文件，包含支持的模型列表和渠道名称。

## 2. 依赖关系
无外部依赖。

## 3. 类型定义

### ModelList
```go
var ModelList = []string{
    "Doubao-pro-128k",
    "Doubao-pro-32k",
    "Doubao-pro-4k",
    "Doubao-lite-128k",
    "Doubao-lite-32k",
    "Doubao-lite-4k",
    "Doubao-embedding",
    "doubao-seedream-4-0-250828",
    "seedream-4-0-250828",
    "doubao-seedance-1-0-pro-250528",
    "seedance-1-0-pro-250528",
    "doubao-seed-1-6-thinking-250715",
    "seed-1-6-thinking-250715",
}
```
支持的模型包括：
- **Doubao Pro**: 高性能版（4k/32k/128k 上下文）
- **Doubao Lite**: 轻量版（4k/32k/128k 上下文）
- **Doubao Embedding**: 嵌入模型
- **Seedream**: 图片生成模型
- **Seedance**: 视频生成模型
- **Seed Thinking**: 思考型模型

### ChannelName
```go
var ChannelName = "volcengine"
```

## 4. 函数详解
无函数定义。

## 5. 关键逻辑分析
- 模型列表覆盖了文本生成、嵌入、图片生成、视频生成等多种能力
- 部分模型有多个名称变体（如 `doubao-seedream-4-0-250828` 和 `seedream-4-0-250828`）

## 6. 关联文件
- `volcengine/adaptor.go` — 使用 ModelList 和 ChannelName
