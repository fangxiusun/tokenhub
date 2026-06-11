# utils.go 代码阅读文档

## 1. 全局总结
utils.go 包含各种通用工具函数，如随机数生成、UUID 生成、时间格式化、数值解析等。

## 2. 依赖关系
### 2.1 导入的包
- `crypto/rand` - 安全随机数
- `encoding/hex` - 十六进制编码
- `fmt` - 格式化
- `math/rand` - 随机数
- `os/exec` - 系统命令执行
- `strconv` - 字符串转换
- `time` - 时间处理
- `uuid` - UUID 生成

### 2.2 被引用的文件
- 整个项目的多个模块

## 3. 类型定义
### 3.1 函数
- `GetUUID() string` - 生成 UUID
- `GetRandomString(n int) string` - 生成随机字符串
- `GetRandomNumber(n int) string` - 生成随机数字
- `GetRandomID() string` - 生成随机 ID
- `TimestampToString(timestamp int64) string` - 时间戳转字符串
- `String2Int(s string) (int, error)` - 字符串转整数

## 4. 函数详解
### 4.1 GetUUID
- **签名**: `func GetUUID() string`
- **职责**: 生成唯一标识符
- **返回值**: UUID 字符串

### 4.2 GetRandomString
- **签名**: `func GetRandomString(n int) string`
- **职责**: 生成指定长度的随机字符串
- **参数**: n - 字符串长度

## 5. 关键逻辑分析
- **安全随机数**: 使用 crypto/rand 生成安全的随机数
- **字符集**: 随机字符串使用字母和数字

## 6. 关联文件
- `model/token.go` - Token 密钥生成
- `model/redemption.go` - 兑换码生成
