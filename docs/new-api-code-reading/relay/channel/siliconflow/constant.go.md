# constant.go 代码阅读文档

## 1. 全局总结
本文件定义了 SiliconFlow 渠道的常量配置，包括支持的模型列表和渠道名称。模型列表涵盖聊天、图像生成、嵌入和 Rerank 等多种模型。

## 2. 依赖关系
无外部依赖。

## 3. 类型定义

### 变量定义
| 变量名 | 类型 | 说明 |
|--------|------|------|
| `ModelList` | `[]string` | SiliconFlow 支持的模型列表 |
| `ChannelName` | `string` | 渠道名称，值为 `"siliconflow"` |

### 模型分类
| 类别 | 模型 |
|------|------|
| **聊天模型** | GLM-4-9B-Chat、DeepSeek-LLM-67B-Chat、Qwen1.5/Qwen2 系列、Yi-1.5 系列、ChatGLM3-6B、InternLM2.5 系列 |
| **图像生成** | InstantID、SDXL-Lightning、FLUX.1-schnell |
| **嵌入模型** | BAAI/bge-large-en-v1.5、BAAI/bge-large-zh-v1.5、BAAI/bge-m3、netease-youdao/bce-embedding-base_v1 |
| **Rerank 模型** | netease-youdao/bce-reranker-base_v1、BAAI/bge-reranker-v2-m3 |
| **Pro 版本** | Qwen2-7B/1.5B、GLM-4-9B、ChatGLM3-6B、Yi-1.5-9B/6B、Gemma-2-9B、InternLM2.5-7B、Llama-3-8B、Mistral-7B |
| **语音模型** | FunAudioLLM/SenseVoiceSmall |
| **数学模型** | Qwen2-Math-72B-Instruct |

## 4. 函数详解
本文件无函数定义。

## 5. 关键逻辑分析
- SiliconFlow 提供了丰富的模型生态，覆盖聊天、图像、嵌入、Rerank 等多种 AI 能力。
- Pro 版本模型是经过优化的推理版本，通常性能更好。
- 部分模型（如 Stable Diffusion 系列）被注释掉，可能是暂时下线或不再支持。

## 6. 关联文件
- `relay/channel/siliconflow/adaptor.go` — 适配器实现
- `relay/channel/siliconflow/dto.go` — 数据结构定义
- `relay/channel/siliconflow/relay-siliconflow.go` — Rerank 响应处理
