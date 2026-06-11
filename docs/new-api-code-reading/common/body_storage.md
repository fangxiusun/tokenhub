# body_storage.go 代码阅读文档

## 1. 全局总结

`body_storage.go` 实现了一个灵活的请求体存储系统，支持内存存储和磁盘存储两种模式。该系统是 API 网关中处理大请求体的关键基础设施，能够根据数据大小自动选择存储方式（内存或磁盘），并在内存压力大时提供磁盘缓存作为备选方案。主要服务于请求重试、日志记录等需要多次读取请求体的场景。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `bytes` | 内存存储的字节缓冲 |
| `fmt` | 错误格式化 |
| `io` | I/O 接口和工具函数 |
| `os` | 磁盘文件操作 |
| `sync` | 互斥锁保证并发安全 |
| `sync/atomic` | 原子操作管理关闭状态 |

同包内部依赖：
- `CreateDiskCacheFile` — 创建磁盘缓存文件
- `DiskCacheTypeBody` — 磁盘缓存类型常量
- `GetDiskCacheThresholdBytes` — 获取磁盘缓存阈值
- `IsDiskCacheEnabled` — 检查磁盘缓存是否启用
- `IsDiskCacheAvailable` — 检查磁盘缓存是否有足够空间
- `GetDiskCachePath` — 获取磁盘缓存路径
- `IncrementMemoryBuffers/DecrementMemoryBuffers` — 内存缓冲区计数
- `IncrementDiskFiles/DecrementDiskFiles` — 磁盘文件计数
- `IncrementDiskCacheHits/IncrementMemoryCacheHits` — 缓存命中计数
- `CleanupOldDiskCacheFiles` — 清理旧缓存文件
- `ErrRequestBodyTooLarge/IsRequestBodyTooLargeError` — 请求体过大错误

## 3. 类型定义

### `BodyStorage` 接口

```go
type BodyStorage interface {
    io.ReadSeeker
    io.Closer
    Bytes() ([]byte, error)
    Size() int64
    IsDisk() bool
}
```

请求体存储的核心接口，继承 `io.ReadSeeker` 和 `io.Closer`，额外提供：
- `Bytes()` — 获取全部内容的字节切片
- `Size()` — 获取数据大小
- `IsDisk()` — 判断是否为磁盘存储

### `memoryStorage` 结构体（未导出）

```go
type memoryStorage struct {
    data   []byte        // 存储的数据
    reader *bytes.Reader // 内存读取器
    size   int64         // 数据大小
    closed int32         // 关闭状态（原子操作）
    mu     sync.Mutex    // 并发保护
}
```

### `diskStorage` 结构体（未导出）

```go
type diskStorage struct {
    file     *os.File    // 磁盘文件句柄
    filePath string      // 文件路径
    size     int64       // 文件大小
    closed   int32       // 关闭状态（原子操作）
    mu       sync.Mutex  // 并发保护
}
```

## 4. 函数详解

### 构造函数

#### `newMemoryStorage(data []byte) *memoryStorage`

- **功能**：从字节切片创建内存存储
- **副作用**：调用 `IncrementMemoryBuffers(size)` 更新内存使用统计

#### `newDiskStorage(data []byte, cachePath string) (*diskStorage, error)`

- **功能**：将字节切片写入磁盘缓存文件
- **流程**：创建文件 → 写入数据 → 重置文件指针
- **错误处理**：写入失败时关闭文件并删除临时文件

#### `newDiskStorageFromReader(reader io.Reader, maxBytes int64, cachePath string) (*diskStorage, error)`

- **功能**：从 Reader 流式创建磁盘存储
- **特殊处理**：
  - 使用 `io.LimitReader` 限制读取大小为 `maxBytes+1`
  - 超过限制时返回 `ErrRequestBodyTooLarge`

### `memoryStorage` 方法

#### `Read(p []byte) (n int, err error)`

- 加锁读取，检查关闭状态后委托给 `bytes.Reader`

#### `Seek(offset int64, whence int) (int64, error)`

- 加锁 Seek，检查关闭状态后委托给 `bytes.Reader`

#### `Close() error`

- 使用 `atomic.CompareAndSwapInt32` 确保只关闭一次
- 关闭时调用 `DecrementMemoryBuffers(m.size)` 更新统计

#### `Bytes() ([]byte, error)`

- 加锁读取，检查关闭状态后返回 `m.data`

#### `Size() int64` / `IsDisk() bool`

- 简单返回存储大小和类型标识

### `diskStorage` 方法

#### `Read(p []byte) (n int, err error)`

- 加锁读取，检查关闭状态后委托给 `os.File`

#### `Seek(offset int64, whence int) (int64, error)`

- 加锁 Seek，检查关闭状态后委托给 `os.File`

#### `Close() error`

- 使用 `atomic.CompareAndSwapInt32` 确保只关闭一次
- 关闭文件、删除临时文件、调用 `DecrementDiskFiles(d.size)`

#### `Bytes() ([]byte, error)`

- 保存当前文件指针位置 → 移动到开头 → 读取全部 → 恢复位置
- 确保 `Bytes()` 调用不影响后续的 `Read`/`Seek` 操作

#### `Size() int64` / `IsDisk() bool`

- 简单返回存储大小和类型标识（返回 `true`）

### 工厂函数

#### `CreateBodyStorage(data []byte) (BodyStorage, error)`

- **功能**：根据数据大小和系统配置自动选择存储方式
- **决策逻辑**：
  1. 如果磁盘缓存启用 且 数据大小 ≥ 阈值 且 磁盘空间足够 → 使用磁盘存储
  2. 磁盘创建失败时回退到内存存储（记录错误日志）
  3. 其他情况使用内存存储

#### `CreateBodyStorageFromReader(reader io.Reader, contentLength int64, maxBytes int64) (BodyStorage, error)`

- **功能**：从 Reader 创建存储，支持流式处理大请求
- **决策逻辑**：
  1. 如果磁盘缓存启用 且 `contentLength` ≥ 阈值 且 磁盘空间足够 → 使用磁盘存储
  2. 否则使用内存读取
- **注意**：磁盘存储失败时**无法回退**，因为 Reader 数据已被消费

#### `ReaderOnly(r io.Reader) io.Reader`

- **功能**：包装 Reader 隐藏 `io.Closer` 接口
- **目的**：防止 `http.NewRequest` 类型断言 `io.ReadCloser` 后关闭底层 `BodyStorage`

### 清理函数

#### `CleanupOldCacheFiles()`

- **功能**：清理超过 5 分钟的旧缓存文件
- **调用时机**：应用启动时

## 5. 关键逻辑分析

1. **并发安全设计**：两种存储实现都使用 `sync.Mutex` 保护所有读写操作，使用 `atomic.CompareAndSwapInt32` 确保 `Close()` 只执行一次，避免重复释放资源。

2. **关闭状态检查**：所有读写操作在执行前都检查 `closed` 状态，返回 `ErrStorageClosed` 错误，防止对已关闭存储的操作。

3. **内存管理**：通过 `IncrementMemoryBuffers`/`DecrementMemoryBuffers` 跟踪内存使用量，为系统级内存监控提供数据。

4. **磁盘存储优势**：
   - 大请求体不占用内存
   - `Close()` 时自动清理临时文件
   - 支持流式写入（`newDiskStorageFromReader`）

5. **Reader 消费问题**：`CreateBodyStorageFromReader` 中，磁盘存储路径直接消费 Reader，如果失败则无法回退到内存存储，因为数据已丢失。内存路径通过 `io.ReadAll` 读取全部数据后再创建存储。

6. **ReaderOnly 包装**：`ReaderOnly` 函数通过匿名结构体嵌入 `io.Reader`，巧妙地隐藏了 `Close()` 方法，防止 HTTP 客户端意外关闭底层存储。

7. **缓存清理**：`CleanupOldCacheFiles` 通过统一的 `CleanupOldDiskCacheFiles` 函数清理超过 5 分钟的残留缓存文件，用于启动时的清理工作。

## 6. 关联文件

| 文件 | 关联说明 |
|------|---------|
| `common/disk_cache.go` | 磁盘缓存管理（`CreateDiskCacheFile`、`GetDiskCacheThresholdBytes` 等） |
| `common/metrics.go` | 内存/磁盘使用量统计（`IncrementMemoryBuffers` 等） |
| `middleware/distributor.go` | 分发中间件中可能使用 BodyStorage 进行请求重试 |
| `relay/relay.go` | 中继层可能使用 BodyStorage 存储请求体 |
| `common/errors.go` | 定义 `ErrRequestBodyTooLarge` 错误 |
