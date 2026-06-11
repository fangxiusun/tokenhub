# constant.go 代码阅读文档

## 1. 全局总结

本文件定义了 Cloudflare Workers AI 频道的常量，包括支持的模型列表和频道名称标识。模型列表涵盖了 Cloudflare 平台上部署的各种开源 LLM 模型。

## 2. 依赖关系

无外部依赖，仅使用 Go 内置类型。

## 3. 类型定义

### `ModelList` 变量

```go
var ModelList = []string{...}
```

字符串切片，包含 Cloudflare Workers AI 支持的所有模型 ID。这些模型主要来自以下系列：
- **Meta LLaMA**: `@cf/meta/llama-3.1-8b-instruct`、`@cf/meta/llama-2-7b-chat-fp16` 等
- **Mistral**: `@cf/mistral/mistral-7b-instruct-v0.1`、`@hf/mistralai/mistral-7b-instruct-v0.2` 等
- **DeepSeek**: `@cf/deepseek-ai/deepseek-math-7b-instruct` 等
- **Google Gemma**: `@cf/google/gemma-2b-it-lora`、`@hf/google/gemma-7b-it` 等
- **Qwen**: `@cf/qwen/qwen1.5-14b-chat-awq` 等
- **其他**: Falcon、Phi、TinyLlama、OpenChat 等

模型 ID 格式为 `@cf/{org}/{model}` 或 `@hf/{org}/{model}`，分别表示 Cloudflare 原生部署和 Hugging Face 部署的模型。

### `ChannelName` 常量

```go
var ChannelName = "cloudflare"
```

频道标识名称，用于在系统中注册和识别该频道类型。

## 4. 函数详解

无函数定义。

## 5. 关键逻辑分析

### 模型命名约定
- `@cf/` 前缀：Cloudflare 原生托管的模型，运行在 Cloudflare Workers AI 基础设施上
- `@hf/` 前缀：通过 Cloudflare 运行的 Hugging Face 模型
- 模型名包含组织名、模型名和变体信息（如量化格式 `awq`、`int8`、`fp16`）

### 模型数量
共定义了 33 个模型，覆盖了 7B 到 14B 参数量级的多种任务类型（聊天、代码、数学、多语言等）。

## 6. 关联文件

- `relay/channel/cloudflare/adaptor.go` - 通过 `GetModelList()` 方法返回此模型列表
- `relay/channel/` - 频道注册机制使用 `ChannelName` 和 `ModelList` 进行频道识别
