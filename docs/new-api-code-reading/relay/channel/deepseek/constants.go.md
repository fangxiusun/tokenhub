# constants.go 代码阅读文档

## 1. 全局总结

本文件定义了 DeepSeek 频道的模型列表和频道名称常量。

## 2. 依赖关系

无外部依赖。

## 3. 类型定义

### `ModelList` 变量

```go
var ModelList = []string{
    "deepseek-chat", "deepseek-reasoner",
    "deepseek-v4-flash", "deepseek-v4-flash-none", "deepseek-v4-flash-max",
    "deepseek-v4-pro", "deepseek-v4-pro-none", "deepseek-v4-pro-max",
}
```

DeepSeek 支持的模型：
- **deepseek-chat**: 通用聊天模型
- **deepseek-reasoner**: 推理模型（DeepSeek-R1）
- **deepseek-v4-flash**: V4 快速版，支持 none/max 变体
- **deepseek-v4-pro**: V4 专业版，支持 none/max 变体

### `ChannelName` 变量

```go
var ChannelName = "deepseek"
```

频道标识名称。

## 4. 函数详解

无函数定义。

## 5. 关键逻辑分析

### V4 模型命名模式
V4 模型使用 `{model}-{variant}` 命名：
- `none`: 无思考模式
- `max`: 最大思考模式
- 无后缀: 默认模式

这些变体通过 `adaptor.go` 中的思考后缀解析机制自动转换为 API 参数。

### 模型数量
共 8 个模型，覆盖了聊天、推理和 V4 系列。

## 6. 关联文件

- `relay/channel/deepseek/adaptor.go` - 返回此模型列表
