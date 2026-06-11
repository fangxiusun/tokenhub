# user.go 代码阅读文档

## 1. 全局总结

该文件是用户管理的核心控制器，实现了登录、注册、用户信息管理、密码修改、邮箱验证、OAuth 绑定、管理员用户管理等功能。文件较大（约 1300 行），包含大量用户相关的 HTTP 处理器。

## 2. 依赖关系

- `common` — 通用工具、密码哈希
- `dto` — 数据传输对象
- `i18n` — 国际化
- `logger` — 日志
- `model` — 用户模型
- `service` — 业务逻辑
- `setting` — 系统设置
- `setting/operation_setting` — 运营设置
- `constant` — 常量
- `gin-contrib/sessions` — 会话管理
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `LoginRequest` | 登录请求 |

## 4. 函数详解

### 认证
- `Login` — 密码登录（支持 2FA 检查）
- `setupLogin` — 设置登录 session 和 Cookie

### 用户管理
- 大量用户 CRUD 接口（注册、信息修改、密码修改等）

## 5. 关键逻辑分析

- 登录时检查 2FA 启用状态，已启用则设置 pending session
- `setupLogin` 是所有登录方式的统一入口
- 支持多种 OAuth 提供商的绑定/解绑

## 6. 关联文件

- `model/user.go` — 用户模型
- `controller/twofa.go` — 2FA 验证
- `controller/github.go/discord.go/oidc.go/linuxdo.go` — OAuth 登录
