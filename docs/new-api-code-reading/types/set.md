# set.go 代码阅读文档

## 1. 全局概述

本文件实现了一个泛型集合（Set）数据结构，基于 Go 泛型和 map 实现，提供元素的添加、删除、检查和遍历功能。

## 2. 依赖关系

本文件无外部依赖，仅使用 Go 内置的泛型特性。

## 3. 类型定义

### Set 泛型结构体

```go
type Set[T comparable] struct {
    items map[T]struct{}
}
```

- `T comparable` — 泛型约束，要求元素类型支持比较操作
- 内部使用 `map[T]struct{}` 实现，利用空结构体不占内存的特性

## 4. 函数详情

### NewSet

```go
func NewSet[T comparable]() *Set[T]
```

创建并返回一个新的空 Set。

### Add

```go
func (s *Set[T]) Add(item T)
```

向 Set 中添加一个元素。如果元素已存在，则覆盖（map 语义）。

### Remove

```go
func (s *Set[T]) Remove(item T)
```

从 Set 中移除一个元素。如果元素不存在，操作无效。

### Contains

```go
func (s *Set[T]) Contains(item T) bool
```

检查 Set 是否包含某个元素。

### Len

```go
func (s *Set[T]) Len() int
```

返回 Set 中元素的数量。

### Items

```go
func (s *Set[T]) Items() []T
```

返回 Set 中所有元素组成的切片。由于 map 的无序性，返回的切片元素顺序是随机的。

## 5. 关键逻辑分析

### 实现特点

- **泛型实现**：使用 Go 1.18+ 的泛型特性，支持任意可比较类型
- **内存高效**：使用 `map[T]struct{}` 而非 `map[T]bool`，空结构体不占额外内存
- **线程不安全**：未使用锁机制，需要调用方自行处理并发访问
- **无序遍历**：`Items()` 方法返回的切片顺序不确定

### 使用场景

在系统中用于需要去重的场景，如：
- 模型名称集合
- 渠道标签集合
- 权限标识集合

## 6. 相关文件

- `types/rw_map.go` — 线程安全的 Map 实现
- `model/` — 数据模型层使用 Set 进行数据去重
