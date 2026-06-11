# crypto.go 代码阅读文档

## 1. 全局总结

该文件是项目中加密与哈希相关的工具函数集合，提供了 HMAC-SHA256 签名生成和 bcrypt 密码哈希两个核心功能。文件封装了底层加密库，为项目其他模块提供简洁的加密接口。

## 2. 依赖关系

| 包名 | 用途 |
|------|------|
| `crypto/hmac` | HMAC 消息认证码生成 |
| `crypto/sha256` | SHA-256 哈希算法 |
| `encoding/hex` | 十六进制编码 |
| `golang.org/x/crypto/bcrypt` | bcrypt 密码哈希（加盐哈希） |

内部依赖：`CryptoSecret` 常量（在项目其他文件中定义，作为默认 HMAC 密钥）

## 3. 类型定义

本文件无自定义类型定义，仅使用标准库和扩展库的类型。

## 4. 函数详解

### 4.1 GenerateHMACWithKey(key []byte, data string) string

**功能**：使用指定密钥生成 HMAC-SHA256 签名

**参数**：
- `key`: 用于 HMAC 计算的密钥字节数组
- `data`: 需要签名的字符串数据

**返回值**：十六进制编码的 HMAC 签名字符串

**实现逻辑**：
1. 创建 HMAC-SHA256 实例，绑定密钥
2. 写入数据字节
3. 返回哈希结果的十六进制编码

### 4.2 GenerateHMAC(data string) string

**功能**：使用全局默认密钥 `CryptoSecret` 生成 HMAC-SHA256 签名

**参数**：
- `data`: 需要签名的字符串数据

**返回值**：十六进制编码的 HMAC 签名字符串

**实现逻辑**：
1. 创建 HMAC-SHA256 实例，绑定全局 `CryptoSecret` 密钥
2. 写入数据字节
3. 返回哈希结果的十六进制编码

**注意**：此函数依赖项目中的 `CryptoSecret` 全局常量，该常量应在项目启动时配置

### 4.3 Password2Hash(password string) (string, error)

**功能**：将明文密码转换为 bcrypt 哈希值

**参数**：
- `password`: 明文密码字符串

**返回值**：
- `string`: bcrypt 哈希后的密码字符串
- `error`: 哈希生成过程中的错误

**实现逻辑**：
1. 将密码字符串转换为字节数组
2. 使用 `bcrypt.DefaultCost`（默认成本因子 10）生成哈希
3. 返回哈希字符串

### 4.4 ValidatePasswordAndHash(password string, hash string) bool

**功能**：验证密码与哈希值是否匹配

**参数**：
- `password`: 用户输入的明文密码
- `hash`: 数据库中存储的哈希值

**返回值**：`true` 表示密码匹配，`false` 表示不匹配

**实现逻辑**：
1. 调用 `bcrypt.CompareHashAndPassword` 比较密码与哈希
2. 返回比较结果（err == nil）

## 5. 关键逻辑分析

### 5.1 HMAC 签名的两种模式

文件提供了两种 HMAC 签名方式：
- `GenerateHMACWithKey`：灵活模式，允许调用方指定密钥
- `GenerateHMAC`：简便模式，使用全局默认密钥

这种设计满足了不同场景的需求：
- 需要使用临时密钥或不同密钥时使用 `GenerateHMACWithKey`
- 使用统一密钥的常规场景使用 `GenerateHMAC`

### 5.2 bcrypt 密码哈希的安全性

- 使用 `bcrypt.DefaultCost`（成本因子 10），平衡了安全性和性能
- bcrypt 自动处理加盐，防止彩虹表攻击
- `ValidatePasswordAndHash` 函数返回布尔值，隐藏了具体的错误信息，避免泄露安全信息

### 5.3 HMAC 的十六进制编码

所有 HMAC 函数返回十六进制编码的字符串，这是 JWT、API 签名等场景的标准格式，便于传输和比较。

## 6. 关联文件

- `common/constant.go` 或类似文件：定义 `CryptoSecret` 常量
- `model/user.go`：用户密码存储和验证
- `middleware/auth.go`：JWT 签名验证
- `relay/` 目录：AI API 请求签名
