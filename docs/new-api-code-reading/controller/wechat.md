# wechat.go 代码阅读文档

## 1. 全局总结

该文件实现了微信 OAuth 登录和绑定功能。通过微信服务器代理获取用户信息。

## 2. 依赖关系

- `common` — 通用工具、微信配置
- `model` — 用户模型
- `gin-contrib/sessions` — 会话管理
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `wechatLoginResponse` | 微信登录响应 |

## 4. 函数详解

### `getWeChatIdByCode(code) (string, error)`
通过授权码获取微信用户 ID。调用微信服务器代理 API。

### `WeChatOAuth(c *gin.Context)`
微信 OAuth 登录/注册处理器。

### `WeChatBind(c *gin.Context)`
微信账号绑定处理器。

## 5. 关键逻辑分析

- 通过自建的微信代理服务器获取用户信息
- 使用 `common.WeChatServerAddress` 配置代理地址

## 6. 关联文件

- `controller/user.go` — `setupLogin` 登录设置
