# compact_suffix.go 代码阅读文档

## 1. 全局总结

该文件定义紧凑模型（Compact Model）后缀常量和工具函数，用于识别和处理 OpenAI 压缩模型。

## 2. 依赖关系

- `strings` — 字符串处理

## 3. 类型定义

| 常量名 | 值 | 说明 |
|--------|------|------|
| `CompactModelSuffix` | `"-openai-compact"` | 紧凑模型后缀 |
| `CompactWildcardModelKey` | `"*-openai-compact"` | 紧凑模型通配符键 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `WithCompactModelSuffix` | `func WithCompactModelSuffix(modelName string) string` | 为模型名添加紧凑后缀（如已存在则不重复添加） |

## 5. 关键逻辑分析

- `WithCompactModelSuffix` 检查模型名是否已包含后缀，避免重复添加

## 6. 关联文件

- `setting/ratio_setting/model_ratio.go` — 使用紧凑模型后缀
- `relay/handler.go` — 识别紧凑模型
