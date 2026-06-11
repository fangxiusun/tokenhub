# misc.go 代码阅读文档

## 1. 全局总结

该文件实现了系统状态查询、公告、关于页面、法律文档、邮件验证和密码重置等杂项接口。是前端初始化和用户认证流程的重要支撑。

## 2. 依赖关系

- `common` — 通用工具、邮件发送、验证码
- `constant` — Setup 标志
- `logger` — 日志
- `middleware` — HTTP 统计
- `model` — 用户模型、数据库
- `oauth` — 自定义 OAuth 提供商
- `setting` — 聊天配置
- `setting/console_setting` — 控制台设置
- `setting/operation_setting` — 运营设置
- `setting/system_setting` — 系统设置
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `PasswordResetRequest` | 密码重置请求 |

## 4. 函数详解

### 系统状态
- `TestStatus` — 数据库连接测试 + HTTP 统计
- `GetStatus` — 获取系统完整状态（版本、OAuth 配置、主题、面板开关等）
- `GetNotice` / `GetAbout` — 获取公告/关于内容
- `GetUserAgreement` / `GetPrivacyPolicy` — 获取法律文档
- `GetMidjourney` / `GetHomePageContent` — 获取 MJ 配置/首页内容

### 邮件和密码
- `SendEmailVerification` — 发送邮箱验证邮件（支持域名白名单和别名限制）
- `SendPasswordResetEmail` — 发送密码重置邮件
- `ResetPassword` — 执行密码重置（生成 12 位随机密码）

## 5. 关键逻辑分析

- `GetStatus` 是前端初始化的核心接口，返回 50+ 配置项
- 邮箱验证支持域名白名单（`EmailDomainRestrictionEnabled`）和别名限制（`EmailAliasRestrictionEnabled`）
- 密码重置使用验证码机制，生成 12 位随机密码
- 自定义 OAuth 提供商信息注入到 `GetStatus` 响应中

## 6. 关联文件

- `setting/system_setting/` — 系统配置
- `setting/operation_setting/` — 运营配置
- `setting/console_setting/` — 控制台配置
- `oauth/` — OAuth 提供商
