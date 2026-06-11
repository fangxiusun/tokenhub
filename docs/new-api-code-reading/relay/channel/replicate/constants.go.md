# constants.go 代码阅读文档

## 1. 全局总结
本文件定义了 Replicate 渠道的常量配置，包括渠道名称、默认模型和模型列表。

## 2. 依赖关系
无外部依赖。

## 3. 类型定义

### 常量定义
| 常量名 | 类型 | 值 | 说明 |
|--------|------|----|------|
| `ChannelName` | `string` | `"replicate"` | 渠道名称 |
| `ModelFlux11Pro` | `string` | `"black-forest-labs/flux-1.1-pro"` | 默认图像生成模型 |

### 变量定义
| 变量名 | 类型 | 说明 |
|--------|------|------|
| `ModelList` | `[]string` | 模型列表，当前仅包含 `ModelFlux11Pro` |

## 4. 函数详解
本文件无函数定义。

## 5. 关键逻辑分析
- Replicate 渠道当前仅支持一个模型：FLUX 1.1 Pro（由 Black Forest Labs 开发的图像生成模型）。
- `ModelFlux11Pro` 常量在 `adaptor.go` 中作为默认模型使用。

## 6. 关联文件
- `relay/channel/replicate/adaptor.go` — 使用这些常量
- `relay/channel/replicate/dto.go` — 数据结构定义
