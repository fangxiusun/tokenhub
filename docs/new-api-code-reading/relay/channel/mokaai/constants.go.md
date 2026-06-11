# constants.go 代码阅读文档

## 1. 全局总结
该文件定义了 MokaAI 渠道的常量，包括支持的模型列表和渠道名称标识。

## 2. 依赖关系
无外部依赖。

## 3. 类型定义

### `ModelList` (var []string)
- **作用**：定义 MokaAI 渠道支持的模型列表
- **值**：包含三个 M3E 系列嵌入模型：`m3e-large`、`m3e-base`、`m3e-small`

### `ChannelName` (var string)
- **作用**：渠道名称标识
- **值**：`"mokaai"`

## 4. 函数详解
无函数定义。

## 5. 关键逻辑分析
- MokaAI 渠道目前仅支持嵌入模型（M3E 系列），不支持聊天补全等其他功能
- 模型列表用于渠道配置时的模型选择和验证

## 6. 关联文件
- `relay/channel/mokaai/adaptor.go` — 通过 `GetModelList()` 和 `GetChannelName()` 引用这些常量
