# constants.go 代码阅读文档

## 1. 全局总结
本文件定义了 Perplexity AI 渠道的常量配置，包括支持的模型列表和渠道名称。

## 2. 依赖关系
无外部依赖。

## 3. 类型定义

### 变量定义
| 变量名 | 类型 | 说明 |
|--------|------|------|
| `ModelList` | `[]string` | Perplexity 支持的模型列表 |
| `ChannelName` | `string` | 渠道名称，值为 `"perplexity"` |

### 模型列表
| 模型名称 | 说明 |
|----------|------|
| `llama-3-sonar-small-32k-chat` | Llama 3 Sonar 小型聊天模型 |
| `llama-3-sonar-small-32k-online` | Llama 3 Sonar 小型在线模型 |
| `llama-3-sonar-large-32k-chat` | Llama 3 Sonar 大型聊天模型 |
| `llama-3-sonar-large-32k-online` | Llama 3 Sonar 大型在线模型 |
| `llama-3-8b-instruct` | Llama 3 8B 指令微调模型 |
| `llama-3-70b-instruct` | Llama 3 70B 指令微调模型 |
| `mixtral-8x7b-instruct` | Mixtral 8x7B 指令微调模型 |
| `sonar` | Sonar 基础模型 |
| `sonar-pro` | Sonar 专业版模型 |
| `sonar-reasoning` | Sonar 推理模型 |

## 4. 函数详解
本文件无函数定义。

## 5. 关键逻辑分析
- Perplexity 提供了多个系列的模型：Sonar（搜索增强）、Llama 3（开源微调）、Mixtral（MoE 架构）。
- `sonar` 系列模型支持在线搜索功能（`-online` 后缀）。

## 6. 关联文件
- `relay/channel/perplexity/adaptor.go` — 适配器实现
- `relay/channel/perplexity/relay-perplexity.go` — 请求转换逻辑
