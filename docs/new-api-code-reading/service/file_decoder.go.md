# file_decoder.go 代码阅读文档

## 1. 全局总结

该文件提供文件类型检测和 MIME 类型识别功能。支持从 URL 获取文件类型，以及通过扩展名、Content-Type 头、内容嗅探等多种方式识别文件格式。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | 日志 |
| `logger` | 调试日志 |
| `types` | 文件数据类型 |
| `http` | Content-Type 检测 |
| `image` | 图片格式解码 |

## 3. 函数详解

### `GetFileTypeFromUrl(c, url, reason...) (string, error)`
文件类型检测的多层策略：
1. HTTP Content-Type 头
2. Content-Disposition 文件名
3. URL 路径扩展名
4. 内容嗅探（http.DetectContentType）
5. 图片解码（image.DecodeConfig）
6. HEIF/HEIC 检测

### `GetFileBase64FromUrl(c, url, reason...) (*types.LocalFileData, error)`
已废弃，内部调用 `LoadFileSource` 实现

### `GetMimeTypeByExtension(ext) string`
扩展名到 MIME 类型的映射表

## 5. 关键逻辑分析

1. **渐进式检测**：从 512 字节到 64KB 逐步增加检测数据量
2. **HEIF/HEIC 支持**：Go 标准库不支持，需要特殊处理
3. **多源融合**：结合 HTTP 头、URL、内容多种信息源

## 6. 关联文件

- `file_service.go` — 统一文件服务
- `image.go` — 图片处理
