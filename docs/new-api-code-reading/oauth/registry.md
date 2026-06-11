# registry.go 代码阅读文档

## 1. 全局概述

本文件实现了 OAuth 提供者的注册表（Registry），提供提供者的注册、查询、加载和管理功能。支持内置提供者和自定义提供者两种类型。

## 2. 依赖关系

- `fmt` — 格式化输出
- `sync` — 同步原语
- `github.com/QuantumNous/new-api/common` — 通用工具
- `github.com/QuantumNous/new-api/model` — 数据模型

## 3. 类型定义

本文件无自定义类型定义。

## 4. 函数详情

### Register

```go
func Register(name string, provider Provider)
```

注册一个 OAuth 提供者。用于内置提供者的初始化（如 GitHub、Discord）。

### RegisterCustom

```go
func RegisterCustom(name string, provider Provider)
```

注册一个自定义 OAuth 提供者。自定义提供者可以被注销。

### Unregister

```go
func Unregister(name string)
```

从注册表中移除一个提供者。

### GetProvider

```go
func GetProvider(name string) Provider
```

根据名称获取 OAuth 提供者。

### GetAllProviders

```go
func GetAllProviders() map[string]Provider
```

获取所有已注册的 OAuth 提供者（返回副本）。

### GetEnabledCustomProviders

```go
func GetEnabledCustomProviders() []*GenericOAuthProvider
```

获取所有已启用的自定义 OAuth 提供者。

### IsProviderRegistered

```go
func IsProviderRegistered(name string) bool
```

检查提供者是否已注册。

### IsCustomProvider

```go
func IsCustomProvider(name string) bool
```

检查提供者是否为自定义提供者。

### LoadCustomProviders

```go
func LoadCustomProviders() error
```

从数据库加载所有自定义 OAuth 提供者：
1. 注销所有现有自定义提供者
2. 从数据库查询所有自定义提供者配置
3. 为每个配置创建 GenericOAuthProvider 并注册

### ReloadCustomProviders

```go
func ReloadCustomProviders() error
```

重新加载所有自定义 OAuth 提供者（调用 `LoadCustomProviders`）。

### RegisterOrUpdateCustomProvider

```go
func RegisterOrUpdateCustomProvider(config *model.CustomOAuthProvider)
```

注册或更新单个自定义提供者。

### UnregisterCustomProvider

```go
func UnregisterCustomProvider(slug string)
```

根据 slug 注销自定义提供者。

## 5. 关键逻辑分析

### 注册表结构

```go
var (
    providers = make(map[string]Provider)      // 名称 → 提供者
    customProviderSlugs = make(map[string]bool) // 自定义提供者 slug 集合
)
```

- `providers` 存储所有已注册的提供者
- `customProviderSlugs` 跟踪哪些是自定义提供者（可注销）

### 线程安全

- 使用 `sync.RWMutex` 保护注册表的并发访问
- 读操作使用 `RLock()`，写操作使用 `Lock()`

### 内置 vs 自定义提供者

| 特性 | 内置提供者 | 自定义提供者 |
|------|-----------|-------------|
| 注册方式 | `Register` | `RegisterCustom` |
| 可注销 | 否 | 是 |
| 配置来源 | 代码中硬编码 | 数据库动态配置 |
| 实现 | 独立文件 | `GenericOAuthProvider` |

### 自定义提供者加载流程

```
LoadCustomProviders()
  → 清空现有自定义提供者
  → model.GetAllCustomOAuthProviders()
  → 遍历配置
    → NewGenericOAuthProvider(config)
    → RegisterCustom(config.Slug, provider)
  → 日志记录加载数量
```

## 6. 相关文件

- `oauth/provider.go` — Provider 接口定义
- `oauth/generic.go` — GenericOAuthProvider 实现
- `oauth/github.go` — GitHub 提供者（init 中注册）
- `oauth/discord.go` — Discord 提供者（init 中注册）
- `oauth/oidc.go` — OIDC 提供者（init 中注册）
- `oauth/linuxdo.go` — LinuxDO 提供者（init 中注册）
- `model/custom_oauth.go` — 自定义 OAuth 提供者数据模型
