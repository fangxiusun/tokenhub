# fetch_setting.go 代码阅读文档

## 1. 全局总结

该文件定义网络请求（Fetch）设置，包括 SSRF 防护、IP/域名白名单/黑名单、端口限制等安全配置。

## 2. 依赖关系

- `github.com/QuantumNous/new-api/setting/config` — 全局配置管理器

## 3. 类型定义

| 结构体 | 字段 | 类型 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `FetchSetting` | `EnableSSRFProtection` | `bool` | `true` | 是否启用 SSRF 防护 |
| | `AllowPrivateIp` | `bool` | `false` | 是否允许私有 IP |
| | `DomainFilterMode` | `bool` | `false` | 域名过滤模式（true=白名单） |
| | `IpFilterMode` | `bool` | `false` | IP 过滤模式（true=白名单） |
| | `DomainList` | `[]string` | `[]` | 域名列表 |
| | `IpList` | `[]string` | `[]` | IP 列表（CIDR 格式） |
| | `AllowedPorts` | `[]string` | `["80","443","8080","8443"]` | 允许的端口 |
| | `ApplyIPFilterForDomain` | `bool` | `true` | 对域名启用 IP 过滤 |

## 4. 函数详解

| 函数名 | 签名 | 说明 |
|--------|------|------|
| `GetFetchSetting` | `func GetFetchSetting() *FetchSetting` | 获取 Fetch 设置 |

## 5. 关键逻辑分析

- SSRF 防护默认开启，防止请求内网资源
- 域名支持通配符格式（如 `*.example.com`）
- IP 列表使用 CIDR 格式
- 默认允许端口：80、443、8080、8443

## 6. 关联文件

- `pkg/ionet/` — 网络请求实现
- `relay/handler.go` — 使用 Fetch 设置
