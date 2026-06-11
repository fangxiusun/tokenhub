# codex_oauth.go 代码阅读文档

## 1. 全局总结

该文件实现了 Codex（OpenAI 代码模型）渠道的 OAuth 授权流程。支持通用 OAuth 流程和绑定到已有渠道的 OAuth 流程，通过 PKCE（Proof Key for Code Exchange）增强安全性。

## 2. 依赖关系

- `common` — JSON 序列化
- `constant` — 渠道类型常量
- `model` — 渠道模型
- `relay/channel/codex` — OAuthKey 数据结构
- `service` — OAuth 流程创建、授权码交换、JWT 解析
- `gin-contrib/sessions` — 会话管理
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `codexOAuthCompleteRequest` | 完成 OAuth 的请求体 |

## 4. 函数详解

### `StartCodexOAuth(c *gin.Context)`
启动通用 Codex OAuth 流程（无绑定渠道）。

### `StartCodexOAuthForChannel(c *gin.Context)`
启动绑定到指定渠道的 Codex OAuth 流程。

### `startCodexOAuthWithChannelID(c *gin.Context, channelID int)`
内部实现：创建 OAuth 授权流，将 state 和 verifier 存入 session。

### `CompleteCodexOAuth(c *gin.Context)`
完成通用 Codex OAuth 流程。

### `CompleteCodexOAuthForChannel(c *gin.Context)`
完成绑定到指定渠道的 Codex OAuth 流程。

### `completeCodexOAuthWithChannelID(c *gin.Context, channelID int)`
内部实现：解析授权输入 → 验证 state → 交换授权码 → 提取 account_id 和 email → 生成 OAuthKey → 更新渠道或返回密钥。

### `parseCodexAuthorizationInput(input string) (code, state string, err error)`
解析授权输入，支持多种格式：`code#state`、URL 查询参数、纯 code。

## 5. 关键逻辑分析

- 使用 PKCE 流程（verifier + state）增强安全性
- session 中存储 state、verifier 和创建时间
- 有 channelID 时直接更新渠道的 key 字段
- 无 channelID 时返回生成的密钥 JSON
- 支持代理配置（`channelProxy`）

## 6. 关联文件

- `controller/channel.go` — `RefreshCodexChannelCredential`
- `relay/channel/codex/` — OAuthKey 结构定义
- `service/codex.go` — OAuth 流程服务实现
