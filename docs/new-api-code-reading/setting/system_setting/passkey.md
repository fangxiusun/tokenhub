# passkey.go 代码阅读文档

## 1. 全局总结

该文件定义 WebAuthn/Passkey 登录的配置，支持自动从服务器地址提取 RPID。

## 2. 依赖关系

- `net/url` — URL 解析
- `strings` — 字符串处理
- `github.com/QuantumNous/new-api/common` — 系统名称
- `github.com/QuantumNous/new-api/setting/config` — 全局配置管理器

## 3. 类型定义

| 结构体 | 字段 | 类型 | 说明 |
|--------|------|------|------|
| `PasskeySettings` | `Enabled` | `bool` | 是否启用 Passkey 登录 |
| | `RPDisplayName` | `string` | RP 显示名称 |
| | `RPID` | `string` | RP ID（域名） |
| | `Origins` | `string` | 允许的来源列表 |
| | `AllowInsecureOrigin` | `bool` | 是否允许不安全来源 |
| | `UserVerification` | `string` | 用户验证方式 |
| | `AttachmentPreference` | `string` | 凭据附件偏好 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `GetPasskeySettings` | `func GetPasskeySettings() *PasskeySettings` | 获取 Passkey 配置 |

## 5. 关键逻辑分析

- 如果 RPID 为空，自动从 `ServerAddress` 提取域名
- 如果 Origins 为空，使用 `ServerAddress` 作为默认来源
- 默认用户验证方式为 "preferred"

## 6. 关联文件

- `oauth/passkey.go` — Passkey 实现
- `controller/auth.go` — 认证接口
- `model/passkey.go` — Passkey 模型
