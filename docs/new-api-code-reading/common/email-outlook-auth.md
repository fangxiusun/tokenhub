# email-outlook-auth.go 代码阅读文档

## 1. 全局总结

本文件实现了 Outlook 邮箱的 SMTP 认证机制。由于 Outlook/Office365 邮箱不支持标准的 `PLAIN` 认证方式，需要使用 `LOGIN` 认证方式。文件定义了 `outlookAuth` 结构体来实现 `smtp.Auth` 接口，并提供了一个判断函数来识别 Outlook 服务器。

## 2. 依赖关系

| 包名 | 用途 |
|------|------|
| `errors` | 创建自定义错误 |
| `net/smtp` | SMTP 协议相关类型（`smtp.Auth`, `smtp.ServerInfo`） |
| `strings` | 字符串包含判断 |

**被依赖方**：本文件中的函数被 `email.go` 调用（`LoginAuth` 和 `isOutlookServer`）。

## 3. 类型定义

### `outlookAuth` 结构体

```go
type outlookAuth struct {
    username, password string
}
```

实现 `smtp.Auth` 接口，用于 Outlook 邮箱的 LOGIN 认证流程。

| 字段 | 类型 | 说明 |
|------|------|------|
| `username` | `string` | SMTP 登录用户名 |
| `password` | `string` | SMTP 登录密码/Token |

## 4. 函数详解

### `LoginAuth(username, password string) smtp.Auth`

**功能**：创建并返回一个 Outlook 认证实例。

**参数**：
- `username` — SMTP 登录用户名
- `password` — SMTP 登录密码/Token

**返回值**：`smtp.Auth` 接口实例

**逻辑**：直接构造 `outlookAuth` 结构体并返回指针。

### `(*outlookAuth) Start(_ *smtp.ServerInfo) (string, []byte, error)`

**功能**：`smtp.Auth` 接口方法，启动认证流程。

**参数**：
- `_ *smtp.ServerInfo` — 服务器信息（未使用）

**返回值**：
- `string` — 认证机制名称 `"LOGIN"`
- `[]byte` — 初始数据（空）
- `error` — 始终为 `nil`

**逻辑**：声明使用 `LOGIN` 认证方式，不需要初始数据。

### `(*outlookAuth) Next(fromServer []byte, more bool) ([]byte, error)`

**功能**：`smtp.Auth` 接口方法，响应服务器的认证挑战。

**参数**：
- `fromServer []byte` — 服务器发送的数据
- `more bool` — 是否还有后续挑战

**返回值**：
- `[]byte` — 响应数据（用户名或密码）
- `error` — 未知挑战时返回错误

**逻辑**：
1. 如果 `more` 为 `true`，根据服务器返回的提示信息响应：
   - `"Username:"` → 返回用户名
   - `"Password:"` → 返回密码
   - 其他 → 返回 `"unknown fromServer"` 错误
2. 如果 `more` 为 `false`，返回 `nil, nil` 表示认证完成

### `isOutlookServer(server string) bool`

**功能**：判断给定的 SMTP 服务器是否为 Outlook 服务器。

**参数**：
- `server string` — SMTP 服务器地址

**返回值**：`bool` — 是否为 Outlook 服务器

**逻辑**：检查服务器地址是否包含 `"outlook"` 或 `"onmicrosoft"` 子串，以兼容多地区的 Outlook 邮箱和 Office365 邮箱。

**注意**：代码注释指出，理想情况下应该添加一个 Option 参数来区分是否使用 LOGIN 方式登录，当前是临时兼容方案。

## 5. 关键逻辑分析

### LOGIN 认证流程

SMTP LOGIN 认证是一种非标准的认证方式，流程如下：
1. 客户端发送 `AUTH LOGIN`
2. 服务器返回 `"Username:"` 挑战
3. 客户端发送 Base64 编码的用户名
4. 服务器返回 `"Password:"` 挑战
5. 客户端发送 Base64 编码的密码
6. 服务器确认认证结果

本实现中，Base64 编码由 Go 标准库的 `net/smtp` 包处理，`outlookAuth` 只负责提供原始字符串。

### Outlook 服务器识别

识别逻辑基于简单的字符串匹配，存在以下情况：
- 包含 `"outlook"` → Microsoft Outlook 邮箱
- 包含 `"onmicrosoft"` → Office365/Exchange Online 邮箱

## 6. 关联文件

| 文件 | 关联关系 |
|------|----------|
| `email.go` | 调用 `LoginAuth()` 和 `isOutlookServer()` |
| `setting/` | 提供 SMTP 相关配置（`SMTPAccount`, `SMTPToken`, `SMTPServer` 等） |
