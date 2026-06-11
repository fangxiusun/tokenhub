# constants.go 代码阅读文档

## 1. 全局总结

本文件定义了 Dify 频道的模型列表和频道名称常量。由于 Dify 是平台型服务，模型列表当前为空。

## 2. 依赖关系

无外部依赖。

## 3. 类型定义

### `ModelList` 变量

```go
var ModelList []string
```

空的模型列表切片。Dify 作为 LLM 应用开发平台，其模型由 Dify 内部管理，不直接暴露模型名给外部系统。

### `ChannelName` 变量

```go
var ChannelName = "dify"
```

频道标识名称。

## 4. 函数详解

无函数定义。

## 5. 关键逻辑分析

### 空模型列表的原因
Dify 是一个应用开发平台，而非模型提供商。用户在 Dify 中创建 Bot/Workflow，每个 Bot 有独立的标识。模型选择在 Dify 平台内部完成，外部 API 只需指定 Bot ID。

## 6. 关联文件

- `relay/channel/dify/adaptor.go` - 返回此模型列表
