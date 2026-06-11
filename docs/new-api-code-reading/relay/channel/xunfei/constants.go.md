# constants.go 代码阅读文档

## 1. 全局总结
讯飞星火渠道的常量定义文件，包含支持的模型列表和渠道名称。

## 2. 依赖关系
无外部依赖。

## 3. 类型定义

### ModelList
```go
var ModelList = []string{
    "SparkDesk",
    "SparkDesk-v1.1",
    "SparkDesk-v2.1",
    "SparkDesk-v3.1",
    "SparkDesk-v3.5",
    "SparkDesk-v4.0",
}
```
讯飞星火支持的模型版本。

### ChannelName
```go
var ChannelName = "xunfei"
```

## 4. 函数详解
无函数定义。

## 5. 关键逻辑分析
- 模型名称包含版本号（v1.1 到 v4.0），对应不同的 API 版本
- 版本号用于确定 WebSocket 端点和 domain 参数

## 6. 关联文件
- `xunfei/adaptor.go` — 使用 ModelList 和 ChannelName
- `xunfei/relay-xunfei.go` — 根据模型名提取 API 版本
