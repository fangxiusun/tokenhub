# file_service.go 代码阅读文档

## 1. 全局总结

该文件是统一的文件处理服务，提供文件下载、解码、缓存等功能的统一入口。支持 URL 和 Base64 两种文件源，实现了内存缓存和磁盘缓存的混合策略。

## 2. 依赖关系

| 依赖包 | 用途 |
|--------|------|
| `common` | 缓存操作、日志 |
| `constant` | 上下文键名、文件大小限制 |
| `types` | FileSource、CachedFileData |
| `webp` | WebP 图片解码 |

## 3. 类型定义

无自定义类型，使用 `types` 包中的接口和结构体。

## 4. 核心函数

### `LoadFileSource(c, source, reason...) (*types.CachedFileData, error)`
统一文件加载入口：
1. 快速检查内部缓存
2. 加锁保护加载过程
3. 双重检查缓存
4. 根据来源类型加载（URL/Base64）
5. 设置缓存
6. 注册清理回调

### `GetImageConfig(c, source) (image.Config, string, error)`
获取图片配置（尺寸、格式）

### `GetBase64Data(c, source, reason...) (string, string, error)`
获取 Base64 编码数据

### `GetMimeType(c, source) (string, error)`
获取文件 MIME 类型

### `DetectFileType(mimeType) types.FileType`
根据 MIME 类型检测文件类型（图片/音频/视频/文件）

## 5. 关键逻辑分析

1. **三级缓存**：Source 内部缓存 → gin.Context 缓存 → 磁盘/内存缓存
2. **懒加载**：首次访问时加载，后续直接使用缓存
3. **自动清理**：请求结束时自动清理 FileSource
4. **智能存储**：大文件使用磁盘缓存，小文件使用内存缓存
5. **HEIF/HEIC 支持**：解析 ISOBMFF box 获取图片尺寸

## 6. 关联文件

- `file_decoder.go` — 文件类型检测
- `download.go` — 文件下载
- `types/file_source.go` — FileSource 接口定义
