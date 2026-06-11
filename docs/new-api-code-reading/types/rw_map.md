# rw_map.go 代码阅读文档

## 1. 全局概述

本文件实现了一个线程安全的泛型读写 Map（RWMap），基于 Go 的 `sync.RWMutex` 实现，支持并发读写操作。同时提供了 JSON 序列化/反序列化支持和便捷的加载函数。

## 2. 依赖关系

- `sync` — Go 标准库同步包
- `github.com/QuantumNous/new-api/common` — JSON 序列化封装

## 3. 类型定义

### RWMap 泛型结构体

```go
type RWMap[K comparable, V any] struct {
    data  map[K]V
    mutex sync.RWMutex
}
```

- `K comparable` — 键类型约束
- `V any` — 值类型无约束
- 使用读写锁实现并发安全

## 4. 函数详情

### NewRWMap

```go
func NewRWMap[K comparable, V any]() *RWMap[K, V]
```

创建并返回一个新的空 RWMap。

### Get

```go
func (m *RWMap[K, V]) Get(key K) (V, bool)
```

获取指定键的值。使用读锁，支持并发读取。

### Set

```go
func (m *RWMap[K, V]) Set(key K, value V)
```

设置指定键的值。使用写锁，阻塞其他读写操作。

### AddAll

```go
func (m *RWMap[K, V]) AddAll(other map[K]V)
```

批量添加键值对。使用写锁。

### Clear

```go
func (m *RWMap[K, V]) Clear()
```

清空 Map 中的所有数据。使用写锁。

### ReadAll

```go
func (m *RWMap[K, V]) ReadAll() map[K]V
```

返回 Map 的完整副本。使用读锁，返回的是深拷贝，修改副本不影响原 Map。

### Len

```go
func (m *RWMap[K, V]) Len() int
```

返回 Map 中键值对的数量。使用读锁。

### UnmarshalJSON

```go
func (m *RWMap[K, V]) UnmarshalJSON(b []byte) error
```

从 JSON 数据反序列化到 RWMap。使用写锁，会覆盖现有数据。

### MarshalJSON

```go
func (m *RWMap[K, V]) MarshalJSON() ([]byte, error)
```

将 RWMap 序列化为 JSON 数据。使用读锁。

### MarshalJSONString

```go
func (m *RWMap[K, V]) MarshalJSONString() string
```

将 RWMap 序列化为 JSON 字符串。序列化失败时返回 `"{}"`。

### LoadFromJsonString

```go
func LoadFromJsonString[K comparable, V any](m *RWMap[K, V], jsonStr string) error
```

从 JSON 字符串加载数据到 RWMap。使用写锁，会覆盖现有数据。

### LoadFromJsonStringWithCallback

```go
func LoadFromJsonStringWithCallback[K comparable, V any](m *RWMap[K, V], jsonStr string, onSuccess func()) error
```

从 JSON 字符串加载数据到 RWMap，成功时调用回调函数。

## 5. 关键逻辑分析

### 并发安全设计

- 所有读操作使用 `RLock()`，允许多个并发读
- 所有写操作使用 `Lock()`，独占写锁
- `ReadAll()` 返回深拷贝，避免锁泄露

### JSON 序列化

- 使用 `common.Marshal` 和 `common.Unmarshal` 进行 JSON 操作（遵循项目规范）
- `UnmarshalJSON` 在反序列化前会清空现有数据
- `MarshalJSON` 在序列化时加读锁，确保数据一致性

### 使用场景

在系统中广泛用于需要并发访问的配置数据，如：
- 模型配置映射
- 渠道配置缓存
- 动态设置存储

## 6. 相关文件

- `types/set.go` — 泛型 Set 实现
- `common/json.go` — JSON 序列化封装
- `setting/` — 配置管理中广泛使用 RWMap
