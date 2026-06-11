# custom-event.go 代码阅读文档

## 1. 全局总结

该文件实现了 Server-Sent Events (SSE) 的自定义事件渲染功能，源自 Gin 框架的扩展。文件定义了 `CustomEvent` 结构体和相关的编码/渲染逻辑，用于在 HTTP 连接中推送实时事件流。该实现遵循 W3C SSE 规范（2009年工作草案）。

## 2. 依赖关系

| 包名 | 用途 |
|------|------|
| `fmt` | 格式化数据输出 |
| `io` | I/O 接口定义 |
| `net/http` | HTTP 响应处理 |
| `strings` | 字符串替换处理 |
| `sync` | 互斥锁同步 |

## 3. 类型定义

### 3.1 stringWriter 接口

```go
type stringWriter interface {
    io.Writer
    writeString(string) (int, error)
}
```

**用途**：扩展 `io.Writer` 接口，增加 `writeString` 方法，用于高效地写入字符串而无需转换为字节数组。

### 3.2 stringWrapper 结构体

```go
type stringWrapper struct {
    io.Writer
}
```

**用途**：适配器模式，将普通 `io.Writer` 包装为 `stringWriter` 接口，使不支持原生 `writeString` 的写入器也能使用字符串写入优化。

**方法**：
- `writeString(str string) (int, error)`：将字符串转换为字节数组后调用底层 `Write` 方法

### 3.3 CustomEvent 结构体

```go
type CustomEvent struct {
    Event string
    Id    string
    Retry uint
    Data  interface{}
    Mutex sync.Mutex
}
```

**用途**：表示一个 SSE 事件，包含事件的所有必要字段。

**字段说明**：
- `Event`: 事件类型名称（对应 SSE 的 `event:` 字段）
- `Id`: 事件 ID（对应 SSE 的 `id:` 字段）
- `Retry`: 重连间隔（对应 SSE 的 `retry:` 字段，毫秒）
- `Data`: 事件数据（对应 SSE 的 `data:` 字段）
- `Mutex`: 互斥锁，用于并发写入时的线程安全

## 4. 函数详解

### 4.1 checkWriter(writer io.Writer) stringWriter

**功能**：检查写入器是否支持 `stringWriter` 接口，如果不支持则包装为适配器

**参数**：
- `writer`: 任意 `io.Writer` 实现

**返回值**：`stringWriter` 接口实现

**实现逻辑**：
1. 尝试类型断言将 `writer` 转换为 `stringWriter`
2. 如果成功，直接返回原写入器
3. 如果失败，用 `stringWrapper` 包装后返回

### 4.2 encode(writer io.Writer, event CustomEvent) error

**功能**：将自定义事件编码并写入到输出流

**参数**：
- `writer`: 输出写入器
- `event`: 要编码的自定义事件

**返回值**：编码过程中的错误

**实现逻辑**：
1. 调用 `checkWriter` 获取适配后的写入器
2. 调用 `writeData` 写入事件数据

### 4.3 writeData(w stringWriter, data interface{}) error

**功能**：将事件数据写入输出流，处理特殊字符

**参数**：
- `w`: 字符串写入器
- `data`: 事件数据

**返回值**：写入过程中的错误

**实现逻辑**：
1. 使用 `dataReplacer` 替换数据中的特殊字符（`\r` 替换为 `\\r`，`\n` 保留）
2. 如果数据以 `"data"` 开头，追加双换行符分隔

### 4.4 (r CustomEvent) Render(w http.ResponseWriter) error

**功能**：将自定义事件渲染到 HTTP 响应

**参数**：
- `w`: HTTP 响应写入器

**返回值**：渲染过程中的错误

**实现逻辑**：
1. 调用 `WriteContentType` 设置响应头
2. 调用 `encode` 编码并写入事件数据

### 4.5 (r CustomEvent) WriteContentType(w http.ResponseWriter)

**功能**：设置 SSE 响应的 Content-Type 头

**参数**：
- `w`: HTTP 响应写入器

**实现逻辑**：
1. 加锁保证并发安全
2. 设置 `Content-Type` 为 `text/event-stream`
3. 如果未设置 `Cache-Control`，则设置为 `no-cache`

## 5. 关键逻辑分析

### 5.1 SSE 协议实现

该文件实现了 SSE 协议的核心部分：
- 使用 `text/event-stream` 内容类型
- 支持事件类型（`Event`）、事件 ID（`Id`）、重连间隔（`Retry`）和数据（`Data`）
- 通过双换行符（`\n\n`）分隔事件

### 5.2 字符串写入优化

通过 `stringWriter` 接口和 `stringWrapper` 适配器：
- 避免每次写入字符串时都进行字节数组转换
- 对于原生支持 `writeString` 的写入器（如 `bytes.Buffer`），使用更高效的路径

### 5.3 特殊字符处理

使用两个 `strings.NewReplacer`：
- `fieldReplacer`: 替换 `\n` 和 `\r` 为转义序列（用于字段名）
- `dataReplacer`: 保留 `\n`，只替换 `\r` 为 `\\r`（用于数据内容）

### 5.4 并发安全

`WriteContentType` 方法使用互斥锁保护响应头的写入，确保在并发场景下不会出现数据竞争。

### 5.5 与 Gin 框架的集成

文件中的 `CustomEvent` 结构体实现了类似 Gin 框架的 `Render` 接口，可以方便地在 Gin 路由处理器中使用。

## 6. 关联文件

- `relay/stream.go` 或类似文件：AI API 流式响应处理
- `controller/` 目录：HTTP 请求处理器，使用 SSE 推送实时数据
- `middleware/` 目录：可能使用 SSE 实现实时通知
- `gin/render/` 目录：Gin 框架的渲染接口定义
