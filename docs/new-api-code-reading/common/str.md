# str.go 代码阅读文档

## 1. 全局总结

本文件提供了一系列字符串处理工具函数，涵盖：日志内容截断、随机字符串生成、JSON 互转、类型转换、敏感信息掩码（URL、域名、IP、API Key）等。其中 `MaskSensitiveInfo` 是核心功能，用于在日志中隐藏敏感信息，防止 PII（个人可识别信息）泄露。

## 2. 依赖关系

| 包名 | 用途 |
|------|------|
| `encoding/base64` | Base64 编码 |
| `encoding/json` | JSON 序列化/反序列化 |
| `fmt` | 格式化输出 |
| `net/url` | URL 解析（用于敏感信息掩码） |
| `regexp` | 正则表达式匹配 |
| `strconv` | 字符串与数值转换 |
| `strings` | 字符串操作 |
| `unsafe` | 零拷贝字符串转字节切片 |
| `github.com/samber/lo` | 泛型工具库（随机字符串生成） |

**项目内依赖：**
- `common.Unmarshal` — JSON 反序列化封装
- `common.DebugEnabled` — 调试模式开关

## 3. 类型定义

### 正则表达式常量

```go
var (
    maskURLPattern    = regexp.MustCompile(`(http|https)://[^\s/$.?#].[^\s]*`)
    maskDomainPattern = regexp.MustCompile(`\b(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}\b`)
    maskIPPattern     = regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)
    maskApiKeyPattern = regexp.MustCompile(`(['"]?)api_key:([^\s'"]+)(['"]?)`)
)
```

| 正则变量 | 匹配目标 | 用途 |
|----------|----------|------|
| `maskURLPattern` | HTTP/HTTPS URL | 掩码 URL 中的主机、路径、参数 |
| `maskDomainPattern` | 域名（如 openai.com） | 掩码无协议前缀的域名 |
| `maskIPPattern` | IPv4 地址 | 掩码 IP 地址 |
| `maskApiKeyPattern` | `api_key:xxx` 格式 | 掩码 API 密钥 |

### 常量

```go
const LocalLogContentLimit = 2048
```

本地日志内容最大长度限制（字符数）。

## 4. 函数详解

### LocalLogPreview(content string) string

```go
func LocalLogPreview(content string) string
```

限制日志内容长度：
- 调试模式下返回完整内容
- 非调试模式下，超过 2048 字符截断并添加 `[truncated, ...]` 标记

### GetStringIfEmpty(str string, defaultValue string) string

```go
func GetStringIfEmpty(str string, defaultValue string) string
```

如果字符串为空，返回默认值。

### GetRandomString(length int) string

```go
func GetRandomString(length int) string
```

生成指定长度的随机字母数字字符串。使用 `lo.RandomString` 实现。

### MapToJsonStr(m map[string]interface{}) string

```go
func MapToJsonStr(m map[string]interface{}) string
```

将 map 转换为 JSON 字符串。

### StrToMap(str string) (map[string]interface{}, error)

```go
func StrToMap(str string) (map[string]interface{}, error)
```

将 JSON 字符串解析为 map。

### StrToJsonArray(str string) ([]interface{}, error)

```go
func StrToJsonArray(str string) ([]interface{}, error)
```

将 JSON 字符串解析为数组。

### IsJsonArray(str string) bool

```go
func IsJsonArray(str string) bool
```

判断字符串是否为有效的 JSON 数组。

### IsJsonObject(str string) bool

```go
func IsJsonObject(str string) bool
```

判断字符串是否为有效的 JSON 对象。

### String2Int(str string) int

```go
func String2Int(str string) int
```

将字符串转换为整数，转换失败返回 0。

### StringsContains(strs []string, str string) bool

```go
func StringsContains(strs []string, str string) bool
```

检查字符串切片是否包含指定字符串。

### StringToByteSlice(s string) []byte

```go
func StringToByteSlice(s string) []byte
```

零拷贝将字符串转换为字节切片。**注意：** 返回的切片仅可读，追加操作会导致 panic。

### EncodeBase64(str string) string

```go
func EncodeBase64(str string) string
```

对字符串进行 Base64 标准编码。

### GetJsonString(data any) string

```go
func GetJsonString(data any) string
```

将任意类型序列化为 JSON 字符串。

### NormalizeBillingPreference(pref string) string

```go
func NormalizeBillingPreference(pref string) string
```

规范化计费偏好设置：
- 有效值：`subscription_first`、`wallet_first`、`subscription_only`、`wallet_only`
- 无效值默认返回 `subscription_first`

### MaskEmail(email string) string

```go
func MaskEmail(email string) string
```

掩码邮箱地址，防止 PII 泄露：
- 空邮箱返回 `***masked***`
- 无 `@` 符号返回 `***masked***`
- 正常邮箱只显示域名部分：`***@example.com`

### maskHostTail(parts []string) []string

```go
func maskHostTail(parts []string) []string
```

获取域名尾部应保留的部分：
- 国家代码顶级域名（如 `co.uk`、`com.cn`）保留最后两级
- 其他情况只保留顶级域名

### maskHostForURL(host string) string

```go
func maskHostForURL(host string) string
```

掩码 URL 中的主机名：
- `api.openai.com` → `***.com`
- `sub.domain.co.uk` → `***.co.uk`

### maskHostForPlainDomain(domain string) string

```go
func maskHostForPlainDomain(domain string) string
```

掩码纯域名（无协议前缀）：
- `openai.com` → `***.com`
- `api.openai.com` → `***.***.com`
- `sub.domain.co.uk` → `***.***.co.uk`

### MaskSensitiveInfo(str string) string

```go
func MaskSensitiveInfo(str string) string
```

核心敏感信息掩码函数，处理流程：

1. **URL 掩码**：解析 URL，掩码主机、路径、查询参数
   - `http://example.com` → `http://***.com`
   - `https://api.test.org/v1/users/123?key=secret` → `https://***.org/***/***/?key=***`
2. **域名掩码**：掩码无协议前缀的域名
   - `openai.com` → `***.com`
3. **IP 掩码**：将 IPv4 地址替换为 `***.***.***.***`
4. **API Key 掩码**：将 `api_key:xxx` 格式的密钥值替换为 `***`

## 5. 关键逻辑分析

### 5.1 敏感信息掩码策略

`MaskSensitiveInfo` 采用多层掩码策略：
- **URL 掩码**：保留协议，掩码主机（只显示 TLD），掩码路径和参数
- **域名掩码**：根据子域名深度使用多个 `***` 标记
- **IP 掩码**：完全替换为星号
- **API Key 掩码**：只保留 key 名称，掩码值

### 5.2 零拷贝技术

`StringToByteSlice` 使用 `unsafe.Pointer` 实现零拷贝转换：
- 避免内存分配和复制
- 但返回的切片不可追加，否则会导致数据损坏
- 适用于只读场景（如日志输出、比较操作）

### 5.3 国家代码 TLD 处理

`maskHostTail` 函数智能识别国家代码顶级域名：
- 检测最后两级是否符合国家代码 TLD 模式（如 `co.uk`、`com.cn`）
- 如果是，保留两级；否则只保留一级
- 这确保了掩码后的域名仍有意义

### 5.4 日志安全

通过以下函数组合实现日志安全：
- `LocalLogPreview`：限制日志长度，防止日志注入
- `MaskSensitiveInfo`：掩码敏感信息
- `MaskEmail`：保护用户邮箱

### 5.5 潜在问题

1. **正则性能**：多个正则表达式在每次调用时编译（使用 `regexp.MustCompile` 已优化为包级别编译）
2. **掩码不完全**：可能遗漏某些格式的敏感信息（如自定义格式的 API Key）
3. **零拷贝风险**：`StringToByteSlice` 的使用需要开发者明确知道返回值的限制

## 6. 关联文件

- `new-api/common/json.go` — `Unmarshal` 函数定义
- `new-api/common/log.go` — 日志函数，可能使用 `LocalLogPreview`
- `new-api/middleware/logging/` — 日志中间件，可能使用 `MaskSensitiveInfo`
- `new-api/controller/` — 控制器层，可能使用字符串工具函数
