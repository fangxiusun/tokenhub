# legal.go 代码阅读文档

## 1. 全局总结

该文件定义法律相关设置，包括用户协议和隐私政策内容。

## 2. 依赖关系

- `github.com/QuantumNous/new-api/setting/config` — 全局配置管理器

## 3. 类型定义

| 结构体 | 字段 | 类型 | 说明 |
|--------|------|------|------|
| `LegalSettings` | `UserAgreement` | `string` | 用户协议内容 |
| | `PrivacyPolicy` | `string` | 隐私政策内容 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `GetLegalSettings` | `func GetLegalSettings() *LegalSettings` | 获取法律设置 |

## 5. 关键逻辑分析

- 默认内容为空，需管理员在后台填写

## 6. 关联文件

- `controller/legal.go` — 法律内容接口
- `web/default/src/` — 前端法律页面
