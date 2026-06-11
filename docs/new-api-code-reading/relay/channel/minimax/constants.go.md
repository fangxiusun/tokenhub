# constants.go 代码阅读文档

## 1. 全局总结
本文件定义了 MiniMax 渠道的常量，包括渠道名称和支持的模型列表。MiniMax 提供文本生成、语音合成和图像生成三类模型。

## 2. 依赖关系
无外部依赖。

## 3. 类型定义

### `ModelList` 变量
```go
var ModelList = []string{...}
```
MiniMax 渠道支持的模型列表：

**文本生成模型**:
- `abab6.5-chat`, `abab6.5s-chat`, `abab6-chat`, `abab5.5-chat`, `abab5.5s-chat` — Abab 系列
- `MiniMax-M2.7`, `MiniMax-M2.7-highspeed` — M2.7 系列
- `MiniMax-M2.1`, `MiniMax-M2.1-highspeed` — M2.1 系列
- `MiniMax-M2`, `MiniMax-M2.5`, `MiniMax-M2.5-highspeed` — M2/M2.5 系列

**语音合成模型**:
- `speech-2.5-hd-preview`, `speech-2.5-turbo-preview` — 2.5 版本
- `speech-02-hd`, `speech-02-turbo` — 02 版本
- `speech-01-hd`, `speech-01-turbo` — 01 版本

**图像生成模型**:
- `image-01` — 图像生成
- `image-01-live` — 实时图像生成

### `ChannelName` 变量
```go
var ChannelName = "minimax"
```
渠道标识名称。

## 4. 函数详解
无函数定义。

## 5. 关键逻辑分析

### 模型分类
MiniMax 的模型列表按能力分类：
- Abab 系列为早期文本模型
- M 系列为新一代文本模型，区分标准版和 highspeed 版
- Speech 系列按版本和质量分级（hd/turbo）
- Image 系列支持标准和实时生成

### highspeed 变体
MiniMax 为多个模型提供 `highspeed` 变体（如 `MiniMax-M2.7-highspeed`），这是 MiniMax 特有的快速推理版本。

## 6. 关联文件
- `adaptor.go` — 使用 `ModelList` 和 `ChannelName`
