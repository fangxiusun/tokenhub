# constants.go 代码阅读文档

## 1. 全局总结
智谱 AI 旧版渠道的常量定义文件。

## 2. 依赖关系
无外部依赖。

## 3. 类型定义

### ModelList
```go
var ModelList = []string{
    "chatglm_turbo", "chatglm_pro", "chatglm_std", "chatglm_lite",
}
```
智谱 ChatGLM 系列模型。

### ChannelName
```go
var ChannelName = "zhipu"
```

## 4. 函数详解
无函数定义。

## 5. 关键逻辑分析
- 包含 ChatGLM 的四个版本：turbo、pro、std、lite
- 这是旧版 API，新版使用 `zhipu_4v` 渠道

## 6. 关联文件
- `zhipu/adaptor.go` — 使用 ModelList 和 ChannelName
