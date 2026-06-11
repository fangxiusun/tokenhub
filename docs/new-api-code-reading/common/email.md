# email.go 代码阅读文档

## 1. 全局总结

本文件实现了系统的邮件发送功能，支持 TLS/SSL 加密连接和普通 SMTP 连接两种方式。核心函数 `SendEmail` 负责构建邮件内容并通过 SMTP 协议发送。文件还包含 Message-ID 生成、认证方式选择等辅助功能。

## 2. 依赖关系

| 包名 | 用途 |
|------|------|
| `crypto/tls` | TLS 加密连接 |
| `encoding/base64` | 邮件主题的 UTF-8 Base64 编码 |
| `fmt` | 格式化字符串 |
| `net/smtp` | SMTP 协议客户端 |
| `slices` | 切片包含判断（Go 1.21+） |
| `strings` | 字符串分割与处理 |
| `time` | 时间格式化 |

**内部依赖**：
- `common/json.go` — 未直接依赖，但遵循项目 JSON 规范
- `email-outlook-auth.go` — `LoginAuth()`, `isOutlookServer()`
- `setting/` — SMTP 配置变量（`SMTPFrom`, `SMTPAccount`, `SMTPServer`, `SMTPPort`, `SMTPToken`, `SMTPSSLEnabled`, `SMTPForceAuthLogin`, `SystemName`, `EmailLoginAuthServerList`）
- `common/` — `GetRandomString()`, `SysError()`

## 3. 类型定义

本文件无自定义类型定义，仅使用标准库和内部包的类型。

## 4. 函数详解

### `generateMessageID() (string, error)`

**功能**：生成符合 RFC 标准的 Message-ID。

**返回值**：
- `string` — 格式为 `<时间戳.随机字符串@域名>` 的 Message-ID
- `error` — SMTPFrom 格式无效时返回错误

**逻辑**：
1. 从 `SMTPFrom` 中提取 `@` 后面的域名部分
2. 使用纳秒时间戳和 12 位随机字符串构造 Message-ID
3. 格式：`<1718234567890.abcdefghijkl@example.com>`

### `shouldUseSMTPLoginAuth() bool`

**功能**：判断是否应该使用 LOGIN 认证方式。

**返回值**：`bool`

**逻辑**：
1. 如果 `SMTPForceAuthLogin` 为 `true`，直接返回 `true`
2. 否则检查是否为 Outlook 服务器（`isOutlookServer`）
3. 或者服务器在 `EmailLoginAuthServerList` 白名单中

### `getSMTPAuth() smtp.Auth`

**功能**：根据服务器类型选择合适的 SMTP 认证方式。

**返回值**：`smtp.Auth` — LOGIN 或 PLAIN 认证实例

**逻辑**：
- 如果需要 LOGIN 认证 → 调用 `LoginAuth()` 返回 Outlook 认证
- 否则 → 使用标准的 `smtp.PlainAuth()`

### `SendEmail(subject, receiver, content string) error`

**功能**：发送 HTML 格式的邮件。

**参数**：
- `subject` — 邮件主题
- `receiver` — 收件人（多个用 `;` 分隔）
- `content` — HTML 格式的邮件正文

**返回值**：`error` — 发送失败时返回错误

**逻辑**：
1. **兼容处理**：如果 `SMTPFrom` 为空，使用 `SMTPAccount` 作为发件人
2. **生成 Message-ID**：调用 `generateMessageID()`
3. **配置检查**：SMTP 服务器未配置时返回错误
4. **编码主题**：使用 RFC 2047 编码（`=?UTF-8?B?...?=`）
5. **构建邮件**：拼接标准邮件头（To, From, Subject, Date, Message-ID, Content-Type）
6. **选择连接方式**：
   - **TLS/SSL**（端口 465 或 `SMTPSSLEnabled`）：
     - 创建 TLS 配置（跳过证书验证）
     - 使用 `tls.Dial` 建立加密连接
     - 通过 `smtp.NewClient` 创建客户端
     - 依次执行：认证 → 设置发件人 → 添加收件人 → 发送数据
   - **普通连接**：
     - 使用 `smtp.SendMail` 一步完成
7. **错误处理**：失败时调用 `SysError` 记录日志

## 5. 关键逻辑分析

### 邮件编码

邮件主题使用 RFC 2047 编码，格式为 `=?UTF-8?B?<base64编码>?=`，确保中文等非 ASCII 字符在邮件头中正确传输。

### TLS 连接流程

```
TLS 配置（跳过证书验证）
    ↓
tls.Dial 建立 TCP+TLS 连接
    ↓
smtp.NewClient 创建 SMTP 客户端
    ↓
client.Auth 认证
    ↓
client.Mail 设置发件人
    ↓
client.Rcpt 添加收件人（循环）
    ↓
client.Data 发送邮件内容
    ↓
w.Close 完成发送
```

### 多收件人处理

收件人使用 `;` 分隔符。在 TLS 连接中，通过循环 `client.Rcpt()` 逐个添加；在普通连接中，直接传递分割后的切片给 `smtp.SendMail`。

### 安全注意

TLS 配置中设置了 `InsecureSkipVerify: true`，跳过了服务器证书验证。这在生产环境中存在中间人攻击风险，但在实际使用中（如对接企业邮箱）通常可以接受。

## 6. 关联文件

| 文件 | 关联关系 |
|------|----------|
| `email-outlook-auth.go` | 提供 Outlook LOGIN 认证实现 |
| `setting/` | 提供 SMTP 配置变量 |
| `common/log.go` | `SysError()` 错误日志记录 |
| `common/utils.go` | `GetRandomString()` 随机字符串生成 |
