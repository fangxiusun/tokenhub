# constants.go 代码阅读文档

## 1. 全局总结
该文件定义了 Moonshot（月之暗面/Kimi）渠道的常量，包括支持的模型列表和渠道名称标识。

## 2. 依赖关系
无外部依赖。

## 3. 类型定义

### `ModelList` (var []string)
- **作用**：定义 Moonshot 渠道支持的模型列表
- **值**：包含 5 个 Kimi K2 系列模型：
  - `kimi-k2.5`
  - `kimi-k2-0905-preview`
  - `kimi-k2-turbo-preview`
  - `kimi-k2-thinking`
  - `kimi-k2-thinking-turbo`

### `ChannelName` (var string)
- **作用**：渠道名称标识
- **值**：`"moonshot"`

## 4. 函数详解
无函数定义。

## 5. 关键逻辑分析
- Moonshot 渠道专注于 Kimi K2 系列模型，包括标准版、Turbo 版和思维推理版
- thinking 版本支持思维链推理功能

## 6. 关联文件
- `relay/channel/moonshot/adaptor.go` — 通过 `GetModelList()` 和 `GetChannelName()` 引用这些常量
