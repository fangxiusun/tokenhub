# constants.go 代码阅读文档

## 1. 全局总结
腾讯混元渠道的常量定义文件，包含支持的模型列表和渠道名称。

## 2. 依赖关系
无外部依赖，纯常量定义文件。

## 3. 类型定义

### ModelList
```go
var ModelList = []string{
    "hunyuan-lite",
    "hunyuan-standard",
    "hunyuan-standard-256K",
    "hunyuan-pro",
}
```
腾讯混元支持的模型列表，包含轻量版、标准版、256K 上下文版和专业版。

### ChannelName
```go
var ChannelName = "tencent"
```
渠道标识名称。

## 4. 函数详解
无函数定义。

## 5. 关键逻辑分析
模型列表对应腾讯云混元大模型的不同规格：
- `hunyuan-lite` — 轻量版，适合快速响应
- `hunyuan-standard` — 标准版
- `hunyuan-standard-256K` — 支持 256K 上下文的标准版
- `hunyuan-pro` — 专业版，能力最强

## 6. 关联文件
- `tencent/adaptor.go` — 使用 ModelList 和 ChannelName
