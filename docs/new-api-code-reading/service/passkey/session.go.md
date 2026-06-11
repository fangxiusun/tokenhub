# session.go 代码阅读文档

## 1. 全局总结

本文件提供 Passkey WebAuthn 流程中会话数据的保存与读取功能。通过 Gin Session 中间件将 `webauthn.SessionData` 序列化存储到 Cookie/Session 中，实现注册和登录流程的跨请求状态保持。使用"写入-弹出"（Save/Pop）模式确保会话数据的一次性消费。

## 2. 依赖关系

| 包 | 用途 |
|---|---|
| `encoding/json` | SessionData 的序列化与反序列化 |
| `errors` | 错误创建 |
| `github.com/gin-contrib/sessions` | Gin Session 中间件 |
| `github.com/gin-gonic/gin` | HTTP 上下文 |
| `github.com/go-webauthn/webauthn/webauthn` | `SessionData` 类型 |

## 3. 类型定义

### 错误变量

```go
var errSessionNotFound = errors.New("Passkey 会话不存在或已过期")
```

当会话数据不存在时返回的错误。

## 4. 函数详解

### `SaveSessionData(c *gin.Context, key string, data *webauthn.SessionData) error`

**签名**: `func SaveSessionData(c *gin.Context, key string, data *webauthn.SessionData) error`

**功能**: 将 WebAuthn 会话数据保存到 Gin Session 中。

**逻辑**:
1. 获取当前请求的 Session
2. 若 `data` 为 `nil`，删除对应键并保存（清除会话）
3. 将 `data` 序列化为 JSON 字符串
4. 以 `key` 为键存入 Session
5. 调用 `session.Save()` 持久化

### `PopSessionData(c *gin.Context, key string) (*webauthn.SessionData, error)`

**签名**: `func PopSessionData(c *gin.Context, key string) (*webauthn.SessionData, error)`

**功能**: 从 Gin Session 中读取并删除 WebAuthn 会话数据（一次性消费）。

**逻辑**:
1. 获取 Session 中指定键的原始值
2. 若值为 `nil`，返回 `errSessionNotFound`
3. **立即删除**键并保存 Session（确保不被重放）
4. 类型断言处理两种存储格式：
   - `string` → 直接反序列化
   - `[]byte` → 直接反序列化
   - 其他类型 → 返回格式无效错误
5. 返回反序列化后的 `SessionData` 指针

## 5. 关键逻辑分析

1. **Pop 语义**: `PopSessionData` 实现了严格的"读取并删除"语义，从 Session 中取出数据后立即删除。这防止了 WebAuthn 认证数据被重放攻击利用。

2. **类型兼容**: 同时支持 `string` 和 `[]byte` 两种 Session 存储格式，兼容不同 Session 序列化后端的差异。

3. **错误语义清晰**: `errSessionNotFound` 明确表达"会话不存在或已过期"，帮助前端给出准确的用户提示（如"请重新开始认证"）。

4. **Session 键常量**: 配合 `service.go` 中定义的 `RegistrationSessionKey`、`LoginSessionKey`、`VerifySessionKey` 常量使用，确保键名一致。

## 6. 关联文件

- `new-api/service/passkey/service.go` — 定义 Session 键常量和 WebAuthn 构建逻辑
- `new-api/service/passkey/user.go` — WebAuthn 用户接口
- `new-api/middleware/session.go`（推测） — Session 中间件配置
