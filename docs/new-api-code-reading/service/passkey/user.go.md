# user.go 代码阅读文档

## 1. 全局总结

本文件实现了 `webauthn.User` 接口，将应用的 `model.User` 和 `model.PasskeyCredential` 适配为 WebAuthn 库所需的用户接口。这是连接业务用户模型与 WebAuthn 标准协议的桥梁层。

## 2. 依赖关系

| 包 | 用途 |
|---|---|
| `fmt` | 格式化用户名 |
| `strconv` | 整数转字符串 |
| `strings` | 字符串修剪 |
| `github.com/QuantumNous/new-api/model` | `User` 和 `PasskeyCredential` 数据模型 |
| `github.com/go-webauthn/webauthn/webauthn` | WebAuthn `User` 接口和 `Credential` 类型 |

## 3. 类型定义

### `WebAuthnUser`

```go
type WebAuthnUser struct {
    user       *model.User
    credential *model.PasskeyCredential
}
```

适配器结构体，组合了业务用户模型和 Passkey 凭证模型。

## 4. 函数详解

### `NewWebAuthnUser(user, credential) *WebAuthnUser`

**签名**: `func NewWebAuthnUser(user *model.User, credential *model.PasskeyCredential) *WebAuthnUser`

**功能**: 构造函数，创建 `WebAuthnUser` 实例。

### `WebAuthnID() []byte`

**签名**: `func (u *WebAuthnUser) WebAuthnID() []byte`

**功能**: 返回用户 ID 的字节表示（WebAuthn 接口方法）。

**逻辑**: 将 `user.Id` 转换为字符串再转为 `[]byte`。空值安全（`u` 或 `u.user` 为 `nil` 时返回 `nil`）。

### `WebAuthnName() string`

**签名**: `func (u *WebAuthnUser) WebAuthnName() string`

**功能**: 返回用户名称（WebAuthn 接口方法）。

**逻辑**:
1. 优先使用 `user.Username`（修剪空白）
2. 若为空，返回 `"user-{id}"` 格式的回退名称

### `WebAuthnDisplayName() string`

**签名**: `func (u *WebAuthnUser) WebAuthnDisplayName() string`

**功能**: 返回用户显示名称（WebAuthn 接口方法）。

**逻辑**:
1. 优先使用 `user.DisplayName`（修剪空白）
2. 若为空，回退到 `WebAuthnName()`

### `WebAuthnCredentials() []webauthn.Credential`

**签名**: `func (u *WebAuthnUser) WebAuthnCredentials() []webauthn.Credential`

**功能**: 返回用户的 WebAuthn 凭证列表（WebAuthn 接口方法）。

**逻辑**: 调用 `credential.ToWebAuthnCredential()` 将模型凭证转换为 WebAuthn 库格式，包装为单元素切片返回。

### `ModelUser() *model.User`

**签名**: `func (u *WebAuthnUser) ModelUser() *model.User`

**功能**: 提取内部 `model.User` 的访问方法。

### `PasskeyCredential() *model.PasskeyCredential`

**签名**: `func (u *WebAuthnUser) PasskeyCredential() *model.PasskeyCredential`

**功能**: 提取内部 `model.PasskeyCredential` 的访问方法。

## 5. 关键逻辑分析

1. **适配器模式**: `WebAuthnUser` 是经典的适配器实现，将业务模型适配为 WebAuthn 库要求的接口，解耦了数据模型与第三方库。

2. **空值安全**: 所有方法都包含 `u == nil || u.user == nil` 的防御性检查，避免空指针 panic。

3. **ID 格式约定**: 用户 ID 使用 `strconv.Itoa` 转为字符串再转字节，这要求 `model.User.Id` 为整数类型，且在注册和认证流程中保持一致。

4. **凭证封装**: `WebAuthnCredentials` 返回单元素切片，表明当前设计每个 `WebAuthnUser` 实例对应单个凭证。实际使用中可能通过多次查询构建不同凭证的实例。

## 6. 关联文件

- `new-api/model/user.go` — `User` 数据模型定义
- `new-api/model/passkey.go`（推测） — `PasskeyCredential` 数据模型定义
- `new-api/service/passkey/service.go` — WebAuthn 实例构建
- `new-api/service/passkey/session.go` — 会话数据管理
