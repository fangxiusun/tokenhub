# totp.go 代码阅读文档

## 1. 全局总结

`totp.go` 实现了基于时间的一次性密码（TOTP）功能，用于双因素认证（2FA）。文件包含 TOTP 密钥生成、验证码验证、备用恢复码生成与验证、二维码数据生成等功能。所有函数均依赖 `github.com/pquerna/otp` 库实现核心算法。

## 2. 依赖关系

- `crypto/rand` — 加密安全随机数生成
- `fmt` — 格式化输出
- `os` — 环境变量读取
- `strconv` — 字符串转换
- `strings` — 字符串操作
- `github.com/pquerna/otp` — OTP 核心库
- `github.com/pquerna/otp/totp` — TOTP 实现
- `common` 包内函数：`Get2FAIssuer`、`Password2Hash`、`SystemName`

## 3. 类型定义

### 常量

```go
const (
    BackupCodeLength = 8   // 备用码长度
    BackupCodeCount  = 4   // 生成备用码数量
    MaxFailAttempts  = 5   // 最大失败尝试次数
    LockoutDuration  = 300 // 锁定时间（秒）
)
```

### 无自定义类型

## 4. 函数详解

### GenerateTOTPSecret

```go
func GenerateTOTPSecret(accountName string) (*otp.Key, error)
```

生成 TOTP 密钥和配置：
1. 获取发行者名称（`Get2FAIssuer`）
2. 使用 `totp.Generate` 生成密钥，配置包括：发行者、账户名、周期 30 秒、6 位数字、SHA1 算法
3. 返回 OTP 密钥对象

### ValidateTOTPCode

```go
func ValidateTOTPCode(secret, code string) bool
```

验证 TOTP 验证码：
1. 清理验证码中的空格
2. 检查长度是否为 6 位
3. 调用 `totp.Validate` 进行验证

### GenerateBackupCodes

```go
func GenerateBackupCodes() ([]string, error)
```

生成备用恢复码：
1. 循环生成 `BackupCodeCount`（4）个备用码
2. 每个备用码由 `generateRandomBackupCode` 生成
3. 返回备用码数组

### generateRandomBackupCode

```go
func generateRandomBackupCode() (string, error)
```

生成单个备用码：
1. 使用加密安全随机数从字符集（A-Z, 0-9）中选择字符
2. 生成长度为 `BackupCodeLength`（8）的随机字符串
3. 格式化为 `XXXX-XXXX` 格式

### ValidateBackupCode

```go
func ValidateBackupCode(code string) bool
```

验证备用码格式：
1. 移除所有分隔符并转为大写
2. 检查长度是否为 8 位
3. 验证每个字符是否为字母或数字

### NormalizeBackupCode

```go
func NormalizeBackupCode(code string) string
```

标准化备用码格式：
1. 移除所有分隔符并转为大写
2. 如果长度为 8 位，格式化为 `XXXX-XXXX` 格式
3. 否则返回原字符串

### HashBackupCode

```go
func HashBackupCode(code string) (string, error)
```

对备用码进行哈希：
1. 标准化备用码格式
2. 调用 `Password2Hash` 进行哈希

### Get2FAIssuer

```go
func Get2FAIssuer() string
```

获取 2FA 发行者名称：
1. 返回 `SystemName` 常量

### getEnvOrDefault

```go
func getEnvOrDefault(key, defaultValue string) string
```

获取环境变量或默认值：
1. 使用 `os.LookupEnv` 检查环境变量是否存在
2. 如果存在返回环境变量值，否则返回默认值

### ValidateNumericCode

```go
func ValidateNumericCode(code string) (string, error)
```

验证数字验证码格式：
1. 移除空格
2. 检查长度是否为 6 位
3. 验证是否为纯数字
4. 返回标准化后的验证码或错误

### GenerateQRCodeData

```go
func GenerateQRCodeData(secret, username string) string
```

生成二维码数据：
1. 构建 OTP URI：`otpauth://totp/...`
2. 包含发行者、账户名、密钥、发行者、6 位数字、30 秒周期
3. 返回 URI 字符串

## 5. 关键逻辑分析

### 备用码安全设计

- 备用码使用加密安全随机数生成，字符集为 A-Z 和 0-9
- 备用码格式为 `XXXX-XXXX`，便于用户输入
- 备用码在存储前需要进行哈希处理（`HashBackupCode`）
- 备用码验证后需要手动删除，防止重放攻击

### TOTP 参数配置

| 参数 | 值 | 说明 |
|------|-----|------|
| Period | 30 秒 | 验证码有效期 |
| Digits | 6 | 验证码长度 |
| Algorithm | SHA1 | 哈希算法 |

### 错误处理

- 随机数生成失败时返回错误
- 验证码格式错误时返回 false 或错误信息
- 备用码生成失败时返回错误

### 使用场景

- 用户启用双因素认证时生成 TOTP 密钥
- 用户登录时验证 TOTP 验证码
- 用户备份恢复码用于账户恢复
- 生成二维码供用户扫描添加到认证器应用

## 6. 关联文件

- `common/crypto.go` — 密码哈希函数 `Password2Hash`
- `common/init.go` — 系统名称 `SystemName` 定义
- `middleware/auth.go` — 双因素认证中间件
- `controller/auth.go` — 认证控制器，调用 TOTP 验证