# verification.go 代码阅读文档

## 1. 全局总结

`verification.go` 是一个轻量级的内存验证码管理模块，提供验证码的生成、注册、校验和删除功能。该模块采用内存 Map 存储验证码信息，通过互斥锁（`sync.Mutex`）保证并发安全，主要用于邮箱验证码验证和密码重置令牌管理两种业务场景。

**核心设计特点：**
- 纯内存存储，不依赖外部数据库或缓存（如 Redis）
- 基于 UUID 生成验证码，保证唯一性
- 支持多用途（purpose）隔离，同一 key 在不同用途下互不影响
- 自动过期清理机制，防止内存无限增长
- 最大容量限制（默认 10 条），超出时自动清理过期条目

**局限性：**
- 服务重启后验证码数据全部丢失
- 仅适用于单实例部署，多实例间不共享状态
- 最大容量较小（10 条），不适合高并发场景

## 2. 依赖关系

### 标准库依赖

| 包名 | 用途 |
|------|------|
| `strings` | 字符串操作，用于去除 UUID 中的连字符 |
| `sync` | 提供 `Mutex` 互斥锁，保证并发安全 |
| `time` | 时间处理，用于记录验证码创建时间和过期判断 |

### 第三方依赖

| 包名 | 用途 |
|------|------|
| `github.com/google/uuid` | 生成 UUID 作为验证码的基础数据 |

### 被引用关系

该模块被以下文件引用：

| 文件 | 使用的函数/常量 |
|------|----------------|
| `controller/user.go` | `VerifyCodeWithKey`、`EmailVerificationPurpose` |
| `controller/misc.go` | `GenerateVerificationCode`、`RegisterVerificationCodeWithKey`、`VerifyCodeWithKey`、`DeleteKey`、`EmailVerificationPurpose`、`PasswordResetPurpose`、`VerificationValidMinutes` |

**业务调用场景：**
- **用户注册**：生成 6 位邮箱验证码，注册后校验（`controller/user.go:162`）
- **邮箱验证**：发送验证码邮件，用户提交后校验（`controller/user.go:1033`）
- **密码重置**：生成完整 UUID 作为重置令牌，通过邮件链接发送，用户点击后校验（`controller/misc.go:315-362`）

## 3. 类型定义

### verificationValue 结构体

```go
type verificationValue struct {
    code string    // 验证码或令牌值
    time time.Time // 验证码创建时间
}
```

**说明：** 这是一个未导出的内部结构体，用于存储验证码及其创建时间。`code` 字段存储实际的验证码内容，`time` 字段记录创建时刻，用于后续的过期判断。

### 常量定义

| 常量名 | 值 | 用途 |
|--------|-----|------|
| `EmailVerificationPurpose` | `"v"` | 邮箱验证码用途标识 |
| `PasswordResetPurpose` | `"r"` | 密码重置用途标识 |

**设计说明：** 通过 `purpose` 前缀隔离不同业务场景的验证码。例如，同一个邮箱地址可以同时拥有一个邮箱验证码（key: `"vuser@example.com"`）和一个密码重置令牌（key: `"ruser@example.com"`），两者互不干扰。

### 包级变量

| 变量名 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `verificationMutex` | `sync.Mutex` | - | 保护 `verificationMap` 的互斥锁 |
| `verificationMap` | `map[string]verificationValue` | - | 存储所有验证码的映射表，key 格式为 `purpose + key` |
| `verificationMapMaxSize` | `int` | `10` | Map 最大容量限制，超出时触发过期清理 |
| `VerificationValidMinutes` | `int` | `10` | 验证码有效时间（分钟），可被外部修改 |

**注意：** `VerificationValidMinutes` 是导出变量，允许外部代码动态调整验证码有效期。

## 4. 函数详解

### GenerateVerificationCode(length int) string

**功能：** 生成指定长度的验证码字符串。

**参数：**
- `length int` — 期望的验证码长度。若为 0，则返回完整的 UUID 字符串（去除连字符后为 32 字符）。

**返回值：**
- `string` — 生成的验证码字符串。

**实现逻辑：**
1. 使用 `uuid.New()` 生成一个 v4 UUID
2. 通过 `strings.Replace` 移除所有连字符 `-`
3. 若 `length == 0`，返回完整的 32 字符字符串
4. 否则截取前 `length` 个字符返回

**示例：**
- `GenerateVerificationCode(6)` → `"a1b2c3"` （6 位验证码）
- `GenerateVerificationCode(0)` → `"a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6"` （完整 UUID）
- `GenerateVerificationCode(12)` → `"a1b2c3d4e5f6"` （12 位密码）

**潜在问题：** 当 `length > 32` 时，会导致字符串切片越界 panic。但实际使用场景中不会出现此情况（最大调用为 12 位密码生成）。

---

### RegisterVerificationCodeWithKey(key string, code string, purpose string)

**功能：** 注册一个新的验证码到内存 Map 中。

**参数：**
- `key string` — 标识符（如邮箱地址）
- `code string` — 验证码内容
- `purpose string` — 用途标识（`EmailVerificationPurpose` 或 `PasswordResetPurpose`）

**实现逻辑：**
1. 加锁（`verificationMutex.Lock()`）
2. 使用 `purpose + key` 作为 Map 的 key，存储验证码和当前时间
3. 若 Map 大小超过 `verificationMapMaxSize`（默认 10），调用 `removeExpiredPairs()` 清理过期条目
4. 解锁（`defer verificationMutex.Unlock()`）

**说明：** 同一 `key` 和 `purpose` 组合下，新注册的验证码会覆盖旧的。

---

### VerifyCodeWithKey(key string, code string, purpose string) bool

**功能：** 校验指定 key 的验证码是否正确且未过期。

**参数：**
- `key string` — 标识符（如邮箱地址）
- `code string` — 待校验的验证码
- `purpose string` — 用途标识

**返回值：**
- `bool` — 校验结果。`true` 表示验证码正确且在有效期内；`false` 表示不存在、已过期或验证码错误。

**实现逻辑：**
1. 加锁
2. 从 Map 中查找 `purpose + key` 对应的 `verificationValue`
3. 若不存在（`!okay`），返回 `false`
4. 若验证码已过期（当前时间与创建时间的差值 >= `VerificationValidMinutes * 60` 秒），返回 `false`
5. 比较验证码是否相等，返回比较结果
6. 解锁

**说明：** 校验成功后不会自动删除验证码，需要调用方手动调用 `DeleteKey` 清理。

---

### DeleteKey(key string, purpose string)

**功能：** 删除指定 key 和 purpose 的验证码记录。

**参数：**
- `key string` — 标识符
- `purpose string` — 用途标识

**实现逻辑：**
1. 加锁
2. 从 Map 中删除 `purpose + key` 对应的条目
3. 解锁

**典型使用场景：** 验证码校验成功后清理，防止重复使用。

---

### removeExpiredPairs()

**功能：** 内部辅助函数，清理 Map 中所有已过期的验证码条目。

**注意事项：**
- 这是一个未导出的内部函数
- **调用方必须已持有 `verificationMutex` 锁**，函数内部不再加锁
- 遍历 Map，删除所有过期条目

**实现逻辑：**
1. 获取当前时间
2. 遍历 `verificationMap` 的所有 key
3. 若某条目的创建时间与当前时间的差值 >= `VerificationValidMinutes * 60` 秒，删除该条目

---

### init()

**功能：** 包初始化函数，在 `common` 包被导入时自动执行。

**实现逻辑：**
1. 加锁
2. 初始化 `verificationMap` 为一个空的 `map[string]verificationValue`
3. 解锁

**说明：** 确保 `verificationMap` 在使用前已被初始化，避免 nil map 导致的 panic。

## 5. 关键逻辑分析

### 5.1 并发安全机制

模块使用 `sync.Mutex` 保护对 `verificationMap` 的所有读写操作：

- **写操作（RegisterVerificationCodeWithKey、DeleteKey）：** 完整加锁
- **读操作（VerifyCodeWithKey）：** 完整加锁
- **内部清理（removeExpiredPairs）：** 依赖调用方已加锁，不重复加锁

这种设计避免了 Go 语言中 map 的并发读写 panic 问题，但也意味着所有操作是串行的，在高并发场景下可能成为性能瓶颈。

### 5.2 验证码过期策略

过期判断基于时间差计算：

```go
int(now.Sub(value.time).Seconds()) >= VerificationValidMinutes * 60
```

- `VerificationValidMinutes` 默认为 10 分钟
- 将时间差转换为秒后与 `10 * 60 = 600` 秒比较
- 若验证码创建时间距今超过 600 秒，则判定为过期

**特点：** 过期验证码不会立即删除，仅在以下两种情况下被清理：
1. 新验证码注册时且 Map 超过最大容量限制
2. 手动调用 `removeExpiredPairs()`（目前无外部调用）

### 5.3 Map 容量控制

```go
if len(verificationMap) > verificationMapMaxSize {
    removeExpiredPairs()
}
```

- `verificationMapMaxSize` 默认为 10
- 仅在注册新验证码且 Map 超过限制时触发清理
- 清理所有过期条目，不保证清理后一定低于限制
- 清理后也不立即插入新条目（当前逻辑是先插入后检查）

**潜在问题：** 清理逻辑在插入之后执行，可能导致 Map 短暂超过 `verificationMapMaxSize`。

### 5.4 验证码生成策略

基于 UUID v4 生成，去除连字符后截取：

- UUID v4 有 122 位随机性，截取后的短验证码也有足够的随机性
- 不使用数字验证码（如 6 位纯数字），因此无法通过枚举暴力破解
- 生成的验证码包含字母和数字混合，用户需要准确输入

### 5.5 用途隔离机制

通过 `purpose` 前缀实现不同业务场景的隔离：

```
邮箱验证: key = "v" + "user@example.com" → "vuser@example.com"
密码重置: key = "r" + "user@example.com" → "ruser@example.com"
```

同一邮箱可以同时拥有两种验证码，互不影响。

## 6. 关联文件

### 直接调用方

| 文件路径 | 调用场景 |
|----------|----------|
| `controller/user.go` | 用户注册时校验邮箱验证码（第 162 行）、登录时校验邮箱验证码（第 1033 行） |
| `controller/misc.go` | 发送邮箱验证码（第 287-288 行）、发送密码重置邮件（第 315-316 行）、验证密码重置令牌（第 349 行）、重置密码后删除令牌（第 362 行） |

### 间接关联

| 文件路径 | 关联说明 |
|----------|----------|
| `controller/user.go:356` | 使用 `GenerateVerificationCode(12)` 生成 12 位随机密码，作为用户初始密码 |
| `common/json.go` | 同属 `common` 包，提供 JSON 序列化/反序列化功能 |

### 数据流向

```
用户请求发送验证码
    ↓
controller/misc.go: SendEmailVerification()
    ↓
common.GenerateVerificationCode(6)  →  生成 6 位验证码
common.RegisterVerificationCodeWithKey()  →  存入内存 Map
    ↓
发送邮件（包含验证码）
    ↓
用户提交验证码
    ↓
controller/user.go: Register() / Login()
    ↓
common.VerifyCodeWithKey()  →  校验验证码
    ↓
common.DeleteKey()  →  清理已使用的验证码
```
