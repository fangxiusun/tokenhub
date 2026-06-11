# topup-ratio.go 代码阅读文档

## 1. 全局总结

`topup-ratio.go` 是充值比例管理文件，提供了用户组充值比例的配置和查询功能。文件维护一个全局的 `topupGroupRatio` 映射表，存储不同用户组（如 default、vip、svip）的充值倍率。通过 `sync.RWMutex` 保证并发安全，支持通过 JSON 字符串动态更新配置，并提供按用户组名称查询充值比例的接口。

## 2. 依赖关系

### 标准库依赖
- `encoding/json` — JSON 序列化/反序列化，用于配置的导入导出
- `sync` — 提供 `sync.RWMutex` 读写锁

### 项目内部依赖
- `SysError()` — 系统错误日志输出函数（来自 `sys_log.go`）

## 3. 类型定义

本文件未定义新的类型，使用了以下内置类型：

- `map[string]float64` — 用户组名称到充值比例的映射
- `sync.RWMutex` — 读写锁，保护映射表的并发访问

## 4. 全局变量

### 4.1 `topupGroupRatio`

```go
var topupGroupRatio = map[string]float64{
    "default": 1,
    "vip":     1,
    "svip":    1,
}
```

**功能**：存储用户组充值比例的映射表。

**默认值**：所有用户组的充值比例默认为 1（即 1:1 充值）。

**预定义用户组**：
- `default` — 普通用户组
- `vip` — VIP 用户组
- `svip` — SVIP 用户组

### 4.2 `topupGroupRatioMutex`

```go
var topupGroupRatioMutex sync.RWMutex
```

**功能**：保护 `topupGroupRatio` 映射表的并发安全。

## 5. 函数详解

### 5.1 `TopupGroupRatio2JSONString() string`

**功能**：将当前充值比例配置导出为 JSON 字符串。

**行为**：
1. 获取读锁 `topupGroupRatioMutex.RLock()`
2. 使用 `json.Marshal()` 将映射表序列化为 JSON
3. 释放读锁（通过 `defer`）
4. 如果序列化失败，调用 `SysError()` 记录错误
5. 返回 JSON 字符串

**返回值**：JSON 格式的充值比例配置，例如 `{"default":1,"svip":1,"vip":1}`。

**使用场景**：配置导出、API 响应、配置备份等。

### 5.2 `UpdateTopupGroupRatioByJSONString(jsonStr string) error`

**功能**：通过 JSON 字符串更新充值比例配置。

**行为**：
1. 获取写锁 `topupGroupRatioMutex.Lock()`
2. 创建新的空映射表 `make(map[string]float64)`
3. 使用 `json.Unmarshal()` 将 JSON 字符串解析到新映射表
4. 释放写锁（通过 `defer`）
5. 返回错误（如果有）

**参数**：`jsonStr` — JSON 格式的充值比例配置字符串。

**返回值**：`error` — 解析失败时返回错误，成功时返回 `nil`。

**注意事项**：
- 此函数会完全替换现有配置，而不是合并
- 使用 `make(map[string]float64)` 创建新映射表，避免修改旧数据

### 5.3 `GetTopupGroupRatio(name string) float64`

**功能**：根据用户组名称获取充值比例。

**行为**：
1. 获取读锁 `topupGroupRatioMutex.RLock()`
2. 在映射表中查找指定用户组
3. 释放读锁（通过 `defer`）
4. 如果找到，返回对应比例
5. 如果未找到，调用 `SysError()` 记录错误，返回默认值 1

**参数**：`name` — 用户组名称。

**返回值**：`float64` — 充值比例。未找到时返回默认值 1。

**使用场景**：计费系统计算用户充值金额时，根据用户组获取对应的充值倍率。

## 6. 关键逻辑分析

### 6.1 读写锁并发控制

`topupGroupRatioMutex` 使用 `sync.RWMutex` 保护共享映射表：

- **读锁（RLock）**：`TopupGroupRatio2JSONString()` 和 `GetTopupGroupRatio()` 使用读锁，允许多个读操作并发执行
- **写锁（Lock）**：`UpdateTopupGroupRatioByJSONString()` 使用写锁，确保更新操作独占访问

这种设计在读多写少的场景下提供了良好的并发性能。

### 6.2 配置更新策略

`UpdateTopupGroupRatioByJSONString()` 采用完全替换策略：
1. 创建新的空映射表 `make(map[string]float64)`
2. 将 JSON 解析到新映射表
3. 直接替换全局变量 `topupGroupRatio`

这种策略的优点：
- 简单清晰，避免复杂的合并逻辑
- 保证配置的原子性更新
- 避免部分更新导致的不一致状态

### 6.3 默认值处理

`GetTopupGroupRatio()` 在找不到指定用户组时返回默认值 1：
- 这确保了计费系统不会因为配置缺失而崩溃
- 用户组名称区分大小写（Go map 的特性）
- 错误日志帮助管理员发现配置问题

### 6.4 JSON 序列化错误处理

`TopupGroupRatio2JSONString()` 在序列化失败时：
- 调用 `SysError()` 记录错误（不会导致程序崩溃）
- 返回空字符串 `""`（调用方需处理空字符串情况）

### 6.5 配置数据的原子性

通过 `atomic.Value` 或整体替换的方式，保证配置更新的原子性：
- 读操作总是看到一致的配置状态
- 写操作是原子的，不会出现部分更新的中间状态

## 7. 关联文件

- **`common/sys_log.go`** — 提供 `SysError()` 函数，用于记录配置相关错误
- **`setting/setting.go`**（推测）— 可能包含充值比例的持久化配置和管理接口
- **`model/user.go`**（推测）— 用户模型可能包含用户组字段，用于关联充值比例
- **`controller/topup.go`**（推测）— 充值控制器，可能调用 `GetTopupGroupRatio()` 计算充值金额
- **`middleware/auth.go`**（推测）— 认证中间件，可能解析用户组信息
