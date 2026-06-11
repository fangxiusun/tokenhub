# constants.go 代码阅读文档

## 1. 全局总结

`constants.go` 是系统全局常量和配置变量的集中定义文件，包含了应用运行时所需的各类配置参数、状态常量、功能开关和全局设置。这些变量在系统启动时初始化，部分通过配置文件或环境变量动态更新。该文件是理解系统整体配置和行为的关键入口。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `crypto/tls` | TLS 配置（`InsecureTLSConfig`） |
| `strings` | 字符串前缀匹配（`ThemeAwarePath`） |
| `sync` | 读写锁（`OptionMapRWMutex`） |
| `sync/atomic` | 原子操作（`themeValue`） |
| `time` | 时间相关（`StartTime`、`RateLimitKeyExpirationDuration` 等） |
| `github.com/google/uuid` | 生成唯一标识（`SessionSecret`、`CryptoSecret`） |

## 3. 类型定义

本文件无自定义类型定义，主要使用基础类型和标准库类型。

## 4. 函数详解

### `init()`

- **功能**：初始化 `themeValue` 为 "classic"
- **调用时机**：包首次加载时自动执行

### `GetTheme() string`

- **功能**：获取当前前端主题名称
- **返回值**：`"default"` 或 `"classic"`
- **线程安全**：使用 `atomic.Value` 保证并发读取安全

### `SetTheme(t string)`

- **功能**：设置前端主题
- **参数**：`t string` — 主题名称，仅接受 `"default"` 或 `"classic"`
- **行为**：不接受的值会被静默忽略，不会更新

### `ThemeAwarePath(suffix string) string`

- **功能**：根据当前主题重写路由路径
- **参数**：`suffix string` — 原始路径后缀
- **返回值**：重写后的路径或原始路径
- **路径映射规则**（仅当主题为 `"default"` 时）：
  - `/console/topup` → `/wallet`
  - `/console/log` → `/usage-logs`
  - `/console/personal` → `/profile`
- **设计目的**：兼容 classic 主题的 `/console/*` 路径和 default 主题的新路径

### `IsValidateRole(role int) bool`

- **功能**：验证角色值是否有效
- **有效角色**：`RoleGuestUser(0)`、`RoleCommonUser(1)`、`RoleAdminUser(10)`、`RoleRootUser(100)`

## 5. 关键逻辑分析

### 系统元信息变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `StartTime` | `time.Now().Unix()` | 系统启动时间戳 |
| `Version` | `"v0.0.0"` | 版本号，构建时自动替换 |
| `SystemName` | `"New API"` | 系统名称 |
| `Footer` | `""` | 页脚文本 |
| `Logo` | `""` | Logo 路径 |
| `TopUpLink` | `""` | 充值链接 |

### 前端主题系统

- `themeValue` 使用 `atomic.Value` 存储，支持并发安全的读写
- 默认主题为 `"classic"`，可通过 `SetTheme` 切换
- `ThemeAwarePath` 在请求处理层透明地重写路径，对上层业务代码无感

### 配额与计量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `QuotaPerUnit` | `500 * 1000.0` | 每单位配额对应 $0.002 / 1K tokens |
| `DisplayInCurrencyEnabled` | `true` | 是否以货币形式显示配额 |
| `DisplayTokenStatEnabled` | `true` | 是否显示 Token 统计 |

### 用户注册配置

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `PasswordLoginEnabled` | `true` | 密码登录开关 |
| `PasswordRegisterEnabled` | `true` | 密码注册开关 |
| `EmailVerificationEnabled` | `false` | 邮箱验证开关 |
| `GitHubOAuthEnabled` | `false` | GitHub OAuth 开关 |
| `LinuxDOOAuthEnabled` | `false` | LinuxDO OAuth 开关 |
| `WeChatAuthEnabled` | `false` | 微信认证开关 |
| `TelegramOAuthEnabled` | `false` | Telegram 认证开关 |
| `TurnstileCheckEnabled` | `false` | Cloudflare Turnstile 验证开关 |
| `RegisterEnabled` | `true` | 总注册开关 |

### 邮箱域名限制

- `EmailDomainRestrictionEnabled` — 是否启用邮箱域名限制
- `EmailAliasRestrictionEnabled` — 是否启用邮箱别名限制
- `EmailDomainWhitelist` — 默认白名单包含 9 个常见邮箱域名
- `EmailLoginAuthServerList` — 支持邮箱登录的 SMTP 服务器列表

### 密钥与安全

| 变量 | 说明 |
|------|------|
| `SessionSecret` | 会话密钥，启动时随机生成 UUID |
| `CryptoSecret` | 加密密钥，启动时随机生成 UUID |
| `TLSInsecureSkipVerify` | 是否跳过 TLS 证书验证 |
| `InsecureTLSConfig` | 不安全的 TLS 配置 |

### 邮件服务配置

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `SMTPServer` | `""` | SMTP 服务器地址 |
| `SMTPPort` | `587` | SMTP 端口 |
| `SMTPSSLEnabled` | `false` | 是否启用 SSL |
| `SMTPForceAuthLogin` | `false` | 是否强制认证登录 |

### 第三方 OAuth 配置

包含 GitHub、LinuxDO、WeChat、Telegram、Turnstile 的 Client ID/Secret 和 Token 等配置变量。

### 配额与限制配置

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `QuotaForNewUser` | `0` | 新用户初始配额 |
| `QuotaForInviter` | `0` | 邀请人奖励配额 |
| `QuotaForInvitee` | `0` | 被邀请人奖励配额 |
| `ChannelDisableThreshold` | `5.0` | 渠道自动禁用阈值 |
| `PreConsumedQuota` | `500` | 预消费配额 |
| `RetryTimes` | `0` | 重试次数 |

### 中继配置

| 变量 | 说明 |
|------|------|
| `RelayTimeout` | 中继超时（秒） |
| `RelayIdleConnTimeout` | 空闲连接超时（秒） |
| `RelayMaxIdleConns` | 最大空闲连接数 |
| `RelayMaxIdleConnsPerHost` | 每个主机最大空闲连接数 |
| `RequestInterval` | 请求间隔（Duration） |
| `SyncFrequency` | 同步频率（秒） |

### 速率限制配置

| 变量组 | 说明 |
|--------|------|
| `GlobalApiRateLimit*` | API 全局速率限制 |
| `GlobalWebRateLimit*` | Web 全局速率限制 |
| `CriticalRateLimit*` | 关键操作速率限制（默认 20 次/20 分钟） |
| `UploadRateLimit*` | 上传速率限制（默认 10 次/分钟） |
| `DownloadRateLimit*` | 下载速率限制（默认 10 次/分钟） |
| `SearchRateLimit*` | 搜索速率限制（默认 10 次/分钟） |
| `RateLimitKeyExpirationDuration` | 速率限制 Key 过期时间（20 分钟） |

### 状态常量

#### 用户状态
- `UserStatusEnabled = 1` — 启用
- `UserStatusDisabled = 2` — 禁用

#### Token 状态
- `TokenStatusEnabled = 1` — 启用
- `TokenStatusDisabled = 2` — 禁用
- `TokenStatusExpired = 3` — 过期
- `TokenStatusExhausted = 4` — 用尽

#### 兑换码状态
- `RedemptionCodeStatusEnabled = 1` — 启用
- `RedemptionCodeStatusDisabled = 2` — 禁用
- `RedemptionCodeStatusUsed = 3` — 已使用

#### 渠道状态
- `ChannelStatusUnknown = 0` — 未知
- `ChannelStatusEnabled = 1` — 启用
- `ChannelStatusManuallyDisabled = 2` — 手动禁用
- `ChannelStatusAutoDisabled = 3` — 自动禁用

#### 充值状态
- `TopUpStatusPending = "pending"` — 待处理
- `TopUpStatusSuccess = "success"` — 成功
- `TopUpStatusFailed = "failed"` — 失败
- `TopUpStatusExpired = "expired"` — 过期

### 请求 ID 常量

- `RequestIdKey = "X-Oneapi-Request-Id"` — 本系统请求 ID
- `UpstreamRequestIdKey = "X-Upstream-Request-Id"` — 上游请求 ID

## 6. 关联文件

| 文件 | 关联说明 |
|------|---------|
| `setting/system.go` | 系统设置管理，读写此处的配置变量 |
| `setting/general.go` | 通用设置，管理配额显示等配置 |
| `middleware/rate-limit.go` | 速率限制中间件，使用此处的速率限制变量 |
| `middleware/auth.go` | 认证中间件，使用角色常量进行权限判断 |
| `model/user.go` | 用户模型，使用用户状态常量 |
| `model/token.go` | Token 模型，使用 Token 状态常量 |
| `model/channel.go` | 渠道模型，使用渠道状态常量 |
| `router/web-router.go` | Web 路由，使用 `ThemeAwarePath` 进行路径重写 |
