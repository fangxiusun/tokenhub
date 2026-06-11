# init.go 代码阅读文档

## 1. 全局总结

`init.go` 是 `common` 包的初始化文件，负责处理命令行参数解析、环境变量读取和全局配置变量的初始化。该文件是应用程序启动时的核心配置入口，通过 `InitEnv()` 函数统一管理所有运行时配置项，包括端口、日志目录、会话密钥、数据库路径、缓存配置、速率限制、中继超时等。

## 2. 依赖关系

### 标准库依赖
- `flag` - 命令行参数解析
- `fmt` - 格式化输出
- `log` - 日志输出
- `net/http` - HTTP 传输配置
- `os` - 操作系统环境变量和文件操作
- `path/filepath` - 文件路径处理
- `strconv` - 字符串转换
- `strings` - 字符串处理
- `time` - 时间处理

### 项目内部依赖
- `github.com/QuantumNous/new-api/constant` - 常量定义包

## 3. 类型定义

### 全局变量（命令行参数）
```go
var Port         *int    // 监听端口，默认 3000
var PrintVersion *bool   // 打印版本信息
var PrintHelp    *bool   // 打印帮助信息
var LogDir       *string // 日志目录，默认 ./logs
```

## 4. 函数详解

### printHelp()
```go
func printHelp()
```
打印帮助信息，包含项目名称、版本、维护者和使用说明。

### InitEnv()
```go
func InitEnv()
```
主初始化函数，按顺序执行以下操作：

1. **解析命令行参数** - 调用 `flag.Parse()`
2. **环境变量版本覆盖** - 读取 `VERSION` 环境变量
3. **版本/帮助输出** - 根据命令行参数决定是否退出
4. **会话密钥初始化** - 从 `SESSION_SECRET` 环境变量读取，包含默认值安全检查
5. **加密密钥初始化** - 从 `CRYPTO_SECRET` 环境变量读取，默认使用会话密钥
6. **SQLite 路径配置** - 从 `SQLITE_PATH` 环境变量读取
7. **日志目录初始化** - 将相对路径转为绝对路径，不存在则创建
8. **运行时开关配置** - 包括调试模式、内存缓存、节点类型、TLS 配置等
9. **速率限制配置** - API/Web/关键/搜索四类速率限制参数
10. **初始化常量环境** - 调用 `initConstantEnv()`

### initConstantEnv()
```go
func initConstantEnv()
```
初始化 `constant` 包中的配置变量，包括：

- 流式超时、调试模式、文件下载限制
- 请求体大小限制、匿名请求限制
- 流选项强制、Token 计数、媒体 Token 获取
- 任务更新、Azure API 版本
- 通知限制、任务价格补丁
- 重定向信任域名

## 5. 关键逻辑分析

### 安全检查机制
```go
if ss == "random_string" {
    log.Fatal("Please set SESSION_SECRET to a random string.")
}
```
- 防止使用默认的会话密钥，确保生产环境安全

### TLS 配置
```go
if TLSInsecureSkipVerify {
    if tr, ok := http.DefaultTransport.(*http.Transport); ok && tr != nil {
        if tr.TLSClientConfig != nil {
            tr.TLSClientConfig.InsecureSkipVerify = true
        } else {
            tr.TLSClientConfig = InsecureTLSConfig
        }
    }
}
```
- 仅在环境变量明确设置时跳过 TLS 验证
- 修改默认 HTTP 传输的 TLS 配置

### 配置项分类

| 类别 | 环境变量 | 默认值 | 说明 |
|------|----------|--------|------|
| 基础配置 | `PORT`, `LOG_DIR` | 3000, ./logs | 服务监听和日志 |
| 安全配置 | `SESSION_SECRET`, `CRYPTO_SECRET` | - | 加密密钥 |
| 数据库 | `SQLITE_PATH` | - | SQLite 存储路径 |
| 运行模式 | `DEBUG`, `MEMORY_CACHE_ENABLED` | false, false | 调试和缓存 |
| 节点配置 | `NODE_TYPE`, `NODE_NAME` | master, - | 集群节点角色 |
| 超时配置 | `POLLING_INTERVAL`, `RELAY_TIMEOUT` | - | 轮询和中继超时 |
| 速率限制 | `GLOBAL_API_RATE_LIMIT` 等 | 360 | API 限流配置 |

## 6. 关联文件

- `common/constants.go` - 全局常量定义，被 `initConstantEnv()` 初始化
- `common/utils.go` - 包含 `GetEnvOrDefault`, `GetEnvOrDefaultBool`, `GetEnvOrDefaultString` 辅助函数
- `common/session.go` - 使用 `SessionSecret` 变量
- `common/crypto.go` - 使用 `CryptoSecret` 变量
- `model/main.go` - 数据库初始化使用 `SQLitePath` 变量
- `router/router.go` - 使用 `Port` 变量配置监听端口
