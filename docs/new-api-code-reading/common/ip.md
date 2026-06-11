# ip.go 代码阅读文档

## 1. 全局总结

`ip.go` 是 IP 地址处理工具文件，提供了 IP 地址验证、解析、私有 IP 检测和 CIDR 匹配等功能。这些工具函数主要用于网络安全相关的场景，如访问控制、IP 白名单/黑名单、SSRF 防护等。

## 2. 依赖关系

### 标准库依赖
- `net` - 网络地址解析和处理

### 项目内部依赖
- 无

## 3. 类型定义

该文件未定义新的类型，主要操作 `net.IP` 类型。

## 4. 函数详解

### IsIP(s string) bool
```go
func IsIP(s string) bool
```
判断字符串是否为有效的 IP 地址。

**参数：**
- `s` - 待验证的字符串

**返回值：**
- `bool` - 是否为有效 IP 地址

**实现逻辑：**
- 使用 `net.ParseIP()` 解析字符串
- 返回解析结果是否为 `nil`

### ParseIP(s string) net.IP
```go
func ParseIP(s string) net.IP
```
解析字符串为 `net.IP` 类型。

**参数：**
- `s` - IP 地址字符串

**返回值：**
- `net.IP` - 解析后的 IP 对象，无效时返回 `nil`

### IsPrivateIP(ip net.IP) bool
```go
func IsPrivateIP(ip net.IP) bool
```
判断 IP 地址是否为私有/保留地址。

**参数：**
- `ip` - 待检测的 IP 对象

**返回值：**
- `bool` - 是否为私有 IP 地址

**检测范围：**
1. 回环地址 (`IsLoopback()`)
2. 链路本地单播 (`IsLinkLocalUnicast()`)
3. 链路本地多播 (`IsLinkLocalMulticast()`)
4. 私有网络范围：
   - `10.0.0.0/8` - A 类私有网络
   - `172.16.0.0/12` - B 类私有网络
   - `192.168.0.0/16` - C 类私有网络

### IsIpInCIDRList(ip net.IP, cidrList []string) bool
```go
func IsIpInCIDRList(ip net.IP, cidrList []string) bool
```
判断 IP 地址是否在指定的 CIDR 列表中。

**参数：**
- `ip` - 待检测的 IP 对象
- `cidrList` - CIDR 表示的网络列表

**返回值：**
- `bool` - IP 是否在列表中

**实现逻辑：**
1. 遍历 CIDR 列表
2. 尝试解析为 CIDR 格式
3. 如果解析失败，尝试作为单个 IP 解析
4. 使用 `Contains()` 方法检测 IP 是否在网络范围内

## 5. 关键逻辑分析

### 私有 IP 地址范围

| 地址类别 | CIDR 表示 | 地址范围 |
|---------|-----------|----------|
| 回环地址 | - | 127.0.0.0/8 |
| A 类私有 | 10.0.0.0/8 | 10.0.0.0 - 10.255.255.255 |
| B 类私有 | 172.16.0.0/12 | 172.16.0.0 - 172.31.255.255 |
| C 类私有 | 192.168.0.0/16 | 192.168.0.0 - 192.168.255.255 |
| 链路本地 | - | 169.254.0.0/16 |

### CIDR 列表匹配容错处理
```go
if err != nil {
    // 尝试作为单个IP处理
    if whitelistIP := net.ParseIP(cidr); whitelistIP != nil {
        if ip.Equal(whitelistIP) {
            return true
        }
    }
    continue
}
```
- 支持混合格式：CIDR 表示和单个 IP 地址可以混合在同一列表中
- 无效的 CIDR 条目会被静默跳过，不会导致整个列表失效

## 6. 关联文件

- `common/security.go` - 可能使用 IP 验证函数进行访问控制
- `middleware/auth.go` - 可能使用 IP 白名单进行鉴权
- `controller/ip.go` - 如果存在，可能提供 IP 相关的 API 接口
- `setting/operation.go` - 可能存储 IP 黑白名单配置
