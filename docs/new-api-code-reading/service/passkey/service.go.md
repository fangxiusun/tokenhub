# service.go 代码阅读文档

## 1. 全局总结

本文件是 Passkey（WebAuthn）服务的核心配置模块，负责构建 `webauthn.WebAuthn` 实例。包含 RP（Relying Party）配置的解析、Origin 推导、RPID 确定、协议检测等逻辑。是 Passkey 注册和认证流程的基础配置层。

## 2. 依赖关系

| 包 | 用途 |
|---|---|
| `errors` | 错误创建 |
| `fmt` | 格式化错误信息 |
| `net` | 网络工具（`SplitHostPort`） |
| `net/http` | HTTP 请求处理 |
| `net/url` | URL 解析 |
| `strings` | 字符串处理 |
| `time` | 超时配置 |
| `github.com/QuantumNous/new-api/common` | 系统名称、调试模式标志 |
| `github.com/QuantumNous/new-api/setting/system_setting` | Passkey 设置和服务器地址 |
| `github.com/go-webauthn/webauthn/protocol` | WebAuthn 协议常量和选择器 |
| `github.com/go-webauthn/webauthn/webauthn` | WebAuthn 核心库 |

## 3. 类型定义

### 常量

```go
const (
    RegistrationSessionKey = "passkey_registration_session"
    LoginSessionKey        = "passkey_login_session"
    VerifySessionKey       = "passkey_verify_session"
)
```

用于 Session 存储的键名，分别对应注册、登录、验证三种 Passkey 操作的会话数据。

## 4. 函数详解

### `BuildWebAuthn(r *http.Request) (*webauthn.WebAuthn, error)`

**签名**: `func BuildWebAuthn(r *http.Request) (*webauthn.WebAuthn, error)`

**功能**: 根据系统设置和请求上下文构建 WebAuthn 实例。

**逻辑**:
1. 从 `system_setting.GetPasskeySettings()` 获取 Passkey 配置
2. 解析 `RPDisplayName`，为空时使用 `common.SystemName`
3. 调用 `resolveOrigins` 解析允许的 Origin 列表
4. 调用 `resolveRPID` 确定 Relying Party ID
5. 构建 `AuthenticatorSelection` 配置：
   - `ResidentKey` 设为 `Required`
   - `UserVerification` 从设置读取，默认为 `Preferred`
   - `AuthenticatorAttachment` 从设置读取（可选）
6. 设置登录和注册超时均为 2 分钟
7. 调用 `webauthn.New(config)` 创建实例

### `resolveOrigins(r *http.Request, settings) ([]string, error)`

**签名**: `func resolveOrigins(r *http.Request, settings *system_setting.PasskeySettings) ([]string, error)`

**功能**: 解析并验证允许的 Origin 列表。

**逻辑**:
1. **手动配置优先**: 若 `settings.Origins` 非空，按逗号分割并逐个验证：
   - 检查是否为不安全的 HTTP Origin（非 localhost 时拒绝，除非 `AllowInsecureOrigin` 为 true）
   - 过滤空字符串
   - 过滤后为空则跳转自动推导
2. **自动推导**: 从请求上下文推导 Origin：
   - 检测协议（`X-Forwarded-Proto` → TLS → URL scheme → `X-Forwarded-Protocol` → 默认 HTTP）
   - 非 HTTPS 且非 localhost 时拒绝
   - 从 `r.Host` 或 `system_setting.ServerAddress` 获取主机名
   - 构造 `{scheme}://{host}` 格式的 Origin

### `resolveRPID(r *http.Request, settings, origins) (string, error)`

**签名**: `func resolveRPID(r *http.Request, settings *system_setting.PasskeySettings, origins []string) (string, error)`

**功能**: 确定 Relying Party ID。

**逻辑**:
1. 若 `settings.RPID` 非空，去除端口后返回
2. 否则从第一个 Origin 的主机部分提取（去除端口）
3. Origin 为空时返回错误

### `hostWithoutPort(host string) string`

**签名**: `func hostWithoutPort(host string) string`

**功能**: 从主机名中去除端口号。

### `detectScheme(r *http.Request) string`

**签名**: `func detectScheme(r *http.Request) string`

**功能**: 检测请求的协议方案。

**检测优先级**:
1. `X-Forwarded-Proto` 头
2. TLS 连接状态（`r.TLS != nil`）
3. `r.URL.Scheme`
4. `X-Forwarded-Protocol` 头
5. 默认 `"http"`

## 5. 关键逻辑分析

1. **安全策略**: 默认拒绝 HTTP Origin（除 localhost 外），通过 `AllowInsecureOrigin` 开关控制。这是 WebAuthn 安全规范的要求。

2. **Origin 推导链**: 在用户未手动配置 Origin 时，通过多层回退策略自动推导，适用于反向代理（`X-Forwarded-Proto`）、直连（TLS 检测）等多种部署场景。

3. **RPID 去端口化**: WebAuthn 规范要求 RP ID 为域名（不含端口），`hostWithoutPort` 确保这一约束。

4. **超时配置**: 注册和登录均设置 2 分钟超时并强制执行（`Enforce: true`），防止长时间挂起的认证会话。

## 6. 关联文件

- `new-api/service/passkey/session.go` — Passkey 会话数据管理
- `new-api/service/passkey/user.go` — WebAuthn 用户接口实现
- `new-api/setting/system_setting/` — Passkey 设置定义
- `new-api/controller/passkey.go`（推测） — Passkey HTTP 控制器
