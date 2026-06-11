# github.go 代码阅读文档

## 1. 全局总结

该文件实现了 GitHub OAuth 登录和绑定功能。支持通过 GitHub 账号登录/注册新用户，以及将已有账号与 GitHub 账号绑定。

## 2. 依赖关系

- `common` — 通用工具、GitHub OAuth 配置
- `model` — 用户模型
- `gin-contrib/sessions` — 会话管理
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `GitHubOAuthResponse` | GitHub OAuth Token 响应 |
| `GitHubUser` | GitHub 用户信息 |

## 4. 函数详解

### `getGitHubUserInfoByCode(code string) (*GitHubUser, error)`
通过授权码获取 GitHub 用户信息。流程：交换 Token → 获取用户信息。

### `GitHubOAuth(c *gin.Context)`
GitHub OAuth 登录/注册处理器。验证 state → 检查是否已登录 → 获取用户信息 → 创建或登录用户。支持邀请码（aff code）。

### `GitHubBind(c *gin.Context)`
GitHub 账号绑定处理器。检查 GitHub Login 是否已被绑定 → 更新用户 GitHub ID。

## 5. 关键逻辑分析

- 使用 `common.GitHubClientId` 和 `common.GitHubClientSecret` 配置
- 新用户注册时用户名格式为 `github_{maxUserId+1}`
- 支持邀请码机制（session 中的 `aff` 字段）
- 用户 ID 从 session 获取而非 context（注释标注为 critical bug fix）

## 6. 关联文件

- `controller/user.go` — `setupLogin` 登录设置
- `model/user.go` — 用户模型
