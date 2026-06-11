# quota_setting.go 代码阅读文档

## 1. 全局总结

该文件定义额度设置，控制免费模型是否启用预消耗功能。

## 2. 依赖关系

- `github.com/QuantumNous/new-api/setting/config` — 全局配置管理器

## 3. 类型定义

| 结构体 | 字段 | 类型 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `QuotaSetting` | `EnableFreeModelPreConsume` | `bool` | `true` | 是否对免费模型启用预消耗 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `GetQuotaSetting` | `func GetQuotaSetting() *QuotaSetting` | 获取额度设置 |

## 5. 关键逻辑分析

- 预消耗功能默认启用，即使模型免费也会预先检查额度

## 6. 关联文件

- `relay/handler.go` — 使用预消耗设置
- `service/quota.go` — 额度服务
