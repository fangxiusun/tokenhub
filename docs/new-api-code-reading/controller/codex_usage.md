# codex_usage.go 代码阅读文档

## 1. 全局总结

该文件实现了 Codex 渠道的用量查询功能。通过调用上游 Wham API 获取 Codex 渠道的使用量信息，支持 Token 自动刷新。

## 2. 依赖关系

- `common` — JSON 序列化、日志
- `constant` — 渠道类型
- `model` — 渠道模型
- `relay/channel/codex` — OAuthKey 解析
- `service` — 代理客户端、Token 刷新、用量获取
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

无自定义类型定义。

## 4. 函数详解

### `GetCodexChannelUsage(c *gin.Context)`
获取指定 Codex 渠道的用量信息。流程：验证渠道类型 → 解析 OAuthKey → 调用 Wham API → 若 401/403 则自动刷新 Token → 返回用量数据。

## 5. 关键逻辑分析

- 仅支持单密钥 Codex 渠道（不支持多密钥）
- Token 自动刷新：当上游返回 401/403 时，使用 refresh_token 获取新的 access_token
- 刷新成功后更新数据库中的渠道 key 和缓存
- 上游响应直接透传给前端（不解析具体格式）

## 6. 关联文件

- `controller/codex_oauth.go` — Codex OAuth 流程
- `service/codex.go` — `FetchCodexWhamUsage`、`RefreshCodexOAuthTokenWithProxy`
