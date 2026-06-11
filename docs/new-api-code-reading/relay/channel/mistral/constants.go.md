# constants.go 代码阅读文档

## 1. 全局总结
本文件定义了 Mistral 渠道的常量，包括渠道名称和支持的模型列表。

## 2. 依赖关系
无外部依赖。

## 3. 类型定义

### `ModelList` 变量
```go
var ModelList = []string{
    "open-mistral-7b",
    "open-mixtral-8x7b",
    "mistral-small-latest",
    "mistral-medium-latest",
    "mistral-large-latest",
    "mistral-embed",
}
```
Mistral 渠道支持的模型：
- `open-mistral-7b` — 7B 参数的开源 Mistral 模型
- `open-mixtral-8x7b` — 8x7B MoE 架构的开源 Mixtral 模型
- `mistral-small-latest` — 最新 Small 规模模型
- `mistral-medium-latest` — 最新 Medium 规模模型
- `mistral-large-latest` — 最新 Large 规模模型
- `mistral-embed` — 嵌入模型

### `ChannelName` 变量
```go
var ChannelName = "mistral"
```
渠道标识名称。

## 4. 函数详解
无函数定义。

## 5. 关键逻辑分析

### 模型分类
Mistral 的模型列表覆盖了：
- **开源模型**: `open-mistral-7b`, `open-mixtral-8x7b` — 可自部署的开源模型
- **闭源模型**: `mistral-small/medium/latest` — Mistral 云服务提供的模型
- **嵌入模型**: `mistral-embed` — 文本嵌入模型

### Mixtral MoE
`open-mixtral-8x7b` 是 Mistral 的 Mixture of Experts 架构模型，使用 8 个 7B 专家网络，总参数量约 46B，但推理时只激活部分专家，效率更高。

### latest 标签
使用 `latest` 标签的模型（如 `mistral-large-latest`）指向 Mistral 提供的最新版本，随时间自动更新。

## 6. 关联文件
- `adaptor.go` — 使用 `ModelList` 和 `ChannelName`
