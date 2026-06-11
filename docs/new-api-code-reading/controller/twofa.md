# twofa.go 代码阅读文档

## 1. 全局总结

该文件实现了两步验证（2FA/TOTP）的完整生命周期管理，包括初始化、启用、禁用、状态查询、备用码管理和登录验证。

## 2. 依赖关系

- `common` — TOTP 密钥生成、验证码验证、备用码生成
- `model` — 2FA 和备用码模型
- `gin-contrib/sessions` — 会话管理
- `gin-gonic/gin` — HTTP 框架

## 3. 类型定义

| 类型名 | 说明 |
|--------|------|
| `Setup2FARequest` | 设置 2FA 请求 |
| `Verify2FARequest` | 验证 2FA 请求 |
| `Setup2FAResponse` | 设置 2FA 响应（含密钥、二维码、备用码） |

## 4. 函数详解

### 设置流程
- `Setup2FA` — 初始化 2FA（生成密钥、二维码、备用码）
- `Enable2FA` — 启用 2FA（验证 TOTP 码后启用）

### 管理
- `Disable2FA` — 禁用 2FA（需验证 TOTP 或备用码）
- `Get2FAStatus` — 获取 2FA 状态（启用状态、锁定状态、备用码剩余数量）
- `RegenerateBackupCodes` — 重新生成备用码

### 登录验证
- `Verify2FALogin` — 登录时 2FA 验证（从 session 获取 pending 用户）

### 管理员操作
- `Admin2FAStats` — 获取 2FA 统计信息
- `AdminDisable2FA` — 管理员强制禁用用户 2FA

## 5. 关键逻辑分析

- 2FA 设置分两步：初始化（生成密钥）→ 启用（验证验证码）
- 登录时 2FA 验证使用 pending session 机制
- 支持 TOTP 验证码和备用码两种验证方式
- 管理员操作需要权限检查（`canManageTargetRole`）
- 已禁用的 2FA 记录在重新设置时会被删除

## 6. 关联文件

- `model/twofa.go` — 2FA 模型
- `controller/secure_verification.go` — 安全验证框架
- `controller/passkey.go` — Passkey 验证
