# 公共工具层详细设计 (`common/`)

## 1. 概述
公共工具层包含在整个应用程序中使用的共享工具、常量和辅助函数。这是系统的基础层，被所有其他层依赖。

## 2. 文件详细说明

### 2.1 核心工具
- **`json.go`** -- 项目范围的 JSON 封装（AGENTS.md 规则 1）：
  - `Marshal(v any)`: JSON 序列化
  - `Unmarshal(data []byte, v any)`: JSON 反序列化
  - `UnmarshalJsonStr(data string, v any)`: JSON 字符串反序列化
  - `DecodeJson(reader io.Reader, v any)`: 从 Reader 解码 JSON
  - `GetJsonType(data json.RawMessage)`: 获取 JSON 类型
  - **设计目的**: 所有业务代码必须使用此封装，禁止直接导入 `encoding/json`
- **`str.go`** -- 字符串工具：URL/域名/IP 掩码、base64 编解码、URL 编码、随机字符串生成、UUID 操作。
- **`utils.go`** -- 通用工具：浏览器打开、随机数/字符串/ID 生成、UUID、时间格式化、数值解析、系统命令执行、base64 转换、文件哈希、模板渲染。
- **`validate.go`** -- 初始化全局 `validator.Validate` 实例（go-playground/validator）用于结构体验证。

### 2.2 加密
- **`crypto.go`** -- 加密辅助：HMAC-SHA256 生成（`GenerateHMAC`, `GenerateHMACWithKey`）、bcrypt 密码哈希（`Password2Hash`）、密码验证（`ValidatePasswordAndHash`）。
- **`hash.go`** -- 哈希工具：SHA-256/SHA-1 原始和十六进制编码、HMAC-SHA256。

### 2.3 数据库
- **`database.go`** -- 数据库类型常量（`mysql`, `sqlite`, `postgres`）和全局布尔标志（`UsingSQLite`, `UsingPostgreSQL`, `UsingMySQL`, `UsingClickHouse`）用于跨数据库条件逻辑。

### 2.4 Redis
- **`redis.go`** -- Redis 客户端初始化和辅助：
  - 连接建立：从 `REDIS_CONN_STRING` 环境变量
  - 键操作：带 TTL 的设置/获取/删除
  - 限流：基于 Redis 的分布式限流
  - 分布式锁：互斥操作
  - 缓存操作：get/set/delete
  - 优雅降级：Redis 不可用时回退到内存缓存

### 2.5 HTTP 与网络
- **`gin.go`** -- Gin 框架工具：请求体读取（带大小限制）、MIME 类型检测、多部分表单解析、IP 提取、User-Agent 解析、上下文键常量（`KeyRequestBody`, `KeyBodyStorage`）。
- **`ssrf_protection.go`** -- SSRF 防护引擎：
  - 域名白名单/黑名单过滤
  - IP CIDR 过滤
  - 私有 IP 阻止
  - 允许端口检查
  - 完整的 DNS 解析验证
- **`url_validator.go`** -- 重定向 URL 安全验证：检查方案（仅 http/https）和域名与受信任域名列表的子域名匹配。
- **`ip.go`** -- IP 地址工具：`IsIP()`, `ParseIP()`, `IsPrivateIP()`（检查回环、链路本地、RFC 1918 范围）、`IsPublicIPv4()`。

### 2.6 邮件
- **`email.go`** -- SMTP 邮件发送：生成消息 ID、选择认证方法（plain vs login）、构建带 TLS 的 MIME 消息、通过配置的 SMTP 服务器发送。
- **`email-outlook-auth.go`** -- 自定义 `smtp.Auth` 实现，用于使用 `LOGIN` 认证机制而非 `PLAIN`/`CRAM-MD5` 的 Outlook/Exchange 服务器。

### 2.7 文件与存储
- **`body_storage.go`** -- `BodyStorage` 接口（`io.ReadSeeker` + `Bytes()` + `Size()`）：内存实现和磁盘实现。
- **`disk_cache.go`** -- 磁盘缓存实现：在 `new-api-body-cache` 目录下创建临时文件，带清理、大小管理和 LRU 驱逐。
- **`disk_cache_config.go`** -- 磁盘缓存线程安全配置：启用状态、阈值、最大大小、路径。
- **`embed-file-system.go`** -- 适配器：通过 `gin-contrib/static` 的 `static.ServeFileSystem` 接口提供 Go `embed.FS`。

### 2.8 限流
- **`rate-limit.go`** -- 内存限流器（`InMemoryRateLimiter`）：基于时间戳的滑动窗口，可配置过期和后台清理协程。
- **`limiter/limiter.go`** -- Redis 分布式限流器：使用 Lua 脚本的原子滑动窗口限流，单例模式。

### 2.9 配置
- **`constants.go`** -- 核心应用范围常量和全局变量：`Version`, `SystemName`, `Footer`, `Logo`, `TopUpLink`，主题管理，从环境变量加载的全局配置值。
- **`env.go`** -- 环境变量辅助：`GetEnvOrDefault()` 用于 int/string/bool，带安全解析和默认回退。
- **`init.go`** -- 应用初始化：CLI 标志解析（`--port`, `--version`, `--help`, `--log-dir`）、日志设置（带日志轮转）、HTTP 客户端配置、从 `.env` 文件加载环境变量。

### 2.10 安全
- **`totp.go`** -- TOTP 2FA：密钥生成、二维码数据、代码验证、备份码生成、带锁定的尝试限制。
- **`verification.go`** -- 邮件验证码系统：带过期的内存映射、代码生成、存储、验证，用于邮件验证和密码重置流程。

### 2.11 监控
- **`system_monitor.go`** -- 系统资源监控器：定期收集 CPU、内存、磁盘使用率到 `atomic.Value` 用于按需状态查询。
- **`system_monitor_unix.go`** -- Unix/macOS 磁盘空间信息（`unix.Statfs`）。
- **`system_monitor_windows.go`** -- Windows 磁盘空间信息（`kernel32.dll GetDiskFreeSpaceExW`）。
- **`pprof.go`** -- CPU 分析守护进程：通过 `gopsutil` 监控 CPU 使用率，CPU 超过 80% 时转储 pprof 文件。
- **`pyro.go`** -- Pyroscope 持续分析集成：从环境变量配置运行时 mutex/block 分析率并启动 Pyroscope 代理。
- **`performance_config.go`** -- 原子值的 `PerformanceMonitorConfig` 结构体（CPU/内存/磁盘阈值）带 getter/setter。

### 2.12 模型辅助
- **`model.go`** -- 模型名称分类辅助：硬编码的 OpenAI 响应模型列表、图像生成模型和 OpenAI 文本模型，带前缀/后缀匹配函数。
- **`api_type.go`** -- 将渠道类型（如 `ChannelTypeOpenAI`, `ChannelTypeAnthropic`）映射到 API 类型（`APITypeOpenAI`, `APITypeAnthropic` 等）。
- **`endpoint_defaults.go`** -- 将 `EndpointType` 常量映射到默认上游路径和 HTTP 方法。
- **`endpoint_type.go`** -- `GetEndpointTypesByChannelType()` 返回给定渠道类型支持的端点类型有序列表。

### 2.13 并发
- **`go-channel.go`** -- 安全通道发送/接收包装器（`SafeSendBool`, `SafeSendString`, `SafeStopChan`），防止因关闭通道引起的 panic。
- **`gopool.go`** -- 初始化用于中转请求处理的 `gopool` 协程池（`RelayCtxGo`），带 panic 处理器。

### 2.14 分页
- **`page_info.go`** -- 通用分页结构体（`PageInfo`）带 `Page`/`PageSize`/`Total`/`Items` 和 Gin 上下文绑定辅助。

### 2.15 深拷贝
- **`copy.go`** -- 通用深拷贝工具 `DeepCopy[T]()`，使用 `copier` 库的 `DeepCopy: true` 和 `IgnoreEmpty: true` 选项。

### 2.16 SSE
- **`custom-event.go`** -- 自定义服务器发送事件（SSE）写入器，带流式 JSON 编码器。

### 2.17 日志
- **`sys_log.go`** -- 系统日志函数 `SysLog()` 和 `SysError()`，将带时间戳的消息写入 gin 的默认写入器，带互斥锁保护。

### 2.18 配额
- **`quota.go`** -- `GetTrustQuota()` 返回 `10 * QuotaPerUnit` 作为默认信任级别配额。
- **`topup-ratio.go`** -- 每组充值比例管理：线程安全的组到比例映射，带 JSON 序列化和原子更新。

---

## 关联文件列表

### 公共工具层核心文件
- `common/json.go` - JSON 封装
- `common/str.go` - 字符串工具
- `common/utils.go` - 通用工具
- `common/validate.go` - 结构体验证
- `common/crypto.go` - 加密工具
- `common/hash.go` - 哈希工具
- `common/database.go` - 数据库工具
- `common/redis.go` - Redis 工具
- `common/gin.go` - Gin 工具
- `common/ssrf_protection.go` - SSRF 防护
- `common/url_validator.go` - URL 验证
- `common/ip.go` - IP 工具
- `common/email.go` - 邮件发送
- `common/email-outlook-auth.go` - Outlook 认证
- `common/body_storage.go` - 请求体存储
- `common/disk_cache.go` - 磁盘缓存
- `common/disk_cache_config.go` - 磁盘缓存配置
- `common/embed-file-system.go` - 嵌入式文件系统
- `common/rate-limit.go` - 内存限流器
- `common/limiter/limiter.go` - Redis 限流器
- `common/constants.go` - 常量定义
- `common/env.go` - 环境变量工具
- `common/init.go` - 应用初始化
- `common/totp.go` - TOTP 2FA
- `common/verification.go` - 验证码系统
- `common/system_monitor.go` - 系统监控
- `common/system_monitor_unix.go` - Unix 磁盘信息
- `common/system_monitor_windows.go` - Windows 磁盘信息
- `common/pprof.go` - CPU 分析
- `common/pyro.go` - Pyroscope 集成
- `common/performance_config.go` - 性能配置
- `common/model.go` - 模型分类
- `common/api_type.go` - API 类型映射
- `common/endpoint_defaults.go` - 端点默认值
- `common/endpoint_type.go` - 端点类型
- `common/go-channel.go` - 安全通道
- `common/gopool.go` - 协程池
- `common/page_info.go` - 分页工具
- `common/copy.go` - 深拷贝
- `common/custom-event.go` - SSE 工具
- `common/sys_log.go` - 系统日志
- `common/quota.go` - 配额工具
- `common/topup-ratio.go` - 充值比例

### 依赖的外部库
- `github.com/gin-gonic/gin` - Web 框架
- `github.com/go-redis/redis/v8` - Redis 客户端
- `github.com/jinzhu/copier` - 深拷贝库
- `github.com/go-playground/validator` - 结构体验证
- `github.com/shirou/gopsutil` - 系统监控
- `github.com/pquerna/otp` - TOTP 库
- `gopkg.in/gomail.v2` - 邮件发送
