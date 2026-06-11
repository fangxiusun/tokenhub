# file_source.go 代码阅读文档

## 1. 全局概述

本文件实现了统一的文件来源抽象层，支持 URL 和 Base64 两种文件来源方式。提供了懒加载、缓存机制（内存和磁盘两种模式）以及自动清理功能。

## 2. 依赖关系

- `fmt` — Go 标准库格式化包
- `image` — Go 标准库图像包（用于图片配置）
- `os` — Go 标准库操作系统包（用于磁盘文件操作）
- `strings` — Go 标准库字符串包
- `sync` — Go 标准库同步包

## 3. 类型定义

### FileSource 接口

```go
type FileSource interface {
    IsURL() bool
    GetIdentifier() string
    GetRawData() string
    ClearRawData()
    SetCache(data *CachedFileData)
    GetCache() *CachedFileData
    HasCache() bool
    ClearCache()
    IsRegistered() bool
    SetRegistered(registered bool)
    Mu() *sync.Mutex
}
```

### baseFileSource 结构体

```go
type baseFileSource struct {
    cachedData  *CachedFileData
    cacheLoaded bool
    registered  bool
    mu          sync.Mutex
}
```

共享的基类实现，提供缓存、锁和注册状态管理。

### URLSource 结构体

```go
type URLSource struct {
    baseFileSource
    URL string
}
```

URL 来源的 FileSource 实现。

### Base64Source 结构体

```go
type Base64Source struct {
    baseFileSource
    Base64Data string
    MimeType   string
}
```

Base64 内联数据来源的 FileSource 实现。

### CachedFileData 结构体

```go
type CachedFileData struct {
    base64Data     string
    MimeType       string
    Size           int64
    DiskSize       int64
    ImageConfig    *image.Config
    ImageFormat    string
    diskPath       string
    isDisk         bool
    diskMu         sync.Mutex
    diskClosed     bool
    statDecremented bool
    OnClose        func(size int64)
}
```

缓存的文件数据，支持内存和磁盘两种模式。

## 4. 函数详情

### 构造函数

- `NewURLFileSource(url string) *URLSource` — 创建 URL 来源
- `NewBase64FileSource(base64Data string, mimeType string) *Base64Source` — 创建 Base64 来源
- `NewFileSourceFromData(data string, mimeType string) FileSource` — 根据数据内容自动选择来源类型

### CachedFileData 方法

- `NewMemoryCachedData(base64Data string, mimeType string, size int64) *CachedFileData` — 创建内存缓存
- `NewDiskCachedData(diskPath string, mimeType string, size int64) *CachedFileData` — 创建磁盘缓存
- `GetBase64Data() (string, error)` — 获取 Base64 数据（磁盘模式会从文件读取）
- `SetBase64Data(data string)` — 设置 Base64 数据（仅内存模式）
- `IsDisk() bool` — 是否使用磁盘缓存
- `Close() error` — 关闭缓存并清理资源

## 5. 关键逻辑分析

### 文件来源抽象

`FileSource` 接口统一了 URL 和 Base64 两种文件来源：
- `IsURL()` — 判断来源类型
- `GetRawData()` — 获取原始数据（URL 返回 URL 字符串，Base64 返回编码数据）
- `ClearRawData()` — 清理原始数据（Base64 模式下大于 1KB 时清空以释放内存）

### 缓存机制

#### 内存缓存
- 小文件直接存储在内存中
- 通过 `base64Data` 字段存储
- 访问速度快，但占用内存

#### 磁盘缓存
- 大文件存储在磁盘临时文件中
- 通过 `diskPath` 字段存储文件路径
- 读取时从磁盘加载，减少内存占用
- 关闭时自动删除临时文件

### 自动清理

- `OnClose` 回调函数在缓存关闭时被调用
- `statDecremented` 标记防止重复扣减统计
- 磁盘文件在 `Close()` 时自动删除

### 线程安全

- `baseFileSource` 使用 `sync.Mutex` 保护缓存操作
- `CachedFileData` 的磁盘操作使用独立的 `diskMu` 锁
- `registered` 状态用于跟踪文件是否已注册到清理列表

## 6. 相关文件

- `types/request_meta.go` — FileMeta 使用 FileSource 接口
- `types/file_data.go` — LocalFileData 简单文件数据结构
- `relay/` — 中继层使用 FileSource 处理文件上传
- `middleware/` — 中间件管理文件源的注册和清理
