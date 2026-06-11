# qwen.go 代码阅读文档

## 1. 全局总结

该文件定义 Qwen 模型的配置，主要管理同步图像模型列表。

## 2. 依赖关系

- `strings` — 字符串处理
- `github.com/QuantumNous/new-api/setting/config` — 全局配置管理器

## 3. 类型定义

| 结构体 | 字段 | 类型 | 说明 |
|--------|------|------|------|
| `QwenSettings` | `SyncImageModels` | `[]string` | 同步图像模型列表 |

默认包含：`z-image`、`qwen-image`、`wan2.6`、`wan2.7`、`qwen-image-edit` 等。

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `GetQwenSettings` | `func GetQwenSettings() *QwenSettings` | 获取 Qwen 配置 |
| `IsSyncImageModel` | `func IsSyncImageModel(model string) bool` | 检查模型是否为同步图像模型 |

## 5. 关键逻辑分析

- `IsSyncImageModel` 使用 `strings.Contains` 进行部分匹配
- 注意：bool 字段必须以 `enabled` 结尾才能在编辑界面生效

## 6. 关联文件

- `relay/qwen/` — Qwen 中继适配器
- `setting/model_setting/` — 其他模型配置
