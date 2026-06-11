# grok.go 代码阅读文档

## 1. 全局总结

该文件定义 Grok 模型的配置，包括违规扣款开关和扣款金额。

## 2. 依赖关系

- `github.com/QuantumNous/new-api/setting/config` — 全局配置管理器

## 3. 类型定义

| 结构体 | 字段 | 类型 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `GrokSettings` | `ViolationDeductionEnabled` | `bool` | `true` | 是否启用违规扣款 |
| | `ViolationDeductionAmount` | `float64` | `0.05` | 违规扣款金额 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `GetGrokSettings` | `func GetGrokSettings() *GrokSettings` | 获取 Grok 配置 |

## 5. 关键逻辑分析

- 违规扣款默认启用，每次扣款 0.05 单位
- 通过 `init()` 注册到全局配置管理器

## 6. 关联文件

- `relay/grok/` — Grok 中继适配器
- `setting/model_setting/` — 其他模型配置
