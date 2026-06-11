# passkey.go 代码阅读文档

## 1. 全局总结
该文件实现了 WebAuthn/Passkey 凭据的管理功能，包括凭据的存储、查询、更新和删除。Passkey 是一种基于公钥密码学的无密码认证机制，该文件提供了与 WebAuthn 库的集成，将 WebAuthn 凭据转换为数据库存储格式，并支持凭据的 CRUD 操作。

## 2. 依赖关系
- **encoding/base64**: 用于凭据数据的 Base64 编码和解码。
- **encoding/json**: 用于传输列表（Transports）的 JSON 序列化/反序列化。
- **errors**: 用于错误处理和错误类型判断。
- **fmt**: 用于格式化错误信息和日志。
- **strings**: 用于字符串处理（修剪空白字符）。
- **time**: 用于时间戳处理。
- **github.com/QuantumNous/new-api/common**: 提供系统日志功能。
- **github.com/go-webauthn/webauthn/protocol**: WebAuthn 协议类型定义。
- **github.com/go-webauthn/webauthn/webauthn**: WebAuthn 核心库。
- **gorm.io/gorm**: 用于数据库 ORM 操作。

## 3. 类型定义
### 错误变量
- **ErrPasskeyNotFound**: 凭据未找到的内部错误。
- **ErrFriendlyPasskeyNotFound**: 用户友好的凭据未找到错误。

### PasskeyCredential 结构体
表示 Passkey 凭据的数据模型，包含以下字段：
- `ID`: 凭据 ID（主键）
- `UserID`: 用户 ID（唯一索引）
- `CredentialID`: WebAuthn 凭据 ID（Base64 编码，唯一索引）
- `PublicKey`: 公钥（Base64 编码）
- `AttestationType`: 证明类型
- `AAGUID`: 认证器全局唯一标识符（Base64 编码）
- `SignCount`: 签名计数器
- `CloneWarning`: 克隆警告标志
- `UserPresent`: 用户存在标志
- `UserVerified`: 用户验证标志
- `BackupEligible`: 备份资格标志
- `BackupState`: 备份状态标志
- `Transports`: 传输方式列表（JSON 格式）
- `Attachment`: 认证器附件类型
- `LastUsedAt`: 最后使用时间
- `CreatedAt`: 创建时间
- `UpdatedAt`: 更新时间
- `DeletedAt`: 软删除时间戳

## 4. 函数详解
### 方法函数
- **(p *PasskeyCredential) TransportList() []protocol.AuthenticatorTransport**
  - 将 Transports 字符串解析为传输类型列表。

- **(p *PasskeyCredential) SetTransports(list []protocol.AuthenticatorTransport)**
  - 将传输类型列表序列化为 JSON 字符串并存储。

- **(p *PasskeyCredential) ToWebAuthnCredential() webauthn.Credential**
  - 将数据库存储的凭据转换为 WebAuthn 库需要的格式。

- **(p *PasskeyCredential) ApplyValidatedCredential(credential *webauthn.Credential)**
  - 用验证后的 WebAuthn 凭据更新当前凭据的所有字段。

### 工厂函数
- **NewPasskeyCredentialFromWebAuthn(userID int, credential *webauthn.Credential) *PasskeyCredential**
  - 从 WebAuthn 凭据创建新的 PasskeyCredential 实例。

### 查询函数
- **GetPasskeyByUserID(userID int) (*PasskeyCredential, error)**
  - 根据用户 ID 查询 Passkey 凭据。
  - 未找到记录返回 ErrPasskeyNotFound，数据库错误返回 ErrFriendlyPasskeyNotFound。

- **GetPasskeyByCredentialID(credentialID []byte) (*PasskeyCredential, error)**
  - 根据凭据 ID 查询 Passkey 凭据。
  - 将字节切片转换为 Base64 字符串进行查询。

### 更新函数
- **UpsertPasskeyCredential(credential *PasskeyCredential) error**
  - 插入或更新 Passkey 凭据：在事务中先硬删除用户现有凭据，再创建新凭据。
  - 使用硬删除（Unscoped）避免唯一索引冲突。

### 删除函数
- **DeletePasskeyByUserID(userID int) error**
  - 根据用户 ID 硬删除 Passkey 凭据。

## 5. 关键逻辑分析
- **数据格式转换**: 在数据库存储格式（Base64 字符串）和 WebAuthn 库格式（字节切片）之间进行转换。
- **硬删除策略**: 使用 `Unscoped().Delete()` 进行硬删除，而不是软删除，以避免唯一索引冲突（因为软删除会保留记录，可能导致新凭据无法插入）。
- **事务操作**: `UpsertPasskeyCredential` 使用数据库事务确保删除和创建的原子性。
- **错误处理**: 区分"未找到记录"（正常情况）和"数据库错误"（异常情况），对后者记录日志。
- **用户友好错误**: 对外返回用户友好的错误信息，内部保留详细错误信息。

## 6. 关联文件
- **controller/passkey.go**: 可能包含处理 Passkey 认证 HTTP 请求的控制器。
- **service/passkey.go**: 可能包含 Passkey 业务逻辑服务层。
- **router/passkey.go**: 可能包含 Passkey 相关的路由定义。
- **middleware/auth.go**: 可能包含 Passkey 认证中间件。
- **common/logger.go**: 提供系统日志功能。