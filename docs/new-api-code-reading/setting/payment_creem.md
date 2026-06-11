# payment_creem.go 代码阅读文档

## 1. 全局总结

该文件定义 Creem 支付网关的配置变量。

## 2. 依赖关系

无外部依赖。

## 3. 类型定义

| 变量名 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `CreemApiKey` | `string` | `""` | Creem API 密钥 |
| `CreemProducts` | `string` | `"[]"` | Creem 产品配置（JSON 数组） |
| `CreemTestMode` | `bool` | `false` | 是否启用测试模式 |
| `CreemWebhookSecret` | `string` | `""` | Webhook 签名密钥 |

## 4. 函数详解

无函数定义，仅变量声明。

## 5. 关键逻辑分析

所有变量为包级全局变量，通过环境变量或管理界面配置。

## 6. 关联文件

- `controller/option.go` — 管理界面配置接口
- `service/creem.go` — Creem 支付服务实现
