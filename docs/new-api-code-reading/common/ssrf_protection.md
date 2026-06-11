# ssrf_protection.go 代码阅读文档

## 1. 全局总结

本文件实现了 SSRF（Server-Side Request Forgery，服务端请求伪造）防护功能。提供了完整的 URL 安全验证机制，包括：协议限制（仅允许 HTTP/HTTPS）、私有 IP 地址检测、域名黑白名单过滤、IP 黑白名单过滤、端口范围控制，以及域名解析后的 IP 二次验证。支持灵活的配置组合，可适应不同的安全需求场景。

## 2. 依赖关系

| 包名 | 用途 |
|------|------|
| `fmt` | 格式化错误信息 |
| `net` | IP 地址解析、CIDR 匹配、DNS 解析 |
| `net/url` | URL 解析 |
| `strconv` | 字符串与数值转换 |
| `strings` | 字符串操作（前缀/后缀匹配等） |

本文件是独立的安全工具组件，不依赖项目内其他模块。

## 3. 类型定义

### SSRFProtection

```go
type SSRFProtection struct {
    AllowPrivateIp         bool     // 是否允许访问私有 IP
    DomainFilterMode       bool     // true: 白名单模式, false: 黑名单模式
    DomainList             []string // 域名过滤列表（支持通配符 *.example.com）
    IpFilterMode           bool     // true: 白名单模式, false: 黑名单模式
    IpList                 []string // IP 过滤列表（支持 CIDR 格式）
    AllowedPorts           []int    // 允许的端口列表
    ApplyIPFilterForDomain bool     // 是否对域名启用 IP 过滤（DNS 解析后验证）
}
```

### DefaultSSRFProtection

```go
var DefaultSSRFProtection = &SSRFProtection{
    AllowPrivateIp:   false,
    DomainFilterMode: true,
    DomainList:       []string{},
    IpFilterMode:     true,
    IpList:           []string{},
    AllowedPorts:     []int{},
}
```

默认配置：禁止私有 IP，白名单模式（空列表意味着拒绝所有）。

### 私有地址网段

```go
var privateIPv4Nets []net.IPNet  // IPv4 私有/保留/特殊用途网段
var privateIPv6Nets []net.IPNet  // IPv6 私有/保留/特殊用途网段
```

包含 IANA 标准定义的所有私有和保留地址段：
- **IPv4**：0.0.0.0/8、10.0.0.0/8、100.64.0.0/10、127.0.0.0/8、169.254.0.0/16、172.16.0.0/12、192.0.0.0/24、192.0.2.0/24、192.168.0.0/16、198.18.0.0/15、198.51.100.0/24、203.0.113.0/24、224.0.0.0/4、240.0.0.0/4、255.255.255.255/32
- **IPv6**：未指定地址、回环、IPv4 映射、文档地址、ULA、链路本地、组播等

## 4. 函数详解

### isPrivateIP(ip net.IP) bool

```go
func isPrivateIP(ip net.IP) bool
```

检查 IP 是否为私有/保留/特殊用途地址。

**检查顺序：**
1. nil 检查 → 返回 true
2. 未指定地址（0.0.0.0, ::）→ 返回 true
3. 回环地址 → 返回 true
4. 链路本地单播/组播 → 返回 true
5. 接口本地组播（IPv6）→ 返回 true
6. IPv4：遍历 `privateIPv4Nets` 检查是否在私有网段
7. IPv6：遍历 `privateIPv6Nets` 检查是否在私有网段
8. 兜底：使用 Go 标准库 `ip.IsPrivate()` 检查

### parsePortRanges(portConfigs []string) ([]int, error)

```go
func parsePortRanges(portConfigs []string) ([]int, error)
```

解析端口范围配置字符串。

**支持格式：**
- 单个端口：`"80"`、`"443"`
- 端口范围：`"8000-9000"`
- 混合：`["80", "443", "8000-9000"]`

**验证规则：**
- 端口号必须在 1-65535 范围内
- 范围起始端口不能大于结束端口
- 空字符串会被跳过

### (p *SSRFProtection) isAllowedPort(port int) bool

```go
func (p *SSRFProtection) isAllowedPort(port int) bool
```

检查端口是否在允许列表中。如果列表为空，允许所有端口。

### isDomainListed(domain string, list []string) bool

```go
func isDomainListed(domain string, list []string) bool
```

检查域名是否在指定列表中。

**匹配模式：**
- 精确匹配：`example.com`
- 通配符匹配：`*.example.com` 匹配 `sub.example.com` 和 `example.com`

### (p *SSRFProtection) isDomainAllowed(domain string) bool

```go
func (p *SSRFProtection) isDomainAllowed(domain string) bool
```

根据过滤模式判断域名是否允许：
- 白名单模式：域名必须在列表中
- 黑名单模式：域名不能在列表中

### isIPListed(ip net.IP, list []string) bool

```go
func isIPListed(ip net.IP, list []string) bool
```

检查 IP 是否在指定列表中，使用 CIDR 匹配。

### (p *SSRFProtection) IsIPAccessAllowed(ip net.IP) bool

```go
func (p *SSRFProtection) IsIPAccessAllowed(ip net.IP) bool
```

判断 IP 是否允许访问：
1. 如果是私有 IP 且不允许私有 IP → 拒绝
2. 根据 IP 过滤模式（白名单/黑名单）判断

### (p *SSRFProtection) ValidateURL(urlStr string) error

```go
func (p *SSRFProtection) ValidateURL(urlStr string) error
```

核心 URL 验证方法，完整的验证流程：

1. **协议检查**：仅允许 http/https
2. **端口解析**：提取主机和端口，默认端口为 80（HTTP）或 443（HTTPS）
3. **端口验证**：检查端口是否在允许列表中
4. **IP 直连判断**：如果 host 是 IP 地址，直接进行 IP 访问控制检查
5. **域名过滤**：检查域名是否在域名列表中
6. **IP 解析验证**：如果启用了 `ApplyIPFilterForDomain`，解析域名对应的所有 IP 并逐一验证

### ValidateURLWithFetchSetting(...) error

```go
func ValidateURLWithFetchSetting(urlStr string, enableSSRFProtection bool, allowPrivateIp bool, 
    domainFilterMode bool, ipFilterMode bool, domainList, ipList, allowedPorts []string, 
    applyIPFilterForDomain bool) error
```

使用 FetchSetting 配置参数验证 URL 的便捷函数：
1. 如果 SSRF 防护未启用，直接返回成功
2. 解析端口范围配置
3. 构造 `SSRFProtection` 实例
4. 调用 `ValidateURL` 执行验证

## 5. 关键逻辑分析

### 5.1 多层防护策略

验证采用分层检查策略，每层独立且可配置：
```
URL 解析 → 协议检查 → 端口检查 → IP/域名过滤 → DNS 解析验证
```

### 5.2 白名单/黑名单双模式

域名和 IP 过滤都支持两种模式：
- **白名单模式**（`FilterMode = true`）：只有列表中的才允许，空列表 = 拒绝所有
- **黑名单模式**（`FilterMode = false`）：列表中的被拒绝，空列表 = 允许所有

### 5.3 DNS 重绑定防护

`ApplyIPFilterForDomain` 选项提供了 DNS 重绑定攻击防护：
- 解析域名获取所有关联 IP
- 逐一检查每个 IP 是否通过访问控制
- 防止攻击者通过 DNS 动态解析绕过 IP 过滤

### 5.4 通配符域名匹配

支持 `*.example.com` 格式的通配符：
- 匹配所有子域名：`sub.example.com`、`a.b.example.com`
- 也匹配裸域名本身：`example.com`
- 大小写不敏感

### 5.5 私有地址全覆盖

`isPrivateIP` 函数覆盖了完整的私有/保留地址范围：
- IPv4 和 IPv6 双栈支持
- 包含所有 IANA 定义的特殊用途地址
- 兜底使用 Go 标准库的 `IsPrivate()` 方法

### 5.6 潜在问题

1. **DNS 解析时序**：域名解析在验证时执行，可能存在 TOCTOU（Time-of-Check-Time-of-Use）问题
2. **性能开销**：DNS 解析和多次 IP 检查在高并发场景下可能有性能影响
3. **IPv6 支持**：私有 IPv6 地址列表可能需要根据实际部署环境调整

## 6. 关联文件

- `new-api/middleware/fetch/` — 可能使用本防护机制验证外部请求
- `new-api/setting/` — FetchSetting 配置管理
- `new-api/relay/` — AI API 中继，可能使用本防护机制
