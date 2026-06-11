# discord.go 代码阅读文档

## 1. 全局总结

该文件实现了 Discord OAuth 登录和绑定功能。支持通过 Discord 账号登录/注册新用户，以及将已有账号与 Discord 账号绑定。

## 2. 依赖关系

- `common` — 通用工具
- `model` — 用户模型
- `setting/system_setting` — Discord 配置（ClientId、ClientSecret、启用状态）
- `gin-contrib/sessions` — 会话管理
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `DiscordResponse` | Discord OAuth Token 响应 |
| `DiscordUser` | Discord 用户信息 |

## 4. 函数详解

### `getDiscordUserInfoByCode(code string) (*DiscordUser, error)`
通过授权码获取 Discord 用户信息。流程：交换 Token → 获取用户信息。

### `DiscordOAuth(c *gin.Context)`
Discord OAuth 登录/注册处理器。验证 state → 检查是否已登录（已登录则转绑定）→ 获取用户信息 → 创建或登录用户。

### `DiscordBind(c *gin.Context)`
Discord 账号绑定处理器。检查 Discord ID 是否已被绑定 → 更新用户 Discord ID。

## 5. 关键逻辑分析

- state 参数用于 CSRF 防护，与 session 中存储的值比对
- 已登录用户访问 OAuth 回调时自动转为绑定流程
- 新用户注册时用户名格式为 `discord_{maxUserId+1}`
- 支持显示名称（global_name）回退为 "Discord User"

## 6. 关联文件

- `controller/user.go` — `setupLogin` 登录设置
- `model/user.go` — 用户模型
- `setting/system_setting/` — Discord 配置
