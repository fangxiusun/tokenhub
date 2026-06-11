# constants.go 代码阅读文档

## 1. 全局总结

本文件定义了 Coze 频道的模型列表和频道名称常量。Coze 作为 AI Bot 平台，支持多种第三方模型。

## 2. 依赖关系

无外部依赖。

## 3. 类型定义

### `ModelList` 变量

```go
var ModelList = []string{
    "moonshot-v1-8k", "moonshot-v1-32k", "moonshot-v1-128k",
    "Baichuan4", "abab6.5s-chat-pro", "glm-4-0520",
    "qwen-max", "deepseek-r1", "deepseek-v3",
    "deepseek-r1-distill-qwen-32b", "deepseek-r1-distill-qwen-7b",
    "step-1v-8k", "step-1.5v-mini",
    "Doubao-pro-32k", "Doubao-pro-256k", "Doubao-lite-128k",
    "Doubao-lite-32k", "Doubao-vision-lite-32k", "Doubao-vision-pro-32k",
    "Doubao-1.5-pro-vision-32k", "Doubao-1.5-lite-32k",
    "Doubao-1.5-pro-32k", "Doubao-1.5-thinking-pro", "Doubao-1.5-pro-256k",
}
```

Coze 支持的模型，涵盖多个厂商：
- **Moonshot (月之暗面)**: moonshot-v1 系列
- **百川智能**: Baichuan4
- **MiniMax**: abab6.5s-chat-pro
- **智谱 AI**: glm-4-0520
- **阿里云**: qwen-max
- **DeepSeek**: deepseek-r1、deepseek-v3、distill 系列
- **阶跃星辰**: step 系列
- **字节跳动**: Doubao 系列（豆包）

### `ChannelName` 变量

```go
var ChannelName = "coze"
```

频道标识名称。

## 4. 函数详解

无函数定义。

## 5. 关键逻辑分析

### 多厂商模型聚合
Coze 作为平台聚合了多家 AI 厂商的模型，用户可以通过统一的 Coze API 访问不同厂商的模型。

### 豆包模型系列
Doubao 系列包含多种变体：
- 标准版（pro/lite）
- 视觉版（vision）
- 1.5 版本（新一代）
- thinking-pro（推理增强版）

### 模型数量
共 24 个模型，是所有频道中模型种类最多的之一。

## 6. 关联文件

- `relay/channel/coze/adaptor.go` - 返回此模型列表
