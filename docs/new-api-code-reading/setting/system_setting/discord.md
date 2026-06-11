# discord.go 代码阅读文档

## 1. 全局总结

该文件定义 Discord OAuth 登录的配置。

## 2. 依赖关系

- `github.com/QuantumNous/new-api/setting/config` — 全局配置管理器

## 3. 类型定义

| 结构体 | 字段 | 类型 | 说明 |
|--------|------|------|------|
| `DiscordSettings` | `Enabled` | `bool` | 是否启用 Discord 登录 |
| | `ClientId` | `string` | Discord Client ID |
| | `ClientSecret` | `string` | Discord Client Secret |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `GetDiscordSettings` | `func GetDiscordSettings() *DiscordSettings` | 获取 Discord 配置 |

## 5. 关键逻辑分析

- 默认配置为空，需管理员手动配置

## 6. 关联文件

- `oauth/discord.go` — Discord OAuth 实现
- `controller/auth.go` — 认证接口
