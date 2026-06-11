# constants.go 代码阅读文档

## 1. 全局总结

该文件定义了豆包视频生成频道的常量，包括支持的模型列表、频道名称、视频输入折扣比率映射，以及获取折扣比率的辅助函数。视频输入折扣比率用于在用户提交视频输入时自动调整计费倍率。

## 2. 依赖关系

无外部依赖。

## 3. 类型定义

### 变量

| 变量名 | 类型 | 说明 |
|--------|------|------|
| `ModelList` | `[]string` | 豆包支持的视频生成模型列表 |
| `ChannelName` | `string` | 频道标识，值为 `"doubao-video"` |
| `videoInputRatioMap` | `map[string]float64` | 视频输入折扣比率映射（含视频单价 / 不含视频单价） |

### 模型列表

| 模型名 | 说明 |
|--------|------|
| `doubao-seedance-1-0-pro-250528` | Seedance 1.0 专业版 |
| `doubao-seedance-1-0-lite-t2v` | Seedance 1.0 轻量版（文生视频） |
| `doubao-seedance-1-0-lite-i2v` | Seedance 1.0 轻量版（图生视频） |
| `doubao-seedance-1-5-pro-251215` | Seedance 1.5 专业版 |
| `doubao-seedance-2-0-260128` | Seedance 2.0 |
| `doubao-seedance-2-0-fast-260128` | Seedance 2.0 快速版 |

### 折扣比率

| 模型名 | 比率 | 计算说明 |
|--------|------|----------|
| `doubao-seedance-2-0-260128` | ~0.6087 | 28.0 / 46.0 |
| `doubao-seedance-2-0-fast-260128` | ~0.5946 | 22.0 / 37.0 |

## 4. 函数详解

| 函数签名 | 说明 |
|----------|------|
| `GetVideoInputRatio(modelName string) (float64, bool)` | 根据模型名获取视频输入折扣比率，不存在则返回 false |

## 5. 关键逻辑分析

视频输入折扣比率的设计理念：管理员将 ModelRatio 设置为"不含视频"的较高费率，系统在检测到视频输入时自动乘以此折扣（小于 1.0），从而实现有视频输入时的差异化计费。仅 Seedance 2.0 系列模型支持此折扣。

## 6. 关联文件

- `relay/channel/task/doubao/adaptor.go` — 调用 `GetVideoInputRatio` 进行计费估算
