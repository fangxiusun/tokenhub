# disk_cache.go 代码阅读文档

## 1. 全局总结

该文件实现了磁盘缓存的文件操作功能，提供了缓存文件的创建、读取、写入、删除和清理等核心操作。文件是磁盘缓存子系统的文件 I/O 层，与 `disk_cache_config.go` 配合，实现了完整的磁盘缓存机制。该缓存主要用于存储大体积的请求体数据，减少内存压力。

## 2. 依赖关系

| 包名 | 用途 |
|------|------|
| `fmt` | 格式化字符串和错误信息 |
| `os` | 文件系统操作（创建、读取、删除文件和目录） |
| `path/filepath` | 跨平台路径处理 |
| `time` | 时间计算和文件修改时间判断 |
| `github.com/google/uuid` | 生成唯一文件名 |

内部依赖：
- `GetDiskCachePath()`：获取缓存目录配置
- `IsDiskCacheEnabled()`：检查缓存是否启用
- `GetDiskCacheThresholdBytes()`：获取缓存触发阈值
- `IsDiskCacheAvailable()`：检查缓存空间是否足够
- `DecrementDiskFiles()`：更新缓存统计

## 3. 类型定义

### 3.1 DiskCacheType 类型

```go
type DiskCacheType string

const (
    DiskCacheTypeBody DiskCacheType = "body"
    DiskCacheTypeFile DiskCacheType = "file"
)
```

**用途**：定义磁盘缓存的类型分类，用于文件命名和管理。

**类型说明**：
- `DiskCacheTypeBody`: 请求体缓存，用于存储 AI API 请求的原始数据
- `DiskCacheTypeFile`: 文件数据缓存，用于存储上传的文件内容

## 4. 常量定义

### 4.1 diskCacheDir

```go
const diskCacheDir = "new-api-body-cache"
```

**用途**：定义缓存子目录名称，在配置的缓存路径下创建。

**设计考虑**：
- 使用固定名称便于识别和管理
- 与配置路径组合形成完整的缓存目录

## 5. 函数详解

### 5.1 目录管理函数

#### GetDiskCacheDir() string

**功能**：获取完整的磁盘缓存目录路径

**实现逻辑**：
1. 调用 `GetDiskCachePath()` 获取配置的缓存路径
2. 如果路径为空，使用系统临时目录 `os.TempDir()`
3. 组合缓存路径和子目录名称 `new-api-body-cache`

**返回值**：完整的缓存目录路径

#### EnsureDiskCacheDir() error

**功能**：确保缓存目录存在，不存在则创建

**实现**：调用 `os.MkdirAll` 创建目录（权限 0755）

### 5.2 文件操作函数

#### CreateDiskCacheFile(cacheType DiskCacheType) (string, *os.File, error)

**功能**：创建新的磁盘缓存文件

**参数**：
- `cacheType`: 缓存类型（body 或 file）

**返回值**：
- `string`: 文件完整路径
- `*os.File`: 文件句柄
- `error`: 创建过程中的错误

**实现逻辑**：
1. 调用 `EnsureDiskCacheDir` 确保目录存在
2. 生成唯一文件名：`{类型}-{UUID前8位}-{纳秒时间戳}.tmp`
3. 使用 `os.OpenFile` 创建文件（标志：`O_CREATE|O_RDWR|O_EXCL`，权限：`0600`）
4. 返回文件路径和句柄

**文件命名规则**：
- 包含缓存类型便于识别
- UUID 前 8 位保证唯一性
- 纳秒时间戳提供排序依据
- `.tmp` 后缀表示临时文件

#### WriteDiskCacheFile(cacheType DiskCacheType, data []byte) (string, error)

**功能**：将字节数据写入磁盘缓存文件

**参数**：
- `cacheType`: 缓存类型
- `data`: 要写入的数据

**返回值**：
- `string`: 文件完整路径
- `error`: 写入过程中的错误

**实现逻辑**：
1. 调用 `CreateDiskCacheFile` 创建文件
2. 调用 `file.Write` 写入数据
3. 如果写入失败，关闭文件并删除
4. 关闭文件，如果关闭失败则删除文件
5. 返回文件路径

**错误处理**：
- 写入失败时清理已创建的文件
- 关闭失败时清理已写入的文件
- 确保不会留下部分写入的损坏文件

#### WriteDiskCacheFileString(cacheType DiskCacheType, data string) (string, error)

**功能**：将字符串数据写入磁盘缓存文件

**参数**：
- `cacheType`: 缓存类型
- `data`: 要写入的字符串数据

**返回值**：
- `string`: 文件完整路径
- `error`: 写入过程中的错误

**实现**：调用 `WriteDiskCacheFile`，将字符串转换为字节数组

#### ReadDiskCacheFile(filePath string) ([]byte, error)

**功能**：读取磁盘缓存文件内容

**参数**：
- `filePath`: 文件完整路径

**返回值**：
- `[]byte`: 文件内容字节数组
- `error`: 读取过程中的错误

**实现**：调用 `os.ReadFile` 读取整个文件

#### ReadDiskCacheFileString(filePath string) (string, error)

**功能**：读取磁盘缓存文件内容为字符串

**参数**：
- `filePath`: 文件完整路径

**返回值**：
- `string`: 文件内容字符串
- `error`: 读取过程中的错误

**实现**：
1. 调用 `os.ReadFile` 读取文件
2. 将字节数组转换为字符串

#### RemoveDiskCacheFile(filePath string) error

**功能**：删除磁盘缓存文件

**参数**：
- `filePath`: 文件完整路径

**返回值**：删除过程中的错误

**实现**：调用 `os.Remove` 删除文件

### 5.3 清理和维护函数

#### CleanupOldDiskCacheFiles(maxAge time.Duration) error

**功能**：清理超过指定存活时间的旧缓存文件

**参数**：
- `maxAge`: 文件最大存活时间

**返回值**：清理过程中的错误

**实现逻辑**：
1. 获取缓存目录路径
2. 读取目录内容
3. 遍历所有文件条目
4. 获取文件修改时间
5. 如果 `当前时间 - 修改时间 > maxAge`，删除文件
6. 删除成功后调用 `DecrementDiskFiles` 更新统计

**注意事项**：
- 目录不存在时返回 nil（无需清理）
- 跳过子目录条目
- 获取文件信息失败时跳过该文件
- 删除失败时继续处理其他文件

#### GetDiskCacheInfo() (fileCount int, totalSize int64, err error)

**功能**：获取缓存目录的文件数量和总大小

**返回值**：
- `fileCount`: 文件数量
- `totalSize`: 总大小（字节）
- `error`: 获取过程中的错误

**实现逻辑**：
1. 获取缓存目录路径
2. 读取目录内容
3. 遍历所有文件条目
4. 统计文件数量和总大小

### 5.4 判断函数

#### ShouldUseDiskCache(dataSize int64) bool

**功能**：判断给定大小的数据是否应该使用磁盘缓存

**参数**：
- `dataSize`: 数据大小（字节）

**返回值**：`true` 表示应该使用磁盘缓存

**实现逻辑**：
1. 检查磁盘缓存是否启用
2. 获取缓存触发阈值
3. 如果数据大小小于阈值，返回 false
4. 检查缓存空间是否足够

**判断条件**：
- 磁盘缓存已启用
- 数据大小 >= 阈值
- 缓存空间充足

## 6. 关键逻辑分析

### 6.1 文件命名策略

使用 `{类型}-{UUID前8位}-{纳秒时间戳}.tmp` 的命名策略：
- **类型前缀**：便于区分不同类型的缓存文件
- **UUID**：保证文件名唯一性，避免并发冲突
- **纳秒时间戳**：提供时间排序依据，便于清理旧文件
- **.tmp 后缀**：标识临时文件，避免与永久文件混淆

### 6.2 原子文件创建

使用 `os.O_EXCL` 标志确保文件创建的原子性：
- 如果文件已存在，立即返回错误
- 避免并发创建同名文件时的竞争条件
- 结合 UUID 和时间戳，几乎不可能发生文件名冲突

### 6.3 错误恢复机制

`WriteDiskCacheFile` 函数实现了完整的错误恢复：
1. 创建文件失败：直接返回错误
2. 写入失败：关闭文件并删除
3. 关闭失败：删除已写入的文件
4. 成功：返回文件路径

这种设计确保：
- 不会留下部分写入的损坏文件
- 不会泄露磁盘空间
- 每个操作都是幂等的

### 6.4 清理策略

`CleanupOldDiskCacheFiles` 实现了基于时间的清理：
- 使用文件修改时间判断文件年龄
- 只清理超过 `maxAge` 的文件
- 清理后更新统计信息
- 跳过删除失败的文件（可能被其他进程使用）

### 6.5 缓存决策逻辑

`ShouldUseDiskCache` 函数综合了多个条件：
1. **功能开关**：磁盘缓存必须启用
2. **大小阈值**：小数据不值得使用磁盘缓存
3. **空间限制**：防止缓存无限增长

这种多层判断确保：
- 只对大请求体使用磁盘缓存
- 避免频繁的小文件 I/O
- 控制磁盘使用量

### 6.6 与配置模块的协作

该文件与 `disk_cache_config.go` 紧密协作：
- 使用配置函数获取缓存参数
- 更新统计信息供配置模块查询
- 遵循配置模块定义的限制和阈值

## 7. 关联文件

- `disk_cache_config.go`：配置管理和统计信息
- `setting/performance_setting.go`：磁盘缓存配置的持久化
- `relay/` 目录：使用磁盘缓存处理大请求体
- `controller/` 目录：通过 API 管理缓存（清理、统计）
