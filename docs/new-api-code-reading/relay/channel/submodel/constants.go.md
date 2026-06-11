# constants.go 代码阅读文档

## 1. 全局总结
本文件定义了 SubModel 渠道的常量配置，包括支持的模型列表和渠道名称。模型列表包含多个开源大模型。

## 2. 依赖关系
无外部依赖。

## 3. 类型定义

### 变量定义
| 变量名 | 类型 | 说明 |
|--------|------|------|
| `ModelList` | `[]string` | SubModel 支持的模型列表 |

### 常量定义
| 常量名 | 类型 | 值 | 说明 |
|--------|------|----|------|
| `ChannelName` | `string` | `"submodel"` | 渠道名称 |

### 模型列表
| 模型 | 说明 |
|------|------|
| `NousResearch/Hermes-4-405B-FP8` | Nous Research 的 Hermes 4 405B 模型（FP8 量化） |
| `Qwen/Qwen3-235B-A22B-Thinking-2507` | Qwen3 235B 推理模型（2025-07 版本） |
| `Qwen/Qwen3-Coder-480B-A35B-Instruct-FP8` | Qwen3 480B 编码模型（FP8 量化） |
| `Qwen/Qwen3-235B-A22B-Instruct-2507` | Qwen3 235B 指令微调模型 |
| `zai-org/GLM-4.5-FP8` | 智谱 GLM-4.5 模型（FP8 量化） |
| `openai/gpt-oss-120b` | OpenAI 开源 120B 模型 |
| `deepseek-ai/DeepSeek-R1-0528` | DeepSeek R1 推理模型（2024-05-28 版本） |
| `deepseek-ai/DeepSeek-R1` | DeepSeek R1 推理模型 |
| `deepseek-ai/DeepSeek-V3-0324` | DeepSeek V3 模型（2024-03-24 版本） |
| `deepseek-ai/DeepSeek-V3.1` | DeepSeek V3.1 模型 |

## 4. 函数详解
本文件无函数定义。

## 5. 关键逻辑分析
- SubModel 提供了多个大参数量的开源模型，包括 Qwen3、DeepSeek R1/V3、GLM-4.5 等。
- 多个模型使用 FP8 量化，表明平台注重推理效率。
- 包含推理模型（R1、Thinking）和指令微调模型（Instruct），覆盖不同使用场景。
- `ChannelName` 使用 `const` 而非 `var`，与其他渠道的风格略有不同。

## 6. 关联文件
- `relay/channel/submodel/adaptor.go` — 适配器实现
