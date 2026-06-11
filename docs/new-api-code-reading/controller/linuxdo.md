# linuxdo.go 代码阅读文档

## 1. 全局总结

该文件实现了 LinuxDO（Linux.do 社区）OAuth 登录和绑定功能。支持通过 LinuxDO 账号登录/注册，包含信任等级（Trust Level）验证。

## 2. 依赖关系

- `common` — 通用工具、LinuxDO OAuth 配置
- `model` — 用户模型
- `gin-contrib/sessions` — 会话管理
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `LinuxdoUser` | LinuxDO 用户信息（含 trust_level、silenced 等） |

## 4. 函数详解

### `getLinuxdoUserInfoByCode(code string, c *gin.Context) (*LinuxdoUser, error)`
通过授权码获取 LinuxDO 用户信息。使用 Basic Auth 认证获取 Token，再获取用户信息。

### `LinuxdoOAuth(c *gin.Context)`
LinuxDO OAuth 登录/注册处理器。验证 state → 检查错误码 → 检查登录状态 → 获取用户信息 → 信任等级验证 → 创建或登录用户。

### `LinuxDoBind(c *gin.Context)`
LinuxDO 账号绑定处理器。

## 5. 关键逻辑分析

- 使用 Basic Auth（`base64(clientId:clientSecret)`）获取 Token
- 支持 `LINUX_DO_TOKEN_ENDPOINT` 和 `LINUX_DO_USER_ENDPOINT` 环境变量自定义端点
- 新用户注册需满足最低信任等级（`common.LinuxDOMinimumTrustLevel`）
- redirect_uri 动态构建（根据请求 scheme 和 host）
- 支持邀请码机制

## 6. 关联文件

- `controller/user.go` — `setupLogin` 登录设置
- `model/user.go` — 用户模型
