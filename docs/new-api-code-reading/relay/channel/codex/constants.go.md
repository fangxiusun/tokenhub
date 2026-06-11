# constants.go 代码阅读文档

## 1. 全局总结

本文件定义了 Codex 频道的模型列表和频道名称常量。模型列表基于 GPT-5 系列模型，并通过 `withCompactModelSuffix` 函数自动生成 Compact 变体模型名称。

## 2. 依赖关系

| 包 | 用途 |
|---|------|
| `setting/ratio_setting` | 计费比率设置，提供 Compact 模型后缀生成 |
| `lo` | 工具库（Map、Uniq） |

## 3. 类型定义

### `baseModelList` 变量

```go
var baseModelList = []string{
    "gpt-5", "gpt-5-codex", "gpt-5-codex-mini",
    "gpt-5.1", "gpt-5.1-codex", "gpt-5.1-codex-max", "gpt-5.1-codex-mini",
    "gpt-5.2", "gpt-5.2-codex", "gpt-5.3-codex", "gpt-5.3-codex-spark",
    "gpt-5.4",
}
```

GPT-5 系列基础模型列表。

### `ModelList` 变量

```go
var ModelList = withCompactModelSuffix(baseModelList)
```

最终模型列表，包含基础模型和 Compact 变体。

### `ChannelName` 常量

```go
const ChannelName = "codex"
```

频道标识名称。

## 4. 函数详解

### `withCompactModelSuffix(models []string) []string`
为模型列表生成 Compact 变体：
1. 保留原始模型列表
2. 对每个模型调用 `ratio_setting.WithCompactModelSuffix` 生成 Compact 后缀版本
3. 使用 `lo.Uniq` 去重
4. 返回合并后的唯一模型列表

## 5. 关键逻辑分析

### Compact 模型变体
Codex 支持两种模式：普通 Responses 和 Compact Responses。Compact 模式对应更精简的模型变体，通过 `ratio_setting.WithCompactModelSuffix` 自动添加后缀。这确保了两种模式的模型都能被正确识别和路由。

### 模型命名模式
- `gpt-5`: 基础模型
- `gpt-5-codex`: Codex 专用变体
- `gpt-5-codex-mini`: Codex 轻量变体
- `gpt-5.1-codex-max`: Codex 最大上下文变体
- `gpt-5.3-codex-spark`: Codex Spark 变体

## 6. 关联文件

- `relay/channel/codex/adaptor.go` - 通过 `GetModelList()` 返回此模型列表
- `setting/ratio_setting/` - 提供 Compact 模型后缀生成逻辑
