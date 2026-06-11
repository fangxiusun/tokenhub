# constants.go 代码阅读文档

## 1. 全局总结
该文件定义了 Ollama 渠道的常量，包括支持的模型列表和渠道名称标识。

## 2. 依赖关系
无外部依赖。

## 3. 类型定义

### `ModelList` (var []string)
- **作用**：定义 Ollama 渠道支持的模型列表
- **值**：仅包含一个默认模型 `llama3-7b`

### `ChannelName` (var string)
- **作用**：渠道名称标识
- **值**：`"ollama"`

## 4. 函数详解
无函数定义。

## 5. 关键逻辑分析
- Ollama 是本地部署的 LLM 服务，模型列表仅为默认值，实际可用模型取决于用户本地部署情况
- 与云端 API 不同，Ollama 支持动态拉取和删除模型

## 6. 关联文件
- `relay/channel/ollama/adaptor.go` — 通过 `GetModelList()` 和 `GetChannelName()` 引用这些常量
- `relay/channel/ollama/relay-ollama.go` — 通过 `FetchOllamaModels` 动态获取实际可用模型
