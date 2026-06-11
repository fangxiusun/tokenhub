# constants.go 代码阅读文档

## 1. 全局总结
智谱 AI v4 版渠道的常量定义文件。

## 2. 依赖关系
无外部依赖。

## 3. 类型定义

### ModelList
```go
var ModelList = []string{
    "glm-4", "glm-4v", "glm-3-turbo", "glm-4-alltools",
    "glm-4-plus", "glm-4-0520", "glm-4-air", "glm-4-airx",
    "glm-4-long", "glm-4-flash", "glm-4v-plus",
    "glm-4.6", "glm-4.6v", "glm-4.7", "glm-4.7-flash", "glm-5",
}
```
智谱 GLM 系列模型，包括文本和多模态版本。

### ChannelName
```go
var ChannelName = "zhipu_4v"
```

## 4. 函数详解
无函数定义。

## 5. 关键逻辑分析
- 模型涵盖 GLM-3 到 GLM-5 多个版本
- 多模态模型以 `v` 结尾（如 `glm-4v`、`glm-4.6v`）
- 特殊版本包括 `alltools`（工具增强）、`air`/`airx`（轻量版）、`long`（长上下文）、`flash`（快速版）

## 6. 关联文件
- `zhipu_4v/adaptor.go` — 使用 ModelList 和 ChannelName
