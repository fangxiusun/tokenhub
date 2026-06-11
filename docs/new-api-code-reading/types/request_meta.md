# request_meta.go 代码阅读文档

## 1. 全局概述

本文件定义了请求元数据相关的类型，包括文件类型、Token 类型、Token 计数元数据、文件元数据和请求元数据。这些类型用于在请求处理过程中传递和存储请求的上下文信息。

## 2. 依赖关系

本文件无外部依赖（但引用了同包的 `FileSource` 接口）。

## 3. 类型定义

### FileType 类型

```go
type FileType string
```

文件类型标识符。

### TokenType 类型

```go
type TokenType string
```

Token 类型标识符。

### TokenCountMeta 结构体

```go
type TokenCountMeta struct {
    TokenType       TokenType   `json:"token_type,omitempty"`
    CombineText     string      `json:"combine_text,omitempty"`
    ToolsCount      int         `json:"tools_count,omitempty"`
    NameCount       int         `json:"name_count,omitempty"`
    MessagesCount   int         `json:"messages_count,omitempty"`
    Files           []*FileMeta `json:"files,omitempty"`
    MaxTokens       int         `json:"max_tokens,omitempty"`
    ImagePriceRatio float64     `json:"image_ratio,omitempty"`
}
```

### FileMeta 结构体

```go
type FileMeta struct {
    FileType
    Source FileSource // 统一的文件来源（URL 或 base64）
    Detail string     // 图片细节级别（low/high/auto）
}
```

### RequestMeta 结构体

```go
type RequestMeta struct {
    OriginalModelName string `json:"original_model_name"`
    UserUsingGroup    string `json:"user_using_group"`
    PromptTokens      int    `json:"prompt_tokens"`
    PreConsumedQuota  int    `json:"pre_consumed_quota"`
}
```

## 4. 函数详情

### NewFileMeta

```go
func NewFileMeta(fileType FileType, source FileSource) *FileMeta
```

创建新的 FileMeta 实例。

### NewImageFileMeta

```go
func NewImageFileMeta(source FileSource, detail string) *FileMeta
```

创建图片类型的 FileMeta 实例，包含细节级别。

### GetIdentifier

```go
func (f *FileMeta) GetIdentifier() string
```

获取文件标识符，用于日志输出。

### IsURL

```go
func (f *FileMeta) IsURL() bool
```

判断文件来源是否为 URL。

### GetRawData

```go
func (f *FileMeta) GetRawData() string
```

获取原始数据（已废弃，请使用 `Source.GetRawData()`）。

## 5. 关键逻辑分析

### 文件类型常量

| 常量名 | 值 | 说明 |
|--------|-----|------|
| `FileTypeImage` | `"image"` | 图像文件 |
| `FileTypeAudio` | `"audio"` | 音频文件 |
| `FileTypeVideo` | `"video"` | 视频文件 |
| `FileTypeFile` | `"file"` | 通用文件 |

### Token 类型常量

| 常量名 | 值 | 说明 |
|--------|-----|------|
| `TokenTypeTextNumber` | `"text_number"` | 文本或数字 Token |
| `TokenTypeTokenizer` | `"tokenizer"` | Tokenizer Token |
| `TokenTypeImage` | `"image"` | 图像 Token |

### TokenCountMeta 的作用

`TokenCountMeta` 用于存储请求的 Token 计数元数据：
- `TokenType` — Token 计数方式（文本计数、Tokenizer 计数、图像计数）
- `CombineText` — 所有消息合并后的文本
- `ToolsCount` / `NameCount` / `MessagesCount` — 请求中的工具、名称、消息数量
- `Files` — 请求中的文件列表
- `ImagePriceRatio` — 图像价格比率

### 设计说明

- `FileMeta` 通过嵌入 `FileType` 实现类型继承
- 使用 `FileSource` 接口抽象文件来源，支持 URL 和 Base64 两种方式
- `RequestMeta` 存储请求的全局元数据，用于 Token 计数和配额预扣

## 6. 相关文件

- `types/file_source.go` — FileSource 接口和实现
- `types/file_data.go` — 本地文件数据结构
- `relay/` — 中继层使用这些元数据进行 Token 计数
- `middleware/` — 中间件设置请求元数据
