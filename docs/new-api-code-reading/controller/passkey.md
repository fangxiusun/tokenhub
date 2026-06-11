# passkey.go 代码阅读文档

## 1. 全局总结

该文件实现了 WebAuthn/Passkey 的完整生命周期管理，包括注册、登录、验证和删除。Passkey 是一种无密码认证方式，基于 FIDO2 标准。

## 2. 依赖关系

- `common` — 通用工具
- `model` — Passkey 凭证模型
- `service/passkey` — WebAuthn 服务构建
- `setting/system_setting` — Passkey 配置
- `gin-contrib/sessions` — 会话管理
- `gin-gonic/gin` — HTTP 框架
- `go-webauthn/webauthn` — WebAuthn 库

## 3. 类型定义

无自定义类型定义（使用 `passkeysvc.WebAuthnUser` 等外部类型）。

## 4. 函数详解

### 注册流程
- `PasskeyRegisterBegin` — 开始注册（生成 challenge）
- `PasskeyRegisterFinish` — 完成注册（验证 response，保存凭证）

### 登录流程
- `PasskeyLoginBegin` — 开始登录（生成 discoverable login challenge）
- `PasskeyLoginFinish` — 完成登录（验证 response，更新凭证，设置登录状态）

### 验证流程（安全操作前验证）
- `PasskeyVerifyBegin` — 开始验证
- `PasskeyVerifyFinish` — 完成验证（设置 PasskeyReady 状态）

### 管理
- `PasskeyDelete` — 删除 Passkey（需安全验证）
- `PasskeyStatus` — 查询 Passkey 状态
- `AdminResetPasskey` — 管理员重置用户 Passkey

### 辅助函数
- `getSessionUser` — 从 session 获取用户
- `requirePasskeyRegistrationVerification` — 注册前 2FA 验证
- `requirePasskeyDeleteVerification` — 删除前安全验证
- `requireSecureVerificationMethod` — 验证安全操作的 session 状态

## 5. 关键逻辑分析

- 注册时排除已有凭证（避免重复注册）
- 登录使用 Discoverable Login（无需输入用户名）
- 验证流程用于保护敏感操作（如查看渠道密钥）
- 2FA 启用时，注册/删除 Passkey 需先通过 2FA 验证
- session 中存储 registration/login/verify 三种 session data
- `PasskeyReadySessionKey` 标记 Passkey 验证完成

## 6. 关联文件

- `model/passkey.go` — Passkey 凭证模型
- `service/passkey/` — WebAuthn 服务
- `controller/secure_verification.go` — 安全验证框架
