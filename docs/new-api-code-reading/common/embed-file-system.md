# embed-file-system.go 代码阅读文档

## 1. 全局总结

本文件实现了基于 Go `embed` 包的静态文件服务系统，用于将前端资源嵌入到二进制文件中。提供了两个核心实现：`embedFileSystem`（单主题嵌入文件系统）和 `themeAwareFileSystem`（支持运行时主题切换的文件系统）。这两个实现都符合 `gin-contrib/static` 包的 `ServeFileSystem` 接口。

## 2. 依赖关系

| 包名 | 用途 |
|------|------|
| `embed` | Go 1.16+ 的文件嵌入功能 |
| `io/fs` | 文件系统子目录提取 |
| `net/http` | HTTP 文件系统接口（`http.FileSystem`） |
| `os` | 文件系统错误（`os.ErrNotExist`） |
| `github.com/gin-contrib/static` | Gin 静态文件服务接口 |

**被依赖方**：本文件被项目启动时的静态文件服务初始化代码调用。

## 3. 类型定义

### `embedFileSystem` 结构体

```go
type embedFileSystem struct {
    http.FileSystem
}
```

嵌入文件系统的包装器，实现了 `static.ServeFileSystem` 接口。

| 字段 | 类型 | 说明 |
|------|------|------|
| `http.FileSystem` | 嵌入字段 | Go 标准库的文件系统接口 |

### `themeAwareFileSystem` 结构体

```go
type themeAwareFileSystem struct {
    defaultFS static.ServeFileSystem
    classicFS static.ServeFileSystem
}
```

主题感知的文件系统，根据当前主题动态选择对应的文件系统。

| 字段 | 类型 | 说明 |
|------|------|------|
| `defaultFS` | `static.ServeFileSystem` | 默认主题的文件系统 |
| `classicFS` | `static.ServeFileSystem` | 经典主题的文件系统 |

## 4. 函数详解

### `(*embedFileSystem) Exists(prefix string, path string) bool`

**功能**：检查请求的路径是否存在。

**参数**：
- `prefix string` — URL 前缀
- `path string` — 请求路径

**返回值**：`bool` — 文件是否存在

**逻辑**：尝试打开文件，如果成功则存在，否则不存在。

### `(*embedFileSystem) Open(name string) (http.File, error)`

**功能**：打开指定路径的文件。

**参数**：
- `name string` — 文件路径

**返回值**：
- `http.File` — 文件对象
- `error` — 错误信息

**逻辑**：
1. 如果路径为 `"/"`（根路径），返回 `os.ErrNotExist`
2. 这样做的目的是让根路径请求走到 Gin 的 `NoRouter` 处理器，使用包含分析代码的替换版 index 页面
3. 其他路径正常返回嵌入的文件

### `EmbedFolder(fsEmbed embed.FS, targetPath string) static.ServeFileSystem`

**功能**：将嵌入的文件系统转换为 Gin 可用的静态文件服务。

**参数**：
- `fsEmbed embed.FS` — Go embed 生成的文件系统
- `targetPath string` — 子目录路径

**返回值**：`static.ServeFileSystem` — Gin 静态文件服务接口

**逻辑**：
1. 使用 `fs.Sub` 提取子目录
2. 使用 `http.FS` 包装为 HTTP 文件系统
3. 包装为 `embedFileSystem` 返回
4. 如果子目录不存在则 panic

### `(*themeAwareFileSystem) Exists(prefix string, path string) bool`

**功能**：根据当前主题检查路径是否存在。

**逻辑**：调用 `GetTheme()` 获取当前主题，委托给对应的文件系统。

### `(*themeAwareFileSystem) Open(name string) (http.File, error)`

**功能**：根据当前主题打开文件。

**逻辑**：调用 `GetTheme()` 获取当前主题，委托给对应的文件系统。

### `NewThemeAwareFS(defaultFS, classicFS static.ServeFileSystem) static.ServeFileSystem`

**功能**：创建主题感知的文件系统。

**参数**：
- `defaultFS` — 默认主题文件系统
- `classicFS` — 经典主题文件系统

**返回值**：`static.ServeFileSystem` — 主题感知的文件系统实例

## 5. 关键逻辑分析

### 根路径处理

`embedFileSystem.Open("/")` 返回 `os.ErrNotExist`，这是一个巧妙的设计：
- Gin 的静态文件服务会检查文件是否存在
- 根路径不存在时，请求会被转发到 `NoRouter` 处理器
- `NoRouter` 处理器返回预处理过的 index.html（包含分析代码）
- 这样可以动态注入统计代码，而不是使用嵌入的静态 index.html

### 主题切换机制

`themeAwareFileSystem` 通过 `GetTheme()` 函数实现运行时主题切换：
- 每次请求都会检查当前主题
- 无需重启服务器即可切换主题
- 支持 `"default"` 和 `"classic"` 两种主题

### 嵌入文件系统架构

```
embed.FS (Go 嵌入)
    ↓
fs.Sub (提取子目录)
    ↓
http.FS (转换为 HTTP 文件系统)
    ↓
embedFileSystem (包装为 Gin 接口)
    ↓
themeAwareFileSystem (多主题支持，可选)
```

## 6. 关联文件

| 文件 | 关联关系 |
|------|----------|
| `web/default/` | 默认主题的前端资源 |
| `web/classic/` | 经典主题的前端资源 |
| `router/` | 初始化静态文件服务 |
| `common/setting.go` | `GetTheme()` 函数 |
